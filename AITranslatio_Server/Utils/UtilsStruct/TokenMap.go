package UtilsStruct

import (
	"sync"
	"time"
)

type TokenInfo struct {
	Revoked        bool
	RefreshToken   string
	AccessToken    string
	RegisteredTime string
}

type TokenMap struct {
	TokenMap map[int]*TokenInfo
	ticker   *time.Ticker
	done     chan struct{}
	MU       sync.RWMutex
}

func InitTokenMap() *TokenMap {
	return &TokenMap{
		TokenMap: make(map[int]*TokenInfo),
		ticker:   nil,
		done:     make(chan struct{}),
		MU:       sync.RWMutex{},
	}
}
