package ints

import (
	"crypto/rand"
	"sync"
	"time"
)

var DefaultRandomUuidGenerator = &UuidRealRandomGenerator{}

const UuidSize = 16
const DefaultUuidRealRandomGeneratorSize = 1000
const DefaultUuidRealRandomGeneratorBlocks = 4

// UuidRealRandomGenerator - random uuid generator with random values
type UuidRealRandomGenerator struct {
	nextRandoms       [DefaultUuidRealRandomGeneratorBlocks][]byte
	nextRandomsReload [DefaultUuidRealRandomGeneratorBlocks]bool
	Size              int
	currentBlock      int
	currentPos        int
	mx                sync.Mutex
}

func (urg *UuidRealRandomGenerator) Init() error {
	if urg.Size <= 0 {
		urg.Size = DefaultUuidRealRandomGeneratorSize
	}

	for i := range urg.nextRandoms {
		urg.nextRandoms[i] = make([]byte, 0)
		urg.nextRandomsReload[i] = true
	}

	err := urg.reloadRandoms(0)

	for i := 1; i < DefaultUuidRealRandomGeneratorBlocks; i++ {
		go urg.reloadRandoms(i)
	}

	return err
}

func (urg *UuidRealRandomGenerator) reloadRandoms(i int) error {
	urg.mx.Lock()
	urg.nextRandomsReload[i] = true
	urg.mx.Unlock()

	urg.nextRandoms[i] = make([]byte, urg.Size*UuidSize)

	_, err := rand.Read(urg.nextRandoms[i])

	urg.mx.Lock()
	urg.nextRandomsReload[i] = (err != nil)
	urg.mx.Unlock()

	return err
}

func (urg *UuidRealRandomGenerator) reloadRandomsMust(i int) {
	urg.mx.Lock()
	nrTrue := urg.nextRandomsReload[i]
	urg.mx.Unlock()

	for nrTrue {
		urg.reloadRandoms(i)

		urg.mx.Lock()
		nrTrue = urg.nextRandomsReload[i]
		urg.mx.Unlock()
	}
}

// GetRandBytes returns 16 bytes
func (urg *UuidRealRandomGenerator) GetRandBytes() []byte {
	urg.mx.Lock()

	nrTrue := urg.nextRandomsReload[urg.currentBlock]
	if nrTrue {
		urg.mx.Unlock()
		time.Sleep(1 * time.Microsecond)
		return urg.GetRandBytes()
	}

	defer urg.mx.Unlock()

	res := urg.nextRandoms[urg.currentBlock][urg.currentPos*UuidSize : (urg.currentPos+1)*UuidSize]

	urg.currentPos++
	if len(urg.nextRandoms[urg.currentBlock]) <= urg.currentPos*UuidSize {
		urg.nextRandomsReload[urg.currentBlock] = true
		go urg.reloadRandomsMust(urg.currentBlock)
		urg.currentBlock++
		urg.currentPos = 0
	}

	if urg.currentBlock >= DefaultUuidRealRandomGeneratorBlocks {
		urg.currentBlock = 0
	}

	return res
}

// Next Random id
// xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx
// M 0100 (version 4)
// N 10 (variant 1)
// rrrrrrrr-rrrr-Mrrr-Nrrr-rrrrrrrrrrrr
// M - 0100 (version 4)
// N - 10rr (variant 1)
// r - random value
func (urg *UuidRealRandomGenerator) Next() (res Uuid) {
	bts := urg.GetRandBytes()
	res.UInt128.SetBytes(bts)

	res.UInt128[1] = Nval + res.UInt128[1]>>2

	res.UInt128[0] = (res.UInt128[0]>>16)<<16 + (res.UInt128[0]<<52)>>52 + (1 << 14)

	return res
}

// NextOrDefault gets Next from usg where usg != nil else uses DefaultRandomUuidGenerator.Next()
func (usg *UuidRealRandomGenerator) NextOrDefault() (res Uuid) {
	if usg != nil {
		return usg.Next()
	}

	return DefaultRandomUuidGenerator.Next()
}

// RandUUID - Next Random id from DefaultRandomUuidGenerator
// xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx
// M 0100 (version 4)
// N 10 (variant 1)
// rrrrrrrr-rrrr-Mrrr-Nrrr-rrrrrrrrrrrr
// M - 0100 (version 4)
// N - 10rr (variant 1)
// r - random value
func RandUUID() (res Uuid) {
	return DefaultRandomUuidGenerator.Next()
}
