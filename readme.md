# Realtime Chat (Go) — WebSocket + Redis Pub/Sub


## What it is
A small chat server written in Go that demonstrates:


- WebSocket connections (gorilla/websocket)
- Redis Pub/Sub for cross-instance message propagation (go-redis)
- Simple frontend (single HTML file) that connects with WebSocket
- Docker Compose to run Redis + the app


## Run locally (Docker)


```bash
docker-compose up --build
