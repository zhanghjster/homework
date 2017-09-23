---
title: raft
date: 2017-09-22 09:42:29
tags:
	- distributed consensus
---

对于单节点的系统来说，保存一个数据并保持他的完整性是非常容易的。但对于多个节点来说，就会遇到分布式一致性问题，数据在不同的节点上保存的值可能不同，这就需要算法实现保证数据一致性，Paxo和Raft便是这样的算法。

共识通常出现在复制状态机的上下环境中，它是构建容错系统的一般方法。每个server都有一个状态机和一个日志。状态机就是我们要进行容错的组件，比如hash table。 对于client来说他正在与一个唯一的可靠的状态机进行交互，即便有部分服务器出现故障。 每个状态机从log中拿到更新数据的指令，比如hash table可能就是设置一个key的‘x'的值为3。一致性算法的目的就是要让server中log保持一致，对于任意的状态机如果第n个命令是将hash table的x设置为3，没有其他的状态机的第n个命令是将x设置为其他值，这样做就使每个状态机处理相同的命令序列，产生相同的一些列结果并达到相同的状态

Raft是基于Paxo的一致性算法，它在容错和性能上与Paxo相当，但它比Paxo更容易理解和实现（没有多少人能够理解Paxo，即便是Raft的作者团队前后用了一年时间才理解。Raft实现一致性算法的方法是首先选出一个唯一的leader,  然后由leader负责管理日志复制，leader负责从client哪里接收日志单元然后复制给其他server，并且告诉server是否能够把日志单元保存到状态机里。这种办法简化了日志的管理。这样Raft将一致性问题分解为三个问题 Leader election、log replication和Safty

#### Raft Basic

一个Raft集群典型情况下有5个server，这样可以容忍两个server出现故障。server的状态可以为follower、cadidate或者leader。通常的情况下，只会存在一个leader，其余均为follower。follower不产生请求，只回应来自leader和candidate的请求。leader负责将处理来自client的所有请求（如果client是将请求发给follower，则follower需要将转发给leader。cadidate则是leader election过程的中间状态。

server在三种状态下的转化过程如下图:

[server状态变化图]()

Raft将时间按照‘terms'来分开，用连续性数字ID来表示。每个term从一次选举开始，如果一个cadidate赢得选举后，他将在这个term内一直扮演leader的角色。在有些时候，选举过程会出现分裂(多个candidate), 这种情况下term会以选举无效结束，等待下一term开始。Raft保证在同一时间最多只能有一个leader。

每个server保存了他当前得到的term id,当与其他server(leader或candidate)交互时，会比较其他server的term id，更新为大的值，如果一个candidate或者leader发现自己的term id比其他server的小，则自动降级为follower(server之间交互的信息都带有term id)。

server之间通过RPC协议进行交互，协议的负载只有RequestVoete和AppendEntries两种，参数如下
	
```
RequestVote:
	Arguments {
		term   			# candidate's term
		candidateId 	# candidate requesting vote
		lastLogIndex	# index of candidate's last log entry
		lastLogTerm		# term of candidate's last lot entry
	}
	
	Result {
		term			# currentTerm, for candidate to update itself
		voteGranted		# true means candidate received vote
	}
	
	Receiver implementation:		1. Reply false if term < currentTerm 		2. If votedFor is null or candidateId, and candidate’s log is atleast as up-to-date as receiver’s log, grant vote
	
	
AppendEntries:
	Arguments {
		term 			# leader's item
		leaderId		# follower use it to redirect clients
		prevLogIndex	# index of log entry immediately preceding new ones
		preLogTerm		# term of prevLogIndex entry
		entries[]		# log entries to store(empty for heartbeat)
		leaderCommit	# leader's commitIndex
	}
	
	Result {
		term			# current term, for leader to update itself
		success			# true if follower container entry matching prevLogIndex and prevLogItem
	}
	
```

#### Leader Election

Raft 用心跳机制来出发leader election。共有election timeout和heartbeat timeout两个超时。server的状态变化如下步骤:

1. 所有server的初始状态是follwer, 在此状态下有个超时计时election timeout, 一般为150ms-300ms
2. 如果在election timeout内收到来自leader的心跳，election timeout就重置, 否则server将自己的term id增加并转换为cadidate状态
3. 升级为candidate后节点就开始了选举周期，首先是先投票给自己，然后向其他节点发出投票(votes)请求
4. 其他节点收到请求后如果在新的选举周期内还没有投票，就将投票这个candidate，然后重置election timeout
5. andidate节点收到多数节点的投票后升级为leader，然后持续发送心跳消息(append entries)给follower，发送的时间间隔是 heartbeat timeout
6. follower收到心跳消息后要重置election timeout和更新转矩周期，然后发送回复信息给leader，周而复始，一直到follower不再收到heartbeat变成dadidate
7. 当leader收到新的投票周期大于自己当前周期的投票请求或者心跳消息，自动降级为follower


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