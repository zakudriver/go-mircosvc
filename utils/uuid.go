package utils

import (
	"crypto/rand"
	"encoding/hex"
)

const Size = 16

type UUID [Size]byte

func (u UUID) String() string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}

// 生成uuid
func NewUUID() (*UUID, error) {
	u := &UUID{}

	if _, err := rand.Read(u[:]); err != nil {
		return nil, err
	}

	u[6] = (u[6] & 0x0f) | (4 << 4)
	u[8] = (u[8]&(0xff>>2) | (0x02 << 6))

	return u, nil
}
