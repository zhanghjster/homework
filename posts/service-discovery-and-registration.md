---
title: Consul
date: 2017-09-21 09:57:29
tags:
	- consul
	- etcd
	- raft
---

对于一个分布式系统，会包含多种服务，并部署到不同的服务器上，为了实现服务之间的相互访问，需要有配置来描述不同的服务的IP和端口，在服务数量和服务器较少的情况下，手动维护配置文件就可以了。但如果成千上万，手动方式就成了不可能的任务，需要有系统来支持完成下面两个任务

* 服务注册,服务提供方将服务信息(如主机和端口)注册到存储系统
* 服务发现,服务消费方从存储系统中获得获取所需服务的信息

作为众多服务的中间媒介，通常采用的是健壮高可用的键/值存储系统

<!-- more -->

### ZooKeeper

历史悠久，数据存储格式类似文件系统，通过私有协议访问，集群式架构。优点是成熟稳定，缺点是系统复杂，资源占用高

### etcd

etcd是通过HTTP协议访问的k/v存储系统，采用集群架构，容易部署和使用。但他更多功能是提供存储，要实现服务发现还得配合一些第三方的应用或者自己实现。
​	
* “registrator”, 自动注册工具，将服务提供方的信息存储到etcd, consul这种kv存储系统
* ”confd“，轻量级的配置管理工具，他可以从etcd里取最新的服务信息生成配置文件，服务使用方就可以用它来实时更新配置文件 

### Consul

Consul 提供了高可用的kv存储，集群架构，这点和etcd zookeeper类似。 另外也提供了自动服务发现注册的套件，并且能否对服务进行健康检查。 结合consul-template可以实现服务提供方信息更新(比如增加了API服务器)时，自动生成配置文件给服务使用方自动更新配置。

##### 架构图

<img src="https://www.consul.io/assets/images/consul-arch-420ce04a.png" width="400">

首先从图中可看到的是consul是支持跨数据中心的，在每个数据中心有client和server。一般server是3-5个，少了会影响可靠性，多了会影响速度。client则没有数量的限制。

数据中心的每个node都会参与到gossip协议。这样做，一是不必为client配置server的地址，因为他们会被自动发现。二是对node的检测被分布到每个node上，不必由server来执行。三是作为gossip作为消息层来分发类似选主这种重要的事件。

每个数据中的server通过选主过程来确定leader，来负责集群的管理事务，其他server则将收到的各种RPC请求转发到leader。

server还负责与其他数据中心交互来处理跨数据中心的请求，当server收到这种请求它会将请求转发到相应的数据中心活本地的leader。

#### Raft in Consul

Consul使用Raft算法实现分布式存储的一致性，在consul集群里，只有server结点参与了Raft，所有的client结点将请求转发给server，原因是Raft集群的结点数量不能太多(在3-5)， 如果client也参与到Raft，那么随着集群结点数量增加，在Raft算法下集群效率会下降

在一开始，单结点的Consule serveri进入到“bootstrap”模式，server自动升级为leader,之后其他server可以安全的一致性的加入到成员列表里。最终当新server数量增加到指定数量时“boostrap”模式结束。

参考：

https://www.consul.io/docs/internals/consensus.html














