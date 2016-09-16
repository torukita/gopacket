// Copyright 2016 torukita
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

package layers

import (
    "fmt"
	"encoding/binary"
    "github.com/google/gopacket"
)

/* GTPv2 is not still implemented */


type GTPv1 struct {
	BaseLayer
	// Header Fields
	Version                        uint8   // 3bit
	ProtocolType                   uint8   // 1bit
	Reserved                       uint8   // 1bit
	ExtentionHeaderFlag            uint8   // 1bit 
	SequenceNumberFlag             uint8   // 1bit
	NPDUNumberFlag                 uint8   // 1bit
	MessageType                    uint8   // 8bit
	MessageLength                  uint16  // 16bit
	TEID                           uint32  // 32bit
	SequenceNumber                 uint16  // 16bit
	NPDUNumber                     uint8   // 8bit
    NextExtentionHeaderType        uint8   // 8bit
	restOfData                     []byte
}

/* Layer Interface */
func (g *GTPv1) LayerType() gopacket.LayerType { return LayerTypeGTPv1 }

/* func (g *GTPv1) LayerContents() []byte {} */

//DecodeFunc wraps a function to make it a Decoder.
func decodeGTPv1(data []byte, p gopacket.PacketBuilder) error {
	gtp := &GTPv1{}
	err := gtp.DecodeFromBytes(data, p)
	if err != nil {
		fmt.Println("ERR")
		return nil
	}
    p.AddLayer(gtp)
	/* Should be LayerTypeIPvr? */
    //return p.NextDecoder(gopacket.LayerTypePayload)
    return p.NextDecoder(LayerTypeIPv4)
}

/* for Decoding Interface */
func (g *GTPv1) CanDecode() gopacket.LayerClass {
	return LayerTypeGTPv1
}

func (g *GTPv1) NextLayerType() gopacket.LayerType {
	/* Should be LayerIPv4 ? */
	//return gopacket.LayerTypePayload
	return LayerTypeIPv4
}

func (g *GTPv1) LayerPayload() []byte {
    return g.restOfData
}

func (g *GTPv1) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 8 {
		df.SetTruncated()
		return fmt.Errorf("GTP packet too short")
	}
	g.Version = (data[0] >>5) & 0x07
	g.ProtocolType = (data[0] >> 4) & 0x01
	g.Reserved = 0
	g.ExtentionHeaderFlag = (data[0] >> 2) & 0x01
	g.SequenceNumberFlag = (data[0] >> 1) & 0x01
	g.NPDUNumberFlag = data[0] & 0x01
	g.MessageType = uint8(data[1])
	g.MessageLength = binary.BigEndian.Uint16(data[2:4])
	// fmt.Printf("MessageLength=%d\n", g.MessageLength)
	g.BaseLayer = BaseLayer{Contents: data[:len(data) - int(g.MessageLength)]}
	g.TEID = binary.BigEndian.Uint32(data[4:8])

	if g.ExtentionHeaderFlag >0 || g.SequenceNumberFlag >0 || g.NPDUNumberFlag > 0 {
		g.SequenceNumber = binary.BigEndian.Uint16(data[8:10])
		g.NPDUNumber = uint8(data[10])
		g.NextExtentionHeaderType = uint8(data[11])
		g.restOfData = data[12:]
	} else {
		g.restOfData = data[8:]
	}
	// fmt.Printf("DATA=%v\n", g.restOfData)
	return nil
}
