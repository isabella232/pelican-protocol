package schema

import (
	"io"

	capn "github.com/glycerine/go-capnproto"
)

func (s *PelicanPacket) Save(w io.Writer) error {
	seg := capn.NewBuffer(nil)
	PelicanPacketGoToCapn(seg, s)
	_, err := seg.WriteTo(w)
	return err
}

func (s *PelicanPacket) Load(r io.Reader) error {
	capMsg, err := capn.ReadFromStream(r, nil)
	if err != nil {
		//panic(fmt.Errorf("capn.ReadFromStream error: %s", err))
		return err
	}
	z := ReadRootPelicanPacketCapn(capMsg)
	PelicanPacketCapnToGo(z, s)
	return nil
}

func PelicanPacketCapnToGo(src PelicanPacketCapn, dest *PelicanPacket) *PelicanPacket {
	if dest == nil {
		dest = &PelicanPacket{}
	}
	dest.ResponseSerial = int64(src.ResponseSerial())
	dest.RequestSerial = int64(src.RequestSerial())
	dest.Key = src.Key()

	var n int

	// Mac
	n = src.Mac().Len()
	dest.Mac = make([]byte, n)
	for i := 0; i < n; i++ {
		dest.Mac[i] = byte(src.Mac().At(i))
	}

	// Payload
	n = src.Payload().Len()
	dest.Payload = make([]byte, n)
	for i := 0; i < n; i++ {
		dest.Payload[i] = byte(src.Payload().At(i))
	}

	dest.RequestAbTm = int64(src.RequestAbTm())
	dest.RequestLpTm = int64(src.RequestLpTm())
	dest.ResponseLpTm = int64(src.ResponseLpTm())
	dest.ResponseAbTm = int64(src.ResponseAbTm())

	return dest
}

func PelicanPacketGoToCapn(seg *capn.Segment, src *PelicanPacket) PelicanPacketCapn {
	dest := AutoNewPelicanPacketCapn(seg)
	dest.SetResponseSerial(src.ResponseSerial)
	dest.SetRequestSerial(src.RequestSerial)
	dest.SetKey(src.Key)

	mylist1 := seg.NewUInt8List(len(src.Mac))
	for i := range src.Mac {
		mylist1.Set(i, uint8(src.Mac[i]))
	}
	dest.SetMac(mylist1)

	mylist2 := seg.NewUInt8List(len(src.Payload))
	for i := range src.Payload {
		mylist2.Set(i, uint8(src.Payload[i]))
	}
	dest.SetPayload(mylist2)
	dest.SetRequestAbTm(src.RequestAbTm)
	dest.SetRequestLpTm(src.RequestLpTm)
	dest.SetResponseLpTm(src.ResponseLpTm)
	dest.SetResponseAbTm(src.ResponseAbTm)

	return dest
}

func SliceByteToUInt8List(seg *capn.Segment, m []byte) capn.UInt8List {
	lst := seg.NewUInt8List(len(m))
	for i := range m {
		lst.Set(i, uint8(m[i]))
	}
	return lst
}

func UInt8ListToSliceByte(p capn.UInt8List) []byte {
	v := make([]byte, p.Len())
	for i := range v {
		v[i] = byte(p.At(i))
	}
	return v
}
