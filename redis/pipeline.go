// -*- coding: utf-8 -*-

package redis


//
// Begin a transaction
//
func (client *Client) PipelineBegin() error {
    _, err := client.Do("MULTI")
    return err
}

//
// Send a queued command
//
func (client *Client) Pipeline(cmd string, args ...interface{}) error {
    _, err := client.Do(cmd, args...)
    return err
}

//
// Execute all queued commands and get replies
//
func (client *Client) PipelineCommit() ([]*Variable, error) {
    va, err := client.Do("EXEC")
    if err != nil {
	return nil, err
    }
    return va.Array(), err
}

//
// 
func (client *Client) PipelineRollback() ([]*Variable, error) {
    va, err := client.Do("DISCARD")
    if err != nil {
	return nil, err
    }
    return va.Array(), err
}
