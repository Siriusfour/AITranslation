package Utils

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

func Base64urlToInt64(userHandle string) (int64, error) {
	// 1. base64url 解码（RawURLEncoding 很关键）
	b, err := base64.RawURLEncoding.DecodeString(userHandle)
	if err != nil {
		return 0, fmt.Errorf("base64url decode failed: %w", err)
	}

	// 2. WebAuthn user.id 通常是 8 字节（int64）
	if len(b) != 8 {
		return 0, fmt.Errorf("invalid userHandle length: %d", len(b))
	}

	// 3. 按 big-endian 解析为 int64
	uid := int64(binary.BigEndian.Uint64(b))
	return uid, nil
}
