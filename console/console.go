package console

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func connect(address string) net.Conn {
	c, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return c
}

func RunConsole() {
	fmt.Printf("<-----AEM Console----->\n")
	connection := connect("127.0.0.1:1337")

	fmt.Print(">>")
	for {
		reader := bufio.NewReader(os.Stdin)

		text, _ := reader.ReadString('\n')
		fmt.Fprintf(connection, text+"\n")

		message, _ := bufio.NewReader(connection).ReadString('\n')
		// fmt.Print(message)

		switch strings.TrimSpace(string(message)) {
		case Stop:
			{
				fmt.Println("TCP client exiting...")
				connection.Close()
				return
			}
		case "connect":
			{
				fmt.Println("connecting..")
				connection = connect("127.0.0.1:1337")
			}
		default:
			{
				fmt.Print(">>")
			}
		}
	}
}

// for integration tests
func AutoRunConsole(command string) {
	connection := connect("127.0.0.1:1337")

	for {

		fmt.Fprintf(connection, command+"\n")

		message, _ := bufio.NewReader(connection).ReadString('\n')

		switch strings.TrimSpace(string(message)) {
		case Stop:
			{
				connection.Close()
				return
			}
		case "connect":
			{
				connection = connect("127.0.0.1:1337")
			}
		default:
			{
				return
			}
		}
	}
}
