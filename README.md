# Redis client: implementation of go


## API Docs

[http://godoc.org/github.com/smallfz/small-redis/redis](http://godoc.org/github.com/smallfz/small-redis/redis)

## Samples

**Initialization and clean up**

	rc, _ := redis.NewClient("tcp", "localhost:6379")
    defer rc.Close()

**Array response:**

    result, _ := rc.Do("KEYS", "*")
    for _, a := range result.StringArray() {
		fmt.Fprintln(w, a)
    }

**Numeric response:**

	result, _ = rc.Do("LLEN", "test")
    fmt.Fprintf(w, "test's len: %d", result.Integer())

### Pipeline/Transaction

	rc.PipelineBegin()
    rc.Pipeline("keys", "*")
    rc.Pipeline("lrange", "test", 0, -1)
    rc.Pipeline("expire", "v1", 0)
    vars, _ := rc.PipelineCommit()
    for _, v := range vars {
		fmt.Fprintf(w, "%v\r\n", v)
    }



