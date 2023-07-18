package main

import (
	"fmt"
	"sync"
	"time"
)

var ConnectionPool = sync.Pool{
	New: func() interface{} {
		return &Connection{}
	},
}

func GetConn() *Connection {
	cli := ConnectionPool.Get().(*Connection)
	return cli
}

func PutConn(cli *Connection) {
	cli.reset()
	ConnectionPool.Put(cli) // 放回连接池
}

var redisOpZRangeCost = time.Millisecond * 1000
var redisOpDeleteCost = time.Millisecond * 100

var clients = sync.Map{}

// Connection Operator structure
type Connection struct {
	Id  int
	Uid int
}

func main() {
	go clear()

	fmt.Println("....add....", time.Now().Unix())
	for j := 1; j < 10000; j++ {
		connect(j, j)
	}
	fmt.Println("....delete....", time.Now().Unix())

	go func() {
		for j := 1; j < 10000; j++ {
			disconnect(get(j))
		}
	}()
	select {}
}

func connect(id, uid int) {
	c := GetConn()
	c.Id = id
	c.Uid = uid
	clients.Store(c.Id, c)
}

func disconnect(c *Connection) {
	if c == nil {
		return
	}
	clients.Delete(c.Id)
	// redis 操作
	time.Sleep(redisOpDeleteCost)
	PutConn(c)
}

func get(id int) *Connection {
	val, ok := clients.Load(id)
	if ok {
		return val.(*Connection)
	}
	return nil
}

func clear() {
	for {
		fmt.Println("....clear....", time.Now().Unix())

		clients.Range(func(_, v interface{}) bool {
			client := v.(*Connection)
			zRangeOp(client.Id, client.Uid)
			return true
		})
		time.Sleep(time.Second * 3)
	}
}

func zRangeOp(id, uid int) {
	time.Sleep(redisOpZRangeCost)

	if uid == 0 {
		fmt.Println(id, uid, "..............Fatal Error..................................................", time.Now().Unix())
	}
	fmt.Println(id, uid)

}

func (cli *Connection) reset() {
	cli.Uid = 0
	cli.Id = 0
}
