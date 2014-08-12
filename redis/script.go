// -*- coding: utf-8 -*-

package redis


//
// Load a script into the scripts cache, without executing it.
// Returns the digest of the script.
//
func (client *Client) ScriptLoad(scriptText string) string {
    v, _ := client.Do("SCRIPT", "LOAD", scriptText)
    return v.String()
}

func (client *Client) ScriptExists(scriptHash string) bool {
    v, err := client.Do("SCRIPT", "EXISTS", scriptHash)
    if err == nil {
	if v.Type() == TYPE_ARR {
	    return v.Array()[0].Bool()
	}else{
	    return v.Bool()
	}
    }
    return false
}

func (client *Client) ScriptKill(scriptHash string) {
    client.Do("SCRIPT", "KILL", scriptHash)
}

//
// Flush the Lua scripts cache.
//
func (client *Client) ScriptFlush(){
    client.Do("SCRIPT", "FLUSH")
}

//
// Execute a lua script
//
func (client *Client) Eval(scriptText string, keys, args []interface{}) (*Variable, error) {
    _args := make([]interface{}, 2+len(keys)+len(args))
    _args[0] = scriptText
    _args[1] = len(keys)
    for i:=0; i<len(keys); i+=1 {
	_args[i+2] = keys[i]
    }
    for i:=0; i<len(args); i+=1 {
	_args[i+2+len(keys)] = args[i]
    }
    return client.Do("EVAL", _args...)
}

//
// Evaluates a script cached on the server side by its hash digest.
//
func (client *Client) EvalSha(scriptHash string, keys, args []interface{}) (*Variable, error) {
    _args := make([]interface{}, 2+len(keys)+len(args))
    _args[0] = scriptHash
    _args[1] = len(keys)
    for i:=0; i<len(keys); i+=1 {
	_args[i+2] = keys[i]
    }
    for i:=0; i<len(args); i+=1 {
	_args[i+2+len(keys)] = args[i]
    }
    return client.Do("EVALSHA", _args...)
}

