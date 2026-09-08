package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"regexp"
	"time"
)

var uuidV7Regex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

// NewV7 generates a valid RFC 9562 UUID version 7.
// It encodes a 48-bit millisecond timestamp in the most significant bits,
// followed by a 4-bit version (7), 12 bits of entropy, a 2-bit variant (10),
// and 62 bits of entropy.
func NewV7() string {
	var b [16]byte

	// 48-bit timestamp (milliseconds since Unix epoch)
	milli := uint64(time.Now().UnixMilli())
	b[0] = byte(milli >> 40)
	b[1] = byte(milli >> 32)
	b[2] = byte(milli >> 24)
	b[3] = byte(milli >> 16)
	b[4] = byte(milli >> 8)
	b[5] = byte(milli)

	// Fill remaining 10 bytes with cryptographic randomness
	if _, err := rand.Read(b[6:]); err != nil {
		binary.BigEndian.PutUint64(b[8:], uint64(time.Now().UnixNano()))
	}

	// Version 7: high nibble of byte 6 is 0x7
	b[6] = (b[6] & 0x0f) | 0x70
	// Variant 10xx: high 2 bits of byte 8 are 0x8
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// IsValidV7 checks whether the string is a valid UUIDv7 format.
func IsValidV7(s string) bool {
	return uuidV7Regex.MatchString(s)
}
