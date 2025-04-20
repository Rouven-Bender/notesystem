package main

import (
	"log"
	"net"
	"os"
	"path"

	"grimoire/rpc"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("couldn't get home directory of user: %v", err)
	}
	p := path.Join(home, ".local", "share", "grimoire", "grimoired.sock")
	conn, err := net.Dial("unix", p)
	if err != nil {
		log.Fatalf("couldn't connect to grimoire daemon: %v", err)
	}
	defer conn.Close()

	msg := rpc.EasterMessage{
		Method: "easter",
	}
	s, _ := rpc.EncodeMessage(msg)
	conn.Write([]byte(s))

	r := make([]byte, 128)
	conn.Read(r)
	os.Stdout.Write(r)
}
