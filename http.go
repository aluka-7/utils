package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return e.Message
}

type RemoteError struct {
	Host string
	Err  error
}

func (e *RemoteError) Error() string {
	return e.Err.Error()
}

var CallUserAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/29.0.1541.0 Safari/537.36"

// HttpCall makes HTTP method call.
func HttpCall(client *http.Client, method, url string, header http.Header, body io.Reader) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", CallUserAgent)
	for k, vs := range header {
		req.Header[k] = vs
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 200 {
		return resp.Body, nil
	}
	resp.Body.Close()
	if resp.StatusCode == 404 { // 403 can be rate limit error.  || resp.StatusCode == 403 {
		err = fmt.Errorf("resource not found: %s", url)
	} else {
		err = fmt.Errorf("%s %s -> %d", method, url, resp.StatusCode)
	}
	return nil, err
}

// HttpGet gets the specified resource.
// ErrNotFound is returned if the server responds with status 404.
func HttpGet(client *http.Client, url string, header http.Header) (io.ReadCloser, error) {
	return HttpCall(client, "GET", url, header, nil)
}

// HttpPost posts the specified resource.
// ErrNotFound is returned if the server responds with status 404.
func HttpPost(client *http.Client, url string, header http.Header, body []byte) (io.ReadCloser, error) {
	return HttpCall(client, "POST", url, header, bytes.NewBuffer(body))
}

// HttpGetToFile gets the specified resource and writes to file.
// ErrNotFound is returned if the server responds with status 404.
func HttpGetToFile(client *http.Client, url string, header http.Header, fileName string) error {
	rc, err := HttpGet(client, url, header)
	if err != nil {
		return err
	}
	defer rc.Close()

	os.MkdirAll(path.Dir(fileName), os.ModePerm)
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, rc)
	return err
}

// HttpGetBytes gets the specified resource. ErrNotFound is returned if the server
// responds with status 404.
func HttpGetBytes(client *http.Client, url string, header http.Header) ([]byte, error) {
	rc, err := HttpGet(client, url, header)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return ioutil.ReadAll(rc)
}

// HttpGetJSON gets the specified resource and mapping to struct.
// ErrNotFound is returned if the server responds with status 404.
func HttpGetJSON(client *http.Client, url string, v interface{}) error {
	rc, err := HttpGet(client, url, nil)
	if err != nil {
		return err
	}
	defer rc.Close()
	err = json.NewDecoder(rc).Decode(v)
	if _, ok := err.(*json.SyntaxError); ok {
		return fmt.Errorf("JSON syntax error at %s", url)
	}
	return nil
}

// HttpPostJSON posts the specified resource with struct values,
// and maps results to struct.
// ErrNotFound is returned if the server responds with status 404.
func HttpPostJSON(client *http.Client, url string, body, v interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	rc, err := HttpPost(client, url, http.Header{"content-type": []string{"application/json"}}, data)
	if err != nil {
		return err
	}
	defer rc.Close()
	err = json.NewDecoder(rc).Decode(v)
	if _, ok := err.(*json.SyntaxError); ok {
		return fmt.Errorf("JSON syntax error at %s", url)
	}
	return nil
}

// A RawFile describes a file that can be downloaded.
type RawFile interface {
	Name() string
	RawUrl() string
	Data() []byte
	SetData([]byte)
}

// FetchFiles fetches files specified by the rawURL field in parallel.
func FetchFiles(client *http.Client, files []RawFile, header http.Header) error {
	ch := make(chan error, len(files))
	for i := range files {
		go func(i int) {
			p, err := HttpGetBytes(client, files[i].RawUrl(), nil)
			if err != nil {
				ch <- err
				return
			}
			files[i].SetData(p)
			ch <- nil
		}(i)
	}
	for _ = range files {
		if err := <-ch; err != nil {
			return err
		}
	}
	return nil
}

// FetchFiles uses command `curl` to fetch files specified by the rawURL field in parallel.
func FetchFilesCurl(files []RawFile, curlOptions ...string) error {
	ch := make(chan error, len(files))
	for i := range files {
		go func(i int) {
			stdout, _, err := ExecCmd("curl", append(curlOptions, files[i].RawUrl())...)
			if err != nil {
				ch <- err
				return
			}

			files[i].SetData([]byte(stdout))
			ch <- nil
		}(i)
	}
	for _ = range files {
		if err := <-ch; err != nil {
			return err
		}
	}
	return nil
}

// "section"包含产品名称，其版本和可选注释.
type section struct {
	name    string
	version string
	comment []string
}

// 包含可以从User-Agent字符串中提取的所有信息。
type UserAgent struct {
	ua           string
	mozilla      string
	platform     string
	os           string
	localization string
	browser      Browser
	bot          bool
	mobile       bool
	undecided    bool
}

// 从给定的字符串中读取，直到到达给定的分隔符或字符串的末尾。
//
// 第一个参数是要解析的用户代理字符串。
// 第二个参数是指向用户代理字符串当前索引的引用。
// delimiter参数指定哪个字符是定界符，cat参数确定是否应忽略嵌套的'('。
//
// 返回包含已读取内容的字节数组。
func readUntil(ua string, index *int, delimiter byte, cat bool) []byte {
	var buffer []byte

	i := *index
	catalan := 0
	for ; i < len(ua); i = i + 1 {
		if ua[i] == delimiter {
			if catalan == 0 {
				*index = i + 1
				return buffer
			}
			catalan--
		} else if cat && ua[i] == '(' {
			catalan++
		}
		buffer = append(buffer, ua[i])
	}
	*index = i + 1
	return buffer
}

// Parse the given product, that is, just a name or a string
// formatted as Name/Version.
//
// It returns two strings. The first string is the name of the product and the
// second string contains the version of the product.
func parseProduct(product []byte) (string, string) {
	prod := strings.SplitN(string(product), "/", 2)
	if len(prod) == 2 {
		return prod[0], prod[1]
	}
	return string(product), ""
}

// 解析部分.节的格式通常如下"名称/版本(注释)".注释和版本都是可选的。
//
// 第一个参数是要解析的用户代理字符串。
// 第二个参数是指向用户代理字符串当前索引的引用。
//
// 返回一个节，其中包含我们可以从最后一个已分析节中提取的信息。
func parseSection(ua string, index *int) (s section) {
	buffer := readUntil(ua, index, ' ', false)

	s.name, s.version = parseProduct(buffer)
	if *index < len(ua) && ua[*index] == '(' {
		*index++
		buffer = readUntil(ua, index, ')', true)
		s.comment = strings.Split(string(buffer), "; ")
		*index++
	}
	return s
}

// Initialize the parser.
func (p *UserAgent) initialize() {
	p.ua = ""
	p.mozilla = ""
	p.platform = ""
	p.os = ""
	p.localization = ""
	p.browser.Engine = ""
	p.browser.EngineVersion = ""
	p.browser.Name = ""
	p.browser.Version = ""
	p.bot = false
	p.mobile = false
	p.undecided = false
}

// 解析给定的User-Agent字符串并获取结果UserAgent对象。
//
// 返回解析给定User-Agent字符串后已初始化的UserAgent对象。
func New(ua string) *UserAgent {
	o := &UserAgent{}
	o.Parse(ua)
	return o
}

// Parse the given User-Agent string. After calling this function, the
// receiver will be setted up with all the information that we've extracted.
func (p *UserAgent) Parse(ua string) {
	var sections []section

	p.initialize()
	p.ua = ua
	for index, limit := 0, len(ua); index < limit; {
		s := parseSection(ua, &index)
		if !p.mobile && s.name == "Mobile" {
			p.mobile = true
		}
		sections = append(sections, s)
	}

	if len(sections) > 0 {
		if sections[0].name == "Mozilla" {
			p.mozilla = sections[0].version
		}

		p.detectBrowser(sections)
		p.detectOS(sections[0])

		if p.undecided {
			p.checkBot(sections)
		}
	}
}

// Returns the mozilla version (it's how the User Agent string begins:
// "Mozilla/5.0 ...", unless we're dealing with Opera, of course).
func (p *UserAgent) Mozilla() string {
	return p.mozilla
}

// Returns true if it's a bot, false otherwise.
func (p *UserAgent) Bot() bool {
	return p.bot
}

// Returns true if it's a mobile device, false otherwise.
func (p *UserAgent) Mobile() bool {
	return p.mobile
}

// Returns the original given user agent.
func (p *UserAgent) UA() string {
	return p.ua
}

// Represents full information on the operating system extracted from the user agent.
type OSInfo struct {
	// Full name of the operating system. This is identical to the output of ua.OS()
	FullName string

	// Name of the operating system. This is sometimes a shorter version of the
	// operating system name, e.g. "Mac OS X" instead of "Intel Mac OS X"
	Name string

	// Operating system version, e.g. 7 for Windows 7 or 10.8 for Max OS X Mountain Lion
	Version string
}

// Normalize the name of the operating system. By now, this just
// affects to Windows NT.
//
// Returns a string containing the normalized name for the Operating System.
func normalizeOS(name string) string {
	sp := strings.SplitN(name, " ", 3)
	if len(sp) != 3 || sp[1] != "NT" {
		return name
	}

	switch sp[2] {
	case "5.0":
		return "Windows 2000"
	case "5.01":
		return "Windows 2000, Service Pack 1 (SP1)"
	case "5.1":
		return "Windows XP"
	case "5.2":
		return "Windows XP x64 Edition"
	case "6.0":
		return "Windows Vista"
	case "6.1":
		return "Windows 7"
	case "6.2":
		return "Windows 8"
	case "6.3":
		return "Windows 8.1"
	case "10.0":
		return "Windows 10"
	}
	return name
}

// Guess the OS, the localization and if this is a mobile device for a
// Webkit-powered browser.
//
// The first argument p is a reference to the current UserAgent and the second
// argument is a slice of strings containing the comment.
func webkit(p *UserAgent, comment []string) {
	if p.platform == "webOS" {
		p.browser.Name = p.platform
		p.os = "Palm"
		if len(comment) > 2 {
			p.localization = comment[2]
		}
		p.mobile = true
	} else if p.platform == "Symbian" {
		p.mobile = true
		p.browser.Name = p.platform
		p.os = comment[0]
	} else if p.platform == "Linux" {
		p.mobile = true
		if p.browser.Name == "Safari" {
			p.browser.Name = "Android"
		}
		if len(comment) > 1 {
			if comment[1] == "U" {
				if len(comment) > 2 {
					p.os = comment[2]
				} else {
					p.mobile = false
					p.os = comment[0]
				}
			} else {
				p.os = comment[1]
			}
		}
		if len(comment) > 3 {
			p.localization = comment[3]
		} else if len(comment) == 3 {
			_ = p.googleBot()
		}
	} else if len(comment) > 0 {
		if len(comment) > 3 {
			p.localization = comment[3]
		}
		if strings.HasPrefix(comment[0], "Windows NT") {
			p.os = normalizeOS(comment[0])
		} else if len(comment) < 2 {
			p.localization = comment[0]
		} else if len(comment) < 3 {
			if !p.googleBot() {
				p.os = normalizeOS(comment[1])
			}
		} else {
			p.os = normalizeOS(comment[2])
		}
		if p.platform == "BlackBerry" {
			p.browser.Name = p.platform
			if p.os == "Touch" {
				p.os = p.platform
			}
		}
	}
}

// Guess the OS, the localization and if this is a mobile device
// for a Gecko-powered browser.
//
// The first argument p is a reference to the current UserAgent and the second
// argument is a slice of strings containing the comment.
func gecko(p *UserAgent, comment []string) {
	if len(comment) > 1 {
		if comment[1] == "U" {
			if len(comment) > 2 {
				p.os = normalizeOS(comment[2])
			} else {
				p.os = normalizeOS(comment[1])
			}
		} else {
			if p.platform == "Android" {
				p.mobile = true
				p.platform, p.os = normalizeOS(comment[1]), p.platform
			} else if comment[0] == "Mobile" || comment[0] == "Tablet" {
				p.mobile = true
				p.os = "FirefoxOS"
			} else {
				if p.os == "" {
					p.os = normalizeOS(comment[1])
				}
			}
		}
		// Only parse 4th comment as localization if it doesn't start with rv:.
		// For example Firefox on Ubuntu contains "rv:XX.X" in this field.
		if len(comment) > 3 && !strings.HasPrefix(comment[3], "rv:") {
			p.localization = comment[3]
		}
	}
}

// Guess the OS, the localization and if this is a mobile device
// for Internet Explorer.
//
// The first argument p is a reference to the current UserAgent and the second
// argument is a slice of strings containing the comment.
func trident(p *UserAgent, comment []string) {
	// Internet Explorer only runs on Windows.
	p.platform = "Windows"

	// The OS can be set before to handle a new case in IE11.
	if p.os == "" {
		if len(comment) > 2 {
			p.os = normalizeOS(comment[2])
		} else {
			p.os = "Windows NT 4.0"
		}
	}

	// Last but not least, let's detect if it comes from a mobile device.
	for _, v := range comment {
		if strings.HasPrefix(v, "IEMobile") {
			p.mobile = true
			return
		}
	}
}

// Guess the OS, the localization and if this is a mobile device
// for Opera.
//
// The first argument p is a reference to the current UserAgent and the second
// argument is a slice of strings containing the comment.
func opera(p *UserAgent, comment []string) {
	slen := len(comment)

	if strings.HasPrefix(comment[0], "Windows") {
		p.platform = "Windows"
		p.os = normalizeOS(comment[0])
		if slen > 2 {
			if slen > 3 && strings.HasPrefix(comment[2], "MRA") {
				p.localization = comment[3]
			} else {
				p.localization = comment[2]
			}
		}
	} else {
		if strings.HasPrefix(comment[0], "Android") {
			p.mobile = true
		}
		p.platform = comment[0]
		if slen > 1 {
			p.os = comment[1]
			if slen > 3 {
				p.localization = comment[3]
			}
		} else {
			p.os = comment[0]
		}
	}
}

// Guess the OS. Android browsers send Dalvik as the user agent in the
// request header.
//
// The first argument p is a reference to the current UserAgent and the second
// argument is a slice of strings containing the comment.
func dalvik(p *UserAgent, comment []string) {
	slen := len(comment)

	if strings.HasPrefix(comment[0], "Linux") {
		p.platform = comment[0]
		if slen > 2 {
			p.os = comment[2]
		}
		p.mobile = true
	}
}

// Given the comment of the first section of the UserAgent string,
// get the platform.
func getPlatform(comment []string) string {
	if len(comment) > 0 {
		if comment[0] != "compatible" {
			if strings.HasPrefix(comment[0], "Windows") {
				return "Windows"
			} else if strings.HasPrefix(comment[0], "Symbian") {
				return "Symbian"
			} else if strings.HasPrefix(comment[0], "webOS") {
				return "webOS"
			} else if comment[0] == "BB10" {
				return "BlackBerry"
			}
			return comment[0]
		}
	}
	return ""
}

// Detect some properties of the OS from the given section.
func (p *UserAgent) detectOS(s section) {
	if s.name == "Mozilla" {
		// Get the platform here. Be aware that IE11 provides a new format
		// that is not backwards-compatible with previous versions of IE.
		p.platform = getPlatform(s.comment)
		if p.platform == "Windows" && len(s.comment) > 0 {
			p.os = normalizeOS(s.comment[0])
		}

		// And finally get the OS depending on the engine.
		switch p.browser.Engine {
		case "":
			p.undecided = true
		case "Gecko":
			gecko(p, s.comment)
		case "AppleWebKit":
			webkit(p, s.comment)
		case "Trident":
			trident(p, s.comment)
		}
	} else if s.name == "Opera" {
		if len(s.comment) > 0 {
			opera(p, s.comment)
		}
	} else if s.name == "Dalvik" {
		if len(s.comment) > 0 {
			dalvik(p, s.comment)
		}
	} else {
		// Check whether this is a bot or just a weird browser.
		p.undecided = true
	}
}

// Returns a string containing the platform..
func (p *UserAgent) Platform() string {
	return p.platform
}

// Returns a string containing the name of the Operating System.
func (p *UserAgent) OS() string {
	return p.os
}

// Returns a string containing the localization.
func (p *UserAgent) Localization() string {
	return p.localization
}

// Return OS name and version from a slice of strings created from the full name of the OS.
func osName(osSplit []string) (name, version string) {
	if len(osSplit) == 1 {
		name = osSplit[0]
		version = ""
	} else {
		// Assume version is stored in the last part of the array.
		nameSplit := osSplit[:len(osSplit)-1]
		version = osSplit[len(osSplit)-1]

		// Nicer looking Mac OS X
		if len(nameSplit) >= 2 && nameSplit[0] == "Intel" && nameSplit[1] == "Mac" {
			nameSplit = nameSplit[1:]
		}
		name = strings.Join(nameSplit, " ")

		if strings.Contains(version, "x86") || strings.Contains(version, "i686") {
			// x86_64 and i868 are not Linux versions but architectures
			version = ""
		} else if version == "X" && name == "Mac OS" {
			// X is not a version for Mac OS.
			name = name + " " + version
			version = ""
		}
	}
	return name, version
}

// Returns combined information for the operating system.
func (p *UserAgent) OSInfo() OSInfo {
	// Special case for iPhone weirdness
	os := strings.Replace(p.os, "like Mac OS X", "", 1)
	os = strings.Replace(os, "CPU", "", 1)
	os = strings.Trim(os, " ")

	osSplit := strings.Split(os, " ")

	// Special case for x64 edition of Windows
	if os == "Windows XP x64 Edition" {
		osSplit = osSplit[:len(osSplit)-2]
	}

	name, version := osName(osSplit)

	// Special case for names that contain a forward slash version separator.
	if strings.Contains(name, "/") {
		s := strings.Split(name, "/")
		name = s[0]
		version = s[1]
	}

	// Special case for versions that use underscores
	version = strings.Replace(version, "_", ".", -1)

	return OSInfo{
		FullName: p.os,
		Name:     name,
		Version:  version,
	}
}

var ie11Regexp = regexp.MustCompile("^rv:(.+)$")

// A struct containing all the information that we might be
// interested from the browser.
type Browser struct {
	// The name of the browser's engine.
	Engine string

	// The version of the browser's engine.
	EngineVersion string

	// The name of the browser.
	Name string

	// The version of the browser.
	Version string
}

// Extract all the information that we can get from the User-Agent string
// about the browser and update the receiver with this information.
//
// The function receives just one argument "sections", that contains the
// sections from the User-Agent string after being parsed.
func (p *UserAgent) detectBrowser(sections []section) {
	slen := len(sections)

	if sections[0].name == "Opera" {
		p.browser.Name = "Opera"
		p.browser.Version = sections[0].version
		p.browser.Engine = "Presto"
		if slen > 1 {
			p.browser.EngineVersion = sections[1].version
		}
	} else if sections[0].name == "Dalvik" {
		// When Dalvik VM is in use, there is no browser info attached to ua.
		// Although browser is still a Mozilla/5.0 compatible.
		p.mozilla = "5.0"
	} else if slen > 1 {
		engine := sections[1]
		p.browser.Engine = engine.name
		p.browser.EngineVersion = engine.version
		if slen > 2 {
			sectionIndex := 2
			// The version after the engine comment is empty on e.g. Ubuntu
			// platforms so if this is the case, let's use the next in line.
			if sections[2].version == "" && slen > 3 {
				sectionIndex = 3
			}
			p.browser.Version = sections[sectionIndex].version
			if engine.name == "AppleWebKit" {
				switch sections[slen-1].name {
				case "Edge":
					p.browser.Name = "Edge"
					p.browser.Version = sections[slen-1].version
					p.browser.Engine = "EdgeHTML"
					p.browser.EngineVersion = ""
				case "OPR":
					p.browser.Name = "Opera"
					p.browser.Version = sections[slen-1].version
				default:
					if sections[sectionIndex].name == "Chrome" {
						p.browser.Name = "Chrome"
					} else if sections[sectionIndex].name == "Chromium" {
						p.browser.Name = "Chromium"
					} else {
						p.browser.Name = "Safari"
					}
				}
			} else if engine.name == "Gecko" {
				name := sections[2].name
				if name == "MRA" && slen > 4 {
					name = sections[4].name
					p.browser.Version = sections[4].version
				}
				p.browser.Name = name
			} else if engine.name == "like" && sections[2].name == "Gecko" {
				// This is the new user agent from Internet Explorer 11.
				p.browser.Engine = "Trident"
				p.browser.Name = "Internet Explorer"
				for _, c := range sections[0].comment {
					version := ie11Regexp.FindStringSubmatch(c)
					if len(version) > 0 {
						p.browser.Version = version[1]
						return
					}
				}
				p.browser.Version = ""
			}
		}
	} else if slen == 1 && len(sections[0].comment) > 1 {
		comment := sections[0].comment
		if comment[0] == "compatible" && strings.HasPrefix(comment[1], "MSIE") {
			p.browser.Engine = "Trident"
			p.browser.Name = "Internet Explorer"
			// The MSIE version may be reported as the compatibility version.
			// For IE 8 through 10, the Trident token is more accurate.
			// http://msdn.microsoft.com/en-us/library/ie/ms537503(v=vs.85).aspx#VerToken
			for _, v := range comment {
				if strings.HasPrefix(v, "Trident/") {
					switch v[8:] {
					case "4.0":
						p.browser.Version = "8.0"
					case "5.0":
						p.browser.Version = "9.0"
					case "6.0":
						p.browser.Version = "10.0"
					}
					break
				}
			}
			// If the Trident token is not provided, fall back to MSIE token.
			if p.browser.Version == "" {
				p.browser.Version = strings.TrimSpace(comment[1][4:])
			}
		}
	}
}

// Returns two strings. The first string is the name of the engine and the
// second one is the version of the engine.
func (p *UserAgent) Engine() (string, string) {
	return p.browser.Engine, p.browser.EngineVersion
}

// Returns two strings. The first string is the name of the browser and the
// second one is the version of the browser.
func (p *UserAgent) Browser() (string, string) {
	return p.browser.Name, p.browser.Version
}

var botFromSiteRegexp = regexp.MustCompile("http://.+\\.\\w+")

// Get the name of the bot from the website that may be in the given comment. If
// there is no website in the comment, then an empty string is returned.
func getFromSite(comment []string) string {
	if len(comment) == 0 {
		return ""
	}

	// Where we should check the website.
	idx := 2
	if len(comment) < 3 {
		idx = 0
	}

	// Pick the site.
	results := botFromSiteRegexp.FindStringSubmatch(comment[idx])
	if len(results) == 1 {
		// If it's a simple comment, just return the name of the site.
		if idx == 0 {
			return results[0]
		}

		// This is a large comment, usually the name will be in the previous
		// field of the comment.
		return strings.TrimSpace(comment[1])
	}
	return ""
}

// Returns true if the info that we currently have corresponds to the Google
// mobile bot. This function also modifies some attributes in the receiver
// accordingly.
func (p *UserAgent) googleBot() bool {
	// This is a hackish way to detect Google's mobile bot.
	if strings.Index(p.ua, "Googlebot") != -1 {
		p.platform = ""
		p.undecided = true
	}
	return p.undecided
}

// Set the attributes of the receiver as given by the parameters. All the other
// parameters are set to empty.
func (p *UserAgent) setSimple(name, version string, bot bool) {
	p.bot = bot
	if !bot {
		p.mozilla = ""
	}
	p.browser.Name = name
	p.browser.Version = version
	p.browser.Engine = ""
	p.browser.EngineVersion = ""
	p.os = ""
	p.localization = ""
}

// Fix some values for some weird browsers.
func (p *UserAgent) fixOther(sections []section) {
	if len(sections) > 0 {
		p.browser.Name = sections[0].name
		p.browser.Version = sections[0].version
		p.mozilla = ""
	}
}

var botRegex = regexp.MustCompile("(?i)(bot|crawler|sp(i|y)der|search|worm|fetch|nutch)")

// Check if we're dealing with a bot or with some weird browser. If that is the
// case, the receiver will be modified accordingly.
func (p *UserAgent) checkBot(sections []section) {
	// If there's only one element, and it's doesn't have the Mozilla string,
	// check whether this is a bot or not.
	if len(sections) == 1 && sections[0].name != "Mozilla" {
		p.mozilla = ""

		// Check whether the name has some suspicious "bot" or "crawler" in his name.
		if botRegex.Match([]byte(sections[0].name)) {
			p.setSimple(sections[0].name, "", true)
			return
		}

		// Tough luck, let's try to see if it has a website in his comment.
		if name := getFromSite(sections[0].comment); name != "" {
			// First of all, this is a bot. Moreover, since it doesn't have the
			// Mozilla string, we can assume that the name and the version are
			// the ones from the first section.
			p.setSimple(sections[0].name, sections[0].version, true)
			return
		}

		// At this point we are sure that this is not a bot, but some weirdo.
		p.setSimple(sections[0].name, sections[0].version, false)
	} else {
		// Let's iterate over the available comments and check for a website.
		for _, v := range sections {
			if name := getFromSite(v.comment); name != "" {
				// Ok, we've got a bot name.
				results := strings.SplitN(name, "/", 2)
				version := ""
				if len(results) == 2 {
					version = results[1]
				}
				p.setSimple(results[0], version, true)
				return
			}
		}

		// We will assume that this is some other weird browser.
		p.fixOther(sections)
	}
}
