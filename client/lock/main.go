package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
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

	// 创建一个会话
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer s1.Close()

	// 创建另一个会话
	s2, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer s2.Close()

	m1 := concurrency.NewMutex(s1, "/my-lock/")
	m2 := concurrency.NewMutex(s2, "/my-lock/")

	// 会话1获取锁
	if err := m1.Lock(context.TODO()); err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("s1 获取锁")

	m2locked := make(chan struct{})
	go func() {
		defer close(m2locked)
		if err := m2.Lock(context.TODO()); err != nil {
			log.Fatal(err)
			return
		}
	}()

	if err := m1.Unlock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("s1 释放锁")

	<-m2locked
	fmt.Println("s2获得锁")

}
