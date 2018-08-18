package timeid

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"
)

var minimumDate = time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()

const uint32MaxValue = 4294967295

// NewTimeID creates a new ID based on the current time and a random number.
func NewTimeID(rand *rand.Rand, now func() time.Time) (tid TimeID, err error) {
	n := now()
	unixTime := n.Unix() - minimumDate
	if unixTime > uint32MaxValue {
		err = fmt.Errorf("timeid: cannot create an ID for time %v, time value %d exceeds maximum value of %d", n, unixTime, uint32MaxValue)
		return
	}
	tid = TimeID{
		Time:   uint32(unixTime),
		Random: rand.Uint32(),
	}
	return
}

// TimeID is an ID based on time and a random number.
type TimeID struct {
	Random uint32
	Time   uint32
}

// Bytes which make up the underlying time ID.
func (tid TimeID) Bytes() (op []byte, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, tid)
	op = buf.Bytes()
	return
}

// FromInt64 recovers the underlying data from the time ID.
func FromInt64(id int64) (tid TimeID, err error) {
	w := bytes.NewBuffer([]byte{})
	err = binary.Write(w, binary.LittleEndian, id)
	if err != nil {
		return
	}
	r := bytes.NewBuffer(w.Bytes())
	err = binary.Read(r, binary.LittleEndian, &tid)
	return
}

// Int64 creates a single integer which represents the time and random value.
func (tid TimeID) Int64() (op int64, err error) {
	b, err := tid.Bytes()
	if err != nil {
		return
	}
	r := bytes.NewReader(b)
	err = binary.Read(r, binary.LittleEndian, &op)
	return
}

// GetTime returns the time that the ID was created.
func (tid TimeID) GetTime() time.Time {
	return time.Unix(int64(tid.Time)+minimumDate, 0)
}
