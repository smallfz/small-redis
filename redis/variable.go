// -*- coding: utf-8 -*-

package redis

import (
//    "bytes"
    "strconv"
    "reflect"
)

var (
    crlf = []byte{'\r','\n'}
)

// type Buffer interface {
//     Bytes() []byte
//     Write(a []byte) (n int, err error)
//     WriteByte(b byte) (err error)
//     WriteString(s string) (n int, err error)
// }

// func NewBuffer() Buffer {
//     return &bytes.Buffer{}
// }

func _String(t string) []byte {
    var b []byte
    bs := []byte(t)
    b = append(b, '$')
    b = append(b, strconv.Itoa(len(bs))...)
    b = append(b, crlf...)
    b = append(b, bs...)
    b = append(b, crlf...)
    return b
}

func Variable(item interface{}) []byte {
    var b []byte
    switch v := item.(type) {
    case int:
	b = append(b, ':')
	b = append(b, strconv.FormatInt(int64(v), 10)...)
	b = append(b, crlf...)
	break
    case uint:
	b = append(b, ':')
	b = append(b, strconv.FormatUint(uint64(v), 10)...)
	b = append(b, crlf...)
	break
    case string:
	b = append(b, _String(v)...)
	break
    case []byte:
	b = append(b, '$')
	b = append(b, strconv.Itoa(len(v))...)
	b = append(b, crlf...)
	b = append(b, v...)
	b = append(b, crlf...)
	break
    default:
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
	    fa := reflect.ValueOf(v)
	    var a []interface{}
	    for i:=0; i<fa.Len(); i+=1 {
		a = append(a, fa.Index(i).Interface())
	    }
	    b = append(b, Array(a)...)
	default:
	    panic("Unable to serialize the argument.")
	    break
	}
	break
    }
    return b
}

func _ArrayBody(items []interface{}) []byte {
    var b []byte
    for _, item := range items {
	b = append(b, Variable(item)...)
    }
    return b
}

func Array(items []interface{}) []byte {
    var b []byte
    cnt := len(items)
    b = append(b, '*')
    b = append(b, strconv.Itoa(cnt)...)
    b = append(b, crlf...)
    b = append(b, _ArrayBody(items)...)
    b = append(b, crlf...)
    return b
}

func NormalizeCommand(cmd string, args []interface{}) []byte {
    var b []byte
    cnt := 1 + len(args)
    b = append(b, '*')
    b = append(b, strconv.Itoa(cnt)...)
    b = append(b, crlf...)
    b = append(b, _String(cmd)...)
    b = append(b, _ArrayBody(args)...)
    return b
}

