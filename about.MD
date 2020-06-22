# 王培仲 Go开发工程师

---
- 男 29岁 
- 学历：深圳大学 本科
- 居住地：深圳市福田区&emsp;当前状态：离职
- 手机：18682126930&emsp;&emsp;邮箱：w.xy.z@live.com&emsp;&emsp;微信：wxyzyxw
- Github：github.com/peizhong

---
# 概况
- 7年软件开发工作经验，近2年从事微服务项目的开发，1年多的Go开发经验
- 熟悉Go语言编程及生态，懂得Go Runtime、调优。阅读过自带和第三方源码
- 以往项目中使用过C++、Python、C#开发。熟悉常用软件设计模式、数据结构和算法
- 熟悉微服务模式的后台开发，具备架构能力，熟悉TCP、HTTP、MySQL、Redis、etcd、RabbitMQ、gRPC。使用框架包括gin、gorm、evio、Django、ASP .NET Core
- 具有分布式解决方案的经验：Docker、Kubernetes、MongoDB、监控(Prometheus)、服务注册(Consul)、网关(Nginx、Ocelot)、身份认证(IdenetityServer)、日志(ELK)、链路追踪(Jaeger)、分布式任务(celery)。了解分布式应用协议(gossip, raft)
- Linux(Ubuntu、WSL)开发经验，了解操作系统、虚拟化和容器技术
- 较好的英语能力，可以口语交流，高效地阅读Github项目文档，查阅英文资料

---
# 工作经历
## 深圳市零成本科技股份公司 (2019年3月 ~ 2020年4月) Go工程师
### 项目：珠宝零售云平台
- 将珠宝行业Saas的营销模块后端代码从C#迁移至Go。通过WebAPI和gRPC与其他模块交互，持续功能迭代发布
- 优化框架：在HTTP、gRPC、MySQL、RabbitMQ等由框架集成的模块中增加Jaeger链路追踪，使用etcd集中管理配置
- 优化业务：pprof分析性能消耗、合理使用goroutine、调整微服务模块的划分、优化SQL、数据库读写分离、缓存、延时任务，提高系统响应速度
- 其他用到的技术和框架包括：Redis、ELK、Prometheus、gin、gorm、logrus

## 深圳市法本信息技术有限公司 (2018年8月 ~ 2019年2月) 全栈工程师
### 项目：招商银行IT资产采购系统 
- 开发前后端分离(ASP. Net Core + React)微服务架构的系统，负责采购申请、状态追踪等模块
- 优化框架：加入日志NLog、故障处理Polly、对象映射AutoMapper、数据库映射Dapper、单元测试xUnit等第三方库，提高开发效率
- 系统还包括服务注册、消息队列、缓存等模块

## 深圳市康拓普信息技术有限公司 (2016年11月 ~ 2018年7月) C#工程师
### 项目：南方电网设备资产管理软件
- C#开发Windows桌面程序，优化、添加新功能包括数据校验、批量操作等功能
- 改进程序架构，利用多线程提高软件响应速度，用orm代替sql

## 深圳市华腾生物医疗电子有限公司 (2016年3月 ~ 2016年9月) C#工程师
### 项目：便携式医疗设备软件
- C#开发Windows平板应用，实现蓝牙、USB传输，数据与服务器同步
- C#开发服务器程序，为客户端提供WCF服务，制作图形界面实现数多设备同时监控

## 深圳市科瑞康实业有限公司 (2013年6月 ~ 2015年12月) C++工程师
### 项目：病人监护仪软件
- C++开发Linux平台监护仪，负责实时数据打印、声音控制模块
- C#开发Windows桌面程序，负责数据回放、报表打印功能

---
# Github项目
github.com/peizhong
  - Letsgo: 学习Go开发。阅读任务调度、内存分配、GC等源码。使用第三方库的demo
    1. playground/gossip: gossip协议的分布式缓存，基于hashicorp/memberlist，阅读了源码了解gossip实现方式
    2. playground/net: 学习go自带的tcp库及第三方库evio，阅读源码了解go在Linux平台下epoll操作
    3. playground/mix: go与c混合编程，通过cgo在go程序调用c函数，也可通过mmap实现c和go的数据共享
  - pystuff: 学习Python开发，制作爬虫程序和展示爬虫结果的网站
    1. 爬虫：基于celery的定时任务，定时检查网络电台上的更新，下载新数据到本地、写入数据库
    2. 网站：Django后端，前端页面展示爬虫结果、控制爬虫执行时间
    3. 发布：编写Dockerfile，在云服务器上以Docker方式部署
  - NetApp: 学习 .Net Core开发微服务，包括网关、身份认证、数据库读写等服务
    1. 服务注册及网关：请求到达Ocelot网关后，查询Conusl服务注册信息，将数据转发给指定后端。阅读Ocelot源码，了解网关实现
    2. 身份认证：基于IdenetityServer4的认证服务器，为其他模块提供OpenId认证和OAuth授权功能

---
# 自我评价
- 良好的学习能力，乐于探索新技术和开源项目，应用到实际项目中
- 良好的编程习惯，有代码洁癖，注重代码质量，认同代码走查和单元测试
- 有责任心，良好的团队合作意识，也具有独立分析和解决问题的能力
- 对分布式系统、集群、云端等技术感兴趣」