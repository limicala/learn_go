package rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"learnset/rpc/codec"
	"log"
	"net"
	"reflect"
)

/*
	Option: Json Type {MagicNumber: xxx, CodecType: xxx}
	Header: Fixed Length, | ServerMethod | Seq Num | ... |
	Body: Contains Request and Response

	rpc message datagram
	| Option | Header1 | Body1 | Header2 | Body2 | ...
*/

/* ref
- https://research.swtch.com/interfaces
*/

type A interface {
	Call()
	Deep()
}

type Option struct {
	MagicNumber int
	CodecType   codec.Type
}

type Server struct {
}

type request struct {
	h            *codec.Header
	argv, replyv reflect.Value
}

func NewServer() *Server {
	return &Server{}
}

var DefaultOption = &Option{
	MagicNumber: 0x4b3d22,
	CodecType:   codec.GobType,
}

func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		return nil, err
	}
	req := &request{h: &h}
	req.argv = reflect.New(reflect.TypeOf(""))
	if err := cc.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err:", err)
	}
	return req, nil
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}) {
	if err := cc.Write(h, body); err != nil {
		log.Fatal(err)
	}
	log.Println("rpc server: resp", h, body)
}

func (server *Server) handleRequest(cc codec.Codec, r *request) {
	r.replyv = reflect.ValueOf(fmt.Sprintf("rpc resp %d", r.h.Seq))
	server.sendResponse(cc, r.h, r.replyv.Interface())
}

func (server *Server) serveCodec(cc codec.Codec) {
	for {
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, struct{}{})
			continue
		}
		server.handleRequest(cc, req)
	}
}

func (server *Server) serveConn(conn io.ReadWriteCloser) {
	var option Option
	if err := json.NewDecoder(conn).Decode(&option); err != nil {
		log.Fatal(err)
		return
	}
	f := codec.NewCodecFuncMap[option.CodecType]
	if f == nil {
		return
	}
	server.serveCodec(f(conn))
}

func (server *Server) Accept(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			return
		}
		go server.serveConn(conn)
	}
}

var s = NewServer()

func Accept(listener net.Listener) {
	s.Accept(listener)
}
