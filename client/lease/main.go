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
		log.Fatal(err)
		return
	}
	defer cli.Close()

	// client实现了lease接口，所以有grant方法
	leaseGrantResp, err := cli.Grant(context.TODO(), 5) // 创建一个5秒的租约,返回的是一个租约
	if err != nil {
		log.Fatal("创建租约失败：", err)
		return
	}
	fmt.Printf("租约id：%v , 租期：%v \n", leaseGrantResp.ID, leaseGrantResp.TTL)

	// 在put的时候加入租约， 租约到期后这个key就会被移除
	ops := []clientv3.OpOption{clientv3.WithLease(leaseGrantResp.ID)}
	_, err = cli.Put(context.TODO(), "/wb/study", "127.0.0.1:80888", ops...)
	_, err = cli.Put(context.TODO(), "/wb/study/info", "127.111.0.1:80899", ops...) // 租约可以作用多个key
	if err != nil {
		log.Fatal(err)
	}

	// 可以使用keepAlive向服务端发请求，来使租约续租为永久, 方法返回的是一个通道
	leaseKeepAliveResponse, err := cli.KeepAlive(context.TODO(), leaseGrantResp.ID)
	if err != nil {
		log.Fatal(err)
		return
	}

	for { // 循环监听返回的通道， 监听服务端的keepAlive响应信息
		v := <-leaseKeepAliveResponse
		fmt.Printf("%v \n", v.TTL)
		fmt.Printf("%v \n", v)
		if v != nil { // 收到keepAlive响应后跳出
			break
		}
	}

	// 睡眠6秒等待租约到期
	time.Sleep(6 * time.Second)

	getResp, err := cli.Get(context.TODO(), "/wb/study", clientv3.WithPrefix())
	for _, kv := range getResp.Kvs {
		fmt.Printf("%s : %s \n", kv.Key, kv.Value)
	}
}
