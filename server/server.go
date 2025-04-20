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
		<-c
		os.Remove(socketFilePath)
		os.Exit(0)
	}()

	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Println("accepting incoming connection failed: %v", err)
			continue
		}
		go func (conn net.Conn) {
			log.Println("connection accepted")
			defer conn.Close()

			scanner := bufio.NewScanner(conn)
			scanner.Split(rpc.Split)

			for scanner.Scan() {
				msg := scanner.Bytes()
				method, content, err := rpc.DecodeMessage(msg)
				if err != nil {
					log.Printf("error: %v", err)
					continue
				}
				handleMessage(conn, method, content)
			}
		}(conn)
	}
}
func handleMessage(conn net.Conn, method string, content []byte) {
	switch method {
	case "easter":
		rsp := rpc.BaseResponse {
			Response: "Happy Easter",
		}
		writeResponse(conn, rsp)
	default:
		// TODO: make error for invalid method
	}
}

func writeResponse(conn net.Conn, msg any) {
	reply, _ := rpc.EncodeMessage(msg)
	conn.Write([]byte(reply))
}
