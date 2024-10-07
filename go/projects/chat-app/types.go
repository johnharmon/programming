package main

import (
	"net"
	"sync"
)

type GlobalBroadcaster struct {
	globalProducer chan ChatMessage
	//globalConsumers map[string]chan ChatMessage // map of all channels clients are consuming from, whenever anything is sent to the global producer, it replicates to all global consumers
	ClientUpdates         chan ClientUpdate     // channel that takes connection states, used to update the internal map of client channeld to send to
	globalClients         map[uint64]ChatClient // map of chat clients, try to coordinate ChatClient.ConnectionID with the map key
	globalClientUpdateMut sync.Mutex            // Mutex to lock client map as clients are added/removed
}

type ChatMessage struct { // All messages will be wrapped/unwrapped via this struct on either side
	Username     string `json:"username"`
	Message      string `json:"message"`
	Timestamp    int64  `json:"timestamp"`
	ConnectionID uint64 `json:"connection_id"`
}

type ChatConnectionState struct { // Will be used to update the server that connections are being made/terminated
	Username     string `json:"username"`
	State        string `json:"connection_state"`
	ConnectionID uint64 `json:"connection_id"`
}

type ChatClient struct { // Represents a connected chat client
	Username       string           `json:"client_username"`
	ConnectionID   uint64           `json:"client_connection_id"`
	Connection     net.Conn         // Connection to read from
	LocalAddress   net.Addr         `json:"local_connection_address"`
	RemoteAddress  net.Addr         `json:"remote_connection_address"`
	MessageBuffer  []byte           // Any bytes that have been read since the last mesage delimiter was detected // Any bytes that have been read since the last mesage delimiter was detected // Any bytes that have been read since the last mesage delimiter was detected
	ClientToServer chan ChatMessage // Messages coming FROM the client. Messages will be unmarshalled from the raw connection then put on this channel. SHOULD be the same as the global producer
	ServerToClient chan []byte      // Messages going TO the client. Global messages will be marshalled first then put on this channel to be written to the raw conection
	ReadErrors     chan error
	WriteErrors    chan error
}

type ClientUpdate struct { // Represents a state change in chat clients, either connecting a new one or disconnecting an existing one
	ClientID     string `json:"client_id"`
	ConnectionID uint64
	Channel      chan ChatMessage
	Action       string `json:"action"` // "connect" or "disconnect"
	RemoteAddr   net.Addr
	Connection   net.Conn
}
