// -*- coding: utf-8 -*-

package redis

import (
    "strconv"
    "reflect"
)

var (
    crlf = []byte{'\r','\n'}
    delim byte = '\n'
)

type Variable struct {
    typeCode byte
    data []byte
    items []*Variable
    nilArray bool
    nilString bool
}

func (v *Variable) Bytes() ([]byte, error) {
    var b []byte
    switch v.typeCode {
    case ':', '+', '-':
	b = append(b, v.typeCode)
	b = append(b, crlf...)
	b = append(b, v.data...)
	b = append(b, crlf...)
	break
    case '$':
	if v.nilString {
	    b = append(b, []byte{'$','-','1','\r','\n'}...)
	    break
	}
	b = append(b, v.typeCode)
	b = append(b, strconv.Itoa(len(v.data))...)
	b = append(b, crlf...)
	b = append(b, v.data...)
	b = append(b, crlf...)
	break
    case '*':
	if v.nilArray {
	    b = append(b, []byte{'*','-','1','\r','\n'}...)
	    break
	}
	b = append(b, v.typeCode)
	b = append(b, strconv.Itoa(len(v.items))...)
	b = append(b, crlf...)
	for _, item := range v.items {
	    a, _ := item.Bytes()
	    b = append(b, a...)
	}
	break
    default:
	return nil, Err("Invalid variable type code.")
    }
    return b, nil
}

func (v *Variable) LPush(item *Variable) error {
    if v.typeCode == '*' {
	var items []*Variable
	items = append(items, item)
	items = append(items, v.items...)
	v.items = items
	return nil
    }
    return Err("LPush can only happen on array.")
}

func NewVariableSimpleString(t string) (*Variable) {
    va := &Variable{
	typeCode: '+',
	data: []byte(t),
    }
    return va
}

func NewVariableErrorString(t string) (*Variable) {
    va := &Variable{
	typeCode: '-',
	data: []byte(t),
    }
    return va
}

func NewVariable(item interface{}) (*Variable, error) {
    va := &Variable{}
    switch v := item.(type) {
    case int:
	va.typeCode = ':'
	va.data = append(va.data, strconv.FormatInt(int64(v), 10)...)
	break
    case uint:
	va.typeCode = ':'
	va.data = append(va.data, strconv.FormatUint(uint64(v), 10)...)
	break
    case string:
	va.typeCode = '$'
	va.data = append(va.data, v...)
	break
    case []byte:
	va.typeCode = '$'
	va.data = append(va.data, v...)
    default:
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
	    va.typeCode = '*'
	    fa := reflect.ValueOf(v)
	    for i:=0; i<fa.Len(); i+=1 {
		subVa, err := NewVariable(fa.Index(i).Interface())
		if err != nil {
		    return nil, err
		}
		va.items = append(va.items, subVa)
	    }
	    break
	default:
	    return nil, Err("Invalid type for new variable.")
	}
	break
    }
    return va, nil
}

type Reader interface {
    Read([]byte) (n int, err error)
    ReadByte()(c byte, err error)
    ReadBytes(delim byte)(b []byte, err error)
}

func ReadByCount(reader Reader, count int) []byte {
    b := make([]byte, count)
    n, _ := reader.Read(b)
    return b[:n]
}

func NewVariableFromReader(reader Reader) (*Variable, error) {
    va := &Variable{}
    t, _ := reader.ReadByte()
    switch t {
    case ':', '+', '-':
	line, _ := reader.ReadBytes(delim)
	va.typeCode = t
	va.data = line[:len(line)-2]
	break
    case '$':
	va.typeCode = t
	line, _ := reader.ReadBytes(delim)
	strLen, _ := strconv.Atoi(string(line[:len(line)-2]))
	b := ReadByCount(reader, strLen + 2)
	va.data = b[:len(b)-2]
	break
    }
    return va, nil
}

func NormalizeCommand(cmd string, args []interface{}) ([]byte, error) {
    va, err := NewVariable(args)
    if err != nil {
	return nil, err
    }
    a, err := NewVariable(cmd)
    if err != nil {
	return nil, err
    }
    err = va.LPush(a)
    if err != nil {
	return nil, err
    }
    b, err := va.Bytes()
    if err != nil {
	return nil, err
    }
    return b, nil
}

