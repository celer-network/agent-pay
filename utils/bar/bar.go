package bar

import (
	"encoding/base64"
	"encoding/json"

	"google.golang.org/protobuf/types/known/anypb"
)

// MarshalOpaqueAnyJSON returns a JSON object that preserves the Any type URL and
// base64-encodes the raw payload without requiring the embedded message type to
// be linked into the binary. Intended for logging/debugging.
func MarshalOpaqueAnyJSON(a *anypb.Any) ([]byte, error) {
	if a == nil {
		return []byte("null"), nil
	}
	return json.Marshal(map[string]any{
		"@type": a.TypeUrl,
		"value": base64.StdEncoding.EncodeToString(a.Value),
	})
}
