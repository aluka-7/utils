package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

// 给指定的字符串进行MD5加密
func MD5(origData string) string {
	m := md5.New()
	m.Write([]byte(origData))
	return hex.EncodeToString(m.Sum(nil))
}

/**
 * 对输入的密码进行验证
 *
 * @param attemptedPassword 待验证的密码
 * @param encryptedPassword 密文
 * @param salt              盐值
 * @return 是否验证成功
 */
func Authenticate(attemptedPassword, encryptedPassword, salt string) bool {
	// 用相同的盐值对用户输入的密码进行加密
	eap := EncryptedPassword(attemptedPassword, salt)
	// 把加密后的密文和原密文进行比较，相同则验证成功，否则失败
	return eap == encryptedPassword
}

/**
 * 通过提供加密的强随机数生成器 生成盐
 *
 * @return
 */
func GenerateSalt() string {
	keyword := RandString(10)
	data := []byte(keyword)
	hash := SHA1(SHA1(data))
	key := hash[0:16]
	return hex.EncodeToString(key)
}
func SHA1(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

/**
 * 生成密文
 *
 * @param rawPwd 明文密码
 * @param salt     盐值
 * @return
 */
func EncryptedPassword(rawPwd string, salt string) string {
	pwd := PBKDF2([]byte(rawPwd), []byte(salt), 10000, 50, sha256.New)
	return hex.EncodeToString(pwd)
}
func PBKDF2(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
}

func Encrypt(origData, key []byte) ([]byte, error) {
	keyBytes := getKeyBytes(string(key))
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, keyBytes[:blockSize])
	cryptic := make([]byte, len(origData))
	blockMode.CryptBlocks(cryptic, origData)
	return cryptic, nil
}

func Decrypt(cryptic, key []byte) ([]byte, error) {
	keyBytes := getKeyBytes(string(key))
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, keyBytes[:blockSize])
	origData := make([]byte, len(cryptic))
	blockMode.CryptBlocks(origData, cryptic)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}
func getKeyBytes(key string) []byte {
	keyBytes := []byte(key)
	switch l := len(keyBytes); {
	case l < 16:
		keyBytes = append(keyBytes, make([]byte, 16-l)...)
	case l > 16:
		keyBytes = keyBytes[:16]
	}
	return keyBytes
}
func PKCS5Padding(cipher []byte, blockSize int) []byte {
	padding := blockSize - len(cipher)%blockSize
	text := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipher, text...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	un := int(origData[length-1])
	return origData[:(length - un)]
}
