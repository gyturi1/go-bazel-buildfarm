package pkg

import (
	"errors"
	"testing"
)

func TestFormat(t *testing.T) {
	type args struct {
		e error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{errors.New("")}, "couldn't accept: "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.args.e); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}
