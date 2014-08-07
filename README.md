# Redis client: implementation of go

### Samples

Array response:

	rc, _ := redis.NewClient("tcp", "localhost:6379")
    defer rc.Close()
    result, _ := rc.Do("KEYS", "*")
    for _, a := range result.StringArray() {
		fmt.Fprintln(w, a)
    }

Numeric response:

	result, _ = rc.Do("LLEN", "test")
    fmt.Fprintf(w, "test's len: %d", result.Integer())

### API
[http://godoc.org/github.com/smallfz/small-redis/redis](http://godoc.org/github.com/smallfz/small-redis/redis)


