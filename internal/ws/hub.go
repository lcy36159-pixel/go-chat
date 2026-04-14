package ws

var clients = make(map[uint]*Client)

// 註冊 client
func Register(client *Client) {
	clients[client.UserID] = client
}

// 傳送訊息給指定 user
func SendToUser(userID uint, data []byte) {
	if client, ok := clients[userID]; ok {
		client.Conn.WriteMessage(1, data)
	}
}
