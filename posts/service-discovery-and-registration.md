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

从图中可看到的是Consul支持跨数据中心，在每个数据中心有client和server。一般server是3-5个，少了会影响可靠性，多了会影响速度。client则没有数量的限制。

集群通过goosip协议管理所有node成员和广播集群消息，client需要配置文件就能发现急群众所有的server，节点故障检测被分不到整个集群，不需要server的参与。所有server通过Raft协议进行选主(leader), leader负责处理所有的查询和事务。

server还负责与其他数据中心交互来处理跨数据中心的请求，当server收到这种请求它会将请求转发到相应的数据中心活本地的leader。

#### Raft in Consul

Consul使用Raft算法实现分布式存储的一致性。Raft节点在集群中有follower,candidate或leader三种状态，初始为follower，在这个状态下节点可以接收来自leader的日志单元和投票选leader。如果一段时间内收不到来自leader的心跳，则升级为candidate想集群所有node发送投票邀请，如果candidate收到过半节点的投票则升级为leader，leader负责接收日志单元并复制给所有follower

集群当选出leader后才能接收新的日志单元，节点收到日志时把日志转发给leader，leader将日志保存到本地存储复制给所有follower，当超过半数的follower成功复制后日志单元状态变为‘committed’，进而交给状态机去执行。然后leader会将日志单元的committ状态发送给所有follower让他们来更新到各自状态机。

Consul自动的通过快照机制将Raft日志压缩来避免其无限增长

在consul集群里，只有server结点通过Raft算法进行数据一致性管理，原因是Raft集群的结点数量不能太多(在3-5)，如果client也参与到Raft那么随着集群结点数量增加，在Raft算法下集群效率会下降，client结点只是将数据请求转发给server

Raft集群的启动方式有多种，其中一种直接的办法是用个配置文件记录所有server的列表，每个server启动后用这个静态的列表作为Raft的server，但这需要所有server维护一个静态的配置文件，比如下面的一个YAML格式内容：

```YAML
servers:
	- "192.168.1.11:11011"
	- "192.168.1.12:11011"
	- "192.168.1.13:11011"
```

另外可以让server启动时自动的维护Raft的server列表，这需要避免split-brain(双主)的情况，Consul采用了下面两种方式实现：

* '-bootstrap'参数，如果集群先启动一个server并且他优先成为leader(Consul在这里做了特殊设置，让这个bootstrap server从log中恢复成leader状态)，之后加入的所有server一直是follower状态。但这还有个问题，就是还是得固定一台server成为第一个leader，启动参数必须与其他server不同(必须带有-bootstrap）
* '-bootstrap-expect'参数，所有server都用同一参数'-boostrap-expect N'说明集群的server个数为N，在有N个server join进来后，cluster开始启动raft逻辑选主。注意，N一定要与server总数相同，否则会出现split-brain问题，比如N=3 儿集群server总数为7，就很可能出现两个leader的情况。

对于读操作，consul支持多种一致性模式：

* 'default', 这种模式下如果leader因为与成员分离导致新的leader被选出的情况下，虽然旧的leader不能再提交日志，但还可以负责进行读取操作，这回导致部分读取的数据是过期的。但这种状态不会维持多长时间，因为旧的leader的状态很快就会降级为follower
* ‘consistent’，这种模式下要求leader给follower发消息确认仍然是leader，这就造成了读取操作时多了一步操作增加了延迟
* ‘stale’，这个模式下每个server都可以负责读取，这就导致读出来的数据因为还未在leader把数据复制到全部follower时被读出儿造成不一致，但这种情况也是少数，并且效率会非常快，并且即便没有leader的情况下还能够提供读操作

#### Gossip in Consul

Consul使用[gossip](https://en.wikipedia.org/wiki/Gossip_protocol)协议来管理节点和集群内发送消息，这个功能由底层的[Serf](https://www.serf.io/)库提供。集群中有LAN和WAN两个gossip pool。前者包含集群中所有的节点，用于client自动发现server、故障检测、消息群发等目的。“WAN pool”用于跨数据中心的server交互

### Consul功能

#### 注册服务(Service)

要做到“服务发现”就需要有一个存有所有“service”的目录，consul提供了HTTP API和配置文件两种机制来向目录中注册service

下面是配置文件里service的例子

```
{
	"service": {
	    "name": "redis",
	    "tags": ["v2.0"，"primary"],
	    "address": "",
	    "port": 8000,
	    "checks": [
	        {
	            "id": "redis check",
	            "script": "/usr/local/bin/check_redis.pl",
	            "interval": "5s"
	        },
	        {
	            "id": "redis check",
	            "script": "/usr/local/bin/check_redis.pl",
	            "interval": "5s"
	        }
	    ]
	  },
	}
}	
```

* name，必须填写，在同一个node下不可重名
* id，不填写默认为name
* tags，对service进行标签，可用于watch时通过tag过滤
* address，服务的地址
* port，端口
* tocken，ACL token
* checks, ‘name’默认使用'service:<service-id>'，如果重复多个则用service:<service-id>:<num>(num自增)。检查失败的时候DNS会自动将service去除，如果check里不设置默认的status，则默认为‘critial'防止service自动的被加到服务池里

#### 本地DNS

DNS是Consul提供服务查询的主要接口，程序可以通过DNS来获得自己感兴趣的服务的地址，而不用通过与consul的的HTTP API交互

默认情况下consul在127.0.0.1:8600监听DNS请求，默认不支持递归查询

Consul支持两种类型的DNS查询，"node lookup"和"service lookup", query的格式也比较特别

* node lookup，格式如下：

```
	<node>.node[.datacenter].<domain>
```

要查询'node1'这个节点可以使用node1.node.dc1.consul, 其中‘cd1’为datacenter，是FQDN(Fully Qualified Domain Name)的可选项，如果不包含则为Agent所处的datacenter。比如 ‘node1.node.consul'就是要查Agent所在datacenter的node1这个节点

node查询只会返回它的A记录：

```
[root@node4 consul.d]# dig @127.0.0.1 -p 8600  node1.node.consul. ANY

; <<>> DiG 9.9.4-RedHat-9.9.4-51.el7 <<>> @127.0.0.1 -p 8600 node1.node.consul. ANY
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 49521
;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1
;; WARNING: recursion requested but not available

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
;; QUESTION SECTION:
;node1.node.consul.		IN	ANY

;; ANSWER SECTION:
node1.node.consul.	0	IN	A	172.19.0.2

;; Query time: 1 msec
;; SERVER: 127.0.0.1#8600(127.0.0.1)
;; WHEN: Thu Sep 28 03:53:19 UTC 2017
;; MSG SIZE  rcvd: 62

```

* service lookup, 格式如下

```
[tag.]<service>.service[.datacenter].<domain>
```

比如查询"redis"这个service的结果如下:

```
[root@node4 consul.d]# dig @127.0.0.1 -p 8600 redis.service.consul

; <<>> DiG 9.9.4-RedHat-9.9.4-51.el7 <<>> @127.0.0.1 -p 8600 redis.service.consul
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 10593
;; flags: qr aa rd; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 1
;; WARNING: recursion requested but not available

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
;; QUESTION SECTION:
;redis.service.consul.		IN	A

;; ANSWER SECTION:
redis.service.consul.	0	IN	A	172.19.0.4
redis.service.consul.	0	IN	A	172.19.0.6

;; Query time: 2 msec
;; SERVER: 127.0.0.1#8600(127.0.0.1)
;; WHEN: Thu Sep 28 06:15:04 UTC 2017
;; MSG SIZE  rcvd: 81
```

当有service通不过健康检查，consul会自动将其从DNS里去除, 比如将172.19.0.6上的check结果设置为‘critial’

```
[root@node4 consul.d]# dig @127.0.0.1 -p 8600 redis.service.consul

; <<>> DiG 9.9.4-RedHat-9.9.4-51.el7 <<>> @127.0.0.1 -p 8600 redis.service.consul
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 40502
;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1
;; WARNING: recursion requested but not available

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
;; QUESTION SECTION:
;redis.service.consul.		IN	A

;; ANSWER SECTION:
redis.service.consul.	0	IN	A	172.19.0.6

;; Query time: 3 msec
;; SERVER: 127.0.0.1#8600(127.0.0.1)
;; WHEN: Thu Sep 28 06:22:38 UTC 2017
;; MSG SIZE  rcvd: 65
```

还可以通过[SRV](https://en.wikipedia.org/wiki/SRV_record)记录来查询‘service’设置的端口, 比如下面的记录里返回了个redis service的端口是‘8000’

```
[root@node4 consul.d]# dig @127.0.0.1 -p 8600 redis.service.consul SRV

; <<>> DiG 9.9.4-RedHat-9.9.4-51.el7 <<>> @127.0.0.1 -p 8600 redis.service.consul SRV
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 25116
;; flags: qr aa rd; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 3
;; WARNING: recursion requested but not available

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
;; QUESTION SECTION:
;redis.service.consul.		IN	SRV

;; ANSWER SECTION:
redis.service.consul.	0	IN	SRV	1 1 8000 node5.node.dc1.consul.
redis.service.consul.	0	IN	SRV	1 1 8000 node4.node.dc1.consul.

;; ADDITIONAL SECTION:
node5.node.dc1.consul.	0	IN	A	172.19.0.4
node4.node.dc1.consul.	0	IN	A	172.19.0.6

;; Query time: 1 msec
;; SERVER: 127.0.0.1#8600(127.0.0.1)
;; WHEN: Thu Sep 28 06:29:33 UTC 2017
;; MSG SIZE  rcvd: 142
```

* Service RFC 2782 Lookup, 格式如下

```
	_<service>._<protocol>[.service][.datacenter][.domain]
	
	consul文档里这么写，但感觉应该是下面这样
   
   _<service>._<protocol>.service[.datacenter].<domain>
```
**protocol**用service的tag，如果没有就用'tcp'


#### 健康检查(Checks)

Agent的一个重要任务是做系统级和程序级的健康检测，检测通过配置文件或HTTP接口来定义并被保存到Agent所运行的节点

检查的类型

* 定时脚本检查(script+Interval), 调用一个第三方的程序进行健康检测，使用退出码或标准输出来表明检测结果(类似Nagios), 输出内容不能超过4K，默认的检查超时时间是30秒，可以通过配置文件里的"timeout"来定义。Agent使用enable_script_checks参数来打开这种检查

  脚本退出码与检测结果对应关系：
  * 0, passing
  * 1, warning
  * other, failing

  脚本的任何输出内容会被保存到消息的‘Output’字段， Agent启动后，
* 定时HTTP检查(HTTP+Interval), 定时通过HTTP GET调用指定的URL，根绝HTTP返回码判断服务状态。‘2xx’表示通过检查，‘429’表示警告，其余均表示故障。默认情况下http请求的timeout与检查的interval相同，最长为10秒。可以通过‘timeout'来在配置文件里设置。检查返回的内容大小不能超过4K，超过部分会被截断。默认情况下支持SSL，可以通过tls_skip_verify来关掉。
* 定时TCP检查(TCP+Interval), 定时通过TCP连接检查指定host/IP和端口，根据是否能够建立连接成功与否判定service状态，成功连接表示service正常，否则表示事态危急. 默认超时时间10秒。
* TTL检查，通过服务主动的定时更新TTL，超过时间的定位service故障。
* 定时Docker检查，通过调用常驻docker里的一个检查程序来进行，这个程序通过调用Docker Exec API来启动，需要consul agent具有调用Docker HTTP API或Unix Socket的权限。consul用 DOCKER_HOST 来定位Docker API端点，检查程序运行完后要返回适当的退出码和输出，输出内容不能超过4K。Agent需要通过enable_script_check来打开这种检查

默认会将check的状态设置为‘critical’，这是防止服务在没有被检查前就被加入到调用这个服务的系统里,下面是一个check的配置项例子：
​	
```
{
	 "check": {
	    "id": "redis check",
	    "script": "/usr/local/bin/check_redis.pl",
	    "interval":"5s",
	    "status":"passing"
	    "service_id": "redis" # 检查绑定到指定的service，只会影响指定service
	}
}
```

多个检查用数组表示
​	
```
{
	"checks": [
		{
			"id": "redis check",
			"script": "/usr/local/bin/check_redis.pl",
			"interval":"5s",
			"service_id": "redis",
			"status":"critial"
		},
		{
			"id": "ssh",
			"name": "ssh port check",
			"tcp": "benx:22",
			"interval": "10s",
			"timeout": "5s"
		}
	]
}
```

下面是监控上文的‘checks’的例子，

#### 状态监视(Watch)

使用Whach可以监视KV、nodes、service、checks等对象的变化，当有更新时会触发watch定义的可执行的handler。

Watch通过阻塞的HTTP API实现，Agent会根据调用相应api请求自动监视指定内容，当有变化时通知handler

Watch还可以加入到Agent的配置文件中的watches生成，下面是Agent配置文件的内容

```
{
  "watches": [
      {
		  "type": "key",
		  "key": "foo/bar",
		  "handler": "/usr/bin/my-watch-handler.pl"
		}
   ]
}

```

watch还可以通过‘consule watch’命令来直接运行并把结果输出到处理程序

当watch监控到数据的更新，可以调用一个handler还做后续处理，watch会将监控的结果以JSON格式发送给handler的标准输入(stdin), watch不同的对象结果的格式会不同

watch的类型：

 * key
 * keprefix
 * services
 * nodes
 * service
 * checks
 * event

下面是监控key"foo/bar"的例子，用脚本/opt/consul/script/kw.pl为handler，它会将watch传给他的内容打印

保存"foo/bar"后启动watch, 他会先将kv的现有状态输出给handler

<img src="http://owo5nif4b.bkt.clouddn.com/kv1.png" width="400">

更新"foo/bar"

<img src="http://owo5nif4b.bkt.clouddn.com/kv2.png" width="400">

“foo/bar"被更新后, watch输出了kv的新的状态给handler

<img src="http://owo5nif4b.bkt.clouddn.com/kv3.png" width="400">

下面是监控上文在配置文件里设置的‘redis‘这个service的check, handler将watch的结果输出， 命令行里如果不做service限制则监控所有checks

<img src="http://owo5nif4b.bkt.clouddn.com/kv6.png" width="400">

当前状态为warning, 以为此时check_redis.pl测试脚本退出码是1

<img src="http://owo5nif4b.bkt.clouddn.com/kv5.png" width="400">

当把exit(1)改为exit(0)后， handler打印了新的check结果

<img src="http://owo5nif4b.bkt.clouddn.com/kv7.png" width="400">


参考：

https://www.consul.io














