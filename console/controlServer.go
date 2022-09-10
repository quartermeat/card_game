package console

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func StartClient() {
	
}

//StartServer starts the control server
func StartServer() {
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

		fmt.Print("-> ", command, "\n")

		var response string
		switch command {
		case "exit":
			{
				fmt.Println("Exiting TCP server!")
				response = fmt.Sprintln("exit", command)
				connection.Write([]byte(response))
				connection.Close()
				return
			}
		default:
			{
				response = fmt.Sprintln("unknown command:", command)
				connection.Write([]byte(response))
			}
		}
	}
}
