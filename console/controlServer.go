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

type ConsoleCommand struct {
	Command string
}

// StartServer starts the control server
func StartServer(writeInputHandler chan<- ConsoleCommand) {
	inputHandlerCommand := ConsoleCommand{}

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
				response = fmt.Sprintln(Poke)
				inputHandlerCommand.Command = Poke
				select {
				case writeInputHandler <- inputHandlerCommand:
					{
						connection.Write([]byte(response))
					}
				default:
					{
						// don't do anything
					}
				}
			}
		case Stop:
			{
				response = fmt.Sprintln(Stop)
				inputHandlerCommand.Command = Stop
				select {
				case writeInputHandler <- inputHandlerCommand:
					{
						connection.Write([]byte(response))
					}
				default:
					{
						// don't do anything
					}
				}
			}
		default:
			{
				response = fmt.Sprintln("unknown command:", command)
				connection.Write([]byte(response))
			}
		}
	}
}
