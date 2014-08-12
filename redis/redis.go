// -*- coding: utf-8 -*-

package redis


import (
    "time"
    "net"
    "bufio"
    "strings"
    "net/url"
    "regexp"
    "strconv"
    "fmt"
)

var (
    DefaultMaxTimeCommand time.Duration = 20 * time.Second
    DefaultRedisServerPort = 6379
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
    reader *bufio.Reader
    maxTimeCommand time.Duration
    lastIOTime time.Time
    // pipeline bool
    // pipelineCmds []*_Cmd
    // pipelineVars []*Variable
}

type _Cmd struct {
    command string
    args []interface{}
    va *Variable
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
// Initialize a new redis client from an URL.
// eg. "redis://localhost:6379/0"
//
func NewClientFromUrl(hostUrl string) (*Client, error) {
    _url, err := url.Parse(hostUrl)
    if err != nil {
	return nil, err
    }

    var db, port int

    ptPath := regexp.MustCompile("^\\/(\\d+)$")
    parts := ptPath.FindStringSubmatch(_url.Path)
    if parts == nil {
	db = 0
    }else{
	_db, err := strconv.Atoi(parts[1])
	if err == nil {
	    db = _db
	}
    }

    ptPort := regexp.MustCompile("\\:(\\d+)$")
    parts = ptPort.FindStringSubmatch(_url.Host)
    var hostStr string
    if parts == nil {
	// port not specified. 
	scheme := strings.ToLower(_url.Scheme)
	if scheme == "redis" {
	    port = DefaultRedisServerPort
	    hostStr = fmt.Sprintf("%s:%d", _url.Host, port)
	}else{
	    _errMsg := fmt.Sprintf("Protocol %s not supported.", scheme)
	    return nil, Err(_errMsg)
	}
    }else{
	hostStr = _url.Host
    }
    rc, err := NewClient("tcp", hostStr)
    if err != nil {
	return nil, err
    }
    if db != 0 {
	rc.Select(db)
    }
    return rc, nil
}

//
// Close client connection, release all resources
//
func (client *Client) Close() {
    if client.Conn != nil {
	client.Conn.Close()
    }
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
    client.reader = bufio.NewReader(conn)
    return conn, err
}

func (client *Client) _Send(cmd string, args ...interface{}) error {
    timeout := client.maxTimeCommand
    if timeout <= 0 {
    	timeout = DefaultMaxTimeCommand
    }
    client._Connect()
    // client.Conn.SetDeadline(time.Now().Add( timeout))
    cmdBytes, err := Command(cmd, args)
    if err != nil {
	return err
    }
    client.Conn.Write(cmdBytes)
    return nil
}

//
// Send command to redis and get reply
//
func (client *Client) Do(cmd string, args ...interface{}) (*Variable, error) {
    cmd = strings.ToUpper(cmd)

    err := client._Send(cmd, args...)
    if err != nil {
	return nil, err
    }

    reader := client.reader
    ra, err := NewVariableFromReader(reader)

    if err == nil {
	if ra.Type() == TYPE_ERR {
	    err = Err(ra.String())
	    panic(err)
	    return ra, err
	}
    }
    
    return ra, err
}

func (client *Client) Select(db int) {
    client.Do("SELECT", db)
}


