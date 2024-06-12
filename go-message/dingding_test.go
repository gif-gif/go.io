package gomessage

import "testing"

func TestDingDing(t *testing.T) {
	InitDingDing("bb96f055f83a0ad78b3112ca849f29d37de5f1bfecc5d1e6f205e4f63e6b0e93", "SEC6fd7250c7e489eeda966719217d8bc45136819cfd73bee38e776769c563cacfe")
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
			text: "test",
		}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DingDing(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("DingDing() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
