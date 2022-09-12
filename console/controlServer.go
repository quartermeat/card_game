package console

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// COMMAND IDs
const (
	Test string = "test"
	Poke string = "poke"
	Stop string = "stop"
)

// ConsoleTxTopic is the structure to hold a write
// channel and commandId to get info to other go routines
type ConsoleTxTopic struct {
	TopicId            string
	consoleToInputChan chan<- IConsoleTxTopic
}

// IConsoleTxTopic interface allows reading topic Id and async send the topic
type IConsoleTxTopic interface {
	SendTopic(topicId string, connection net.Conn)
	GetTopicId() string
}

// GetTopicId returns the command Id string
func (command ConsoleTxTopic) GetTopicId() string {
	return command.TopicId
}

// SendCommand takes a command id, and a connection
// it will asynchronously send a topic
func (command ConsoleTxTopic) SendTopic(topicId string, connection net.Conn) {
	inputHandlerCommand := ConsoleTxTopic{TopicId: topicId}
	select {
	case command.consoleToInputChan <- inputHandlerCommand:
		{
			response := fmt.Sprintln(topicId)
			connection.Write([]byte(response))
		}
	default:
		{
			// don't do anything
		}
	}
}

// StartServer starts a tcp server
func StartServer(writeInputHandler chan<- IConsoleTxTopic) {
	//register the sender
	inputHandlerCommand := ConsoleTxTopic{consoleToInputChan: writeInputHandler}

	PORT := ":" + "1337"
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	connection, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		netData, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		command := strings.TrimSpace(string(netData))

		var response string
		switch command {
		case Test:
			{
				response = fmt.Sprintln("executing on server")
				connection.Write([]byte(response))
			}
		case Poke:
			{
				inputHandlerCommand.SendTopic(Poke, connection)
			}
		case Stop:
			{
				inputHandlerCommand.SendTopic(Stop, connection)
			}
		default:
			{
				response = fmt.Sprintln("unknown command:", command)
				connection.Write([]byte(response))
			}
		}
	}
}
