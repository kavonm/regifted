package mpeg4file

import (
		"strconv"
		"regifted/mpeg4file/moofSubBox"
		)

type moof struct{
	size uint32
	largeSize uint64
	boxType uint32
	containedBoxes moofSubBox.MoofLevel2 
}

func NewMoof(s uint64, box uint32) *moof{
	newMoof:=new(moof)
	newMoof.SetSize(s)
	newMoof.boxType = box
	return newMoof
}

func (m *moof) SetSize(s uint64){
	if s==0 {
		m.size = 0
	}else if s>4294967295 {
		m.size = 1
		m.largeSize = s
	} else {
		m.size = uint32(s)
	}
}

func (m *moof) String() string{
	return strconv.FormatUint(uint64(m.size),10)
}

func (m *moof) Write() []byte{
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