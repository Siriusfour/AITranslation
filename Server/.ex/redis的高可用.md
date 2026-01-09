主从复制：

一个master节点，一个slave节点，主从复制，可能会有数据丢失，切换需要手动

哨兵模式：

一个master节点，多个slave节点 多个sentinel节点
正常模式下master节点负责正常的读和写，并把自身的数据复制到多个slave节点备份（第一次全量，后面复制增量），每个sentinel节点都会向master节点确认心跳（轮询）
如果有超过半数的sentinel确认当前的master已经挂了，那他们就会一raft共识算法投票在多个slave节点选出一个新的master
Q1：主从切换期间大概会有10s的不可读，怎么处理？
A1：根据请求的种类区分不同的处理方式，

分布式集群：

数据结构：
string：存储单个kv对， SDS,就是在一个结构体显式维护容量cap，已用长度len，数组本体
hash：存储对象的信息，例如一个用户的所有信息，
set：唯一集合，
zset：排行榜，延时任务
list ： 
Bitmap：签到统计
GEO：

