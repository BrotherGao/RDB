// This is a very basic example of a program that implements rdb.decoder and
// outputs a human readable diffable dump of the rdb file.
package main

import (
	"fmt"
	"os"
	"github.com/BrotherGao/RDB"

	"github.com/BrotherGao/RDB/nopdecoder"
)

type decoder struct {
	db int
	i  int
	nopdecoder.NopDecoder
}

func (p *decoder) StartDatabase(n int) {
	p.db = n
	fmt.Printf("Start parsing DB%d\n", p.db)
}

func (p *decoder) EndDatabase(n int) {
	fmt.Printf("Finish parsing DB%d\n", p.db)
}

func (p *decoder) Aux(key, value []byte) {
	fmt.Printf("aux_key=%q\n", key)
	fmt.Printf("aux_value=%q\n", value)
}

func (p *decoder) Set(key, value []byte, expiry int64) {
	fmt.Printf("db=%d %q -> %q ttl=%d\n", p.db, key, value, expiry)
}

func (p *decoder) Hset(key, field, value []byte) {
	fmt.Printf("db=%d %q . %q -> %q\n", p.db, key, field, value)
}

func (p *decoder) Sadd(key, member []byte) {
	fmt.Printf("db=%d %q { %q }\n", p.db, key, member)
}

func (p *decoder) StartList(key []byte, length, expiry int64) {
	p.i = 0
}

func (p *decoder) Rpush(key, value []byte) {
	fmt.Printf("db=%d %q[%d] -> %q\n", p.db, key, p.i, value)
	p.i++
}

func (p *decoder) StartZSet(key []byte, cardinality, expiry int64) {
	p.i = 0
}

func (p *decoder) Zadd(key []byte, score float64, member []byte) {
	fmt.Printf("db=%d %q[%d] -> {%q, score=%g}\n", p.db, key, p.i, member, score)
	p.i++
}

func (p *decoder) StartRDB() {
	fmt.Println("Start parsing RDB")
}

func (p *decoder) EndRDB() {
	fmt.Println("Finish parsing RDB")
}

func maybeFatal(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	f, err := os.Open(os.Args[1])
	maybeFatal(err)
	err = rdb.Decode(f, &decoder{})
	maybeFatal(err)
}
