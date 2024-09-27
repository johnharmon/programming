package main

import (
	"net"
	"sync"
)

type GlobalBroadcaster struct {
	globalProducer chan ChatMessage
	//globalConsumers map[string]chan ChatMessage // map of all channels clients are consuming from, whenever anything is sent to the global producer, it replicates to all global consumers
	ClientUpdates         chan ClientUpdate  // channel that takes connection states, used to update the internal map of client channeld to send to
	globalClients         map[int]ChatClient // map of chat clients, try to coordinate ChatClient.ConnectionID with the map key
	globalClientUpdateMut sync.Mutex
}

type ChatMessage struct { // All messages will be wrapped/unwrapped via this struct on either side
	Username     string `json:"username"`
	Message      string `json:"message"`
	Timestamp    int64  `json:"timestamp"`
	ConnectionID int    `json:"connection_id"`
}

type ChatConnectionState struct { // Will be used to update the server that connections are being made/terminated
	Username     string `json:"username"`
	State        string `json:"connection_state"`
	ConnectionID int    `json:"connection_id"`
}

type ChatClient struct { // Represents a connected chat client
	Username          string
	ConnectionAddress net.IP
	ConnectionChannel chan ChatMessage
	ConnectionID      int
}

type ClientUpdate struct { // Represents a state change in chat clients, either connecting a new one or disconnecting an existing one
	ClientID     string
	ConnectionID int
	Channel      chan ChatMessage
	Action       string // "connect" or "disconnect"
}
