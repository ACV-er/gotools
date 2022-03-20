package compress

import (
	"encoding/base64"
	"reflect"
	"testing"
)

func TestCompress(t *testing.T) {
	type args struct {
		mode int
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantRet string
		wantErr bool
	}{
		{
			name: "正常压缩数据",
			args: args{
				COMPRESS_GZIP,
				[]byte("1351451451512"),
			},
			wantRet: "H4sIAAAAAAAA/zI0NjU0ASNTQyNAAAAA//9lZlNKDQAAAA==",
			wantErr: false,
		},
		{
			name: "压缩空字符串",
			args: args{
				COMPRESS_GZIP,
				[]byte(""),
			},
			wantRet: "H4sIAAAAAAAA/wEAAP//AAAAAAAAAAA=",
			wantErr: false,
		},
		{
			name: "压缩空字节数组",
			args: args{
				COMPRESS_GZIP,
				[]byte{},
			},
			wantRet: "H4sIAAAAAAAA/wEAAP//AAAAAAAAAAA=",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, err := Compress(tt.args.mode, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRet != tt.wantRet {
				t.Errorf("Compress() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestUnCompress(t *testing.T) {
	common_base64 := base64.StdEncoding.EncodeToString([]byte(`{"msg":"这是一个json,不能被gzip解压"}`))
	type args struct {
		mode int
		data string
	}
	tests := []struct {
		name    string
		args    args
		wantRet []byte
		wantErr bool
	}{
		{
			name: "正常压缩数据",
			args: args{
				COMPRESS_GZIP,
				"H4sIAAAAAAAA/zI0NjU0ASNTQyNAAAAA//9lZlNKDQAAAA==",
			},
			wantRet: []byte("1351451451512"),
			wantErr: false,
		},
		{
			name: "测试非base64",
			args: args{
				COMPRESS_GZIP,
				`{"msg":"这是一个json"}`,
			},
			wantErr: true,
		},
		{
			name: "测试非zip的base64",
			args: args{
				COMPRESS_GZIP,
				common_base64,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, err := UnCompress(tt.args.mode, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnCompress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("UnCompress() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}
