package ints

import (
	"sync"
	"time"
)

var DefaultSnowflakeGenerator = &SnowflakeGenerator{
	TimeShift: DefaultTimeShift,
}

const (
	// DefaultTimeShift 2010-11-04, 01:42:54 UTC
	// ~69 years in 4 bit
	DefaultTimeShift = 1288834974657
)

// SnowflakeGenerator - generator IDs
// 1 - 0 | 41 - time in milliseconds | 5 - data center | 5 server | 12 sequence
// ~69 years in 4 bit
type SnowflakeGenerator struct {
	Mx sync.Mutex
	// TimeShift - time shift / Use 1288834974657 for default generator
	TimeShift int64
	// DataCenterID - data center id 0-31 (valid values)
	DataCenterID int64
	// ServerID - server id 0-31 (valid values)
	ServerID int64
	// sequence - values between 0 - 4095
	sequence int64
	// lastTime - last call time
	lastTime int64
}

// Next - gets next id: 1 - 0 | 41 - time in milliseconds | 5 - data center | 5 server | 12 sequence
// ~69 years in 4 bit
func (fg *SnowflakeGenerator) Next() int64 {
	id, owerload := fg.NextLock()

	for owerload {
		id, owerload = fg.NextLock()
	}

	return id
}

// Next - gets next id: 1 - 0 | 41 - time in milliseconds | 5 - data center | 5 server | 12 sequence
// ~69 years in 4 bit
func (fg *SnowflakeGenerator) NextLock() (id int64, owerload bool) {
	fg.Mx.Lock()
	defer fg.Mx.Unlock()

	return fg.NextRaw()
}

// NextRaw - gets next id: 1 - 0 | 41 - time in milliseconds | 5 - data center | 5 server | 12 sequence
// ~69 years in 4 bit
func (fg *SnowflakeGenerator) NextRaw() (id int64, owerload bool) {
	m := time.Now().UnixMilli() - fg.TimeShift

	if fg.lastTime != m {
		fg.lastTime = m
		fg.sequence = 0
	}

	if fg.sequence >= 4095 {
		owerload = true
	}

	fg.sequence++

	m = m<<22 + fg.DataCenterID<<17 + fg.ServerID<<12 + fg.sequence

	return m, owerload
}

// LimitID - generates min id for time
func (fg *SnowflakeGenerator) LimitID(t time.Time) (id int64) {
	m := t.UnixMilli() - fg.TimeShift
	return m << 22
}

// GetTimeFromID - get time from id
func (fg *SnowflakeGenerator) GetTimeFromID(id int64) time.Time {
	t := id >> 22
	return time.UnixMilli(t + fg.TimeShift)
}

// LimitID - generates min id for time
func LimitID(t time.Time) (id int64) {
	return DefaultSnowflakeGenerator.LimitID(t)
}

// GetTimeFromID - get time from id
func GetTimeFromID(id int64) time.Time {
	return DefaultSnowflakeGenerator.GetTimeFromID(id)
}

// NextID - gets next id: 1 - 0 | 41 - time in milliseconds | 5 - data center | 5 server | 12 sequence
// ~69 years in 4 bit
func NextID() int64 {
	return DefaultSnowflakeGenerator.Next()
}
