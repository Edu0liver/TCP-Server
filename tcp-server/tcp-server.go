package tcpServer

import (
	"fmt"
	"net"
)

type Message struct {
	From    string
	Payload []byte
}

type Server struct {
	address string
	ln      net.Listener
	quitch  chan struct{}
	Msgch   chan Message
}

func NewServer(address string) *Server {
	return &Server{
		address: address,
		quitch:  make(chan struct{}),
		Msgch:   make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()

	<-s.quitch
	close(s.Msgch)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Accept Error:", err)
			continue
		}

		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read Error:", err)
			continue
		}

		s.Msgch <- Message{
			From:    conn.RemoteAddr().String(),
			Payload: buf[:n],
		}

		conn.Write([]byte(fmt.Sprintf("Thank you for the message (%s)", conn.RemoteAddr().String())))
	}
}
