package main

import (
	"log"
	"net"
	"os"
	"path"
	"flag"

	"grimoire/rpc"
)
type grimoireClient struct {
	Connection net.Conn
}
func (c *grimoireClient) Disconnect() {
	msgB, _ := rpc.EncodeMessage(rpc.CloseMessage)
	c.Connection.Write([]byte(msgB))
	log.Println("closed connection")
	c.Connection.Close()
}

var client grimoireClient

func main() {
	var msg = flag.String("msg", "", "Which message to send")
	flag.Parse()

	client = grimoireClient {
		Connection: connectToServer(),
	}
	defer client.Disconnect()

	switch *msg {
	case "easter":
		m := rpc.EasterMessage{
			Method: "easter",
		}
		s, _ := rpc.EncodeMessage(m)
		client.Connection.Write([]byte(s))
		r := make([]byte, 128)
		client.Connection.Read(r)
		os.Stdout.Write(r)
	case "shutdown":
		m := rpc.ShutdownMessage
		s, _ := rpc.EncodeMessage(m)
		client.Connection.Write([]byte(s))
	}
}

func connectToServer() (net.Conn) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("couldn't get home directory of user: %v", err)
	}
	p := path.Join(home, ".local", "share", "grimoire", "grimoired.sock")
	conn, err := net.Dial("unix", p)
	if err != nil {
		log.Fatalf("couldn't connect to grimoire daemon: %v", err)
	}
	return conn
}
