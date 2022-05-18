package main

import (
	"fmt"
	"sync"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/gotomicro/ego/core/elog"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	dns = ""
)

func main() {

	query := ""
	query2 := ""

	fmt.Println(query)
	timer := 1

	tn := time.Now()
	for i := 0; i < timer; i++ {
		// do()
		doQuery(query)

	}
	tnc := time.Since(tn)

	wg := sync.WaitGroup{}
	t := time.Now()
	for i := 0; i < timer; i++ {
		wg.Add(1)
		go func(k int, timeline time.Time) {
			doQuery(query2)
			wg.Done()
		}(i, t)
	}
	wg.Wait()
	elog.Info("optimize", elog.String("step", "done"), elog.Any("costCon", time.Since(t)), elog.Any("costSerial", tnc))

}

const clickHouseDriverName = "clickhouse"

func New(datasource string, opts ...sqlx.SqlOption) sqlx.SqlConn {
	return sqlx.NewSqlConn(clickHouseDriverName, datasource, opts...)
}

var (
	clientxPool = sync.Map{}
)

func TakeClientx(dsn string) sqlx.SqlConn {
	c, ok := clientxPool.Load(dsn)
	if !ok {
		c := New(dsn)
		clientxPool.Store(dsn, c)
		return c
	}
	return c.(sqlx.SqlConn)
}

func doQuery(query string) (res interface{}, err error) {
	t := time.Now()
	elog.Info("optimize", elog.String("step", "start"), elog.Any("cost", time.Since(t)))
	connx := TakeClientx("dns")
	err = connx.QueryRow(&res, query)
	if err != nil {
		elog.Error("ClickHouse", elog.Any("step", "rows"), elog.Any("error", err.Error()))
		return
	}
	elog.Info("optimize", elog.String("step", "finish"), elog.Any("cost", time.Since(t)))
	return
}
