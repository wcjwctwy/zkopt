package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"os"
	"strings"
	"time"
)

//将老集群的数据复制到新集群上
func CopyZkData(oldzk *zk.Conn, path string, newzk *zk.Conn) {
	children, _, e := oldzk.Children(path)
	bytes, _, _ := oldzk.Get(path)
	fmt.Printf("|%-120s|%s\n", path, bytes)
	if !strings.Contains(path, "/sysconfig/bigdata/system/status") {
		//flags有4种取值：
		// 0:永久，除非手动删除
		// zk.FlagEphemeral = 1:短暂，session断开则改节点也被删除
		// zk.FlagSequence  = 2:会自动在节点后面添加序号
		// 3:Ephemeral和Sequence，即，短暂且自动添加序号
		var acls = zk.WorldACL(zk.PermAll) //控制访问权限模式
		create, err := newzk.Create(path, bytes, 0, acls)
		fmt.Println(create)
		if err != nil {
			return
		}
	}
	if e != nil || len(children) == 0 {
		return
	}
	for _, cPath := range children {
		CopyZkData(oldzk, path+"/"+cPath, newzk)
	}
}

//读取节点数据
func QueryNode(zk *zk.Conn, path string) {
	children, _, e := zk.Children(path)
	bytes, _, _ := zk.Get(path)
	fmt.Printf("|%-120s|%s\n", path, bytes)

	if e != nil || len(children) == 0 {
		return
	}
	for _, cPath := range children {
		QueryNode(zk, path+"/"+cPath)
	}
}

//读取节点数据
func Delete(zk *zk.Conn, path string) {
	children, _, e := zk.Children(path)
	bytes, _, _ := zk.Get(path)
	fmt.Printf("|%-120s|%s\n", path, bytes)

	if e != nil || len(children) == 0 {
		return
	}
	for _, cPath := range children {
		QueryNode(zk, path+"/"+cPath)
	}
}

func main() {
	fmt.Println(os.Args)
	if len(os.Args) < 4 {

		return
	}
	c, _, err := zk.Connect([]string{os.Args[1]}, time.Second) //*10)
	if err != nil {
		panic(err)
	}
	c1, _, err := zk.Connect([]string{os.Args[3]}, time.Second) //*10)
	if err != nil {
		panic(err)
	}
	CopyZkData(c, os.Args[2], c1)

	QueryNode(c1, os.Args[2])
}
