package main

import (
	"encoding/json"
	"fmt"
	"net"
)

func NewClientUpdate(clientID string, connectionID int, chatChan chan ChatMessage, action string) ClientUpdate {

}

func (gb *GlobalBroadcaster) addConnection(conn *net.Conn) error {
	return nil
}

func (gb *GlobalBroadcaster) Broadcast(cm ChatMessage) {
	for connID, chatClient := range gb.globalClients {
		if connID != cm.ConnectionID {
			chatClient.ConnectionChannel <- cm
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
		ConnectionChannel: make(chan ChatMessage),
		ConnectionID:      cu.ConnectionID,
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

func ManageGlobalBroadcaster(gb *GlobalBroadcaster) {
	for {
		select {
		case message := <-gb.globalProducer:
			go gb.Broadcast(message)
			//write message to all consumers
		case clientUpdate := <-gb.ClientUpdates:
			gb.UpdateClient(clientUpdate)
			// method for updating client map, closing/opening channels, etc
		}
	}
}

func handleConnection(conn net.Conn, gb *GlobalBroadcaster, lci *int) error {
	clientUpdate := ClientUpdate{
		ClientID:     "",
		ConnectionID: *lci + 1,
		Channel:      make(chan ChatMessage),
		Action:       "connect",
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

	return nil

}

func main() {

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
