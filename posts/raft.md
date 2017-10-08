---
title: Raft
date: 2017-09-22 09:42:29
tags:
    - distributed consensus
    - 分布式一致
categories:
    - 算法
---

对于单台服务器的系统来说，保存一条数据并保持它的完整性是非常容易的。但对于多个服务器组成的集群，就会遇到分布式一致性问题，数据在不同的节点上保存的值可能不同，这就需要算法实现保证数据一致性，Paxos和Raft便是这样的算法

Paxos曾经是一致性算法的标杆，大多数算法都基于它或者受它的影响，但它难于理解，在实际应用中难于实现

Raft是人们经过Paxos痛苦的折磨后设计出的一个更好的一致性算法，它在容错和性能上与Paxos相当，但比Paxos更容易理解和实现

<!-- more -->

#### 一致性(consensus)

一致性是分布式系统容错要解决的一个根本性问题，它要求集群中的多数server对数据的正确与否形成共识。典型的一致性算法要求在集群中大多数(过半)服务器可用时才会认为系统可用，比如一个5台server的集群在2台故障的情况下任然会被认为可靠，可以继续提供服务，但超过2台则停止服务

[复制状态机](https://www.cs.cornell.edu/fbs/publications/SMSurvey.pdf)是一致性算法产生的背景，它是集群容错的关键要素。在一个集群中每个server都有状态机和日志。状态机就是要进行容错的组件，比如hash table。 对于client来说他正在与一个唯一的可靠的状态机进行交互，即便有部分服务器出现故障。 每个状态机从log中拿到更新数据的指令，比如hash table可能就是设置一个key的‘x'的值为3。一致性算法的目的就是要让server中log保持一致，对于任意的状态机如果第n个命令是将hash table的x设置为3，没有其他的状态机的第n个命令是将x设置为其他值，这样做就使每个状态机处理相同的命令序列，产生相同的一些列结果并达到相同的状态。简单描述就是一致性算法要每台server上的状态机拿到相同顺序的指令以达到相同的状态结果的目的, 下图是复制状态机的架构图：

<img src="http://owo5nif4b.bkt.clouddn.com/smc.png" width="400">

#### Raft

Raft通过leader这个角色来实现一致性，集群需要选出leader，然后由它来管理日志，包括从client端获取日志，向其他server复制日志，确认日志在大多数server保存后通知server将日志同步到他们的状态机。

通过建立leader这个角色，Raft将一致性问题分解为下面三个问题:

* 选主(Leader election): 集群需要在server中选出一个leader
* 日志复制(Log replication): leader必须从client端接受日志并复制到集群中，并强制要求其他服务器的日志保持和自己相同
* 安全性(Safty): 如果一个服务器已经将给定索引位置的日志条目应用到本机状态机中，则所有其他服务器在该索引位置必须具有相同的条目

#### 基础(Basic)

一个Raft集群典型情况下有5台server，系统可以容忍两个server出现故障依然保持正常工作。server的状态分为为leader、follower、cadidate三种。正常的情况下，集群只存在一个leader，其余均为follower。follower是被动的不会产生任何请求，只回应来自leader和candidate的请求。leader负责将处理来自client的所有请求（如果client是将请求发给follower，则follower需要将转发给leader)。cadidate则是leader election过程的中间状态。

server在三种状态下的转化过程如下图:

<img src="http://owo5nif4b.bkt.clouddn.com/littlestatus.png" width="400" >

Raft将时间按照任期(terms)来分开，用连续性数字ID来表示。每个term从一次选举开始，当一个cadidate赢得选举后，它将在这个term内一直扮演leader的角色。在有些时候，选举过程会出现分裂(多个candidate), 这种情况下term会以选举无效结束，等待下一term开始。Raft保证在同一时间最多只能有一个leader

时间按照terms分开后如下图：

<img src="http://owo5nif4b.bkt.clouddn.com/terms1.png" width="400">

每个server保存了他当前得到的term id，当与其他server(leader或candidate)交互时，会比较其他server的term id，更新为大的值，如果一个candidate或者leader发现自己的term id比其他server的小，则自动降级为follower(server之间交互的信息都带有term id)。

server之间通过RPC协议进行交互，协议的负载只有RequestVoete和AppendEntries两种，

##### RequestVote RPC:
* Arguments:

|    Field     |             Description             |
| :----------: | :---------------------------------: |
|     term     |          candidates' term           |
| candidateId  |      dandidate requesting vote      |
| lastLogIndex | index of candidate's last log entry |
| lastLogTerm  |  term of dadidate's last log entry  |

* Results:

|    Field    |               Description                |
| :---------: | :--------------------------------------: |
|    term     | currentTerm, for candidate to update itself |
| voteGranted |    true means candidate received vote    |

* Receiver implementation:

  1. Reply false if term < currentTerm
  2. If votedFor is null or cadidatedId, and candidat's log is at least as up-to-date as receiver's log, grant the vote

##### AppendEntries RPC
* Arguments

|    Field     |               Description                |
| :----------: | :--------------------------------------: |
|     term     |              leader's term               |
|   leaderId   |     so follower can redirect clients     |
| prevLogIndex | index of entry immediatey preceding new ones |
|  preLogTerm  |        term of prevLogIndex entry        |
|  entries[]   | log entries to store(empty for heartbeat;may send more than one for efficiency |
| leaderCommit |           leader's commitIndex           |

* Results

|  Field  |               Description                |
| :-----: | :--------------------------------------: |
|  term   | currentTerm, for leader to update itself |
| success | true if follower container entry mating prevLogIndex and prevLogTerm |

* Reciver implementation
  1. Reply false if term < currentTerm
  2. Reply false if log dosen's contain and emtry at prevLogIndex whose term matches prevLogTerm
  3. If an existing entry conflicts with a new one(same index but different terms), delete the existing entry and all that follow it
  4. Append any new entries not already in the log 
  5. If leaderCommit > commitIndex, set commitIndex = min(leaderCommit, index of last new entry)

#### 选主(Leader Election)

Raft 用心跳机制来触发leader election, 共有election timeout和heartbeat timeout两个超时。选主过程与状态变化如下步骤:

1. 所有server的初始状态是follwer, 在此状态下有个超时计时election timeout, 一般为150ms-300ms
2. 如果在election timeout内收到来自leader的心跳，election timeout就重置, 否则server将自己的term id增加并转换为cadidate状态
3. 升级为candidate后节点就开始了新的term，首先是先投票给自己，然后向其他节点发出投票(votes)请求
4. 其他节点收到请求后如果在新的term内还没有投票，就将投票这个candidate，然后重置election timeout
5. andidate节点收到多数节点的投票后升级为leader，然后持续发送心跳消息(append entries)给follower，发送的时间间隔是 heartbeat timeout
6. follower收到心跳消息后要重置election timeout和更新term，并回复信息给leader，周而复始，一直到follower不再收到heartbeat变成dadidate
7. 当leader收到新的term大于自己当前term的投票请求或者心跳消息，自动降级为follower

<img src="http://owo5nif4b.bkt.clouddn.com/leaderelection1.png" width="400">

一个server在candidate状态会遇到三种情况：

1. 得到大多数server的投票，成为leader
2. 收到其他server的投票请求，如果对方的term id大于自己的id，则自己降为follower，重置timeout
3. 收到其他server的投票请求，并且term id相同，则这种情况下是split vote，即同时有多个server在candidate状态进行选主。如果server最终拿到的票数不超过所有server的半数，则选举失效。timeout后term id增加，重新发出vote邀请，用随机的election timeout也尽可能的减少了这种情况的发生

#### 日志复制(Log Replication)

集群中每个server都有一个由日志条目目组成的日志。日志条目是由leader在收到来自client的请求后将请求包装而成，内容包含“复制状态机”要执行的指令和当时leader的term。日志条目还存在”提交“(commited)和“未提交”(uncommited)两种状态，初始状态为“未提交”。 

下图是在某一时刻集群每个server日志的结构:

<img src="http://owo5nif4b.bkt.clouddn.com/logs.png" width="400" title="日志">

下面是server将日志复制到集群的过程：

1. leader将放到它的日志中，然后并发的通过AppendEntries RPC请求发送给follower要求他去复制
2. follower节点保存了日志条目后返回"成功给"leader
3. leader收到过半follower节点的"成功"信号后将‘日志条目’的状态更新为‘committed’
4. leader将日志条目里的命令交给它的状态机执行并将执行结果交给client
5. leader跟踪它所知道的最大的日志条目的索引，并在之后的所有AppendEntries RPC均带有这个值(leaderCommit)，目的是让follower都知道这条条目已经提交，它门将条目应用到本地状态机
6. follower出现崩溃或者运行缓慢或者网络丢包的情况出现，通过leader无限重试AppendEntries RPC直到所有follower最终存储了所有日志条目

Raft如何保证不同服务器上的日志的一致性呢，它通过保证以下特性来实现：

* 如果在不同的日志中两个日志条目有相同的索引和term，他们所存储的的命令相同
* 如果在不同的日志中两个日志条目具相同的索引和term，他们之前的所有条目都完全相同

第一条特性通过同一leader在一个任期内在给定的日志索引位置最多创建一条日志条目，并且日志的索引不会改变。第二条特性通过AppendEntries的一致性检查实现，Append RPC的参数里包含leader最新条目之前的条目索引位置(preLogIndex)和term(preLogTerm)，如果follower没有在日志中找到相同索引和term，则拒绝新的日志条目，这种方式确保了AppendEntries返回成功时leader就知道follower和它的日志是一致的

leader通过强制follower来复制他们的日志来处理他们之间的不一致, leader为每个follower维护一个nextIndex，它表示leader将要发送给follower的下一条日志条目的索引，当一个server成为leader时，它会将nextIndex设置为它的最新日志条目索引数+1，如果follower日志和leader的不同，会在AppendEntries一致性检查时返回失败，leader将nextIndex递减，然后重试AppendEntries RPC直至返回成功，follower上的冲突日志都移除，然后leader再从最后计算到的nextIndex开始逐步把日志同步给follower直至一致。

#### 安全性(Safty)

强制leader和follower之间的日志的复制还不能完全满足一致性，因为当出现leader的更替，如果新的leader是刚从故障恢复的follower里选上来的，那么它和上任leader之间可能会缺失很多日志条目，这样再强制follower和leader日志同步，他的follower从上任leader那里复制过来的日志会有一些被覆盖，因为日志条目录的index被新的leader同步到了之前的某个值

这要求leader必须满足“存储全部已经提交的日志条目”这个条件，即Leader Completeness原则(一个已提交的日志条目，一定会出现在term更大的leader的日志中)，Raft通过选主环节的一些限制条件来实现这个目的。这个限制是，如果candidate的日志和大多数服务器上的一样新，那么他一定包含所有已提交的日志条目。因为如果一个日志条目是被提交的，他一定在大多数服务器上是保存的了。实现方法是：在RequestVote RPC时候，投票人如果发现自己的日志比候选人的日志新(比较最后日志条目的索引和term)，则拒绝投票请求

如果follower活着candidate出现故障，leader在发送的RequestVote和AppendEntries RPC会失败，在恢复后，会再次受到leader发来的同一个RPC，直至日志与leader一致

参考：

https://raft.github.io

https://raft.github.io/slides/uiuc2016.pdf

http://www.infoq.com/cn/articles/raft-paper