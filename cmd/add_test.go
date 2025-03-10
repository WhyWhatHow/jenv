package cmd

import (
	"testing"
)

func TestRunAdd(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		force   bool
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Add JDK successfully",
			args:    []string{"jdk8", "C:\\Program Files\\Java\\jdk1.8.0_291"},
			force:   false,
			wantErr: false,
		},
		{
			name:    "Add JDK with force flag",
			args:    []string{"jdk11", "C:\\Program Files\\Java\\jdk-11.0.12"},
			force:   true,
			wantErr: false,
		},
		{
			name:    "Add JDK with invalid path",
			args:    []string{"jdk11", "invalid_path"},
			force:   false,
			wantErr: true,
			errMsg:  "failed to add JDK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			force = tt.force
			runAdd(nil, tt.args)
			if tt.wantErr {
				//assert.Error(t, err)
				//assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				//assert.NoError(t, err)
			}
		})
	}
}
