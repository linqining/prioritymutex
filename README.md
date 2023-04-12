# PriorityMutex
PriorityMutex provides a priority lock implementation in go, in some case, we expect to do some operations to resources with priority, such as collection some items and fire a batch resource request.
## Installation
Install PriorityMutes using go get command:

    $go get github.com/linqining/prioritymutex

## Example

```go
package main

import (
	redis "github.com/gomodule/redigo/redis"
	"github.com/linqining/prioritymutex"
)

var rdsClient redis.Conn

func init() {
	ConnectRedis()
}

func ConnectRedis() redis.Conn {
	var err error
	rdsClient, err = redis.Dial("tcp", "localhost:6379", redis.DialPassword("123456"))
	if err != nil {
		panic(err)
	}
	return rdsClient
}

var elms []int64

func Foo(p *prioritymutex.PriorityMutex, uids []int64) {
	p.PLock()
	defer p.PUnlock()
	elms = append(elms, uids...)
}

func Bar(p *prioritymutex.PriorityMutex) {
	p.Lock()
	defer p.Unlock()
	tmp := elms
	elms = []int64{}
	args := []interface{}{}
	args = append(args, "KEY")
	for _, elm := range tmp {
		args = append(args, elm)
	}
	rdsClient.Do("HGET", args...)
}

func main() {
	p := &prioritymutex.PriorityMutex{}
	go Foo(p, []int64{1, 2, 3, 4})
	Bar(p)
}
```