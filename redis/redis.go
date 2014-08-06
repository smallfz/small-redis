// -*- coding: utf-8 -*-

package redis


import (
//    "fmt"
    "time"
    "net"
    "bufio"
)

var (
    DefaultMaxTimeCommand time.Duration = 20 * time.Second
)


type Client struct {
    network, host string
    db int
    Conn net.Conn
    maxTimeCommand time.Duration
    lastIOTime time.Time
}

//
// Connect(tcp) to redis server, initialize a new client
//
func NewClient(network, host string) (*Client, error) {
    client := &Client{
	network: network,
	host: host,
	db: 0,
	maxTimeCommand: DefaultMaxTimeCommand,
    }
    return client, nil
}

func (client *Client) Close() {
    if client.Conn != nil {
	client.Conn.Close()
    }
}

func _Timeout(ch chan int, seconds time.Duration){
    time.Sleep(seconds * time.Second)
    ch <- 0
}

func (client *Client) _Connect() (net.Conn, error){
    if client.Conn != nil {
	return client.Conn, nil
    }
    conn, err := net.Dial(client.network, client.host)
    if err != nil {
	return nil, err
    }
    client.Conn = conn
    return conn, err
}

func (client *Client) Do(cmd string, args ...interface{}) (*Reply, error) {
    timeout := client.maxTimeCommand
    if timeout <= 0 {
    	timeout = DefaultMaxTimeCommand
    }
    client._Connect()
    client.Conn.SetDeadline(time.Now().Add( timeout))
    client.Conn.Write(NormalizeCommand(cmd, args))
    reader := bufio.NewReader(client.Conn)
    reply, err := NewReplyFromReader(reader)
    return reply, err
}


