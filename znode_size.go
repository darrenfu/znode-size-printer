package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func printSize(cli *zk.Conn, node string, threshold int32) {
	exists, stat, err := cli.Exists(node)
	if err != nil {
		fmt.Printf("error retrieving znode stat: %v\n", err)
		os.Exit(1)
	}
	if !exists {
		fmt.Println("node does not exist")
		os.Exit(1)
	}
	if stat == nil {
		fmt.Println("stat is nil :(")
		os.Exit(1)
	}
	if stat.DataLength > threshold {
		fmt.Printf("%10d %s\n", stat.DataLength, node)
	}
}

func recurse(cli *zk.Conn, node string, threshold int32) {
	printSize(cli, node, threshold)
	children, _, err := cli.Children(node)
	if err != nil {
		fmt.Printf("could not get children for node %s: %v\n", node, err)
		return
	}
	for _, child := range children {
		recurse(cli, node+"/"+child, threshold)
	}
}

func main() {
	key := flag.String("dir", "/universe/mon-marathon-service/state", "path to find sizes for")
	host := flag.String("zk", "master.mesos:2181", "zk host:port")
	threshold := flag.Int("min-sz", 0, "minimum size to print")
	flag.Parse()

	cli, _, err := zk.Connect([]string{*host}, 5*time.Second)
	if err != nil {
		fmt.Printf("error connecting to zk: %v\n", err)
		os.Exit(1)
	}
	recurse(cli, *key, int32(*threshold))
}
