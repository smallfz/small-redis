// -*- coding: utf-8 -*-

package redis


import (
    "time"
    "net"
    "bufio"
)

var (
    DefaultMaxTimeCommand time.Duration = 20 * time.Second
)

type _Err struct{
    msg string
}

func (e *_Err) Error() string{
    return e.msg
}

func Err(msg string) (*_Err){
    return &_Err{msg: msg}
}

type Client struct {
    network, host string
    db int
    Conn net.Conn
    maxTimeCommand time.Duration
    lastIOTime time.Time
}

//
// Initialize a new redis client
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

//
// Close client connection, release all resources
//
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

//
// Send command to redis and get reply
//
func (client *Client) Do(cmd string, args ...interface{}) (*Variable, error) {
    timeout := client.maxTimeCommand
    if timeout <= 0 {
    	timeout = DefaultMaxTimeCommand
    }
    client._Connect()
    client.Conn.SetDeadline(time.Now().Add( timeout))
    cmdBytes, err := Command(cmd, args)
    if err != nil {
	return nil, err
    }
    client.Conn.Write(cmdBytes)
    reader := bufio.NewReader(client.Conn)
    ra, err := NewVariableFromReader(reader)
    return ra, err
}


