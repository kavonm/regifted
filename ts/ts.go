package ts

import (
	"regifted/data"

	"fmt"
	"log"
	"os"
)

const TS_PACKET_SIZE = 188

const (
	PACKET_TYPE_ERROR   = 0
	PACKET_TYPE_PAT     = 2
	PACKET_TYPE_PCR     = 3
	PACKET_TYPE_PES     = 4
	PACKET_TYPE_ES      = 1
	PACKET_TYPE_PMT     = 5
	PACKET_TYPE_PROGRAM = 6
	PACKET_TYPE_TS      = 7
)

type TSState struct {
	globals_initialized bool
	// the keys in these maps are pids
	pesCollector           map[uint]Pes
	pmtConstructors        map[uint]Pmt
	entryConstructors      map[uint]PmtEntry
	elementaryConstructors map[uint]ElementaryStreamPacket
	types                  map[uint]uint
	pat                    Pat

	bytes  []byte
	reader *data.Reader
	// reader *data.BufferedReader
	pcr uint

	// pes.streamtype -> pes[]
	pesMap map[uint][]Pes
}

// todo( mathew guest ) add error return
func Load(fh *os.File) *TSState {
	fmt.Println("load()")

	var state *TSState
	state = &TSState{}
	rc := state.init()
	if rc != true {
		log.Printf("could not initialize state\n")
		// return 71
		return nil
	}
	state.reader = data.NewReaderFromStream(fh)
	// state.reader = data.NewBufferedReaderFromStream(fh)
	state.main()
	return state
}

func (state *TSState) main() {
	reader := state.reader

	for reader.Cursor < reader.Size {
		var pesData *Pes
		byteChunk := reader.ReadBytes(TS_PACKET_SIZE)
		tsPacket := TsPacket{}
		tsPacket.byteChunk = byteChunk
		packetType, packetReader := state.ReadTSPacket(&tsPacket)

		switch {
		case packetType == PACKET_TYPE_PAT:
			state.readPat(&tsPacket, packetReader)

		case packetType == PACKET_TYPE_PMT:
			state.readPMT(&tsPacket, packetReader)

		case packetType == PACKET_TYPE_ES:
			pesData = state.readES(&tsPacket, packetReader)

			if pesData != nil {
				if state.pesMap[pesData.streamType] != nil {
					state.pesMap[pesData.streamType] = make([]Pes, 1, 1)

				}

				state.pesMap[pesData.streamType] = append(state.pesMap[pesData.streamType], *pesData)

			}
		}

		if tsPacket.hasAdaptation && tsPacket.adaptation.hasPCR {
			state.pcr = tsPacket.adaptation.pcr.pcr
		}
	}

	for key := range state.pesCollector {
		state.CreateAndDispensePes(key, state.types[key])
	}

}

//CreateAndDispensePes
//Dump the remaining PES
func (state *TSState) CreateAndDispensePes(pid uint, streamType uint) {
	pes := state.pesCollector[pid]
	pes.pid = pid
	pes.streamType = streamType
	pes.Read()
	pes.Print()
}

func (state *TSState) readPat(tsPacket *TsPacket, reader *data.Reader) {
	state.pat.byteChunk = reader.ReadBytes(reader.Size - reader.Cursor)
	state.pat.unitStart = tsPacket.unitStart
	state.pat.Read()

	state.loadPAT(&state.pat)
}

func (state *TSState) readPMT(tsPacket *TsPacket, reader *data.Reader) {
	pmt, _ := state.pmtConstructors[tsPacket.pid]
	pmt.unitStart = tsPacket.unitStart
	pmt.byteChunk = reader.ReadBytes(reader.Size - reader.Cursor)
	pmt.Read()

	state.loadPMT(&pmt)
}

func (state *TSState) readES(tsPacket *TsPacket, reader *data.Reader) *Pes {
	var pesData *Pes
	elementaryStreamPacket, _ := state.elementaryConstructors[tsPacket.pid]
	elementaryStreamPacket.pid = tsPacket.pid
	elementaryStreamPacket.unitStart = tsPacket.unitStart

	if tsPacket.hasAdaptation {
		elementaryStreamPacket.payload = tsPacket.adaptation.payload
	} else {
		elementaryStreamPacket.payload = reader.ReadBytes(reader.Size - reader.Cursor)
	}

	pesData = state.dispatch(&elementaryStreamPacket)
	elementaryStreamPacket.Print()
	return pesData
}

//Init
//Initialize the constructors
func (state *TSState) init() bool {
	if state.globals_initialized == true {
		log.Printf("EE attempted to initialize globals twice\n")
		return false
	}
	state.pmtConstructors = make(map[uint]Pmt)
	state.entryConstructors = make(map[uint]PmtEntry)
	state.types = make(map[uint]uint)
	state.pesCollector = make(map[uint]Pes)
	state.elementaryConstructors = make(map[uint]ElementaryStreamPacket)
	state.pat = Pat{}
	state.pat.tableId = 0
	state.pesMap = make(map[uint][]Pes)
	state.globals_initialized = true
	return true
}

func (state *TSState) DeleteState() {
	if state.globals_initialized == false {
		return
	}
	state.globals_initialized = false
	state.init()
	state.globals_initialized = false
}

func ReadHeaderData(bytes []byte) uint {
	reader := data.NewReader(bytes)
	var a uint = (reader.Read(1) >> 1) & 0x07
	var b uint = reader.Read(1)
	var c uint = (reader.Read(1) >> 1) & 0x7f
	var d uint = reader.Read(1)
	var e uint = (reader.Read(1) >> 1) & 0x7f
	var timestamp uint = (a << 30) | (b << 22) | (c << 15) | (d << 7) | e
	return timestamp
}
