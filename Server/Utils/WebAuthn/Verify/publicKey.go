package Verify

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/fxamacker/cbor/v2"
	"math/big"
)

// COSE Key label constants (per RFC 8152 / WebAuthn)
const (
	coseKeyLabelKty = 1 // key type
	coseKeyLabelAlg = 3 // algorithm

	// EC2 / OKP specific labels
	coseKeyLabelCrv = -1 // curve (EC2/OKP)
	coseKeyLabelX   = -2 // x coord (EC2) or public key bytes (OKP)
	coseKeyLabelY   = -3 // y coord (EC2)
	// RSA specific labels
	coseKeyLabelN = -1 // modulus n
	coseKeyLabelE = -2 // exponent e
)

// kty values
const (
	ktyOKP = 1 // Octet Key Pair (e.g., Ed25519)
	ktyEC2 = 2 // EC2 (e.g., P-256)
	ktyRSA = 3 // RSA
)

// crv values
const (
	crvP256    = 1 // EC2 P-256
	crvEd25519 = 6 // OKP Ed25519
)

func extractCredentialData(authData []byte) (*AttestedCredentialData, error) {

	//校验是否有变长部分
	flags := authData[32]
	if (flags & 0x40) == 0 {
		return nil, errors.New("AT flag 未设置，没有凭证数据")
	}

	offset := 37
	aaguid := authData[offset : offset+16]
	offset += 16

	credIdLen := binary.BigEndian.Uint16(authData[offset : offset+2])
	offset += 2

	credentialID := authData[offset : offset+int(credIdLen)]
	offset += int(credIdLen)

	cosePublicKey := authData[offset:]

	return &AttestedCredentialData{
		AAGUID:              aaguid,
		CredentialID:        credentialID,
		CredentialPublicKey: cosePublicKey,
	}, nil

}

func parseCOSEPublicKey(coseBytes []byte) (interface{}, int, error) {
	var m map[int]interface{}
	if err := cbor.Unmarshal(coseBytes, &m); err != nil {
		return nil, 0, fmt.Errorf("cbor unmarshal COSE key: %w", err)
	}
	return parseCOSEPublicKeyMap(m)
}

// parseCOSEPublicKeyMap parses an already-decoded COSE map (useful if你上游已解了 CBOR).
func parseCOSEPublicKeyMap(m map[int]interface{}) (interface{}, int, error) {
	// kty / alg
	kty, ok := asInt(m[coseKeyLabelKty])
	if !ok {
		return nil, 0, errors.New("COSE key missing/invalid kty")
	}
	alg, _ := asInt(m[coseKeyLabelAlg]) // alg 可缺省，但尽量取到

	switch kty {
	case ktyEC2:
		// EC2: need crv, x, y
		crv, ok := asInt(m[coseKeyLabelCrv])
		if !ok {
			return nil, alg, errors.New("EC2 key missing crv")
		}
		xBytes, ok := asBytes(m[coseKeyLabelX])
		if !ok {
			return nil, alg, errors.New("EC2 key missing x")
		}
		yBytes, ok := asBytes(m[coseKeyLabelY])
		if !ok {
			return nil, alg, errors.New("EC2 key missing y")
		}
		switch crv {
		case crvP256:
			pub := &ecdsa.PublicKey{
				Curve: elliptic.P256(),
				X:     new(big.Int).SetBytes(xBytes),
				Y:     new(big.Int).SetBytes(yBytes),
			}
			if !pub.Curve.IsOnCurve(pub.X, pub.Y) {
				return nil, alg, errors.New("EC2 P-256 point not on curve")
			}
			return pub, alg, nil
		default:
			return nil, alg, fmt.Errorf("unsupported EC2 crv=%d (only P-256 supported)", crv)
		}

	case ktyOKP:
		// OKP: crv + x (no y)
		crv, ok := asInt(m[coseKeyLabelCrv])
		if !ok {
			return nil, alg, errors.New("OKP key missing crv")
		}
		xBytes, ok := asBytes(m[coseKeyLabelX])
		if !ok {
			return nil, alg, errors.New("OKP key missing x")
		}
		switch crv {
		case crvEd25519:
			// x is the 32-byte Ed25519 public key
			if l := len(xBytes); l != ed25519.PublicKeySize {
				return nil, alg, fmt.Errorf("Ed25519 public key length invalid: %d", l)
			}
			return ed25519.PublicKey(xBytes), alg, nil
		default:
			return nil, alg, fmt.Errorf("unsupported OKP crv=%d (only Ed25519 supported)", crv)
		}

	case ktyRSA:
		// RSA: n (modulus), e (exponent) as big-endian byte strings
		nBytes, okN := asBytes(m[coseKeyLabelN])
		eBytes, okE := asBytes(m[coseKeyLabelE])
		if !okN || !okE {
			return nil, alg, errors.New("RSA key missing n or e")
		}
		n := new(big.Int).SetBytes(nBytes)
		if n.Sign() <= 0 {
			return nil, alg, errors.New("RSA modulus n must be positive")
		}
		var eb big.Int
		eb.SetBytes(eBytes)
		e := int(eb.Int64())
		if e <= 0 {
			return nil, alg, errors.New("RSA exponent e must be positive")
		}
		return &rsa.PublicKey{N: n, E: e}, alg, nil

	default:
		return nil, alg, fmt.Errorf("unsupported kty=%d", kty)
	}
}

func asInt(v interface{}) (int, bool) {
	switch t := v.(type) {
	case int:
		return t, true
	case int64:
		return int(t), true
	case uint64:
		return int(t), true
	case uint32:
		return int(t), true
	case uint:
		return int(t), true
	default:
		return 0, false
	}
}

func asBytes(v interface{}) ([]byte, bool) {
	b, ok := v.([]byte)
	return b, ok
}
