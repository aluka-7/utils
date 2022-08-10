package utils

import "testing"

func Test_EncryptedPassword(t *testing.T) {
	expected := "b8c9780d6775dae3080de04f28730c40fafc4cea8a33f9e4d5027f5696babc313929925d6def2379c068ec5cf7e497840bac"
	salt := "4322f9a718b6b117c842bf6d52ca4eda"
	actual := EncryptedPassword("forchange", salt)

	if expected != actual {
		t.Error("生成的密码不匹配", len(actual), "|", "预期:", expected, "|", "实际:", actual)
	}
}

func Test_Authenticate(t *testing.T) {
	salt := GenerateSalt()
	attemptedPassword := "forchange"
	actual := EncryptedPassword(attemptedPassword, salt)
	if !Authenticate(attemptedPassword, actual, salt) {
		t.Error("生成的密码不匹配", len(actual), actual, salt)
	}
}
