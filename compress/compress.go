package compress

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
)

const COMPRESS_GZIP = 1

// gzip压缩，[]byte -> gzip(bin) -> base64(string)
// 暂时只支持gzip
func Compress(mode int, data []byte) (ret string, err error) {
	var in bytes.Buffer
	gz := gzip.NewWriter(&in)
	if _, err = gz.Write(data); err != nil {
		return
	}
	if err = gz.Close(); err != nil {
		return
	}

	ret = base64.StdEncoding.EncodeToString(in.Bytes())

	return
}

// gzip解压缩，base64(string) -> gzip(bin) -> []byte
func UnCompress(mode int, data string) (ret []byte, err error) {
	src_data, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return
	}

	gz_reader, err := gzip.NewReader(bytes.NewReader(src_data))
	if err != nil {
		return
	}

	ret, err = ioutil.ReadAll(gz_reader)

	return
}
