package ratlogparser

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestSimpleParser_parseLine(t *testing.T) {
	p := SimpleParser{}

	tests := []struct {
		name    string
		input   string
		want    Entry
		wantErr bool
	}{
		{
			name:  "ok",
			input: `[request|info] GET | ip: \:\:ffff\:172.18.0.3 | url: /api/v1/about | method: GET | xhr: true`,
			want: Entry{
				Tags:    []Tag{"request", "info"},
				Message: "GET",
				Fields: map[string]fmt.Stringer{
					"ip":     BasicField("::ffff:172.18.0.3"),
					"url":    BasicField("/api/v1/about"),
					"method": BasicField("GET"),
					"xhr":    BasicField("true"),
				},
			},
			wantErr: false,
		},
		{
			name:  "ok",
			input: `[graphql] operation-responsetime | operation: onDigitalConsultationEvent | duration: 27.303113068 | operationType: subscription`,
			want: Entry{
				Tags:    []Tag{"graphql"},
				Message: "operation-responsetime",
				Fields: map[string]fmt.Stringer{
					"operation":     BasicField("onDigitalConsultationEvent"),
					"duration":      BasicField("27.303113068"),
					"operationType": BasicField("subscription"),
				},
			},
			wantErr: false,
		},
		{
			name:  "message only",
			input: `message only`,
			want: Entry{
				Tags:    nil,
				Message: "message only",
				Fields:  map[string]fmt.Stringer{},
			},
			wantErr: false,
		},
		{
			name:  "message and fields",
			input: `File not found | path: /tmp/notfound.txt`,
			want: Entry{
				Tags:    nil,
				Message: "File not found",
				Fields: map[string]fmt.Stringer{
					"path": BasicField("/tmp/notfound.txt"),
				},
			},
			wantErr: false,
		},
		{
			name:  "message with escaped slash",
			input: `path not found \\`,
			want: Entry{
				Tags:    nil,
				Message: `path not found \`,
				Fields:  map[string]fmt.Stringer{},
			},
			wantErr: false,
		},
		{
			name:  "line break in message",
			input: `error\\nin file`,
			want: Entry{
				Tags:    nil,
				Message: `error\nin file`,
				Fields:  map[string]fmt.Stringer{},
			},
			wantErr: false,
		},
		{
			name:  "escaped bracket in tag",
			input: `[nginx\[] GET | url: /v1/`,
			want: Entry{
				Tags:    []Tag{"nginx["},
				Message: `GET`,
				Fields: map[string]fmt.Stringer{
					"url": BasicField("/v1/"),
				},
			},
			wantErr: false,
		},
		{
			name:  "escaped pipe in tag",
			input: `[nginx\|] GET | url: /v1/`,
			want: Entry{
				Tags:    []Tag{"nginx|"},
				Message: `GET`,
				Fields: map[string]fmt.Stringer{
					"url": BasicField("/v1/"),
				},
			},
			wantErr: false,
		},
		{
			name:  "escaped pipe in tags",
			input: `[nginx\||api] GET | url: /v1/`,
			want: Entry{
				Tags:    []Tag{"nginx|", "api"},
				Message: `GET`,
				Fields: map[string]fmt.Stringer{
					"url": BasicField("/v1/"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.parseLine(bytes.NewBufferString(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLine() got = %v, want %v", got, tt.want)
			}
		})
	}
}
