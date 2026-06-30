// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rpc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestServerRegisterName(t *testing.T) {
	server := NewServer()
	service := new(testService)

	if err := server.RegisterName("test", service); err != nil {
		t.Fatalf("%v", err)
	}

	if len(server.services.services) != 2 {
		t.Fatalf("Expected 2 service entries, got %d", len(server.services.services))
	}

	svc, ok := server.services.services["test"]
	if !ok {
		t.Fatalf("Expected service calc to be registered")
	}

	wantCallbacks := 9
	if len(svc.callbacks) != wantCallbacks {
		t.Errorf("Expected %d callbacks for service 'service', got %d", wantCallbacks, len(svc.callbacks))
	}
}

func TestNewServerLimitDefaults(t *testing.T) {
	server := NewServer()

	if server.batchItemLimit != 0 {
		t.Fatalf("batch item limit: got %d, want 0", server.batchItemLimit)
	}
	if server.batchResponseLimit != 0 {
		t.Fatalf("batch response limit: got %d, want 0", server.batchResponseLimit)
	}
	if server.httpBodyLimit != defaultBodyLimit {
		t.Fatalf("http body limit: got %d, want %d", server.httpBodyLimit, defaultBodyLimit)
	}
	if server.wsReadLimit != wsMessageSizeLimit {
		t.Fatalf("websocket read limit: got %d, want %d", server.wsReadLimit, wsMessageSizeLimit)
	}
}

func TestServerSetBatchLimits(t *testing.T) {
	server := NewServer()

	server.SetBatchLimits(2, 128)

	if server.batchItemLimit != 2 {
		t.Fatalf("batch item limit: got %d, want 2", server.batchItemLimit)
	}
	if server.batchResponseLimit != 128 {
		t.Fatalf("batch response limit: got %d, want 128", server.batchResponseLimit)
	}
}

func TestServerBatchItemLimit(t *testing.T) {
	server := NewServer()
	service := new(batchLimitTestService)
	if err := server.RegisterName("batch", service); err != nil {
		t.Fatal(err)
	}
	server.SetBatchLimits(1, 0)

	resp := serveBatchRequest(t, server, `[
		{"jsonrpc":"2.0","id":1,"method":"batch_touch"},
		{"jsonrpc":"2.0","id":2,"method":"batch_touch"}
	]`)

	if service.calls.Load() != 0 {
		t.Fatalf("batch call count: got %d, want 0", service.calls.Load())
	}
	if len(resp) != 1 {
		t.Fatalf("response count: got %d, want 1", len(resp))
	}
	if string(resp[0].ID) != "1" {
		t.Fatalf("response id: got %s, want 1", resp[0].ID)
	}
	if resp[0].Error == nil {
		t.Fatalf("response error: got nil, want batch too large")
	}
	if resp[0].Error.Code != -32600 {
		t.Fatalf("response error code: got %d, want -32600", resp[0].Error.Code)
	}
	if resp[0].Error.Message != errMsgBatchTooLarge {
		t.Fatalf("response error message: got %q, want %q", resp[0].Error.Message, errMsgBatchTooLarge)
	}
}

func TestServerBatchResponseLimit(t *testing.T) {
	server := NewServer()
	if err := server.RegisterName("batch", new(batchLimitTestService)); err != nil {
		t.Fatal(err)
	}
	server.SetBatchLimits(0, 64)

	resp := serveBatchRequest(t, server, `[
		{"jsonrpc":"2.0","id":1,"method":"batch_large"}
	]`)

	if len(resp) != 1 {
		t.Fatalf("response count: got %d, want 1", len(resp))
	}
	if string(resp[0].ID) != "1" {
		t.Fatalf("response id: got %s, want 1", resp[0].ID)
	}
	if resp[0].Error == nil {
		t.Fatalf("response error: got nil, want batch response too large")
	}
	if resp[0].Error.Message != errMsgBatchResponseTooLarge {
		t.Fatalf("response error message: got %q, want %q", resp[0].Error.Message, errMsgBatchResponseTooLarge)
	}
}

type batchLimitTestService struct {
	calls atomic.Int64
}

func (s *batchLimitTestService) Touch() string {
	s.calls.Add(1)
	return "ok"
}

func (s *batchLimitTestService) Large() string {
	return strings.Repeat("x", 256)
}

func serveBatchRequest(t *testing.T, server *Server, request string) []*jsonrpcMessage {
	t.Helper()

	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	go server.ServeCodec(NewCodec(serverConn), 0)

	clientConn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if _, err := io.WriteString(clientConn, strings.TrimSpace(request)+"\n"); err != nil {
		t.Fatalf("write error: %v", err)
	}

	clientConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	line, err := bufio.NewReader(clientConn).ReadString('\n')
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	var resp []*jsonrpcMessage
	if err := json.Unmarshal([]byte(line), &resp); err != nil {
		t.Fatalf("decode response %q: %v", line, err)
	}
	return resp
}

func TestServer(t *testing.T) {
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatal("where'd my testdata go?")
	}
	for _, f := range files {
		if f.IsDir() || strings.HasPrefix(f.Name(), ".") {
			continue
		}
		path := filepath.Join("testdata", f.Name())
		name := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		t.Run(name, func(t *testing.T) {
			runTestScript(t, path)
		})
	}
}

func runTestScript(t *testing.T, file string) {
	server := newTestServer()
	content, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	go server.ServeCodec(NewCodec(serverConn), 0)
	readbuf := bufio.NewReader(clientConn)
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		switch {
		case len(line) == 0 || strings.HasPrefix(line, "//"):
			// skip comments, blank lines
			continue
		case strings.HasPrefix(line, "--> "):
			t.Log(line)
			// write to connection
			clientConn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if _, err := io.WriteString(clientConn, line[4:]+"\n"); err != nil {
				t.Fatalf("write error: %v", err)
			}
		case strings.HasPrefix(line, "<-- "):
			t.Log(line)
			want := line[4:]
			// read line from connection and compare text
			clientConn.SetReadDeadline(time.Now().Add(5 * time.Second))
			sent, err := readbuf.ReadString('\n')
			if err != nil {
				t.Fatalf("read error: %v", err)
			}
			sent = strings.TrimRight(sent, "\r\n")
			if sent != want {
				t.Errorf("wrong line from server\ngot:  %s\nwant: %s", sent, want)
			}
		default:
			panic("invalid line in test script: " + line)
		}
	}
}

// This test checks that responses are delivered for very short-lived connections that
// only carry a single request.
func TestServerShortLivedConn(t *testing.T) {
	server := newTestServer()
	defer server.Stop()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal("can't listen:", err)
	}
	defer listener.Close()
	go server.ServeListener(listener)

	var (
		request  = `{"jsonrpc":"2.0","id":1,"method":"rpc_modules"}` + "\n"
		wantResp = `{"jsonrpc":"2.0","id":1,"result":{"nftest":"1.0","rpc":"1.0","test":"1.0"}}` + "\n"
		deadline = time.Now().Add(10 * time.Second)
	)
	for i := 0; i < 20; i++ {
		conn, err := net.Dial("tcp", listener.Addr().String())
		if err != nil {
			t.Fatal("can't dial:", err)
		}
		defer conn.Close()
		conn.SetDeadline(deadline)
		// Write the request, then half-close the connection so the server stops reading.
		conn.Write([]byte(request))
		conn.(*net.TCPConn).CloseWrite()
		// Now try to get the response.
		buf := make([]byte, 2000)
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal("read error:", err)
		}
		if !bytes.Equal(buf[:n], []byte(wantResp)) {
			t.Fatalf("wrong response: %s", buf[:n])
		}
	}
}
