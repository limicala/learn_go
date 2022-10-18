package rpc

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"learnset/rpc/codec"
	"log"
	"net"
	"reflect"
	"testing"
	"time"
)

// https://learnku.com/docs/the-way-to-go/1211-uses-gob-to-transmit-data/3671
type P struct {
	X, Y, Z int
	Name    string
}

type Q struct {
	X, Y *int
	Name string
}

/* additional link
- https://go.dev/doc/faq#Functions_methods
*/

func Test_Gob(t *testing.T) {
	// Initialize the encoder and decoder.  Normally enc and dec would be
	// bound to network connections and the encoder and decoder would
	// run in different processes.
	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	dec := gob.NewDecoder(&network) // Will read from network.
	// Encode (send) the value.
	err := enc.Encode(P{3, 4, 5, "Pythagoras"})
	if err != nil {
		log.Fatal("encode error:", err)
	}
	// Decode (receive) the value.
	var q Q
	err = dec.Decode(&q)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	fmt.Printf("%q: {%d,%d}\n", q.Name, *q.X, *q.Y)
}

func Test_Gob2(t *testing.T) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	dec := gob.NewDecoder(&network)

	err := enc.Encode("test code")
	if err != nil {
		log.Fatal("encode error", err)
	}

	r := reflect.New(reflect.TypeOf(""))
	err = dec.Decode(r.Interface())
	log.Println(err, r.Elem())
}

func startServer(addr chan string) {
	/* https://pkg.go.dev/net#Listen
	:0 meaning port number is automatically chosen
	*/
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", listener.Addr())
	addr <- listener.Addr().String()
	Accept(listener)
}

// go test -count=1 -v -run Test_Rpc learnset/rpc
func Test_Rpc(t *testing.T) {
	addr := make(chan string)
	go startServer(addr)

	time.Sleep(time.Second)
	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	_ = json.NewEncoder(conn).Encode(DefaultOption)
	cc := codec.NewGobCodec(conn)

	// Client send Option
	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("rpc req %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply:", reply)
	}
}
