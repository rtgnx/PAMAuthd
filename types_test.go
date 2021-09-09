package main

import "testing"

func TestPasswdEntry_Decode(t *testing.T) {
	type fields struct {
		Name     string
		Password string
		UID      int64
		GID      int64
		Fullname string
		Home     string
		Shell    string
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Test valid passwd entry",
			args:    args{b: []byte("rtgnx:x:1001:1001:,,,:/home/rtgnx:/bin/bash")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PasswdLine{
				Name:     tt.fields.Name,
				Password: tt.fields.Password,
				UID:      tt.fields.UID,
				GID:      tt.fields.GID,
				Fullname: tt.fields.Fullname,
				Home:     tt.fields.Home,
				Shell:    tt.fields.Shell,
			}
			if err := p.UnmarshalText(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("PasswdLine.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
