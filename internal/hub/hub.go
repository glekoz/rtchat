package hub
return h, nil
}


// Register a client
func (h *Hub) Register(c *Client) {
h.mu.Lock()
h.clients[c] = true
h.mu.Unlock()
}


// Unregister a client
func (h *Hub) Unregister(c *Client) {
h.mu.Lock()
if _, ok := h.clients[c]; ok {
delete(h.clients, c)
close(c.Send)
}
h.mu.Unlock()
}


// Broadcast publishes a message to Redis. The pubsubLoop will receive it and broadcast to local clients.
func (h *Hub) Broadcast(msg Message) error {
bs, err := json.Marshal(msg)
if err != nil {
return err
}
return h.redis.Publish(context.Background(), h.channel, bs).Err()
}


// pubsubLoop subscribes to Redis and broadcasts incoming messages to connected clients
func (h *Hub) pubsubLoop(ctx context.Context) {
pubsub := h.redis.Subscribe(ctx, h.channel)
defer pubsub.Close()


ch := pubsub.Channel()
for {
select {
case <-ctx.Done():
return
case m, ok := <-ch:
if !ok {
return
}
var msg Message
if err := json.Unmarshal([]byte(m.Payload), &msg); err != nil {
// ignore malformed
continue
}
h.mu.RLock()
for c := range h.clients {
// try non-blocking send
select {
case c.Send <- msg:
default:
// if client send buffer is full, skip it
}
}
h.mu.RUnlock()
}
}
}


// Close stops the hub
func (h *Hub) Close() {
h.cancel()
_ = h.redis.Close()
}
