package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

// watch 可以监听key的行为， etcd本质上是一个存储k-v的玩意儿，有一点像redis
func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("获取etcd client失败， 连接etcd失败; ", err)
	}
	defer cli.Close()

	// watch 返回一个监听指定key的通道
	watchChan := cli.Watch(context.TODO(), "/wb/info", clientv3.WithPrefix()) // 监听所有前缀是/wb/info的key
	//watchChan := cli.Watch(context.TODO(), "/wb/info")

	for watchResponse := range watchChan { // 监听通道，循环监听
		for _, event := range watchResponse.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", event.Type, event.Kv.Key, event.Kv.Value)
		}
	}

}
