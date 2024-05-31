package gomessage

import "testing"

func TestFeiShu(t *testing.T) {
	type args struct {
		hookUrl string
		text    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "test", args: args{
			hookUrl: "https://open.feishu.cn/open-apis/bot/v2/hook/aa0f28f1-1663-421b-9fa9-af9b0bbe2ca4",
			text:    "test",
		}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := FeiShu(tt.args.hookUrl, tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("FeiShu() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
