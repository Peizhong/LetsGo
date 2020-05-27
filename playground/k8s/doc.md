# Docker
http://dockone.io/article/2941
- 快速部署：基于镜像
- 高效资源利用：不用硬件虚拟化，操作系统上的虚拟化，直接在cpu运行本地指令
- 迁移和扩展：不同平台运行，程序运行的环境也纳入到版本控制中
## 隔离
### Namespaces
命名空间（namespaces）是Linux提供的用于分离进程树、网络接口、挂载点以及进程间通信等资源的方法。如果不区分namespaces，每一个服务都能看到其他服务的进程，也可以访问宿主机器上的任意文件，一旦服务器上的某一个服务被入侵，那么入侵者就能够访问当前机器上的所有服务和文件。Linux 的命名空间机制提供了以下七种不同的命名空间，包括 CLONE_NEWCGROUP、CLONE_NEWIPC、CLONE_NEWNET、CLONE_NEWNS、CLONE_NEWPID、CLONE_NEWUSER 和 CLONE_NEWUTS

docker容器启动时，设置进程、用户、网络、IPC 以及 UTS 相关的命名空间
### 网络
每一个使用docker run启动的容器具有单独的网络命名空间，Docker提供了四种网络模式，Host、Container、None 和 Bridge 模式。Docker通过Linux的命名空间实现了网络的隔离，又通过iptables进行数据包转发
- Bridge：默认的网络设置模式，除了分配隔离的网络命名空间之外，Docker还会为所有的容器设置IP地址。当Docker服务器在主机上启动之后会创建新的虚拟网桥docker0，随后在该主机上启动的全部容器在默认情况下都与该网桥相连，容器的网关设为docker0，添加iptables规则。docker run -p 6379:6379，发完宿主6379都会通过iptable转发到对应容器
> iptables： NAT PREROUTING，NAT POSTROUTING

> 网桥：工作在数据链路层，查看本地的mac转发表，匹配的mac地址条目对应的端口

> TCP/IP网关：一个网络通向其他网络的IP地址，是具有路由功能的设备的IP地址
### 挂载点
CLONE_NEWNS创建隔离的挂载点命名空间。提供一个根文件系统（rootfs）
### CGroup
Control Groups，隔离宿主机器上的物理资源，例如CPU、内存、磁盘I/O和网络带宽
> 查看当前子系统 lssubsys -m
# k8s 概念
## 声明性配置
通过声明理想的状态来生产结果，yaml/json
## 协调与控制器
- 查询：获取理想的世界状态
- 观察：观察世界
- 调整：找出观察到的和理想的世界的差异，采取措施
## Workload
### Pod
集群中调度的最小单元，包含一个或多个容器，共享网络ip、hostname、存储资源、进程间通信命名空间
### ReplicaSet
副本集：无状态的pod副本，降低故障、横向扩展
### StatefulSet
Pod创建按顺序索引从低到高。Pod销毁按顺序索引从高到低。副本0可用于引导，后续副本都可认为副本0存在
### DaemonSet
在节点上运行指定pod。node selector和node label进行选择，将哪些Node包含进来。用于一些节点常驻服务
## Discovery and Load Balance
### Service
最小的网络通信单元，代表TCP或UDP负载均衡服务。定义了Pod的逻辑组，标签选择器选择pod，这些Pod提供相同的功能服务，包含3个数据
1. 自己的ip地址
2. dns的名称
3. 负载均衡规则
集群中的其他容器根据dns的名称，找到该service的负载均衡器地址，通过负载均衡规则，将流量代理到实现该service的pod上
### Ingress
基于http的负载均衡
## Deployment
对ReplicaSet模型中Pod的更新进行支持，从一个版本部署到另一个版本。策略包括：Recreate(全部销毁)，RollingUpdate(滚动升级)
### Job/Cron Jobs
Job允许Pod执行后退出，不会重启。Cron Jobs基于时间来调度和创建 Job 
## Config and Storage
### ConfigMap
配置文件、pod的容器公用、支持在线修改
### Sercet
数据库密码、证书等
### Persistent Volume Claims
Persistent Volume不直接绑定到pod中，通过Persistent Volume Claims将存储资源挂载到pod
## 标记
### Namespace
Namespace的目录存放集群中的大多数对象。提供RoleBaseAccessControl(角色访问控制)。service的dns路径会包含namespace
### Label
标记对象，通过
# k8s架构
## 多组件
可被替换
## api 驱动
通过api 交互，兼容不同组件
## master 组件
- api：服务器：对外暴露 API，其他组件都与它进行通信，也负责鉴权和授权
- etcd：保存集群数据，分布式一致性算法(raft)多个副本容灾
   - optimistic concurrency：乐观并发，compare and swap
   - watch：监视某个建值的变化
- scheduler：调度器，将pod部署到节点
- control manager：控制管理器，执行协调控制循环，一旦不满足状态，采取操作(ReplicaSets、Deployment、Service)
## node 组件
- kubulet：管理节点上pod的生命周期，同时也负责Volume（CVI）和网络（CNI）的管理。通过与api服务器通信，查找应该在其节点上运行的pod，调度和报告pod的运行状态，apis服务器接口接受到这些信息后将节点状态更新到ectd中
- kube-proxy：管理节点上的网络规则并执行连接转发
- container runtime：docker
- supervisord：保持kubelet和kube-proxy运行
- fluentd：提供集群日志
## 计划组件
- kube-dns：service的ip注册到dns，提供命名服务，方便其他组件发现
- Heapster：收集容器的cpu、网络、磁盘使用情况等指标，推送到监控系统。利用指标实现pod自动扩展
# k8s API服务器
## api管理
核心api路径：/api/v1/namespaces/<命名空间名>/<资源类型>/<资源>

api分组路径：/apis/<api分组>/<api版本>namespaces/<命名空间名>/<资源类型>/<资源>
## 请求处理
### 编码格式
content-type指定，yaml, json, protobuf
### 生命周期
1. 身份认证：客户端证书、令牌、openid
2. RBAC授权：身份具有角色权限
3. 准入控制：判断请求内容是否合法，可以对请求进行一定修改。可以通过基于webhook的准入控制动态地添加到api服务器
4. 验证：根据资源类型，自定义
### 特殊请求
流数据，开放的长期连接
- logs：日志流
- websocket：
### 监视操作
watch监视某个api路径上发生的更新，添加?wathc=true，服务器和客户端保持连接
### 乐观并法更新
拒绝较晚的并法写入
# k8s 调度器
刚创建的pod没有指定node，调度器通过watch监视api服务器上没有nodename的pod。调度器为pod选择合适的node更新nodename。node的kubectl会知道要执行pod
> 可以直接制定pod的nodename，如DaemonSet
## 调度算法
1. 谓词：硬性约束，如cpu、内存、节点的标签(nodeSelector)
2. 优先级：如已有镜像的节点优先级高
3. 亲和性：节点选择器的复杂表达式，可以非强制，作为优先级
4. 污点容忍：toleration允许谓词检查失败
## 冲突
调度时间和容器执行时间存在滞后，调度决策可能失败。节点会将pod标记失败。如果pod是由RepicaSet创建的，会创建新的安排到其他节点。一个pod也应该用RepicaSet部署
# k8s 网络
## 网络接口
- CNI规范：实现容器编排平台通过底层网络接连容器。K8s的容器运行时(Docker)调用CNI插件(Calico)向容器的网络命名空间添加或删除接口
## kube-proxy
操控每个节点的iptable，将发往某个服务的流量，重定位到后备容器端口
## 服务发现
DNS：发现Service端点的位置，kube-dns/CoreDNS组件。每个Service创建时，生成虚拟IP和DNS记录：<服务名>.<命名空间>.svc.cluster.local
## 网络策略
NetworkPolicy：命名空间、Pod的网络出口和入口规则，允许某些标签的Namespace或Pod访问制定端口等
## 服务网格
微服务环境，流量通常通过Ingress进入集群，Ingress由Service支持，Serice由任意个Pod支持，Pod可以连接到其他Service及其后备Pod，网络流向复杂。通过DaemonSet或sidecar部署到节点，将Pod的流量代理到服务网格上，让Pod成为网格一部分。服务网格提供的功能
- 流量管理：管理Service上的请求
- 可观测性：分布式追踪
- 安全性：底层网络不提供默认加密，服务网格为所有流量提供TLS
# k8s 监控
## 监控目标
可靠性和观察性
- 白盒监控：查看应用程序产生的信号
- 黑盒监控：探测监控，通过接口实施动作，如果结果与预期不一致，触发报警
## 监控栈
从集群和应用程序中获取数据
- Prometheus：拉取数据
- Fluentd：推送数据，读取本地的文件或流
## 储存监控数据
- InfluxDB：时间序列数据，成对的时间和值
- ElasticSearch：日志结构数据，非结构化或半结构化的日在文件
## 可视化
- Grafana：指标
- Kibana：日志
# 灾难恢复
## 高可用性
关键组件有高可用性，etcd多个节点
## 状态
恢复到之前定义的状态，保存在etcd
## 应用数据
持续卷

# TODO
- 普罗米修斯
- 自动扩展