// Package 'console' provides a TCP server, a TCP client console to interact with the server,
// and the interface to send messaging to the input handler
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

// RunConsole starts the tcp listener and handles user input from the console
// TODO: need to redirect all fmt.Print throughout to the Errors -> may have to reformat
// nomenclature on Errors to DebugLog
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

// AutoRunConsole is a stripped down console, not really ment for user input,
// but to send topics
func AutoRunConsole(topicId string) {
	connection := connect("127.0.0.1:1337")

	for {

		fmt.Fprintf(connection, topicId+"\n")

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
