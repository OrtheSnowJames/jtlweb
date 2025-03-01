package jtltp

import (
	"net"
	"testing"
	"time"
	"sync"
)

func TestJtltpFetch(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	server := NewJtltpServer(listener, "localhost", []string{"test"})
	address := listener.Addr().String()
	
	var wg sync.WaitGroup
	wg.Add(1)
	
	// Start server in goroutine
	go func() {
		defer wg.Done()
		if err := server.AwaitConnection(); err != nil {
			t.Errorf("Server connection error: %v", err)
			return
		}
		msg := server.AwaitMessage()
		if msg["JTLTP-GET"] == "test" {
			server.SendGood("test document", "jtl")
		} else {
			server.Send404()
		}
	}()

	// Test successful fetch
	result, err := JtltpFetch(address, "test")
	if err != nil {
		t.Fatalf("Failed to fetch document: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result but got nil")
	}
	if result["JTLTP-STATUS"] != "200" || result["JTLTP-TYPE"] != "jtl" || result["JTLTP"] != "test document" {
		t.Errorf("Unexpected result for valid request: %v", result)
	}

	wg.Wait()
	
	// Test 404 case
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.AwaitConnection(); err != nil {
			t.Errorf("Server connection error: %v", err)
			return
		}
		msg := server.AwaitMessage()
		if msg["JTLTP-GET"] == "test" {
			server.SendGood("test document", "jtl")
		} else {
			server.Send404()
		}
	}()

	result, err = JtltpFetch(address, "unknown")
	if err != nil {
		t.Fatalf("Failed to fetch document for invalid request: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result but got nil")
	}
	if result["JTLTP-STATUS"] != "404" || result["JTLTP-TYPE"] != "jtl" || result["JTLTP"] != "Not Found" {
		t.Errorf("Unexpected result for invalid request: %v", result)
	}

	wg.Wait()
}

func TestJtltpServer(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	server := NewJtltpServer(listener, "localhost", []string{"test"})
	address := listener.Addr().String()

	var wg sync.WaitGroup
	wg.Add(1)

	// Start server handler
	go func() {
		defer wg.Done()
		if err := server.AwaitConnection(); err != nil {
			t.Errorf("Server connection error: %v", err)
			return
		}
		msg := server.AwaitMessage()
		if msg["JTLTP-GET"] == "test" {
			server.SendGood("test document", "jtl")
		} else {
			server.Send404()
		}
	}()

	// Small delay to ensure server is ready
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}

	// Test good request
	conn.Write([]byte("JTLTP-GET=[test]"))
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		t.Fatalf("Failed to read from connection: %v", err)
	}

	response := string(buffer[:n])
	expected := "JTLTP-STATUS=[200] JTLTP-TYPE=[jtl] JTLTP=MSG=[test document]"
	if response != expected {
		t.Errorf("Expected %s but got %s", expected, response)
	}
	conn.Close()

	wg.Wait()

	wg.Add(1)
	// Start server again for bad request test
	go func() {
		defer wg.Done()
		if err := server.AwaitConnection(); err != nil {
			t.Errorf("Server connection error: %v", err)
			return
		}
		msg := server.AwaitMessage()
		if msg["JTLTP-GET"] == "test" {
			server.SendGood("test document", "jtl")
		} else {
			server.Send404()
		}
	}()

	// Test bad request
	conn, _ = net.Dial("tcp", address)
	conn.Write([]byte("JTLTP-GET=[unknown]"))
	n, _ = conn.Read(buffer)
	response = string(buffer[:n])
	expected = "JTLTP-STATUS=[404] JTLTP-TYPE=[jtl] JTLTP=MSG=[Not Found]"
	if response != expected {
		t.Errorf("Expected %s but got %s", expected, response)
	}
	conn.Close()

	wg.Wait()
}
