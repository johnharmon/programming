package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	messageEnd = 0x00
	//messageSep = byte(messageEnd)
)

func websocketHandler(w http.ResponseWriter, r *http.Request) {
}

func NewClientUpdate(clientID string, connectionID int, chatChan chan ChatMessage, action string) ClientUpdate {
	return ClientUpdate{}
}

func (gb *GlobalBroadcaster) GetChatClient(ConnectionID int) *ChatClient {
	for connID, client := range gb.globalClients {
		if connID == ConnectionID {
			return &client
		}
	}
	return nil
}

func (gb *GlobalBroadcaster) addConnection(conn *net.Conn) error {
	return nil
}

func (gb *GlobalBroadcaster) Broadcast(cm ChatMessage) {
	chatMessageBytes, err := json.Marshal(cm)
	if err != nil {
		fmt.Printf("Error mashalling message: %s\n", err)
	}
	for _, chatClient := range gb.globalClients {
		if chatClient.ConnectionID != cm.ConnectionID {
			chatClient.ServerToClient <- cm
		}
	}
}

func (gb *GlobalBroadcaster) RemoveClient(cu ClientUpdate) bool {
	for k, v := range gb.globalClients {
		if v.ConnectionID == cu.ConnectionID {
			delete(gb.globalClients, k)
			return true
		}
	}
	return false
}
func NewChatClient(cu ClientUpdate) ChatClient {
	cc := ChatClient{
		Username:          cu.ClientID,
		ConnectionAddress: nil,
		ServerToClient:    make(chan []byte),
		ClientToServer:    make(chan ChatMessage),
		ConnectionID:      cu.ConnectionID,
		Connection:        cu.Connection,
	}
	return cc
}

func (gb *GlobalBroadcaster) UpdateClient(cu ClientUpdate) bool {
	gb.globalClientUpdateMut.Lock()
	defer gb.globalClientUpdateMut.Unlock()
	if cu.Action == "connect" {
		//gb.globalClients[cu.ConnectionID] = NewChatClient(cu)
		gb.globalClients[cu.ConnectionID] = NewChatClient(cu)
		return true
	} else if cu.Action == "disconnect" {
		gb.RemoveClient(cu)
		return false
	}
	return false
}

func (cc *ChatClient) ReadMessage() (cm ChatMessage, funcErr error) {
	messageBytes := make([]byte, 1024)
	messageSep := make([]byte, 1)
	messageSep[0] = byte(messageEnd)
	_, err := cc.Connection.Read(messageBytes)
	if err != nil {
		fmt.Printf("Error reading from client: %s\n", err)
		return ChatMessage{}, err
	}
	endOfMessage := bytes.Index(messageBytes, messageSep)
	if endOfMessage != -1 {
		messageTail := messageBytes[endOfMessage:]
		messageBytes = append(cc.MessageBuffer, messageBytes[:endOfMessage]...)
		cc.MessageBuffer = messageTail
	} else {
		cc.MessageBuffer = append(cc.MessageBuffer, messageBytes...)
		//json.Marshal()
	}

}

func ManageGlobalBroadcaster(gb *GlobalBroadcaster) {
	for {
		select {
		case message := <-gb.globalProducer:
			go gb.Broadcast(message) //write message to all consumers
		case clientUpdate := <-gb.ClientUpdates:
			gb.UpdateClient(clientUpdate) // method for updating client map, closing/opening channels, etc
		}
	}
}

func ReadClient(cc *ChatClient) (message []byte, tail []byte) {
	buf := []byte{}
	messageSep := append([]byte{}, byte(messageEnd))
	_, err := cc.Connection.Read(buf)
	if err != nil {
		fmt.Printf("Error reading from client: %s\n", err)
		return message, tail
	}
	endOfMessage := bytes.Index(buf, messageSep)
	if endOfMessage != -1 {
		message = buf[:endOfMessage]
		tail = buf[endOfMessage:]
	} else {
		tail = buf[:]
		//json.Marshal()
	}
	return message, tail

}

func ClientListener(cc *ChatClient) {
	inboundBuffer := []byte{}
	previousBuffer := []byte{}
	messageSep := append([]byte{}, byte(messageEnd))
	message := []byte{}
	for {
		_, err := cc.Connection.Read(inboundBuffer)
		if err != nil {
			fmt.Printf("Error reading from connection: %s", err)
		}
		endOfMessage := bytes.Index(inboundBuffer, messageSep)
		if endOfMessage != -1 {
			message = append(previousBuffer, inboundBuffer[:endOfMessage]...)
			previousBuffer = inboundBuffer[endOfMessage:]

		} else {
			previousBuffer = append(previousBuffer, inboundBuffer...)
		}
		//copy()
	}
}

func ManageChatClient(cc *ChatClient, gb *GlobalBroadcaster) {
	for {
		select {
		case clientMessage := <-cc.ClientToServer:
			gb.globalProducer <- clientMessage
		}
	}
}

func handleConnection(conn net.Conn, gb *GlobalBroadcaster, lci *int) error {
	*lci++                          // increment last connection id by 1
	remoteAddr := conn.RemoteAddr() // Get remote address of connection
	clientUpdate := ClientUpdate{
		ClientID:     "",
		ConnectionID: *lci,
		Channel:      make(chan ChatMessage),
		Action:       "connect",
		RemoteAddr:   remoteAddr,
		Connection:   conn,
	}
	dataBuf := []byte{}
	_, err := conn.Read(dataBuf)
	if err != nil {
		return fmt.Errorf("error reading from connnection: %w", err)
	}
	uErr := json.Unmarshal(dataBuf, &clientUpdate)
	if uErr != nil {
		return fmt.Errorf("error unmarshalling data into ChatClient instance: %w", uErr)
	}
	gb.ClientUpdates <- clientUpdate
	newClient := gb.GetChatClient(*lci)
	ManageChatClient(newClient, gb)

	return nil
}

func main() {
	http.HandleFunc("/ws", websocketHandler)

	// Listen on TCP port 8000 on all interfaces.
	ActiveConnections := []*ChatClient{}
	if len(ActiveConnections) < 1 {
		fmt.Printf("\r")
	}

	lastConnectionID := 0
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	broadcaster := GlobalBroadcaster{}

	fmt.Println("Chat server started on port 8000")
	go ManageGlobalBroadcaster(&broadcaster)
	// Accept new connections in a loop.
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		// Handle each connection in a new goroutine.
		go handleConnection(conn, &broadcaster, &lastConnectionID)
	}
}
