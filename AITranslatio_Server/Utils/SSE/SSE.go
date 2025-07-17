package SSE

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"sync"
)

type SSEClients struct {
	Clients map[int]chan string
	mutex   sync.Mutex
}

func InitSSEClients() *SSEClients {
	return &SSEClients{
		Clients: make(map[int]chan string),
		mutex:   sync.Mutex{},
	}
}

func (clients *SSEClients) SendNotify(UserID int, message string) {
	clients.mutex.Lock()
	defer clients.mutex.Unlock()
	clients.Clients[UserID] <- message
}

func (clients *SSEClients) RemoveClient(UserID int) {
	clients.mutex.Lock()
	delete(clients.Clients, UserID)
	clients.mutex.Unlock()
}

func (SSEClients *SSEClients) CreateSSE(CreateSSEctx *gin.Context, UserID int) error {

	CreateSSEctx.Writer.Header().Set("Content-Type", "text/event-stream")
	CreateSSEctx.Writer.Header().Set("Cache-Control", "no-cache")
	CreateSSEctx.Writer.Header().Set("Connection", "keep-alive")
	CreateSSEctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	notifyChan := make(chan string, 10)
	//把该chan注册到map里
	SSEClients.mutex.Lock()
	SSEClients.Clients[UserID] = notifyChan
	SSEClients.mutex.Unlock()

	CreateSSEctx.Stream(func(w io.Writer) bool {
		select {
		case msg := <-notifyChan:
			//  写入 SSE 格式数据
			data := fmt.Sprintf("data: %s", msg)
			_, err := w.Write([]byte(data))
			if err != nil {
				fmt.Println("Write error for UserID", UserID, ":", err)
				return false
			}

			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
		case <-CreateSSEctx.Request.Context().Done():
			fmt.Println("Client disconnected with UserID:", UserID)
			return false
		}
		return true
	})

	return nil

}
