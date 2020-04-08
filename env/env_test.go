package env

import (
	"reflect"
	"testing"
)

func TestGetEnv(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name             string
		args             args
		wantEnvironments map[string]string
		wantErr          bool
	}{
		{
			name: "test env",
			args: args{filename: PWDFile(".env")},
			wantEnvironments: map[string]string{
				"NAME":     "OhYee",
				"secret":   "123456",
				"password": "A1b2c=!",
			},
			wantErr: false,
		},
		{
			name:             "error",
			args:             args{filename: PWDFile(".error")},
			wantEnvironments: map[string]string{},
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEnvironment, err := GetEnv(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotEnvironment, tt.wantEnvironments) {
				t.Errorf("GetEnv() = %v, want %v", gotEnvironment, tt.wantEnvironments)
			}
		})
	}
}
