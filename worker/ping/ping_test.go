package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestPlugin_Disable(t *testing.T) {
	type fields struct {
		enabled bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "enabled", fields: fields{enabled: true}},
		{name: "disabled", fields: fields{enabled: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugin{
				enabled: tt.fields.enabled,
			}
			p.Disable()
			if p.enabled != false {
				t.Errorf("expected %v, got %v", false, tt.fields.enabled)
			}
		})
	}
}

func TestPlugin_Enable(t *testing.T) {
	type fields struct {
		enabled bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "enabled", fields: fields{enabled: true}},
		{name: "disabled", fields: fields{enabled: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugin{
				enabled: tt.fields.enabled,
			}
			p.Enable()
			if p.enabled != true {
				t.Errorf("expected %v, got %v", true, tt.fields.enabled)
			}
		})
	}
}

func TestPlugin_Run(t *testing.T) {
	type fields struct {
		enabled bool
	}
	type args struct {
		params []string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStdout string
		wantStderr string
		wantErr    bool
	}{
		{name: "empty enabled", fields: fields{enabled: true}, args: args{params: nil}, wantStdout: "", wantStderr: "Usage: ping", wantErr: true},
		{name: "empty disabled", fields: fields{enabled: false}, args: args{params: nil}, wantStdout: "", wantStderr: "", wantErr: true},
		{name: "Hello", fields: fields{enabled: true}, args: args{params: []string{"-V"}}, wantStdout: "ping utility", wantStderr: "", wantErr: false},
		{name: "transmit", fields: fields{enabled: true}, args: args{params: []string{"-c", "4", "8.8.8.8"}}, wantStdout: "packets transmitted", wantStderr: "", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugin{
				enabled: tt.fields.enabled,
			}
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			err := p.Run(stdout, stderr, tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStdout := stdout.String(); !strings.Contains(gotStdout, tt.wantStdout) {
				t.Errorf("Run() gotStdout = %v, want %v", gotStdout, tt.wantStdout)
			}
			if gotStderr := stderr.String(); !strings.Contains(gotStderr, tt.wantStderr) {
				t.Errorf("Run() gotStderr = %v, want %v", gotStderr, tt.wantStderr)
			}
		})
	}
}

func TestPlugin_Status(t *testing.T) {
	type fields struct {
		enabled bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "init disabled", fields: fields{enabled: false}, want: false},
		{name: "init enabled", fields: fields{enabled: true}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugin{
				enabled: tt.fields.enabled,
			}
			if got := p.Status(); got != tt.want {
				t.Errorf("Status() = %v, want %v", got, tt.want)
			}
		})
	}
}
