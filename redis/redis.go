// -*- coding: utf-8 -*-

package redis


import (
//    "fmt"
    "time"
    "net"
)

const (
    DefaultMaxTimeCommand = 20
)


type Client struct {
    Conn net.Conn
    Db int
    MaxTimeCommand time.Duration
}

type Reply struct{}

//
// Connect(tcp) to redis server, initialize a new client
//
func NewClient(network, host string) (*Client, error) {
    conn, err := net.Dial(network, host)
    if err != nil {
	return nil, err
    }
    client := &Client{
	Conn: conn,
	Db: 0,
	MaxTimeCommand: DefaultMaxTimeCommand,
    }
    return client, nil
}

func (client *Client) Close() {
    client.Conn.Close()
}

func _Timeout(ch chan int, seconds time.Duration){
    time.Sleep(seconds * time.Second)
    ch <- 0
}

func (client *Client) Do(cmd string, args ...interface{}) (*Reply, error) {
    timeout := client.MaxTimeCommand
    if timeout <= 0 {
    	timeout = DefaultMaxTimeCommand
    }
    client.Conn.SetDeadline(time.Now().Add(time.Second * timeout))
    client.Conn.Write(_NormalizeCommand(cmd, args))
    return nil, nil
}


