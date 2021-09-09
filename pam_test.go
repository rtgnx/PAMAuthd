package main

import "testing"

func TestPAMAuth(t *testing.T) {
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test Login",
			args: args{username: "testuser", password: ""},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PAMAuth(tt.args.username, tt.args.password); got != tt.want {
				t.Errorf("PAMAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}
