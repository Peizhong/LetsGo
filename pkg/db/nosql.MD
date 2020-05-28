# todo
- MapReduce
# NoSQL 
数据之间关联较少，不用事先定义数据字段，容易分散储存、实现并发处理、横向扩展
- 非结构化数据：图像、声音
- 半结构化数据：json、xml
## BASE特性
- BA：基本可用
- S：可伸缩
- E：最终一致
## NoSQL数据库类型
- 文档型(mongodb)：实体的所有数据在一个文档中，文档存在集合中
- 键值对(redis)：key和value，针对key查找时快
- 列存储
- 图存储
# MongoDB概念
偏向于大数据规模下的CRUD，不支持事务
半结构化的非关系型数据库，分布式的存储架构
## 特性
1. 存储结构：BSON文档(二进制格式)、集合储存文档、原子性修改文档
2. 数据查询：单表查询、聚合查询、MapReduce。支持建立索引(复合、文本、地理空间)
3. 数据库架构：副本集，主张转移。数据块自动切分，横向扩展，数据访问负载均衡
## 对象
1. 数据库：use xx
2. 集合
3. 文档：{}
4. 视图：查询结果
5. 索引
6. 存储过程：javascript编写
## 数据类型
BSON：比JSON多了Data、BinaryData、ObjectId等
ObjectId：包含时间戳_机器名_随机值，_id字段为文档默认组件
子文档：嵌套
数据：元素可以不同类型、文档
# Mongo集群
mongos(路由)、config server、shard
## 副本集
## 分片集
# Mongo操作
```
use xxdb
show collections
db.xxCollection.findOne()._id.getTimestamp()
db.xxCollection.find().pretty()
```
## 插入
- 如果没有数据库，默认test
- 如果集合不存在，自动创建
- 如果文档没有_id，自动创建

