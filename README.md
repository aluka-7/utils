# Utils

## Arrays
返回数组'array'中'val'的索引位置,返回值为:-1 表示不包含
* Contains(array interface{}, val interface{}) (index int)
* ContainsString(array []string, val string) (index int)
* ContainsInt(array []int64, val int64) (index int)
* ContainsUint(array []uint64, val uint64) (index int)
* ContainsBool(array []bool, val bool) (index int)
* ContainsFloat(array []float64, val float64) (index int)
* ContainsComplex(array []complex128, val complex128) (index int)

## cmd
执行系统命令返回stdout,stderr,以及可能的错误
* ExecCmdDirBytes(dir, cmdName string, args ...string) ([]byte, []byte, error)  //在给定目录中执行系统命令,并以字节类型返回stdout,stderr,以及可能的错误
* ExecCmdBytes(cmdName string, args ...string) ([]byte, []byte, error)          //执行系统命令并返回stdout,以字节类型表示的stderr以及可能的错误
* ExecCmdDir(dir, cmdName string, args ...string) (string, string, error)       //在给定目录中执行系统命令并返回stdout,字符串类型的stderr以及可能的错误.
* ExecCmd(cmdName string, args ...string) (string, string, error)               //执行系统命令并返回stdout,字符串类型的stderr以及可能的错误

## Convert
提供各种类型还转工具
* StrTo(string)                                     //转换字符串以指定类型(Uint8,Int,Int64,MustUint8,MustInt,MustInt64,Float64,String,Hash,Exist)。
* ToStr(value interface{}, args ...int) (s string)  //将任何类型转换为字符串。    
* HexStr2int(hexStr string) (int, error)            //将十六进制格式字符串转换为十进制数
* Int2HexStr(num int) (hex string)                  //将十进制数格式转换为十六进制字符串
* ArrayToString(A []int, denim string) string       //整形数组转换成字符串

## Crypt
建立一个go,java,python通用的加解密实现包。
* MD5(origData string) string                                                   //给指定的字符串进行MD5加密
* Authenticate(attemptedPassword, encryptedPassword, salt string) bool          //对输入的密码进行验证
* GenerateSalt() string                                                         //通过提供加密的强随机数生成器 生成盐
* EncryptedPassword(rawPwd string, salt string) string                          //生成密文
* PBKDF2(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte    //基于PBKDF2算法加密
* Encrypt(origData, key []byte) ([]byte, error)                                 //基于PKCS5Padding算法加密
* Decrypt(cryptic, key []byte) ([]byte, error)                                  //基于PKCS5Padding算法解密

## File
提供文件操作相关工具
* HumaneBytes(s uint64, base float64, sizes []string) string        //个性化文件大小计算文件大小
* HumaneFileSize(s uint64) string                                   //个性化文件大小计算文件大小并生成用户友好的字符串
* FileMTime(file string) (int64, error)                             //获取文件的修改时间
* FileSize(file string) (int64, error)                              //获取文件大小
* Copy(src, dest string) error                                      //从源地址复制到目标地址
* WriteFile(filename string, data []byte) error                     //将数据写入文件名指定文件,如果文件不存在,Write File将创建它及其上层路径
* IsFile(filePath string) bool                                      //判断给定路径是不是文件以及是否存在,如果给定的路径是文件，则返回true，或者当它是目录或不存在时返回false
* IsExist(path string) bool                                         //检查文件或目录是否存在,当文件或者目录不存在时返回false
* GetGOPATH() []string                                              // 返回GOPATH变量中的所有路径.
* IsDir(dir string) bool                                            // 如果给定路径是目录,则返回true;如果文件或目录不存在;则返回false.
* StatDir(rootPath string, includeDir ...bool) ([]string, error)    //通过深度优先收集给定目录的信息
* GetAllSubDirs(rootPath string) ([]string, error)                  // 返回给定根路径的所有子目录,返回值不包含给定的路径.

## http
提供http远程访问工具
* HttpCall(client *http.Client, method, url string, header http.Header, body io.Reader) (io.ReadCloser, error)
* HttpGet(client *http.Client, url string, header http.Header) (io.ReadCloser, error)
* HttpPost(client *http.Client, url string, header http.Header, body []byte) (io.ReadCloser, error)
* HttpGetToFile(client *http.Client, url string, header http.Header, fileName string) error
* HttpGetBytes(client *http.Client, url string, header http.Header) ([]byte, error)
* HttpGetJSON(client *http.Client, url string, v interface{}) error
* HttpPostJSON(client *http.Client, url string, body, v interface{}) error
* FetchFiles(client *http.Client, files []RawFile, header http.Header) error
* FetchFilesCurl(files []RawFile, curlOptions ...string) error
* New(ua string) *UserAgent
* (p *UserAgent) Parse(ua string)

## IP
* ExternalIP() 获取外部IP
* InternalIP() 获取内部IP

## SyncMap
提供一个同步map操作工具

## Math
提供了基于随机数生成值的工具
* Div(n, b float64) float64                             //浮点数除法
* RandInt(start int, end int) int                       //随机int
* RandInt64(start int64, end int64) int64               //随机int64
* GenerateRandomCode() string                           //随机获取6位数字符串
* GenFixedLengthChineseChars(length int) string         //指定长度随机中文字符(包含复杂字符)
* GenRandomLengthChineseChars(start, end int) string    //指定范围随机中文字符
* RandStr(len int) string                               //随机英文小写字母
* RandString(n int) string                              //生成指定长度的随机字母和数字字符串，包括0-9、a-z、A-Z的所有字符

## string
* TokenizeToStringArray(str, delimiters string, trimTokens, ignoreEmptyTokens bool) []*string   //根据分隔符进行分割处理，形成包路径数组。默认分割符为：",; \t\n"
* Str2Bytes(s string) []byte                                                                    //字符串到字节
* Bytes2Str(b []byte) string                                                                    //字节到字符串
* StartsWith(str, prefix string, offset int) bool                                               //判断str是以prefix开始
* IsBlank(source string) bool                                                                   //判断是否存在空格
* HasText                                                                                       //判断是否有值
* AppendStr(strs []string, str string) []string                                                 //将字符串追加到数组中,且没有重复

## time
* Date(ti int64, format string) string      //将unix时间整型格式化为字符串
* DateS(ts string, format string) string    //将unix时间字符串格式化为字符串
* DateT(t time.Time, format string) string  //提供java一样格式的日期格式化方式
* StringTime time.Time                      //用于提供字符日期类型json序列化互转
* NumberTime time.Time                      //用于提供unix日期类型json序列化互转
* NewDuration(str string) (dur Duration)    //提供将("ns", "us" (or "µs"), "ms", "s", "m", "h")转换成time.Duration
* NewTime(t float64) time.Time              //从float64创建一个新的time.Time

## License
This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.
