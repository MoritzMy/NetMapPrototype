package icmp

import (
	"fmt"
	"testing"
)

func TestMarshal(t *testing.T) {
	packet := NewEchoICMPPacket(0, 0, []byte(""))
	cp := append([]EchoICMPPacket{}, packet)[0]
	b, err := Marshal(&packet)

	if err != nil {
		t.Error()
	}

	var u EchoICMPPacket

	err = Unmarshal(b, &u)
	if err != nil {
		t.Error()
	}

	if u.Equal(cp) {
		return
	}

	fmt.Println(u.String(), "\n", cp.String())

	t.Fail()

}
