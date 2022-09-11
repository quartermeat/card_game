package console

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func RunConsole() {

	fmt.Printf("<-----AEM Console----->\n")
	connection := Connect("127.0.0.1:1337")

	fmt.Print(">>")
	for {
		reader := bufio.NewReader(os.Stdin)

		text, _ := reader.ReadString('\n')
		fmt.Fprintf(connection, text+"\n")

		message, _ := bufio.NewReader(connection).ReadString('\n')
		// fmt.Print(message)

		switch strings.TrimSpace(string(message)) {
		case "exit":
			{
				fmt.Println("TCP client exiting...")
				connection.Close()
			}
		case "connect":
			{
				fmt.Println("connecting..")
				connection = Connect("127.0.0.1:1337")
			}
		default:
			{
				fmt.Print(">>")
			}
		}
	}
}

func Connect(address string) net.Conn {
	c, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return c
}
