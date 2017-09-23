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
	
* “registrator”, 自动注册工具，将服务提供方的信息存储到etcd, consul这种kv存储系统
* ”confd“，轻量级的配置管理工具，他可以从etcd里取最新的服务信息生成配置文件，服务使用方就可以用它来实时更新配置文件 

### Consul

Consul 提供了高可用的kv存储，集群架构，这点和etcd zookeeper类似。 另外也提供了自动服务发现注册的套件，并且能否对服务进行健康检查。 结合consul-template可以实现服务提供方信息更新(比如增加了API服务器)时，自动生成配置文件给服务使用方自动更新配置。

#### 概念

1. Agent, 在集群的每个node上运行的后台程序。运行’consul agent‘启动，运行在‘client’或者‘server’模式下。所有的agent可以运行DNS或HTTP接口并且负责检查和同步服务。
2. Client, 运行在client模式下的Agent，负责将所有的RPC请求转发到Server。Client是相对无状态的，他唯一的后台活动是LAN gossip池，只占用有限的资源。
3. Server，运行在server模式下的client。任务包括维护集群状态、处理RPC请求、与其他datacenter交换WAN信息、转发查询到leader或远程数据中心等。
4. Consensus(共识), 通过选主达成的协议。
5. Gossip，node之间通过UDP交流所使用的协议
6. LAN Gossip，一个数据中心里同一个本地网络里的node组成的池
7. WAN Gossip，数据中心与外界交互的node。


#### 架构图

![图片](https://www.consul.io/assets/images/consul-arch-420ce04a.png)

首先从图中可看到的是consul是支持跨数据中心的，在每个数据中心有client和server。一般server是3-5个，少了会影响可靠性，多了会影响速度。client则没有数量的限制。

数据中心的每个node都会参与到gossip协议。这样做，一是不必为client配置server的地址，因为他们会被自动发现。二是对node的检测被分布到每个node上，不必由server来执行。三是作为gossip作为消息层来分发类似选主这种重要的事件。

每个数据中的server通过选主过程来确定leader，来负责集群的管理事务，其他server则将收到的各种RPC请求转发到leader。

server还负责与其他数据中心交互来处理跨数据中心的请求，当server收到这种请求它会将请求转发到相应的数据中心活本地的leader。

#### Consensus Protocol（共识协议）

Consul使用[共识协议](https://en.wikipedia.org/wiki/Consensus_(computer_science))实现[一致性](https://en.wikipedia.org/wiki/CAP_theorem)(Consistency)，这个协议是基于Raft算法实现，[这里](http://thesecretlivesofdata.com/raft/)有Raft的动画演示。

##### Raft算法

Raft是基于[Paxos](https://en.wikipedia.org/wiki/Paxos_%28computer_science%29)的一个共识算法。与Paxos相比，Raft具有更少的状态以及简单易懂的算法。

这里有一些在讨论Raft时候用到的关键词：

 * Log - Raft系统的主要工作单元是日志条目。一致性问题被分解成日志条目的复制问题。日志是条目的有序序列，如果所有的成员对日志条目内容和顺序都认可，则我们认为日志是一致的
 * FSM - [有限状态机](https://en.wikipedia.org/wiki/Finite-state_machine)(Finite State Machine). 一个FSM是他们之间具有过度的有限状态所组成的集合。当心的日志被应用时，FSM允许在状态件过度。具有相同日志顺序的应用必须具有相同的状态，即，行为必须具有确定性。
 * Peer set - 参与日志复制的所有成员的集合。对于consul，所有的node都在本地数据中心的peer set里。
 * Quorum - peer set里的一些主要成员组成的集合。如果一个peer set的大小是你n，则quorum至少要(n/2)+1. 如果所有node的quorum是无效的，则集群变成无效的，不会再有新的log被提交
 * Committed Entry - 当一个条目被持久存储到结点的Quorom nodes里时，条目被认为被提交了，进而才可以被应用
 * Leader - 在任何时间，peer set会选出一个node作为leader，负责提取新的log条目，复制条目到‘follower’。

 
 。。。。。。

参考：

https://www.consul.io/docs/internals/consensus.html














