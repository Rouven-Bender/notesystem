package main

import (
	"net"
	"log"
	"os"
	"os/signal"
	"syscall"
	"path"
	"bufio"
	"encoding/json"
	"fmt"
	"strings"

	"grimoire/rpc"
)

type daemon struct {
	ShutdownChannel chan bool
	NotebasePath string
	Socket net.Listener
	SqlDB *sqliteStore
}

func (d *daemon) sendShutdownEvent() {
	d.ShutdownChannel<-true
}

func (d *daemon) setupSocket() {
	// Setup of Socket
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
	d.Socket = socket

	// Cleanup the socket file
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-c:
		case <-d.ShutdownChannel:
		}
		os.Remove(socketFilePath)
		log.Println("cleaned up the socket")
		os.Exit(0)
	}()
}
func (d *daemon) listenforconnections() {
	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Printf("accepting incoming connection failed: %v", err)
			continue
		}
		go d.handleconnection(conn)
	}
}

func InitServer() *daemon {
	grimoire := &daemon{
		ShutdownChannel: make(chan bool, 1),
		NotebasePath: os.Getenv("NOTEBASEPATH"),
	}
	if grimoire.NotebasePath == "" {
		log.Fatalf("$NOTEBASEPATH is empty")
	}
	grimoire.SqlDB, err = NewSqliteStore(path.Join(grimoire.NotebasePath, "db.sqlite"))
	if err != nil {
		log.Fatal(err)
	}
	grimoire.setupSocket()
	grimoire.listenforconnections()
}

func main() {
	grimoire := InitServer()
}

func (d *daemon) handleconnection(conn net.Conn) {
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
		if c := d.handleMessage(conn, method, content); c {
			break
		}
	}
}

func (d *daemon) handleMessage(conn net.Conn, method rpc.MessageName, content []byte) (closeQ bool){
	switch method {
	case rpc.CLOSE:
		log.Println("recieved close message")
		conn.Close()
		return true
	case rpc.SHUTDOWN:
		log.Println("recieved shutdown message")
		d.sendShutdownEvent()
	case rpc.INDEX:
		d.handleIndexRequest(conn, content)
	case rpc.SEARCH:
		d.handleSearchRequest(conn, content)
	default:
		rsp := rpc.NewErrorResponse(fmt.Sprintf("Unknown Method: %s", method))
		writeResponse(conn, rsp)
	}
	return false
}

func (d *daemon) handleSearchRequest(conn net.Conn, content []byte) {
	//m := rpc.SearchMessage{}
	//json.Unmarshal(content, &m)
	//log.Printf("recieved search query: \"%s\"", m.Query)
	//q := bleve.NewMatchQuery(m.Query)
	//r := bleve.NewSearchRequest(q)
	//res, err := grimoire.Index.Search(r)
	//if err != nil {
	//	ersp := rpc.NewErrorResponse(fmt.Sprintf("ran into error while searching index: %s", err))
	//	writeResponse(conn, ersp)
	//	return
	//}
	//hitData := rpc.SearchResults{}
	//hitData.Type = rpc.SUCCESS
	//for _, hit := range res.Hits {
	//	data := rpc.Hit{}
	//	data.Filename = hit.ID
	//	log.Println(len(hit.Fragments))
	//	for fragmentField, fragments := range hit.Fragments {
	//		f := rpc.Fragment{}
	//		f.Field = fragmentField
	//		f.Matches = fragments
	//		data.Fragments = append(data.Fragments, f)
	//	}
	//	hitData.Hits = append(hitData.Hits, data)
	//}
	//writeResponse(conn, hitData)
}

func (d *daemon) handleIndexRequest(conn net.Conn, content []byte) {
	m := rpc.IndexMessage{}
	json.Unmarshal(content, &m)
	log.Printf("recieved index message for file: \"%s\"", m.Filename)
	err := func() error {
		if !strings.HasPrefix(m.Filename, grimoire.NotebasePath) {
			log.Printf("invalid filename: %s\n", m.Filename)
			return fmt.Errorf("invalid filename: %s", m.Filename)
		}
		log.Printf("valid path recieved. starting parsing")
		n, err := ParseFile(m.Filename)
		if err != nil {
			log.Printf("ran into error while parsing file(%s): %s", m.Filename, err)
			return fmt.Errorf("ran into error while parsing file(%s): %s", m.Filename, err)
		}
		err = d.SqlDB.index(m.Filename, n)
		if err != nil {
			log.Println("ran into error while indexing")
			return fmt.Errorf("ran into error while indexing: %s", err)
		}
		rsp := rpc.NewSuccessResponse("finished parsing note")
		writeResponse(conn, rsp)
		log.Printf("finished parsing: %s", m.Filename)
		return nil
	}()
	if err != nil {
		log.Println("got error:", err)
		ersp := rpc.NewErrorResponse(fmt.Sprint(err))
		writeResponse(conn, ersp)
	}
}

func writeResponse(conn net.Conn, msg any) {
	reply, _ := rpc.EncodeMessage(msg)
	conn.Write([]byte(reply))
}
