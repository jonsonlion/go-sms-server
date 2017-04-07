package convert

import (
)
import (
	"strings"
	"strconv"
)

/**
  uint8数组转为字符串
*/
func B2S(bs []uint8) string {
	b := make([]string, len(bs))
	for i, v := range bs {
		b[i] = strconv.Itoa(int(v))
	}
	return strings.Join(b, "")
}


func BytesToString(bs []byte) string {
	switch interface{}(bs).(type) {
	case []uint8:
		return B2S(bs)
	default:
		return string(bs)
	}
}

func StringToBytes(s string) []byte{
	b := make([]byte, 0, len(s))
	for _ , c := range s{
		v := int(c) - 48
		b = append(b, byte(v))
	}
	return b
}
