package server

import (
	"bytes"
	"net"
	"testing"

	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutlv"
)

func TestServer(t *testing.T) {
	s := NewServer("localhost:0")
	defer s.Close()

	c, err := net.Dial("tcp", s.Addr())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	rw := newConn(c)

	// bind
	p := pdu.NewBindTransmitter()
	f := p.Fields()
	f.Set(pdufield.SystemID, "user")
	f.Set(pdufield.Password, "password")
	f.Set(pdufield.InterfaceVersion, 0x34)
	if err = rw.Write(p); err != nil {
		t.Fatal(err)
	}

	// bind resp
	resp, err := rw.Read()
	if err != nil {
		t.Fatal(err)
	}
	id, ok := resp.Fields()[pdufield.SystemID]
	if !ok {
		t.Fatalf("missing system_id field: %#v", resp)
	}
	if id.String() != "smpptest" {
		t.Fatalf("unexpected system_id: want smpptest, have %q", id)
	}

	// submit_sm
	p = pdu.NewSubmitSM(nil)
	f = p.Fields()
	f.Set(pdufield.SourceAddr, "777")
	f.Set(pdufield.DestinationAddr, "380671112222")
	f.Set(pdufield.ShortMessage, pdutext.Latin1("Lorem ipsum"))
	if err = rw.Write(p); err != nil {
		t.Fatal(err)
	}

	// submit_sm_resp
	resp = pdu.NewSubmitSMResp()
	resp.Header().Seq = p.Header().Seq
	resp.Header().Len = 0x19

	r, err := rw.Read()
	if err != nil {
		t.Fatal(err)
	}
	match(t, resp, r)

	// submit_sm + tlv field
	p = pdu.NewSubmitSM(pdutlv.Fields{pdutlv.TagReceiptedMessageID: pdutlv.CString("xyz123")})
	f = p.Fields()
	f.Set(pdufield.SourceAddr, "foobar")
	f.Set(pdufield.DestinationAddr, "bozo")
	f.Set(pdufield.ShortMessage, pdutext.Latin1("Lorem ipsum"))
	if err = rw.Write(p); err != nil {
		t.Fatal(err)
	}
	// same submit_sm
	r, err = rw.Read()
	if err != nil {
		t.Fatal(err)
	}

	resp = pdu.NewSubmitSMResp()
	resp.Header().Seq = p.Header().Seq
	resp.Header().Len = 0x19

	match(t, resp, r)
}

func match(t *testing.T, resp pdu.Body, r pdu.Body) (pdu.Header, pdu.Header) {
	want, have := *resp.Header(), *r.Header()
	if want != have {
		t.Fatalf("unexpected header: want %#v, have %#v", want, have)
	}
	for k, v := range resp.Fields() {
		vv, exists := r.Fields()[k]
		if !exists {
			t.Fatalf("unexpected fields: want %#v, have %#v",
				resp.Fields(), r.Fields())
		}
		if !bytes.Equal(v.Bytes(), vv.Bytes()) {
			t.Fatalf("unexpected field data: want %#v, have %#v",
				v, vv)
		}
	}
	return want, have
}
