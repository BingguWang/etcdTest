package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("获取etcd client失败， 连接etcd错误： ", err)
	}

	defer cli.Close()

	// put
	putResp, err := cli.Put(context.TODO(), "/wb/info", "123.456.123.123:8080")
	putResp, err = cli.Put(context.TODO(), "/wb/info/etc", "123.456.123.124:8081")
	if err != nil {
		log.Fatal("put key失败：", err)
		return
	}
	fmt.Println(*putResp)

	// get, withPrefix() 把前缀是/wb/info的都找出来
	getResp, err := cli.Get(context.TODO(), "/wb/info", clientv3.WithPrefix())
	if err != nil {
		log.Fatal("get key 失败: ", err)
	}
	fmt.Println(*getResp)
	for _, kv := range getResp.Kvs {
		fmt.Printf("%s : %s \n", kv.Key, kv.Value)
	}

	//delete
	//cli.Delete(context.TODO(),"/wb/info")

}
