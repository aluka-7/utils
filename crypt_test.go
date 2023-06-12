package utils

import "testing"

func Test_EncryptedPassword(t *testing.T) {
	expected := "d31016aa4230cc95bf653b5266ca1dff5a35ebce09a1dd659c704fcb162ef0432137a8ed0d225e8bf67c49fbba93b6e34a6c"
	salt := "4322f9a718b6b117c842bf6d52ca4eda"
	actual := EncryptedPassword("12345678", salt)

	if expected != actual {
		t.Error("生成的密码不匹配", len(actual), "|", "预期:", expected, "|", "实际:", actual)
	}
}

func Test_Authenticate(t *testing.T) {
	salt := GenerateSalt()
	attemptedPassword := "12345678"
	actual := EncryptedPassword(attemptedPassword, salt)
	if !Authenticate(attemptedPassword, actual, salt) {
		t.Error("生成的密码不匹配", len(actual), actual, salt)
	}
}
