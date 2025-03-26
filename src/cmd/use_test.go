package cmd

import (
	"testing"
)

func TestRunUse(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "JDK exists",
			args:    []string{"jdk8"},
			wantErr: false,
		},
		{
			name:    "JDK does not exist",
			args:    []string{"jdk11"},
			wantErr: true,
			errMsg:  "JDK 'jdk11' does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runUse(nil, tt.args)
		})
	}
}
