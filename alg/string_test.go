package alg

import (
	"bytes"
	"strings"
	"testing"
)

var strList = []string{
	"key1", "value1",
	"key2", "value2",
	"key3", "value3",
	"key4", "value4",
	"key5", "value5",
	"key6", "value6",
	"key7", "value7",
	"key8", "value8",
	"key9", "value9",
}

func BenchmarkStrAppendJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var list []string
		for i := 0; i < len(strList); i += 2 {
			list = append(list, strList[i]+"="+strList[i+1])
		}
		strings.Join(list, "&")
	}
}

func BenchmarkStrPlusConcat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var str string
		for i := 0; i < len(strList); i += 2 {
			var gap string
			if i > 0 {
				gap = "&"
			}
			str += gap + strList[i] + "=" + strList[i+1]
		}
	}
}

func BenchmarkStrBuff(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		for i := 0; i < len(strList); i += 2 {
			if i > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(strList[i])
			buf.WriteByte('=')
			buf.WriteString(strList[i+1])
		}
		_ = buf.String()
	}
}
