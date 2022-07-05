package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
)

var (
	rdb    *redis.Client
	ctx    = context.Background()
	server = flag.String("server", "localhost:6379", "redis服务器地址")
	db     = flag.Int("db", 0, "redis选库")
	count  = flag.Int64("count", 1000, "redis count")
	size   = flag.Int64("size", 10240, "redis大键大小")
	pass   = flag.String("pass", "", "redisl连接密码")
	fname  = flag.String("fname", "output.log", "保存存文件名")
	save   = flag.Bool("save", false, "是否到文件")
)

func init() {
	flag.Parse()
	rdb = redis.NewClient(&redis.Options{
		Addr:     *server,
		Password: *pass, // no password set
		DB:       *db,   // use default DB
	})

}

func main() {
	output, _ := os.OpenFile("./"+*fname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer output.Close()
	iter := rdb.Scan(ctx, 0, "*", *count).Iterator()
	for iter.Next(ctx) {
		gsize, _ := rdb.MemoryUsage(ctx, iter.Val()).Result()
		if gsize > *size {
			types, _ := rdb.Type(ctx, iter.Val()).Result()
			if *save {
				fmt.Fprintf(output, "%-50s%-10s%-10d\n", iter.Val(), types, gsize)
			} else {
				fmt.Printf("%-50s%-10s%-10d\n", iter.Val(), types, gsize)
			}
		}
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
}
