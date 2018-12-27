## RDB

&emsp;&emsp;rdb是一个用于解析Redis的RDB文件的Go包。解析RDB格式参照[Redis RDB File Format](https://github.com/sripathikrishnan/redis-rdb-tools/blob/master/docs/RDB_File_Format.textile)   
&emsp;&emsp;该包是在[cupcake/rdb](https://github.com/cupcake/rdb)基础上修改。增加对RDB Version 8的支持，以及bug的修复。

RDB Version 8改动部分：
*  Lua脚本可以持久化到RDB文件中，类型为RDB_OPCODE_AUX，以key-value的形式持久化。其中，key为"lua"，value为对应的脚本内容
*  增加RDB_TYPE_ZSET_2类型，浮点类型不在以字符串的形式保存，而是以binary形式保存到RDB中去
*  增加数据的长度增加RDB_64BITLEN类型
*  增加RDB_TYPE_MODULE类型，Redis 4.0引入Module模块。(该包不支持对该部分的解析)

## 使用
如下是示例程序部分代码
```go
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

func main() {
	f, err := os.Open(os.Args[1])
	//处理err
	err = rdb.Decode(f, &decoder{})
	//处理err
}
```
运行结果如下：
```powershell
Start parsing RDB
aux_key="redis-ver"
aux_value="4.0.10"
aux_key="redis-bits"
aux_value="64"
aux_key="ctime"
aux_value="1545881958"
aux_key="used-mem"
aux_value="2081152"
aux_key="repl-stream-db"
aux_value="0"
aux_key="repl-id"
aux_value="8d4eda6495c7cb8ae7c5ba0626abbcc79c56f9ca"
aux_key="repl-offset"
aux_value="0"
aux_key="aof-preamble"
aux_value="0"
Start parsing DB0
db=0 "klose" -> "11" ttl=0
Finish parsing DB0
Start parsing DB1
db=1 "klose"[0] -> {"b", score=1.123456789}
db=1 "klose"[1] -> {"a", score=11.345}
db=1 "klose"[2] -> {"c", score=123}
Finish parsing DB1
Start parsing DB2
db=2 "klose" . "a" -> "1"
db=2 "klose" . "b" -> "2"
db=2 "klose" . "c" -> "3"
Finish parsing DB2
Start parsing DB11
db=11 "klose" -> "123456" ttl=0
db=11 "list"[0] -> "a"
db=11 "list"[1] -> "b"
db=11 "list"[2] -> "c"
db=11 "list"[3] -> "d"
db=11 "list"[4] -> "e"
aux_key="lua"
aux_value="return redis.call('get', KEYS[1])"
aux_key="lua"
aux_value="return redis.call('llen', KEYS[1])"
Finish parsing DB11
Finish parsing RDB
```
具体参见example目录下的test.go文件。fixture目录下为用于测试的各个版本的rdb文件

## 安装

```
go get github.com/BrotherGao/RDB
```
