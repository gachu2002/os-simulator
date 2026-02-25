package sim

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func TraceHash(trace []TraceEvent) string {
	b := strings.Builder{}
	for _, ev := range trace {
		fmt.Fprintf(&b, "%d|%d|%s|%s\n", ev.Tick, ev.Sequence, ev.Kind, ev.Data)
	}
	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}
