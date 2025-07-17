package UtilsStruct

import "sync"

type SSEClients struct {
	Clients map[int]chan string
	mutex   sync.Mutex
}
