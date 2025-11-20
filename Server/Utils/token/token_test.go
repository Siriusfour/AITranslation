package token

import (
	"AITranslatio/Global"
	"AITranslatio/Utils/SnowFlak"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockConfig struct{}

func (*mockConfig) GetInt(key string) int {
	return 24
}

func MockRedis(t *testing.T) *miniredis.Miniredis {
	s, err := miniredis.Run()
	assert.NoError(t, err)

	Global.RedisClient = redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	return s
}

func TestGeneratedToken(t *testing.T) {

	tests := []struct {
		name    string
		token   *Token
		setup   func(s *miniredis.Miniredis)
		wantErr bool
	}{
		// 1. 正常情况：成功生成 token，redis 正确写入
		{
			name:    "valid_generate",
			setup:   func(s *miniredis.Miniredis) {},
			token:   CreateTokenFactory(0, SnowFlak.CreateSnowflakeFactory().GetID()),
			wantErr: false,
		},
		// 2. 错误用例：签名为空
		{
			name:  "sign_is_null",
			setup: func(s *miniredis.Miniredis) {},
			token: &Token{
				Type:       "RefreshToken",
				UserID:     100,
				TokenID:    2000,
				EncryptKey: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				},
				config: &mockConfig{},
			},
			wantErr: true,
		},
		// 3. 异常： Redis pipeline 执行失败
		{
			name: "RedisError",
			setup: func(s *miniredis.Miniredis) {
				s.Close() // ❗ 直接关闭 Redis，让 Exec 失败
			},
			token:   CreateTokenFactory(0, SnowFlak.CreateSnowflakeFactory().GetID()),
			wantErr: true,
		},

		//4. 边界情况：TokenID，UserID为0
		{name: "id_is_zero",
			setup: func(s *miniredis.Miniredis) {},
			token: &Token{
				Type:       "RefreshToken",
				UserID:     0,
				TokenID:    0,
				EncryptKey: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				},
				config: &mockConfig{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {})

		s := MockRedis(t)
		defer s.Close()

		tt.setup(s)
		tokenStr, err := tt.token.GeneratedToken()

		//如果希望报错 → err 必须非空 ， tokenStr 必须为空（失败不能返回 token
		if tt.wantErr {
			assert.Error(t, err)
			assert.Empty(t, tokenStr)
			return
		}

		//如果情况正常，则不err为空，TokenStr不为空
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenStr)

		// 检查 redis 是否写入
		key := fmt.Sprintf("UserID:%d", tt.token.UserID)

		assert.Equal(t, s.HGet(key, "UserID"), fmt.Sprintf("%d", tt.token.UserID))
		assert.Equal(t, s.HGet(key, "TokenID"), fmt.Sprintf("%d", tt.token.TokenID))
	}
}

func TestParseToken()
