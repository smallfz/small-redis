// -*- coding: utf-8 -*-

package redis


const (
    MSG_TYPE_MESSAGE = "message"
    MSG_TYPE_PMESSAGE = "pmessage"
    MSG_TYPE_SUBSCRIBE = "subscribe"
    MSG_TYPE_UNSUBSCRIBE = "unsubscribe"
)

//
// publish a message to the channel
//
func (client *Client) Publish(channel, message string) error {
    _, err := client.Do("PUBLISH", channel, message)
    return err
}

type PubSub struct {
    client *Client
}

func NewPubSub(network, host string) (*PubSub, error){
    client, err := NewClient(network, host)
    if err != nil {
	return nil, err
    }
    ps := &PubSub{
	client: client,
    }
    return ps, nil
}

func (ps *PubSub) Subscribe(channels ...interface{}) (*Variable, error){
    err := ps.client._Send("SUBSCRIBE", channels...)
    if err != nil {
	return nil, err
    }
    return nil, nil
}

func (ps *PubSub) Unsubscribe(channels ...interface{}) (*Variable, error){
    err := ps.client._Send("SUBSCRIBE", channels...)
    if err != nil {
	return nil, err
    }
    return nil, nil
}

func (ps *PubSub) PSubscribe(patterns ...interface{}) (*Variable, error){
    err := ps.client._Send("PSUBSCRIBE", patterns...)
    if err != nil {
	return nil, err
    }
    return nil, nil
}

func (ps *PubSub) PUnsubscribe(patterns ...interface{}) (*Variable, error){
    err := ps.client._Send("PUNSUBSCRIBE", patterns...)
    if err != nil {
	return nil, err
    }
    return nil, nil
}

func (ps *PubSub) Close() {
    ps.client.Close()
}

type Message struct {
    typeCode string
    message string
    chCount int
}

func (msg *Message) String() string {
    return msg.message
}

func (msg *Message) Type() string {
    return msg.typeCode
}

func (ps *PubSub) _Listen(ch chan *Message) {
    defer close(ch)
    reader := ps.client.reader
    for {
	va, err := NewVariableFromReader(reader)
	if err != nil {
	    break
	}
	arr := va.Array()
	ch <- &Message{
	    typeCode: arr[0].String(),
	    message: arr[1].String(),
	    chCount: arr[2].Integer(),
	}
    }
}

func (ps *PubSub) ListenChan() (chan *Message) {
    ch := make(chan *Message)
    go ps._Listen(ch)
    return ch
}

//
// Listen to channel(s). Handle reply messages with the function.
// Stop listening once the function returns a false.
//
func (ps *PubSub) ListenFunc(handler func(*Message)(bool)) {
    if handler == nil {
	return
    }
    ch := ps.ListenChan()
    for msg := range ch {
	if !handler(msg) {
	    return
	}
    }
}

