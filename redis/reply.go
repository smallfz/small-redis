// -*- coding: utf-8 -*-

package redis

var (
    delim byte = '\n'
)

type Reply struct{
    prefix byte
    data []byte
}

type Reader interface {
    ReadByte()(c byte, err error)
    ReadBytes(delim byte)(b []byte, err error)
}

func NewReplyFromReader(reader Reader) (*Reply, error){
    r := &Reply{}
    t := reader.ReadByte()
    switch t {
    case ':':
	line, _ := reader.ReadBytes(delim)
	r.prefix = t
	r.data = line[:-2]
    }
    return r, nil
}
