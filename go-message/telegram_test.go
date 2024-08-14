package gomessage

import "testing"

func TestTelegram(t *testing.T) {
	InitTelegram("7107568224:AAFgdiEsDqtFvBBScIfWku9IB8jr9Dpl-dw", true)
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "test", args: args{
			text: "test",
		}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Telegram(5562314141, tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("Telegram() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
