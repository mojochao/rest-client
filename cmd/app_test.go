package cmd

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

func Test_parseEnvs(t *testing.T) {
	local := "http://localhost:8080"
	prod := "https://prod.example.com"
	date := "20201031"

	text := fmt.Sprintf(`
{
	"local": {
		"svc": "%s",
    "date": "%s"
  },
	"prod": {
		"svc": "%s",
    "date": "%s"
  }
}`, local, date, prod, date)
	reader := strings.NewReader(text)

	want := Envs{
		"local": Vars{
			"svc":  local,
			"date": date,
		},
		"prod": Vars{
			"svc":  prod,
			"date": date,
		},
	}

	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    Envs
		wantErr bool
	}{
		{
			name:    "test parseEnvs()",
			args:    args{r: reader},
			want:    want,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseEnvs(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseEnvs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseEnvs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

var validReqsText = `
#:name getWithNoParams
#:desc Example get with no query parameters
GET {{svc}}{{api}}
Content-Type: application/json

###

#:name getWithNoParams
#:desc Example get with no query parameters
HEAD {{svc}}{{api}}
Content-Type: application/json

###

#:name getWithNoParams
#:desc Example get with no query parameters
POST {{svc}}{{api}}
Content-Type: application/json

{
  "foo": "Foo",
  "bar": "Bar"
}`

func Test_parseReqs(t *testing.T) {
	reader := strings.NewReader(validReqsText)

	type args struct {
		r    io.Reader
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []*Req
		wantErr bool
	}{
		{
			name: "test parseReqs",
			args: args{r: reader, path: "/some/path"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseReqs(tt.args.r, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseReqs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("parseReqs() got = %v, want %v", got, tt.want)
			//}
			if len(got) != 3 {
				t.Errorf("got %d Reqs, want %d", len(got), 1)
			}
		})
	}
}
