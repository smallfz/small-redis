// -*- coding: utf-8 -*-

package redis


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

// func (ps *PubSub) 
