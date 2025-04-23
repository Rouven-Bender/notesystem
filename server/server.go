package main

import (
	"net"
	"log"
	"os"
	"os/signal"
	"syscall"
	"path"
	"bufio"

	"grimoire/rpc"
)
type daemon struct {
	ShutdownChannel chan bool
}
var grimoire daemon = daemon{
	ShutdownChannel: make(chan bool, 1),
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("couldn't get home directory of user: %v", err)
	}
	grimoireFolder := path.Join(home, ".local", "share", "grimoire")
	err = os.MkdirAll(grimoireFolder, 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("couldn't get the local share folder: %v", err)
	}

	socketFilePath := path.Join(grimoireFolder, "grimoired.sock")
	socket, err := net.Listen("unix", socketFilePath)
	if err != nil {
		log.Fatalf("error listening to unix domain socket: %v", err)
	}

	// Cleanup the socket file
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-c:
		case <-grimoire.ShutdownChannel:
		}
		os.Remove(socketFilePath)
		log.Println("cleaned up the socket")
		os.Exit(0)
	}()

	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Printf("accepting incoming connection failed: %v", err)
			continue
		}
		go func (conn net.Conn) {
			defer conn.Close()
			log.Println("connection accepted")
			defer log.Println("connection closed")

			scanner := bufio.NewScanner(conn)
			scanner.Split(rpc.Split)

			for scanner.Scan() {
				msg := scanner.Bytes()
				method, content, err := rpc.DecodeMessage(msg)
				if err != nil {
					log.Printf("error: %v", err)
					continue
				}
				if c := handleMessage(conn, method, content); c {
					break
				}
			}
		}(conn)
	}
}
func handleMessage(conn net.Conn, method string, content []byte) (closeQ bool){
	switch method {
	case "easter":
		rsp := rpc.BaseResponse {
			Response: "Happy Easter",
		}
		writeResponse(conn, rsp)
	case "close":
		log.Println("recieved close message")
		conn.Close()
		return true
	case "shutdown":
		log.Println("recieved shutdown message")
		sendShutdownEvent()
	default:
		rsp := rpc.BaseError {
			Error: "Unknown Method",
		}
		writeResponse(conn, rsp)
	}
	return false
}

func sendShutdownEvent() {
	grimoire.ShutdownChannel<-true
}

func writeResponse(conn net.Conn, msg any) {
	reply, _ := rpc.EncodeMessage(msg)
	conn.Write([]byte(reply))
}
