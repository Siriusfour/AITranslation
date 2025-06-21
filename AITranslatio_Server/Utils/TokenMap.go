package Utils

import (
	"sync"
	"time"
)

type tokenInfo struct {
	Account string
	OutTime time.Time
}

type TokenMap struct {
	TokenMap map[string]*tokenInfo
	ticker   *time.Ticker
	done     chan struct{}
	MU       sync.RWMutex
}

func InitTokenMap() *TokenMap {
	return &TokenMap{
		TokenMap: make(map[string]*tokenInfo),
		ticker:   nil,
		done:     make(chan struct{}),
		MU:       sync.RWMutex{},
	}
}
