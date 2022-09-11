package console

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const (
	Test string = "test"
	Poke string = "poke"
	Stop string = "stop"
)

type ConsoleTxCommand struct {
	Command            string
	consoleToInputChan chan<- IConsoleTxCommand
}

type IConsoleTxCommand interface {
	SendCommand(commandId string, connection net.Conn)
	GetCommand() string
}

func (command ConsoleTxCommand) GetCommand() string {
	return command.Command
}

func (command ConsoleTxCommand) SendCommand(commandId string, connection net.Conn) {
	inputHandlerCommand := ConsoleTxCommand{Command: commandId}
	select {
	case command.consoleToInputChan <- inputHandlerCommand:
		{
			response := fmt.Sprintln(commandId)
			connection.Write([]byte(response))
		}
	default:
		{
			// don't do anything
		}
	}
}

// StartServer starts the control server
func StartServer(writeInputHandler chan<- IConsoleTxCommand) {
	inputHandlerCommand := ConsoleTxCommand{consoleToInputChan: writeInputHandler}

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
				inputHandlerCommand.SendCommand(Poke, connection)
			}
		case Stop:
			{
				inputHandlerCommand.SendCommand(Stop, connection)
			}
		default:
			{
				response = fmt.Sprintln("unknown command:", command)
				connection.Write([]byte(response))
			}
		}
	}
}
