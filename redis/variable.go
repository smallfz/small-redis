// -*- coding: utf-8 -*-

package redis

import (
    "fmt"
    "strings"
    "strconv"
    "reflect"
)

var (
    crlf = []byte{'\r','\n'}
    delim byte = '\n'
)

const (
    TYPE_INT = "int"
    TYPE_STR = "string"
    TYPE_ARR = "array"
    TYPE_ERR = "string_error"
)

type Variable struct {
    typeCode byte
    data []byte
    items []*Variable
    nilArray bool
    nilString bool
}

//
// Get the type name of the variable. 
//
// TYPE_INT: integer
// TYPE_STR: string
// TYPE_ARR: array
// TYPE_ERR: string contains redis error message
//
func (v *Variable) Type() string {
    switch v.typeCode {
    case '+', '$':
	return TYPE_STR
    case ':':
	return TYPE_INT
    case '*':
	return TYPE_ARR
    case '-':
	return TYPE_ERR
    default:
	return ""
    }
}

//
// Get integer value of the variable.
// if the variable contains a string or other, try to parse it to int.
//
func (v *Variable) Vtoi() (int, error){
    if v.typeCode == ':' {
	return strconv.Atoi(string(v.data))
    }else{
	return strconv.Atoi(v.String())
    }
}

//
// Same as Vtoi, just no errors.
//
func (v *Variable) Integer() int {
    i, _ := v.Vtoi()
    return i
}

func (v *Variable) Array() []*Variable {
    if v.nilArray {
	return nil
    }
    return v.items
}

func (v *Variable) StringArray() []string {
    switch v.typeCode {
    case '+', '-', ':', '$':
	return []string{ v.String() }
    case '*':
	a := make([]string, len(v.items))
	for i:=0; i<len(a); i+=1 {
	    a[i] = v.items[i].String()
	}
	return a
    default:
	return nil
    }
}

func (v *Variable) String() string {
    if v.nilString {
	return ""
    }
    switch v.typeCode {
    case '+', '-':
	return string(v.data)
    case ':':
	return string(v.data)
    case '$':
	return string(v.data)
    case '*':
	a := v.StringArray()
	return strings.Join(a, ", ")
    }
    return ""
}

func (v *Variable) Bytes() ([]byte, error) {
    var b []byte
    switch v.typeCode {
    case ':', '+', '-':
	b = append(b, v.typeCode)
	// b = append(b, crlf...)
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
	if v.nilArray {
	    v.nilArray = false
	}
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
    case bool:
	va.typeCode = ':'
	if v {
	    va.data = append(va.data, '1')
	}else{
	    va.data = append(va.data, '0')
	}
    case string:
	va.typeCode = '$'
	if item == nil {
	    va.nilString = true
	}else{
	    va.data = append(va.data, v...)
	}
	break
    case []byte:
	va.typeCode = '$'
	if item == nil {
	    va.nilString = true
	}else{
	    va.data = append(va.data, v...)
	}
	break
    default:
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
	    va.typeCode = '*'
	    if item == nil {
		va.nilArray = true
	    }else{
		fa := reflect.ValueOf(v)
		for i:=0; i<fa.Len(); i+=1 {
		    subVa, err := NewVariable(fa.Index(i).Interface())
		    if err != nil {
			return nil, err
		    }
		    va.items = append(va.items, subVa)
		}
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

func ReadByCount(reader Reader, count int) ([]byte, error) {
    var all []byte
    for len(all) < count {
	b := make([]byte, count)
	n, err := reader.Read(b)
	if err != nil {
	    return all, err
	}
	all = append(all, b[:n]...)
    }
    return all, nil
}

func NewVariableFromReader(reader Reader) (*Variable, error) {
    va := &Variable{}
    t, _ := reader.ReadByte()
    switch t {
    case ':', '+', '-':
	line, err := reader.ReadBytes(delim)
	if err != nil {
	    return nil, err
	}
	va.typeCode = t
	va.data = line[:len(line)-2]
	break
    case '$':
	va.typeCode = t
	line, err := reader.ReadBytes(delim)
	if err != nil {
	    return nil, err
	}
	strLen, _ := strconv.Atoi(string(line[:len(line)-2]))
	if strLen >= 0 {
	    b, err := ReadByCount(reader, strLen + 2)
	    if err != nil {
		return nil, err
	    }
	    va.data = b[:len(b)-2]
	}else{
	    va.nilString = true
	}
	break
    case '*':
	line, err := reader.ReadBytes(delim)
	if err != nil {
	    return nil, err
	}
	arrLen, _ := strconv.Atoi(string(line[:len(line)-2]))
	va.typeCode = '*'
	if arrLen >= 0 {
	    for i:=0; i<arrLen; i+=1 {
		subVa, _ := NewVariableFromReader(reader)
		va.items = append(va.items, subVa)
	    }
	}else{
	    va.nilArray = true
	}
	break
    default:
	return nil, Err("Invalid leading byte for typeCode.")
    }
    return va, nil
}

func Command(cmd string, args []interface{}) ([]byte, error) {
    strArgs := make([]string, len(args))
    for i:=0; i<len(strArgs); i+=1 {
	arg := args[i]
	switch a := arg.(type) {
	case string:
	    strArgs[i] = a
	    break
	case []byte:
	    strArgs[i] = string(a)
	    break
	case int, uint:
	    strArgs[i] = fmt.Sprintf("%d", a)
	    break
	case bool:
	    if a {
		strArgs[i] = "1"
	    }else{
		strArgs[i] = "0"
	    }
	default:
	    strArgs[i] = fmt.Sprintf("%v", a)
	    break
	}
    }
    va, err := NewVariable(strArgs)
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

