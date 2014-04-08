package moof

import (
		"strconv"
		)

type traf struct{
	size uint32
	largeSize uint64
	boxType uint32
}

func NewTraf(s uint64) *traf{
	newTraf:=new(traf)
	newTraf.SetSize(s)
	newTraf.boxType = 0x74726166 // Hex rep of traf
	return newTraf
}

func (t *traf) SetSize(s uint64){
	if s==0{
		t.size = 0
	}else if s > 4294967295 {
		t.size = 1
		t.largeSize = s
	} else {
		t.size = uint32(s)
	}
}

func (t *traf) String() string{
	return strconv.FormatUint(uint64(t.size),10)
}

func (m *traf) Write() []byte{
	buf := new(bytes.Buffer)
	var err error
	// Size
	err=binary.Write(buf, binary.BigEndian, m.size)
	if err!=nil{
		fmt.Println("binary.Write failed:", err)
	}
	// BoxType
	err = binary.Write(buf,binary.BigEndian,m.boxType)
	if err!=nil{
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}