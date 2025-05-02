package main

import (
	"log"
	"net"
	"os"
	"path"
	"flag"
	"bufio"

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
	var indexFilename = flag.String("f", "", "Which file to index")
	var searchQuery = flag.String("q", "", "search query")
	flag.Parse()

	client = grimoireClient {
		Connection: connectToServer(),
	}
	defer client.Disconnect()

	switch *msg {
	case "shutdown":
		m := rpc.ShutdownMessage
		s, _ := rpc.EncodeMessage(m)
		client.Connection.Write([]byte(s))
	case "index":
		m := rpc.NewIndexMessage(*indexFilename)
		s, _ := rpc.EncodeMessage(m)
		client.Connection.Write([]byte(s))
		_, c := client.readResponse()
		os.Stdout.Write(c)
	case "search":
		m := rpc.NewSearchMessage(*searchQuery)
		s, _ := rpc.EncodeMessage(m)
		client.Connection.Write([]byte(s))
		_, c := client.readResponse()
		os.Stdout.Write(c)
	}
}

func (c *grimoireClient) readResponse() (rpc.ResponseName, []byte){
	scanner := bufio.NewScanner(c.Connection)
	scanner.Split(rpc.Split)
	for scanner.Scan() {
		msg := scanner.Bytes()
		tzpe, content, err := rpc.DecodeResponse(msg)
		if err != nil {
			log.Printf("error: %v", err)
			return "", nil
		}
		return tzpe, content
	}
	return "", nil
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
