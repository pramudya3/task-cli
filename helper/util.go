package helper

import (
	"crypto/sha256"
	"encoding/binary"
	"time"

	"github.com/sqids/sqids-go"
)

func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func GenerateId(description string) string {
	h := sha256.New()
	h.Write([]byte(description))
	h.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	sum := h.Sum(nil)

	num := binary.BigEndian.Uint64(sum[:8])

	s, _ := sqids.New(sqids.Options{
		MinLength: 5,
	})
	id, _ := s.Encode([]uint64{num})

	if len(id) > 5 {
		id = id[:5]
	}

	return id
}
