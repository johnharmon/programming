package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

// region GLOBAL BROADCASTER

type GlobalBroadcaster struct {
	globalProducer chan ChatMessage
	//globalConsumers map[string]chan ChatMessage // map of all channels clients are consuming from, whenever anything is sent to the global producer, it replicates to all global consumers
	ClientUpdates         chan ClientUpdate     // channel that takes connection states, used to update the internal map of client channeld to send to
	globalClients         map[uint64]ChatClient // map of chat clients, try to coordinate ChatClient.ConnectionID with the map key
	globalClientUpdateMut sync.Mutex            // Mutex to lock client map as clients are added/removed
}

func (gb *GlobalBroadcaster) UpdateClient(cu ClientUpdate) bool {
	gb.globalClientUpdateMut.Lock()
	defer gb.globalClientUpdateMut.Unlock()
	if cu.Action == "connect" {
		//gb.globalClients[cu.ConnectionID] = NewChatClient(cu)
		gb.globalClients[cu.ConnectionID] = NewChatClient(cu, gb)
		return true
	} else if cu.Action == "disconnect" {
		gb.RemoveClient(cu)
		return false
	}
	return false
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

func (gb *GlobalBroadcaster) Broadcast(cm ChatMessage) {
	chatMessageBytes, err := json.Marshal(cm)
	if err != nil {
		fmt.Printf("Error mashalling message: %s\n", err)
	}
	for _, chatClient := range gb.globalClients {
		if chatClient.ConnectionID != cm.ConnectionID {
			chatClient.ServerToClient <- chatMessageBytes
		}
	}
}

// endregion

type ChatMessage struct { // All messages will be wrapped/unwrapped via this struct on either side
	Username     string    `json:"username"`
	Message      string    `json:"message"`
	Timestamp    time.Time `json:"timestamp"`
	ConnectionID uint64    `json:"connection_id"`
}

type ChatConnectionState struct { // Will be used to update the server that connections are being made/terminated
	Username     string `json:"username"`
	State        string `json:"connection_state"`
	ConnectionID uint64 `json:"connection_id"`
}

// region CHAT CLIENT
type ChatClient struct { // Represents a connected chat client
	Username              string           `json:"client_username"`
	ConnectionID          uint64           `json:"client_connection_id"`
	Connection            net.Conn         // Connection to read from
	LocalAddress          net.Addr         `json:"local_connection_address"`
	RemoteAddress         net.Addr         `json:"remote_connection_address"`
	MessageBuffer         []byte           // Any bytes that have been read since the last mesage delimiter was detected // Any bytes that have been read since the last mesage delimiter was detected // Any bytes that have been read since the last mesage delimiter was detected
	ClientToServer        chan ChatMessage // Messages coming FROM the client. Messages will be unmarshalled from the raw connection then put on this channel. SHOULD be the same as the global producer
	ServerToClient        chan []byte      // Messages going TO the client. Global messages will be marshalled first then put on this channel to be written to the raw conection
	ServerToClientMessage chan ChatMessage // For holding messages broadcast from the server in the ChatMessage format. I may want to move to all internals being of that format and so may need this channel
	ReadErrors            chan error
	WriteErrors           chan error
	MarshallErrorsOut     chan error    // Used to send errors resultsing from attempts to marshal outgoing chat messages to json
	MarshallErrorsIn      chan error    // Used to send errors resultsing from attempts to marshal incoming chat messages to json
	ClientDone            chan struct{} // when this channel is closed, all processes managing the client exit
}

func (cc *ChatClient) ReadMessage() {
	messageBytes := make([]byte, 1024)
	//messageSep := make([]byte, 1)
	//messageSep[0] = byte(messageEnd)
	for {
		_, err := cc.Connection.Read(messageBytes)
		if err != nil {
			cc.ReadErrors <- fmt.Errorf("error reading from client: %w", err)
		}
		endOfMessage := bytes.IndexByte(messageBytes, MessageSep)
		if endOfMessage != -1 {
			messageTail := messageBytes[endOfMessage+1:]
			messageBytes = append(cc.MessageBuffer, messageBytes[:endOfMessage]...)
			cc.MessageBuffer = messageTail
			chatMessage := ChatMessage{Username: cc.Username, ConnectionID: cc.ConnectionID}
			unmarshallError := json.Unmarshal(messageBytes, &chatMessage)
			if unmarshallError != nil {
				cc.ReadErrors <- fmt.Errorf("error unmarshalling message: %w", unmarshallError)
			}
			cc.ClientToServer <- chatMessage
		} else {
			cc.MessageBuffer = append(cc.MessageBuffer, messageBytes...)
		}
	}
}

func (cc *ChatClient) WriteMessage() {

}

//endregion

type ClientUpdate struct { // Represents a state change in chat clients, either connecting a new one or disconnecting an existing one
	ClientID     string `json:"client_id"`
	ConnectionID uint64
	Channel      chan ChatMessage
	Action       string `json:"action"` // "connect" or "disconnect"
	RemoteAddr   net.Addr
	Connection   net.Conn
}
