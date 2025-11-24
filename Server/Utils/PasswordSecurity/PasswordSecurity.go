package PasswordSecurity

import (
	"AITranslatio/Global/MyErrors"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type PasswordGenerator struct {
	cost int //成本因子，平均性能与安全的最佳实践：12
}

func CreatePasswordGeneratorFactory(cost ...int) *PasswordGenerator {
	if len(cost) == 0 {
		cost = append(cost, 1)
	}
	return &PasswordGenerator{cost[0]}
}

// GenerateSalt 生成随机的盐值
func (Generator *PasswordGenerator) GenerateSalt() (string, error) {
	saltBytes := make([]byte, 32) // 32字节 = 256位
	_, err := rand.Read(saltBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(saltBytes), nil
}

// HashPasswordWithSalt 使用额外盐值进行密码哈希（双重保护）
// 注意：其实bcrypt 本身已经包含盐值，这里提供额外的安全层
// 更新： 放弃自己的盐值，只用bcrypt自带的盐值就好，另外加入pepper机制

func (Generator *PasswordGenerator) HashPasswordWithSalt(Password string, Salt string) (string, error) {

	if Salt == "" || Password == "" {
		return "", errors.New(MyErrors.ErrorPasswordOrSaltIsEmpty)
	}

	HashPasswordWithSalt, err := bcrypt.GenerateFromPassword([]byte(Password+Salt), Generator.cost)
	if err != nil {
		return "", errors.New(MyErrors.ErrorPasswordHashIsFail + err.Error())
	}

	return string(HashPasswordWithSalt), nil
}

// ValidatePasswordWithSalt  验证带盐值的密码
func (Generator *PasswordGenerator) ValidatePasswordWithSalt(HashPassword, Password, Salt string) error {

	if Password == "" || Salt == "" || HashPassword == "" {
		return errors.New(MyErrors.ErrorPasswordOrSaltIsEmpty)
	}

	err := bcrypt.CompareHashAndPassword([]byte(HashPassword), []byte(Password+Salt))
	if err != nil {
		return fmt.Errorf("验证密码失败：%w", err)
	}
	return nil
}
