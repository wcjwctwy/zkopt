## zookeeper数据迁移操作
```
func CopyZkData(oldzk *zk.Conn,path string,newzk *zk.Conn)
```
oldzk 老zk的链接
path 迁移复制的目录
newzk 新zk的链接