package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net"
	"net/http"
	"strings"
	"time"

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
	MessageSep = byte(0x00)
	//messageSep = byte(messageEnd)
)

// region CREATORS

func NewChatClient(cu ClientUpdate, gb *GlobalBroadcaster) ChatClient {
	cc := ChatClient{
		Username:       cu.ClientID,
		LocalAddress:   cu.Connection.LocalAddr(),
		RemoteAddress:  cu.Connection.RemoteAddr(),
		Connection:     cu.Connection,
		ClientToServer: gb.globalProducer,
		ServerToClient: make(chan []byte),
		ReadErrors:     make(chan error),
		WriteErrors:    make(chan error),
	}
	cc.ConnectionID = GenConnectionHash(cu.Connection)
	return cc
}

func GenConnectionHash(conn net.Conn) uint64 {
	// Split local and remote addresses along colons (once)
	localAddress := strings.SplitN(conn.LocalAddr().String(), ":", 1)
	remoteAddress := strings.SplitN(conn.RemoteAddr().String(), ":", 1)
	h := fnv.New64a() // Generate 64bit hash (so it can fit in uint64) to ID the connection
	// Create string from local and remote addresses, using colons to delimit sections
	hashString := fmt.Sprintf("%s:%s-%s:%s",
		localAddress[0], localAddress[1], remoteAddress[0], remoteAddress[1])
	h.Write([]byte(hashString)) // Write the formatted string to the hash buffer
	return h.Sum64()            // Returns a uint64 hash value
}

func MakeDisconnectMessage(cc *ChatClient) ChatMessage {
	disconnectMessage := ChatMessage{
		Username:     cc.Username,
		Message:      fmt.Sprintf("User %s has disconnected\n", cc.Username),
		Timestamp:    time.Now(),
		ConnectionID: cc.ConnectionID,
	}
	return disconnectMessage

}

func NewClientUpdate(clientID string, connectionID int, chatChan chan ChatMessage, action string) ClientUpdate {
	return ClientUpdate{}
}

// endregion

func websocketHandler(w http.ResponseWriter, r *http.Request) {
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

// region CLIENT IO
func ClientWriter(cc *ChatClient) {
	for {
		select {
		case serverMessage := <-cc.ServerToClientMessage:
			serverMessageRaw, mErr := json.Marshal(&serverMessage)
			if mErr != nil {
				cc.MarshallErrorsOut <- mErr
			} else {
				serverMessageRaw = append(serverMessageRaw, MessageSep)
				bytesWritten, wErr := cc.Connection.Write(serverMessageRaw)
				if wErr != nil {
					cc.WriteErrors <- wErr
				}
				fmt.Printf("%d bytes written to %s\n", bytesWritten, cc.Username)
			}
		case <-cc.ClientDone:
			return
		}
	}
}

func ClientReader(cc *ChatClient) {
	for {
		go cc.ReadMessage()
		select {
		case <-cc.ClientDone:
			return
		case serverMessage := <-cc.ServerToClient:
			serverMessage = append(serverMessage, MessageSep)
			_, err := cc.Connection.Write(serverMessage)
			if err != nil {
				cc.WriteErrors <- err
			}
		}
	}
}

func readUntilSep(data []byte, sep byte) (dataRead []byte, dataLeft []byte, sepEncountered bool, dErr error) {
	indexOfSep := bytes.IndexByte(data, sep)
	sepEncountered = false
	dErr = nil
	if indexOfSep == -1 {
		dataRead = data
	} else {
		dataRead = data[:indexOfSep]
		dataLeft = data[indexOfSep:]
		sepEncountered = true
	}
	return dataRead, dataLeft, sepEncountered, dErr

}

func ClientListener(cc *ChatClient) {
	for {
		go cc.ReadMessage()
		select {
		case <-cc.ClientDone:
			return
		case serverMessage := <-cc.ServerToClient:
			serverMessage = append(serverMessage, MessageSep)
			_, err := cc.Connection.Write(serverMessage)
			if err != nil {
				cc.WriteErrors <- err
			}
		}
	}
}

// endregion

// region MANAGERS

func ManageChatClient(cc *ChatClient, gb *GlobalBroadcaster) {
	go ClientListener(cc)
	for {
		select {
		case clientMessage := <-cc.ClientToServer:
			gb.globalProducer <- clientMessage
			continue
		case <-cc.ClientDone:
			gb.globalProducer <- MakeDisconnectMessage(cc)
		}
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

// endregion

func handleConnection(conn net.Conn, gb *GlobalBroadcaster) error {
	remoteAddr := conn.RemoteAddr() // Get remote address of connection
	clientUpdate := ClientUpdate{
		ClientID:     "",
		ConnectionID: GenConnectionHash(conn),
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
	newClient := gb.GetChatClient(clientUpdate.ConnectionID)
	ManageChatClient(newClient, gb)

	return nil
}

// region MAIN
func main() {
	http.HandleFunc("/ws", websocketHandler)

	// Listen on TCP port 8000 on all interfaces.
	ActiveConnections := []*ChatClient{}
	if len(ActiveConnections) < 1 {
		fmt.Printf("\r")
	}
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
		go handleConnection(conn, &broadcaster)
	}
}

// endregion

// region DOCUMENTATION
/*
##########FUNCTION FLOW FOR CONNECTING CLIENTS##########

main() -> handleConnection(net.Conn, *GlobalBroadcaster) { *GlobalBraodcaster.ClientUPdates <- ClientUpdate{}} -> ManageChatClient(newClient) \
{*GlobalBroadcaster.globalProducer <- *ChatClient.ClientToServer }


main()
|
V
handleConnection(net.Conn, *GlobalBroadcaster) {
	*GlobalBraodcaster.ClientUPdates <- ClientUpdate{}
}
|
V
ManageChatClient(ChatClient) {
	*GlobalBroadccaster.globalProducer <- *ChatClient.ClientToServer
}
|
V



########################################################

*/

// endregion
