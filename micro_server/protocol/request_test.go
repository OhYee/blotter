package proto

import (
	"bytes"
	"github.com/OhYee/goutils"
	"testing"
)

func TestRequest(t *testing.T) {
	req := NewRequest("/add", map[string]any{
		"a": 1,
		"b": 2,
	})
	b := req.ToBytes()
	got, err := NewRequestFromBytes(bytes.NewBuffer(b))
	if err != nil {
		t.Errorf("Got error %v", err)
	} else if !goutils.Equal(got, req) {
		t.Errorf("Want %v but got %v", req, got)
	}
}
