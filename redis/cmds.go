// -*- coding: utf-8 -*-

package redis


func (client *Client) Get(key string) string {
    v, err := client.Do("GET", key)
    if err != nil {
	return ""
    }
    return v.String()
}

func (client *Client) Set(key string, value interface{}) {
    client.Do("SET", key, value)
}

//
// Set a timeout on key.
//
func (client *Client) Expire(key string, seconds int) int {
    v, _ := client.Do("EXPIRE", key, seconds)
    return v.Integer()
}

func (client *Client) Del(keys ...string) int {
    args := make([]interface{}, len(keys))
    for i:=0; i<len(args); i+=1 {
	args[i] = keys[i]
    }
    v, _ := client.Do("DEL", args...)
    return v.Integer()
}
