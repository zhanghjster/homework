package main

import (
	"log"

	"github.com/coreos/etcd/clientv3"

	"context"

	"flag"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3/concurrency"
)

var ver string

func init() {
	flag.StringVar(&ver, "version", "v3", "version")
	flag.Parse()
}

func main() {
	switch ver {
	case "v2":
		v2()
	case "v3":
		v3()
	default:
		log.Fatal("wrong version")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func v3() {
	cfg := clientv3.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
	}

	c, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	var ctx = context.Background()

	var key = "/myKey"
	var value = "myKeyValue"

	// set the key
	_, err = c.Put(ctx, key, value)
	if err != nil {
		log.Fatal(err)
	}

	// get the key
	resp, err := c.Get(ctx, key)
	if err != nil {
		log.Fatal(err)
	}
	dumpKvs(resp)

	// prefix
	for i := 0; i < 5; i++ {
		c.Put(ctx, "/myKey2/"+strconv.Itoa(i), strconv.Itoa(i))
	}

	resp, err = c.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
	dumpKvs(resp)

	// watch
	watcher := c.Watch(ctx, key, clientv3.WithPrefix())
	go func() {
		for {
			resp, _ := <-watcher
			for i, event := range resp.Events {
				fmt.Printf(
					"get event: %d, %s, %s \n",
					i,
					string(event.Kv.Key),
					string(event.Kv.Value),
				)
			}
			time.Sleep(3 * time.Second)
		}
	}()
	go func() {
		for i := 0; i < 5; i++ {
			c.Put(ctx, key, strconv.Itoa(i))
			time.Sleep(1 * time.Second)
		}
	}()

	// lease
	leasedKey := "leasedKey"
	lresp, err := c.Grant(ctx, 1)
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.Put(ctx, leasedKey, "leasedValue", clientv3.WithLease(lresp.ID))
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(3 * time.Second)

	gresp, _ := c.Get(ctx, leasedKey)
	log.Println("number of keys ", len(gresp.Kvs))

	// election

	s1, _ := concurrency.NewSession(c)
	defer s1.Close()
	s2, _ := concurrency.NewSession(c)
	defer s2.Close()

	var ePrefix = "/my-election"
	e1 := concurrency.NewElection(s1, ePrefix)
	e2 := concurrency.NewElection(s2, ePrefix)

	go func() {
		if err := e1.Campaign(ctx, "e1"); err != nil {
			log.Fatal("campaign err:" + err.Error())
		}
		log.Println("e1 become leader")

		time.Sleep(2 * time.Second)

		e1.Resign(ctx)

		log.Println("e1 down")
	}()

	go func() {
		if err := e2.Campaign(ctx, "e2"); err != nil {
			log.Fatal("campaign err:" + err.Error())
		}

		log.Println("e2 become leader")

		e2.Resign(ctx)

		log.Println("e2 down")
	}()

}

func v2() {

}

func dumpKvs(resp *clientv3.GetResponse) {
	for i, kv := range resp.Kvs {
		log.Printf("%d, %s, %s", i, string(kv.Key), string(kv.Value))
	}
}
