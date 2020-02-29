package main

import (
	"reflect"
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
		wantStdout []string
		wantStderr []string
		wantErr    bool
	}{
		{name: "empty enabled", fields: fields{enabled: true}, args: args{params: nil}, wantStdout: []string{""}, wantStderr: []string{}, wantErr: false},
		{name: "empty disabled", fields: fields{enabled: false}, args: args{params: nil}, wantStdout: []string{}, wantStderr: []string{}, wantErr: true},
		{name: "Hello enabled", fields: fields{enabled: true}, args: args{params: []string{"Hello"}}, wantStdout: []string{"Hello"}, wantStderr: []string{}, wantErr: false},
		{name: "Hello disabled", fields: fields{enabled: false}, args: args{params: []string{"Hello"}}, wantStdout: []string{}, wantStderr: []string{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugin{
				enabled: tt.fields.enabled,
			}
			c, err := p.Run(tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if c != nil {
				<-c.Start()
				if gotStdout := c.Status().Stdout; !reflect.DeepEqual(tt.wantStdout, gotStdout) {
					t.Errorf("Run() gotStdout = %v, want %v", gotStdout, tt.wantStdout)
				}
				if gotStderr := c.Status().Stderr; !reflect.DeepEqual(tt.wantStderr, gotStderr) {
					t.Errorf("Run() gotStderr = %v, want %v", gotStderr, tt.wantStderr)
				}
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
