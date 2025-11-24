package PasswordSecurity

import (
	"strings"
	"testing"
)

func TestHashPasswordWithSalt(t *testing.T) {
	g := &PasswordGenerator{cost: 12}

	tests := []struct {
		name     string
		password string
		salt     string
		wantErr  bool
	}{
		{
			name:     "正常情况",
			password: "password123",
			salt:     "mysalt",
			wantErr:  false,
		},
		{
			name:     "密码为空",
			password: "",
			salt:     "salt",
			wantErr:  true,
		},
		{
			name:     "盐为空",
			password: "password123",
			salt:     "",
			wantErr:  true,
		},
		{
			name:     "密码和盐都为空",
			password: "",
			salt:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := g.HashPasswordWithSalt(tt.password, tt.salt)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (password=%q, salt=%q)", tt.password, tt.salt)
				}
				// 有错就不再验 hash 了
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if hash == "" {
				t.Fatal("hash should not be empty")
			}

			// 简单验证下是否是合法的 bcrypt 格式
			if !strings.HasPrefix(hash, "$2") {
				t.Fatalf("invalid bcrypt format: %s", hash)
			}
		})
	}
}

func TestValidatePasswordWithSalt(t *testing.T) {
	g := &PasswordGenerator{cost: 10}

	// 为了测试“正确情况”，我们先生成一个合法 hash
	correctPassword := "mypassword"
	correctSalt := "mysalt"
	correctHash, _ := g.HashPasswordWithSalt(correctPassword, correctSalt)

	tests := []struct {
		name         string
		hash         string
		password     string
		salt         string
		wantErr      bool
		wantContains string // 错误信息包含某些关键字（可选）
	}{
		{
			name:     "验证成功",
			hash:     correctHash,
			password: correctPassword,
			salt:     correctSalt,
			wantErr:  false,
		},
		{
			name:         "密码错误",
			hash:         correctHash,
			password:     "wrongpass",
			salt:         correctSalt,
			wantErr:      true,
			wantContains: "password", // bcrypt error message
		},
		{
			name:         "盐值错误",
			hash:         correctHash,
			password:     correctPassword,
			salt:         "wrongsalt",
			wantErr:      true,
			wantContains: "password", // bcrypt 返回密码错误
		},
		{
			name:         "Hash 不合法（无法解析）",
			hash:         "invalid-hash-string",
			password:     correctPassword,
			salt:         correctSalt,
			wantErr:      true,
			wantContains: "hash",
		},
		{
			name:     "密码为空",
			hash:     correctHash,
			password: "",
			salt:     correctSalt,
			wantErr:  true,
		},
		{
			name:     "盐为空",
			hash:     correctHash,
			password: correctPassword,
			salt:     "",
			wantErr:  true,
		},
		{
			name:     "Hash 为空",
			hash:     "",
			password: correctPassword,
			salt:     correctSalt,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := g.ValidatePasswordWithSalt(tt.hash, tt.password, tt.salt)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}

				if tt.wantContains != "" && !strings.Contains(err.Error(), tt.wantContains) {
					t.Fatalf("expected error containing %q, got: %v", tt.wantContains, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
