package uuid // import "sour.is/x/toolbox/uuid"

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"time"
)

// NilUUID is a null uuid.
const NilUUID = "00000000-0000-0000-0000-000000000000"

// V4 return UUIDv4
func V4() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return NilUUID
	}

	// this make sure that the 13th character is "4"
	b[6] = (b[6] | 0x40) & 0x4F

	// this make sure that the 17th is "8", "9", "a", or "b"
	b[8] = (b[8] | 0x80) & 0xBF

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
func fromHexChar(c rune) (byte, bool) {
	switch {
	case 48 <= c && c <= 57:
		return byte(c) - 48, true
	case 97 <= c && c <= 122:
		return byte(c) - 97 + 10, true
	case 65 <= c && c <= 90:
		return byte(c) - 65 + 10, true
	}

	return 0, false
}
// V5 return UUIDv5
func V5(name, namespace string) string {
	n := make([]byte, 0, 16)
	half := false
	for _, c := range namespace {
		x, t := fromHexChar(c)
		if t == false {
			continue
		}

		if half == false {
			half = true
			n = append(n, x<<4)
		} else {
			half = false
			n[len(n)-1] = n[len(n)-1] | x
		}
	}
	n = append(n, []byte(name)...)
	b := sha1.Sum(n[:])

	// this make sure that the 13th character is "5"
	b[6] = (b[6] | 0x50) & 0x5F

	// this make sure that the 17th is "8", "9", "a", or "b"
	b[8] = (b[8] | 0x80) & 0xBF

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
// V6 returns UUIDv6
func V6(namespace string, origin bool) string {
	pfx := sha1.Sum([]byte(namespace))

	ns := make([]byte, 0, 16)
	half := false
	for _, c := range namespace {
		x, t := fromHexChar(c)
		if t == false {
			continue
		}

		if half == false {
			half = true
			ns = append(ns, x<<4)
		} else {
			half = false
			ns[len(ns)-1] = ns[len(ns)-1] | x
		}
	}
	node := sha1.Sum(ns[:])

	var ts uint64
	if !origin {
		ts = uint64(time.Now().UnixNano())>>2 | 0x8000000000000000
	}

	tcode := fmt.Sprintf("%016x", ts)

	b := make([]byte, 16)
	copy(b[0:4], pfx[:])
	copy(b[4:8], node[:])

	// this make sure that the 13th character is "6"
	b[6] = (b[6] | 0x60) & 0x6F

	// this make sure that the 17th is "8", "9", "a", or "b"
	//  b[8] = (b[8] | 0x80) & 0xBF

	return fmt.Sprintf("%x-%x-%x-%s-%s", b[0:4], b[4:6], b[6:8], tcode[0:4], tcode[4:16])
}
// Bytes returns the byte values for a uuid
func Bytes(uuid string) (n []byte) {
	n = make([]byte, 0, 16)

	half := false
	for _, c := range uuid {
		x, t := fromHexChar(c)
		if t == false {
			continue
		}

		if half == false {
			half = true
			n = append(n, x<<4)
		} else {
			half = false
			n[len(n)-1] = n[len(n)-1] | x
		}
	}

	return
}
