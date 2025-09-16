package PasswordSecurity

import (
	"AITranslatio/Global/CustomErrors"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type PasswordGenerator struct {
	cost int //成本因子，平均性能与安全的最佳实践：12
}

func CreatePasswordGeneratorFactory(cost int) *PasswordGenerator {
	return &PasswordGenerator{cost}
}

func (Generator *PasswordGenerator) GenerateSalt(length int) (string, error) {
	saltBytes := make([]byte, 32) // 32字节 = 256位
	_, err := rand.Read(saltBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(saltBytes), nil
}

// HashPasswordWithSalt 使用额外盐值进行密码哈希（双重保护）
// 注意：其实bcrypt 本身已经包含盐值，这里提供额外的安全层
func (Generator *PasswordGenerator) HashPasswordWithSalt(Password string, Salt string) (string, error) {

	if Salt == "" || Password == "" {
		return "", errors.New(CustomErrors.ErrorPasswordOrSaltIsEmpty)
	}

	HashPasswordWithSalt, err := bcrypt.GenerateFromPassword([]byte(Password+Salt), Generator.cost)
	if err != nil {
		return "", errors.New(CustomErrors.ErrorPasswordHashIsFail + err.Error())
	}

	return string(HashPasswordWithSalt), nil
}

// ValidatePasswordWithSalt  验证带盐值的密码
func (Generator *PasswordGenerator) ValidatePasswordWithSalt(HashPassword, Password, Salt string) bool {

	if Password == "" || Salt == "" || HashPassword == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(HashPassword), []byte(Password+Salt))
	return err == nil
}
