package moofSubBox

import (
		"strconv"
		)

type MoofLevel2 interface{
	String() string
}

type mfhd struct{
	size uint32
	largeSize uint64
	boxType uint32
	version uint8
	flags [3]byte
}

func NewMfhd(s uint64, box uint32, ver uint8, flag [3]byte){
	newMfhd:=new(mfhd)
	newMfhd.SetSize(s)
	newMfhd.boxType=box
	newMfhd.version=ver
	newMfhd.flags=flag
}

func (m* mfhd) SetSize(s uint64){
	if s == 0 {
		m.size = 0
	} else if s > 4294967295{
		m.size = 1
		m.largeSize = s
	} else {
		m.size = uint32(s)
	}
}

func (m *mfhd) String() string{
	return strconv.FormatUint(uint64(m.size),10)
}

func (m *mfhd) Write() []byte{
	var data []byte
	// Size
	if m.size!=1{
		data = strconv.AppendUint(data, uint64(m.size), 2)	
	} else {
		data = strconv.AppendUint(data, m.largeSize, 2)
	}	
	// BoxType
	// Contained boxes write
	return data
}