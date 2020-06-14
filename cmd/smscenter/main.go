package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/cbi-sh/smscenter/internal/app/server"
)

func main() {
	s := server.NewServer("192.168.0.2:3736")
	defer s.Close()

	fmt.Println(s.Addr())

	c, err := net.Dial("tcp", s.Addr())
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	time.Sleep(1<<63 - 1)

	// rw := newConn(c)
	// // bind
	// p := pdu.NewBindTransmitter()
	// f := p.Fields()
	// f.Set(pdufield.SystemID, "client")
	// f.Set(pdufield.Password, "secret")
	// f.Set(pdufield.InterfaceVersion, 0x34)
	// if err = rw.Write(p); err != nil {
	// 	t.Fatal(err)
	// }
	// // bind resp
	// resp, err := rw.Read()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// id, ok := resp.Fields()[pdufield.SystemID]
	// if !ok {
	// 	t.Fatalf("missing system_id field: %#v", resp)
	// }
	// if id.String() != "smpptest" {
	// 	t.Fatalf("unexpected system_id: want smpptest, have %q", id)
	// }
	// // submit_sm
	// p = pdu.NewSubmitSM(nil)
	// f = p.Fields()
	// f.Set(pdufield.SourceAddr, "foobar")
	// f.Set(pdufield.DestinationAddr, "bozo")
	// f.Set(pdufield.ShortMessage, pdutext.Latin1("Lorem ipsum"))
	// if err = rw.Write(p); err != nil {
	// 	t.Fatal(err)
	// }
	// // same submit_sm
	// r, err := rw.Read()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// want, have := *p.Header(), *r.Header()
	// if want != have {
	// 	t.Fatalf("unexpected header: want %#v, have %#v", want, have)
	// }
	// for k, v := range p.Fields() {
	// 	vv, exists := r.Fields()[k]
	// 	if !exists {
	// 		t.Fatalf("unexpected fields: want %#v, have %#v",
	// 			p.Fields(), r.Fields())
	// 	}
	// 	if !bytes.Equal(v.Bytes(), vv.Bytes()) {
	// 		t.Fatalf("unexpected field data: want %#v, have %#v",
	// 			v, vv)
	// 	}
	// }
	// // submit_sm + tlv field
	// p = pdu.NewSubmitSM(pdutlv.Fields{ pdutlv.TagReceiptedMessageID: pdutlv.CString("xyz123") })
	// f = p.Fields()
	// f.Set(pdufield.SourceAddr, "foobar")
	// f.Set(pdufield.DestinationAddr, "bozo")
	// f.Set(pdufield.ShortMessage, pdutext.Latin1("Lorem ipsum"))
	// if err = rw.Write(p); err != nil {
	// 	t.Fatal(err)
	// }
	// // same submit_sm
	// r, err = rw.Read()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// want, have = *p.Header(), *r.Header()
	// if want != have {
	// 	t.Fatalf("unexpected header: want %#v, have %#v", want, have)
	// }
	// for k, v := range p.Fields() {
	// 	vv, exists := r.Fields()[k]
	// 	if !exists {
	// 		t.Fatalf("unexpected fields: want %#v, have %#v",
	// 			p.Fields(), r.Fields())
	// 	}
	// 	if !bytes.Equal(v.Bytes(), vv.Bytes()) {
	// 		t.Fatalf("unexpected field data: want %#v, have %#v",
	// 			v, vv)
	// 	}
	// }
}
