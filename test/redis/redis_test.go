package main

import (
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func fakeGenerateToken(int) (string, error) {
	return "Ta5EVtRgUD-HFmRwrujAwKZnx247lFfe", nil
}

func TestLogin(t *testing.T) {
	// mock redis client
	rdb, mock := redismock.NewClientMock()

	// login success
	mock.ExpectGet("sms:captcha:13800138000").SetVal("123456")
	mock.ExpectSet("auth:token:Ta5EVtRgUD-HFmRwrujAwKZnx247lFfe", "13800138000", 24*time.Hour).SetVal("OK")

	// invalid sms code or expired
	mock.ExpectGet("sms:captcha:13900139000").RedisNil()

	// invalid sms code
	mock.ExpectGet("sms:captcha:13700137000").SetVal("123123")

	type args struct {
		mobile  string
		smsCode string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr string
	}{
		{
			name: "login success",
			args: args{
				mobile:  "13800138000",
				smsCode: "123456",
			},
			want: "Ta5EVtRgUD-HFmRwrujAwKZnx247lFfe",
		},
		{
			name: "invalid sms code or expired",
			args: args{
				mobile:  "13900139000",
				smsCode: "123459",
			},
			wantErr: "invalid sms code or expired",
		},
		{
			name: "invalid sms code",
			args: args{
				mobile:  "13700137000",
				smsCode: "123457",
			},
			wantErr: "invalid sms code",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Login(tt.args.mobile, tt.args.smsCode, rdb, fakeGenerateToken)
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
