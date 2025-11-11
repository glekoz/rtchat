package main


import (
"context"
"flag"
"fmt"
"log"
"net/http"
"os"
"os/signal"
"syscall"
"time"


"github.com/you/realtime-chat-go/internal/hub"
"github.com/you/realtime-chat-go/internal/ws"
)


func main() {
var addr string
flag.StringVar(&addr, "addr", ":8080", "http service address")
flag.Parse()


redisAddr := os.Getenv("REDIS_ADDR")
if redisAddr == "" {
redisAddr = "localhost:6379"
}
channel := os.Getenv("REDIS_CHANNEL")
if channel == "" {
channel = "chat_messages"
}


h, err := hub.NewHub(redisAddr, channel)
if err != nil {
log.Fatalf("failed to create hub: %v", err)
}


// serve static frontend
fs := http.FileServer(http.Dir("./static"))
http.Handle("/", fs)


// websocket endpoint
http.HandleFunc("/ws", ws.MakeHandler(h))


srv := &http.Server{
Addr: addr,
Handler: nil,
}


go func() {
log.Printf("listening on %s (redis=%s channel=%s)", addr, redisAddr, channel)
if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
log.Fatalf("listen: %v", err)
}
}()


// graceful shutdown
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit


log.Println("shutting down server...")
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
}
