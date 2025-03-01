package jtltp

import (
	"errors"
	"net"
	"regexp"
	"fmt"
	"strings"
	"time"
)

// client
// timeout takes a function and a timeout in milliseconds, times out the function if it takes too long and returns error

func JtltpFetch(address string, what string) (map[string]string, error) {
	// connect to the server
	conn, _ := net.Dial("tcp", address)

	// send the request
	timeoutnum := 30
	conn.SetReadDeadline(time.Now().Add(time.Duration(timeoutnum) * time.Second))

	// read the response
	buffer := make([]byte, 1024)

	conn.Write([]byte("JTLTP-GET=[" + what + "]"))
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	// parse the response
	message := string(buffer[:n])
	// response format: JTLTP-STATUS=[200] JTLTP-TYPE=[jtl] JTLTP=MSG=["jtl document goes here"]
	fmt.Println(message)
	re_get := regexp.MustCompile(`JTLTP-STATUS=\[([0-9]+)\]\s+JTLTP-TYPE=\[([a-zA-Z0-9_-]+)\]\s+JTLTP=MSG=\[([^\]]+)\]`)
	match := re_get.FindStringSubmatch(message)
	if len(match) < 1 {
		return nil, errors.New("bad response")
	}

	// return the message
	return map[string]string{"JTLTP-STATUS": match[1], "JTLTP-TYPE": match[2], "JTLTP": match[3]}, nil
}

// server
type jtltpServer struct {
	listener   net.Listener
	conn       net.Conn
	clientAddr string
	demand     []string
}

func NewJtltpServer(listener net.Listener, clientAddr string, demand []string) *jtltpServer {
	return &jtltpServer{
		listener:   listener,
		clientAddr: clientAddr,
		demand:     demand,
	}
}

// ex JTLTP-STATUS=[200] JTLTP-TYPE=[jtl] JTLTP=MSG=["jtl document goes here"]
func (connection *jtltpServer) SendGood(message string, doctype string) {
	go func() {
		connection.conn.Write([]byte("JTLTP-STATUS=[200] JTLTP-TYPE=[" + doctype + "] JTLTP=MSG=[" + message + "]"))
		// disconnect
		connection.conn.Close()
	}()

	// move the connection to the next demand
	if len(connection.demand) > 0 {
		connection.demand = connection.demand[1:]
	}
}

func (connection *jtltpServer) SendBad(message string, doctype string) {
	go func() {
		connection.conn.Write([]byte("JTLTP-STATUS=[400] JTLTP-TYPE=[" + doctype + "] JTLTP=MSG=[" + message + "]"))
		// disconnect
		connection.conn.Close()
	}()

	// move the connection to the next demand
	if len(connection.demand) > 0 {
		connection.demand = connection.demand[1:]
	}
}

func (connection *jtltpServer) Send404() {
	go func() {
		connection.conn.Write([]byte("JTLTP-STATUS=[404] JTLTP-TYPE=[jtl] JTLTP=MSG=[Not Found]"))
		// disconnect
		connection.conn.Close()
	}()

	// move the connection to the next demand
	if len(connection.demand) > 0 {
		connection.demand = connection.demand[1:]
	}
}

// not really recommended
func (connection *jtltpServer) SendRaw(message string, status string, msg string) {
	go func() {
		connection.conn.Write([]byte("JTLTP-STATUS=[" + status + "] JTLTP-TYPE=[" + msg + "] JTLTP=MSG=[" + message + "]"))
		// disconnect
		connection.conn.Close()
	}()

	// move the connection to the next demand
	if len(connection.demand) > 0 {
		connection.demand = connection.demand[1:]
	}
}

// not really recommended at all
func (connection *jtltpServer) SendRawer(message string) {
	go func() {
		connection.conn.Write([]byte(message))
		// disconnect
		connection.conn.Close()
	}()

	// move the connection to the next demand
	if len(connection.demand) > 0 {
		connection.demand = connection.demand[1:]
	}
}

func (connection *jtltpServer) AwaitMessage() map[string]string {
	if connection.conn == nil {
		return nil
	}
	// read the message
	buffer := make([]byte, 1024)
	n, err := connection.conn.Read(buffer)
	if err != nil {
		return nil
	}
	// parse the message
	message := string(buffer[:n])
	if !strings.HasPrefix(message, "JTLTP-GET=[") {
		return connection.AwaitMessage()
	}

	re_get := regexp.MustCompile(`JTLTP-GET=\[(.+?)\]`)
	match := re_get.FindStringSubmatch(message)
	if len(match) < 2 {
		connection.Send404()
		return nil
	}

	what_to_get := match[1]

	// return the message
	return map[string]string{"JTLTP-GET": what_to_get}
}

func (connection *jtltpServer) AwaitConnection() error {
	conn, err := connection.listener.Accept()
	if err != nil {
		return err
	}
	connection.conn = conn
	return nil
}
