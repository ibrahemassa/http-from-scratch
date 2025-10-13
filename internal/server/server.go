package server

import (
	"fmt"
	"ibrahemassa/http_bootdev/internal/request"
	"ibrahemassa/http_bootdev/internal/response"
	// "io"
	"log"

	"net"
	// "ibrahemassa/http_bootdev/internal/request"
)

type ServerState int

const (
	StateRunning ServerState = 0
	StateClosed  ServerState = 1
)

type HandlerError struct{
	StatusCode response.StatusCode
	Body string
}

// type Handler func(w io.Writer, req *request.Request) *HandlerError
type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	Listener net.Listener
	State    ServerState
	HandlerFunction Handler
}

func newServer(h Handler) *Server {
	return &Server{
		nil,
		StateRunning,
		h,
	}
}

func Serve(port int, handler Handler) (*Server, error) {
	portStr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", portStr)
	if err != nil {
		return nil, err
	}

	server := newServer(handler)
	server.Listener = listener

	go server.Listen()

	return server, nil
}

func (s *Server) Close() error {
	s.State = StateClosed
	return s.Listener.Close()
}

func (s *Server) Listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil || s.State == StateClosed {
			break
		}
		go s.handle(conn)
	}
	s.Close()
}

// func (s *Server) handle(conn net.Conn) {
// 	defer conn.Close()
// 	// _, err := request.RequestFromReader(conn)
// 	// 	if err != nil{
// 	// 		log.Fatal(err)
// 	// 	}
//
// 	err := response.WriteStatusLine(conn, response.OK)
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}
//
// 	h := response.GetDefaultHeader(0)
// 	err = response.WriteHeaders(conn, h)
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}
// }

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	// h := response.GetDefaultHeader(0)

	req, err := request.RequestFromReader(conn)
	if err != nil{
		log.Fatal(err)
		return
	}
	w := response.NewWriter(conn)
	// if err != nil{
	// 	response.WriteStatusLine(conn, response.BadRequest)
	// 	response.WriteHeaders(conn, h)
	// 	return
	// }
	//
	// buff := bytes.NewBuffer([]byte{})
	s.HandlerFunction(w, req)

	// status := response.OK
	// var body []byte = nil
	//
	// if handleErr != nil{
	// 	status = handleErr.StatusCode
	// 	body = []byte(handleErr.Body)
	// } else {
	// 	body = buff.Bytes()
	// }
	// 
	//
	// length := fmt.Sprintf("%d", len(body))
	// h.Replace("content-length", length)
	// err = response.WriteStatusLine(conn, status)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	//
	// err = response.WriteHeaders(conn, h)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	//
	// conn.Write(body)
	// err = response.WriteBody(conn, *buff)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
}

// func handlerErrorWriter(w io.Writer, err HandlerError){
// 	erro := response.WriteStatusLine(w, err.StatusCode)
// 	if erro != nil {
// 		log.Fatal(err)
// 		return
// 	}
//
// 	buff := bytes.NewBuffer(err.Body)
// 	erro = response.WriteBody(w, *buff)
// 	if erro != nil {
// 		log.Fatal(err)
// 		return
// 	}
// }


