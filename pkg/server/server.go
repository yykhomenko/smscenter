package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
)

// Default settings.
var (
	DefaultUser     = "user"
	DefaultPasswd   = "password"
	DefaultSystemID = "smpptest"
)

// HandlerFunc is the signature of a function passed to Server instances,
// that is called when client PDU messages arrive.
type HandlerFunc func(c smpp.Conn, m pdu.Body)

// Server is an SMPP server for testing purposes. By default it authenticate
// clients with the configured credentials, and echoes any other PDUs
// back to the client.
type Server struct {
	User    string
	Passwd  string
	TLS     *tls.Config
	Handler HandlerFunc

	conns []smpp.Conn
	mu    sync.Mutex
	l     net.Listener
}

// NewServer creates and initializes a new Server. Callers are supposed
// to call Close on that server later.
func NewServer(addr string) *Server {
	s := NewUnstartedServer(addr)
	s.Start()
	return s
}

// NewUnstartedServer creates a new Server with default settings, and
// does not start it. Callers are supposed to call Start and Close later.
func NewUnstartedServer(addr string) *Server {
	return &Server{
		User:    DefaultUser,
		Passwd:  DefaultPasswd,
		Handler: DefaultHandler,
		l:       newLocalListener(addr),
	}
}

func newLocalListener(addr string) net.Listener {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(fmt.Sprintf("server: failed to listen on a port: %v", err))
	}

	return l
}

// Start starts the server.
func (srv *Server) Start() {
	go srv.Serve()
}

// Addr returns the local address of the server, or an empty string
// if the server hasn't been started yet.
func (srv *Server) Addr() string {
	if srv.l == nil {
		return ""
	}
	return srv.l.Addr().String()
}

// Close stops the server, causing the accept loop to break out.
func (srv *Server) Close() {
	if srv.l == nil {
		panic("smpptest: server is not started")
	}
	srv.l.Close()
}

// Serve accepts new clients and handle them by authenticating the
// first PDU, expected to be a Bind PDU, then echoing all other PDUs.
func (srv *Server) Serve() {
	for {
		cli, err := srv.l.Accept()
		if err != nil {
			break // on srv.l.Close
		}

		c := newConn(cli)
		srv.conns = append(srv.conns, c)
		go srv.handle(c)
	}
}

// BroadcastMessage broadcasts a test PDU to the all bound clients
func (srv *Server) BroadcastMessage(p pdu.Body) {
	for i := range srv.conns {
		srv.conns[i].Write(p)
	}
}

// handle new clients.
func (srv *Server) handle(c *conn) {
	defer c.Close()
	if err := srv.auth(c); err != nil {
		if err != io.EOF {
			log.Println("smpptest: server auth failed:", err)
		}
		return
	}
	for {
		p, err := c.Read()
		if err != nil {
			if err != io.EOF {
				log.Println("smpptest: read failed:", err)
			}
			break
		}
		srv.Handler(c, p)
	}
}

// auth authenticate new clients.
func (srv *Server) auth(c *conn) error {
	p, err := c.Read()
	if err != nil {
		return err
	}

	var resp pdu.Body

	switch p.Header().ID {
	case pdu.BindTransmitterID:
		resp = pdu.NewBindTransmitterResp()
	case pdu.BindReceiverID:
		resp = pdu.NewBindReceiverResp()
	case pdu.BindTransceiverID:
		resp = pdu.NewBindTransceiverResp()
	default:
		return errors.New("unexpected pdu, want bind")
	}

	f := p.Fields()
	user := f[pdufield.SystemID]
	passwd := f[pdufield.Password]

	if user == nil || passwd == nil {
		return errors.New("malformed pdu, missing system_id/password")
	}
	if user.String() != srv.User {
		return errors.New("invalid user")
	}
	if passwd.String() != srv.Passwd {
		return errors.New("invalid passwd")
	}

	resp.Fields().Set(pdufield.SystemID, DefaultSystemID)

	return c.Write(resp)
}

// DefaultHandler is the default Server HandlerFunc, and echoes back
// any unexpected PDUs received.
func DefaultHandler(c smpp.Conn, m pdu.Body) {

	switch m.Header().ID {
	case pdu.SubmitSMID:
		r := pdu.NewSubmitSMResp()
		r.Header().Seq = m.Header().Seq
		r.Fields().Set(pdufield.MessageID, "40fe50ab")

		fmt.Println("REQ==================================")
		fmt.Println(
			m.Fields()[pdufield.SourceAddr],
			m.Fields()[pdufield.DestinationAddr],
			m.Fields()[pdufield.ShortMessage],
		)
		fmt.Println("===")
		fmt.Println(m.Header())
		fmt.Println(m.Fields())
		fmt.Println(m.TLVFields())

		fmt.Println("RESP=================================")
		fmt.Println(r.Header())
		fmt.Println(r.Fields())
		fmt.Println(r.TLVFields())
		fmt.Println("=====================================")

		c.Write(r)
	default:
		c.Write(m)
	}
}
