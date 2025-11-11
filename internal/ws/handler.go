package ws
})


for {
_, bs, err := conn.ReadMessage()
if err != nil {
if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
log.Printf("ws read error: %v", err)
}
break
}
var msg hub.Message
if err := json.Unmarshal(bs, &msg); err != nil {
// ignore invalid
continue
}
// publish
if err := h.Broadcast(msg); err != nil {
log.Printf("broadcast err: %v", err)
}
}
}


// writer sends messages from hub to the websocket
func writer(conn *websocket.Conn, c *hub.Client) {
ticker := time.NewTicker(54 * time.Second)
defer func() {
ticker.Stop()
conn.Close()
}()


for {
select {
case msg, ok := <-c.Send:
conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
if !ok {
// hub closed channel
conn.WriteMessage(websocket.CloseMessage, []byte{})
return
}
bs, err := json.Marshal(msg)
if err != nil {
continue
}
if err := conn.WriteMessage(websocket.TextMessage, bs); err != nil {
return
}
case <-ticker.C:
conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
return
}
}
}
}
