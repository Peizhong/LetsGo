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
# 架构
高数据量和吞吐量的数据库应用会对单机的性能造成较大压力,大的查询量会将单机的CPU耗尽,大的数据量对单机的存储压力较大,最终会耗尽系统的内存而将压力转移到磁盘IO上。分片将数据分散在多个服务
## shard
mongod存储应用数据记录，用户的数据会分散存储到各个shard
### chunk
shard以chunk为单位存数据。当一个chunk的大小超过配置中的chunk size时，MongoDB的后台进程会把这个chunk切分成更小的chunk。通过balancer迁移chunk，从而均衡各个shard server的负载
### 分片键
数据的分片以集合为基本单位，
## api 服务器
mongos，提供对外应用访问，所有操作均通过mongos执行，从配置服务器加载元数据，提供服务，将用户的请求正确路由到对应的shard
## config server
储存Sharded cluster(Mongod)的元数据，至少三个实例组成的集群。提供元数据的增删改查，保证一致性与分区容错性
## 对象
1. 数据库：use xx
2. array集合
3. document文档：{}
   - element
4. 视图：查询结果
5. 索引
6. 存储过程：javascript编写
## 数据类型
BSON：比JSON多了Data、BinaryData、ObjectId等
ObjectId：包含时间戳_机器名_随机值，_id字段为文档默认组件
子文档：嵌套
数据：元素可以不同类型、文档
用于存储实际的数据块，shard server可由几台机器组个一个replica set
# Mongo操作
```
use xxdb
show collections
db.xxCollection.findOne()._id.getTimestamp()
db.xxCollection.find().pretty()
db.xxCollection.drop()
```
- D ordered
- E element
- M unordered
- A array
## 插入
insertOne(), insertMany()
- 如果没有数据库，默认test
- 如果集合不存在，自动创建
- 如果文档没有_id，自动创建
## 查找
findOne(), find(), limit(), skip(), sort({"value":-1})
```
{"value":50}
{"value":{$lt:50}} //lt,lte,gt,gte,ne
{"value": {$gt:100, "$lt": 200}, $or: [{"id": 1},{"name": "hello"}]}
```
聚合 aggregate
```
db.mycol.aggregate([{$group : {_id : "$by_user", num_tutorial : {$sum : 1}}}])

{
   "result" : [
      {
         "_id" : "runoob.com",
         "num_tutorial" : 2
      },
      {
         "_id" : "Neo4j",
         "num_tutorial" : 1
      }
   ],
   "ok" : 1
}
```
## 更新
operator
- $set: Sets the value of a field in a document
- $unset: Removes the specified field from a document
- $inc: Increments the value of the field by the specified amount
- $rename: Renames a field
```go
filter := bson.D{{"name", "Ash"}}
update := bson.D{
    {"$inc", bson.D{
        {"age", 1},
    }},
}
```
## 删除
deleteOne(), deleteMany()
## 创建索引
db.col.createIndex({"title":1,"description":-1})
## 关系
引用式关系
```
{
   "_id":ObjectId("52ffc33cd85242f436000001"),
   "contact": "987654321",
   "dob": "01-01-1991",
   "name": "Tom Benzamin",
   "address_ids": [
      ObjectId("52ffc4a5d85242602e000000"),
      ObjectId("52ffc4a5d85242602e000001")
   ]
}
result = db.users.findOne({"name":"Tom Benzamin"},{"address_ids":1})
addresses = db.address.find({"_id":{"$in":result["address_ids"]}})
```
## 数据库引用
dbref go不支持?