package xmail

import (
	"testing"
)

func TestSendMail(t *testing.T) {
	type args struct {
		server   string
		from     string
		password string
		to       string
		subject  string
		body     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				server:   "xx.xx.net:465",
				from:     "xx@xx.xx",
				password: "xx",
				to:       "xx@xx.xx",
				subject:  "This is the email subject",
				body:     "This is an example body.\\n With two lines.",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendMail(tt.args.server, tt.args.from, tt.args.password, tt.args.to, tt.args.subject, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("SendMail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
