package torrent

import (
	"fmt"
	"io"
)

const (
	//保留位8位
	Reserved int = 8
	HsMsgLen int = SHALEN + IDLEN + Reserved
)

type HandshakeMsg struct {
	PreStr  string
	InfoSHA [SHALEN]byte
	PeerId  [IDLEN]byte
}

func NewHandShakeMsg(infoSHA, peerId [IDLEN]byte) *HandshakeMsg {
	return &HandshakeMsg{
		PreStr:  "BitTorrent protocol",
		InfoSHA: infoSHA,
		PeerId:  peerId,
	}
}

func WriteHandShake(w io.Writer, msg *HandshakeMsg) (int, error) {
	buf := make([]byte, len(msg.PreStr)+HsMsgLen+1) // 1 byte for prelen
	buf[0] = byte(len(msg.PreStr))
	curr := 1
	// 右边赋值给左边 返回右边的len
	curr += copy(buf[curr:], []byte(msg.PreStr))
	curr += copy(buf[curr:], make([]byte, Reserved))
	curr += copy(buf[curr:], msg.InfoSHA[:])
	curr += copy(buf[curr:], msg.PeerId[:])
	return w.Write(buf)
}

func ReadHandshake(r io.Reader) (*HandshakeMsg, error) {
	lenBuf := make([]byte, 1)
	// 按照lenBuf 的长度读取字节
	_, err := io.ReadFull(r, lenBuf)
	if err != nil {
		return nil, err
	}
	prelen := int(lenBuf[0])

	if prelen == 0 {
		err := fmt.Errorf("prelen cannot be 0")
		return nil, err
	}

	msgBuf := make([]byte, HsMsgLen+prelen)
	_, err = io.ReadFull(r, msgBuf)
	if err != nil {
		return nil, err
	}

	var peerId [IDLEN]byte
	var infoSHA [SHALEN]byte

	copy(infoSHA[:], msgBuf[prelen+Reserved:prelen+Reserved+SHALEN])
	copy(peerId[:], msgBuf[prelen+Reserved+SHALEN:])

	return &HandshakeMsg{
		PreStr:  string(msgBuf[0:prelen]),
		InfoSHA: infoSHA,
		PeerId:  peerId,
	}, nil
}
