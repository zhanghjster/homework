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

Raft集群的启动方式：

一种直接的办法是用个配置文件记录所有server的列表，每个server启动后用这个静态的列表作为Raft的server，但这需要所有server维护一个静态的配置文件，比如下面的一个YAML格式内容：

```YAML
servers:
	- "192.168.1.11:11011"
	- "192.168.1.12:11011"
	- "192.168.1.13:11011"
```

另外可以让server启动时自动的维护Raft的server列表，这需要避免多个split-brain(双主)的情况

'-bootstrap'参数，如果集群先启动一个server并且他优先成为leader(Consul在这里做了特殊设置，让这个bootstrap server从log中恢复成leader状态)，之后加入的所有server一直是follower状态。但这还有个问题，就是还是得固定一台server成为第一个leader，启动参数必须与其他server不同(必须带有-bootstrap'。'-bootstrap-expect'参数则避免了这个配置，所有server都用一个参数 '-boostrap-expect N'说明集群的server个数为N，在有N个server join进来后，cluster开始启动raft逻辑选主。注意，N一定要与server总数相同，否则会出现split-brain问题，比如N=3 儿集群server总数为7，就很可能出现两个leader的情况。


#### 检查(checks)

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

#### 监视(Watch)

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

```
[root@node1 script]# consul kv put foo/bar 2
[root@node1 script]# consul watch -type key -key foo/bar /opt/consul/script/kw.pl
$VAR1 = {
          'Value' => 'Mg==',
          'LockIndex' => 0,
          'CreateIndex' => 2640,
          'Session' => '',
          'Flags' => 0,
          'ModifyIndex' => 2721,
          'Key' => 'foo/bar'
        };

```

更新"foo/bar"

```
[root@node4 consul.d]# consul kv put foo/bar 6
Success! Data written to: foo/bar
```

“foo/bar"被更新后, watch输出了kv的新的状态给handler

```
[root@node1 script]# consul watch -type key -key foo/bar /opt/consul/script/kw.pl
$VAR1 = {
          'Value' => 'Mg==',
          'LockIndex' => 0,
          'CreateIndex' => 2640,
          'Session' => '',
          'Flags' => 0,
          'ModifyIndex' => 2721,
          'Key' => 'foo/bar'
        };
$VAR1 = {
          'Value' => 'Ng==',
          'LockIndex' => 0,
          'CreateIndex' => 2640,
          'Session' => '',
          'Flags' => 0,
          'ModifyIndex' => 2735,
          'Key' => 'foo/bar'
        };

```

下面是监控上文在配置文件里设置的‘redis‘这个service的check, handler将watch的结果输出， 命令行里如果不做service限制则监控所有checks

````
[root@node1 script]# consul watch -type checks -service redis /opt/consul/script/kw.pl
$VAR1 = [
          {
            'Notes' => '',
            'ServiceName' => 'redis',
            'Status' => 'warning',
            'ServiceID' => 'redis',
            'Output' => 'OK',
            'ServiceTags' => [
                               'primary'
                             ],
            'CheckID' => 'redis check',
            'Node' => 'node4',
            'Name' => ''
          }
        ];

````
当前状态为warning, 以为此时check_redis.pl测试脚本退出码是1

```
[root@node4 consul.d]# cat /usr/local/bin/check_redis.pl 
	#!/usr/bin/perl
	use 5.010;
	print 'OK';
	exit(1);      
```
当把exit(1)改为exit(0)后， handler打印了新的check结果

```
[root@node1 script]# consul watch -type checks -service redis /opt/consul/script/kw.pl
$VAR1 = [
          {
            'Notes' => '',
            'ServiceName' => 'redis',
            'Status' => 'warning',
            'ServiceID' => 'redis',
            'Output' => 'OK',
            'ServiceTags' => [
                               'primary'
                             ],
            'CheckID' => 'redis check',
            'Node' => 'node4',
            'Name' => ''
          }
        ];
$VAR1 = [
          {
            'Notes' => '',
            'ServiceName' => 'redis',
            'Status' => 'passing',
            'ServiceID' => 'redis',
            'Output' => 'OK',
            'ServiceTags' => [
                               'primary'
                             ],
            'CheckID' => 'redis check',
            'Node' => 'node4',
            'Name' => ''
          }
        ];

```

参考：

https://www.consul.io/docs/internals/consensus.html














