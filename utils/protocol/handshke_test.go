package proto

import (
	"bytes"
	"testing"
	// gb "github.com/OhYee/goutils/bytes"
)

func TestHandshake(t *testing.T) {
	tests := []struct {
		name    string
		b       []byte
		wantErr bool
		valid   bool
	}{
		{"test handshake", nil, false, true},
		{"test handshake", []byte{}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			if err := SendHandshake(buf); (err != nil) != tt.wantErr {
				t.Errorf("SendHandshake() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.b != nil {
				buf.Reset()
				buf.Write(tt.b)
			}

			if VarifyHandshake(buf) != tt.valid {
				t.Errorf("handshake error want %v %v", tt.valid, buf.Bytes())
			}
		})
	}
}
