package codec

import (
	"encoding/binary"
	"sync"
)

//FtPacket protocol
type FtPacket struct {
	gID        uint64 //整个fec组的id
	index      uint16 //该包在组中的位置
	realLength uint16 //去掉补齐部分的真实长度
	data       []byte
	len        int //在pool里的时候记录长度用
}

//Encode encode
func (ftp *FtPacket) Encode(result []byte) int {
	binary.BigEndian.PutUint64(result, ftp.gID)
	binary.BigEndian.PutUint16(result[8:], ftp.index)
	binary.BigEndian.PutUint16(result[10:], ftp.realLength)
	copy(result[12:], ftp.data)
	return 12 + len(ftp.data)
}

//Decode decode
func (ftp *FtPacket) Decode(data []byte) {
	ftp.gID = binary.BigEndian.Uint64(data[0:])
	ftp.index = binary.BigEndian.Uint16(data[8:])
	ftp.realLength = binary.BigEndian.Uint16(data[10:])
	ftp.data = data[12:]
}

//////

var mFtPool = newFtPool()

type ftPool struct {
	dataPool sync.Pool
}

func newFtPool() *ftPool {
	ftpool := new(ftPool)
	ftpool.dataPool = sync.Pool{
		New: func() interface{} {
			data := new(FtPacket)
			data.data = make([]byte, 2000)
			return data
		},
	}
	return ftpool
}

func (p *ftPool) poolGet() *FtPacket {
	item := p.dataPool.Get()
	return item.(*FtPacket)
}

func (p *ftPool) poolPut(item *FtPacket) {
	p.dataPool.Put(item)
}
