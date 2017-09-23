---
title: raft
date: 2017-09-22 09:42:29
tags:
	- distributed consensus
---

对于单台服务器的系统来说，保存一条数据并保持它的完整性是非常容易的。但对于多个服务器组成的集群，就会遇到分布式一致性问题，数据在不同的节点上保存的值可能不同，这就需要算法实现保证数据一致性，Paxos和Raft便是这样的算法。

Paxos曾经是一致性算法的标杆，大多数算法都基于它或者受它的影响，但它难于理解，在实际应用中难于实现

Raft是人们经过Paxos痛苦的折磨后设计出的一个更好的一致性算法，它在容错和性能上与Paxos相当，但比Paxos更容易理解和实现

#### 一致性(consensus)

一致性是分布式系统容错要解决的一个根本性问题，它要求集群中的多数server对数据的正确与否形成共识。典型的一致性算法要求在集群中大多数(过半)服务器可用时才会认为系统可用，比如一个5台server的集群在2台故障的情况下任然会被认为可靠，可以继续提供服务，但超过2台则停止服务

[复制状态机](https://www.cs.cornell.edu/fbs/publications/SMSurvey.pdf)是一致性算法产生的背景，它是集群容错的关键要素。在一个集群中每个server都有一个状态机和一个日志。状态机就是要进行容错的组件，比如hash table。 对于client来说他正在与一个唯一的可靠的状态机进行交互，即便有部分服务器出现故障。 每个状态机从log中拿到更新数据的指令，比如hash table可能就是设置一个key的‘x'的值为3。一致性算法的目的就是要让server中log保持一致，对于任意的状态机如果第n个命令是将hash table的x设置为3，没有其他的状态机的第n个命令是将x设置为其他值，这样做就使每个状态机处理相同的命令序列，产生相同的一些列结果并达到相同的状态。简单描述就是一致性算法要每台server上的状态机拿到相同顺序的指令以达到相同的状态结果的目的, 下图是复制状态机的架构图：

![复制状态机](http://owo5nif4b.bkt.clouddn.com/statemachine.png)

#### Raft

Raft通过leader这个角色来实现一致性，集群需要选出leader，然后由它来管理日志，包括从client端获取日志，向其他server复制日志，确认日志在大多数server保存后通知server将日志同步到他们的状态机。

通过建立leader这个角色，Raft将一致性问题分解为下面三个问题:

* 选主(Leader election): 集群需要在server中选出一个leader
* 日志复制(Log replication): leader必须从client端接受日志并复制到集群中，并强制要求其他服务器的日志保持和自己相同
* 安全性(Safty): 如果一个服务器已经将给定索引位置的日志条目应用到本机状态机中，则所有其他服务器在该索引位置必须具有相同的条目

##### Basic

一个Raft集群典型情况下有5台server，系统可以容忍两个server出现故障依然保持正常工作。server的状态分为为leader、follower、cadidate三种。正常的情况下，集群只存在一个leader，其余均为follower。follower是被动的不会产生任何请求，只回应来自leader和candidate的请求。leader负责将处理来自client的所有请求（如果client是将请求发给follower，则follower需要将转发给leader)。cadidate则是leader election过程的中间状态。

server在三种状态下的转化过程如下图:

![server状态变化图](http://owo5nif4b.bkt.clouddn.com/serverstatesandrpc.png)

Raft将时间按照任期(terms)来分开，用连续性数字ID来表示。每个term从一次选举开始，如果一个cadidate赢得选举后，他将在这个term内一直扮演leader的角色。在有些时候，选举过程会出现分裂(多个candidate), 这种情况下term会以选举无效结束，等待下一term开始。Raft保证在同一时间最多只能有一个leader

时间按照terms分开后如下图：

![terms](http://owo5nif4b.bkt.clouddn.com/terms.png)

每个server保存了他当前得到的term id，当与其他server(leader或candidate)交互时，会比较其他server的term id，更新为大的值，如果一个candidate或者leader发现自己的term id比其他server的小，则自动降级为follower(server之间交互的信息都带有term id)。

server之间通过RPC协议进行交互，协议的负载只有RequestVoete和AppendEntries两种，参数如下图：
![vote](http://owo5nif4b.bkt.clouddn.com/vote.png)

![request](http://owo5nif4b.bkt.clouddn.com/request.png)

##### 选主Leader Election)

Raft 用心跳机制来触发leader election。共有election timeout和heartbeat timeout两个超时。选主过程与状态变化如下步骤:

1. 所有server的初始状态是follwer, 在此状态下有个超时计时election timeout, 一般为150ms-300ms
2. 如果在election timeout内收到来自leader的心跳，election timeout就重置, 否则server将自己的term id增加并转换为cadidate状态
3. 升级为candidate后节点就开始了选举周期，首先是先投票给自己，然后向其他节点发出投票(votes)请求
4. 其他节点收到请求后如果在新的选举周期内还没有投票，就将投票这个candidate，然后重置election timeout
5. andidate节点收到多数节点的投票后升级为leader，然后持续发送心跳消息(append entries)给follower，发送的时间间隔是 heartbeat timeout
6. follower收到心跳消息后要重置election timeout和更新转矩周期，然后发送回复信息给leader，周而复始，一直到follower不再收到heartbeat变成dadidate
7. 当leader收到新的投票周期大于自己当前周期的投票请求或者心跳消息，自动降级为follower

![选主](http://owo5nif4b.bkt.clouddn.com/leaderelection.png)

一个server在candidate状态下会遇到三种情况：
 
1. 得到大多数server的投票，成为leader
2. 收到其他server的投票请求，如果对方的term id大于自己的id，则自己降为follower，重置timeout
3. 收到其他server的投票请求，并且term id相同，则这种情况下是split vote，即同时有多个server在candidate状态进行选主。如果server最终拿到的票数不超过所有server的半数，则选举失效。timeout后term id增加，重新发出vote邀请，用随机的election timeout也尽可能的减少了这种情况的发生

#### Raft日志复制过程(Log Replication)

当一个server被选为leader后，所有来自client的请求都会被转发到leader, 之后的处理过程如下:

1. 每个请求都被封装为‘日志条目’(log entry)。日志条目存在两种状态‘uncommited’和‘commiitted’。初始状态为’未提交'(uncommited)
2. leader通过AppendEntries RPC请求发送给follower节点，当follower节点保存了日志条目后返回成功给leader
3. 当leader收到过半followr节点的成功信号后将‘日志条目’的状态更新为‘committed’,
4. leader保存日志，
5. leader再次发送消息给follower们标示‘日志条目'状态为‘commited', 这样数据在系统就保持一致了。

每个日志条目都包含当时的term id用来判断是否一致，还包含了一个id用于标示在log中的位置

总结过程要经过leader和follower的两次交互，第一次的往返交互确认日志条目被足够多的节点保存，第二次的发送确定日志条目状态为‘commited'

参考：

https://raft.github.io

https://raft.github.io/slides/uiuc2016.pdf