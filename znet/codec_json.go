package znet

import (
	"encoding/json"
	"myzinx/ziface"
)

var _ ziface.ICodec = &JsonCodec{}

// JsonCodec implements the Codec interface.
// JsonCodec encodes and decodes data in json way.
type JsonCodec struct{}

// Encode implements the Codec Encode method.
func (c *JsonCodec) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Decode implements the Codec Decode method.
func (c *JsonCodec) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
