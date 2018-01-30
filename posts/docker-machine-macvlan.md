---
title: docker-machine 测试 macvlan 环境
date: 2018-01-22 17:56:36
tags: 
    - docker-machine
    - macvlan
---

最近研究了一下docker的几种网络驱动，其中之一macvlan模式能够让container直接使用host的物理网络，配置与host同一网段的IP，最终的效果是container就像一台和host处于同一网段的物理机

下面用docker-machine配置maclvan的实验步骤

<!-- more -->

##### 创建两个docker-machine

```shell
docker-machine create --driver virtualbox node1
docker-machine create --driver virtualbox node2
```

在virtualbox里将两个虚拟机停掉后，在网络设置里分别增加“网卡3“配置，使用“桥接网卡“，混杂模式为‘全部允许’, 之后重新启动两个虚拟机

<img src="http://owo5nif4b.bkt.clouddn.com/QQ20180122-181844@2x.png" width="400">

注意： “网卡1”和网卡2“不要动，否则docker不能正常运行

#### 在docker-machine里创建macvlane

1. 进入node1

   ~~~shell
   docker-machine ssh node1
   ~~~

2. 查看node1的ip，  我的电脑上结果如下

   ~~~shell
   docker@node1:~$ ip addr | grep inet
       inet 127.0.0.1/8 scope host lo
       inet6 ::1/128 scope host 
       inet 10.0.2.15/24 brd 10.0.2.255 scope global eth0
       inet6 fe80::a00:27ff:fe93:3624/64 scope link 
       inet 192.168.1.113/24 brd 255.255.255.255 scope global eth1
       inet6 fe80::a00:27ff:feb4:f61a/64 scope link 
       inet 192.168.99.100/24 brd 192.168.99.255 scope global eth2
       inet6 fe80::a00:27ff:fe48:4740/64 scope link 
       inet 172.17.0.1/16 brd 172.17.255.255 scope global docker0

   ~~~

   其中 ”192.168.1.113“和我的电脑是在同一网段，查看网关为 192.168.1.1 ，命令如下

   ~~~shell
   docker@node1:~$ ip r
   default via 10.0.2.2 dev eth0  metric 1 
   default via 192.168.1.1 dev eth1  metric 1 
   10.0.2.0/24 dev eth0  proto kernel  scope link  src 10.0.2.15 
   127.0.0.1 dev lo  scope link 
   172.17.0.0/16 dev docker0  proto kernel  scope link  src 172.17.0.1 
   192.168.1.0/24 dev eth1  proto kernel  scope link  src 192.168.1.113 
   192.168.99.0/24 dev eth2  proto kernel  scope link  src 192.168.99.100 
   ~~~

3. 设置node1的eth1为混杂模式

   ~~~shell
   docker@node1:~$ sudo ip link set eth1 promisc on
   ~~~

4. 创建macvlan network

   ~~~shell
   docker@node1:~$ docker network create --subnet 192.168.1.0/24 --gateway 192.168.1.1 -o parent=eth1 -d macvlan macvlan1
   ~~~

5. 运行一个container，指定网络为上面创建的’macvlan1‘，ip为 192.168.1.123

   ~~~shell
   docker@node1:~$ docker run -itd --name c1 --network macvlan1 --ip 192.168.1.123 alpine /bin/sh 
   ~~~

6. 查看 上面刚创建的container, 'c1'的ip， 为192.168.1.123， 配置成功

   ~~~shell
   docker@node1:~$ docker exec -it c1 ip a 
   1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1
       link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
       inet 127.0.0.1/8 scope host lo
          valid_lft forever preferred_lft forever
   7: eth0@if4: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1500 qdisc noqueue state UNKNOWN 
       link/ether 02:42:c0:a8:01:7b brd ff:ff:ff:ff:ff:ff
       inet 192.168.1.123/24 brd 192.168.1.255 scope global eth0
          valid_lft forever preferred_lft forever
   ~~~

7. 同样的办法在node2里创建macvlan和创建一个container ‘c2’，最终node1 c2 node2 c2 的ip如下

   ~~~
   node1: 192.168.1.113
   c1: 192.168.1.123
   node2: 192.168.1.114
   c2: 192.168.1.124
   ~~~

   另外我的电脑的ip是192.168.1.101， 

8. 两个container c1 c2 分别运行在两个node上，检查他们之间是否ping 通

   ~~~shell
   docker@node1:~$ docker exec -it c1 ping 192.168.1.124
   PING 192.168.1.124 (192.168.1.124): 56 data bytes
   64 bytes from 192.168.1.124: seq=0 ttl=64 time=5.530 ms
   64 bytes from 192.168.1.124: seq=1 ttl=64 time=0.450 ms
   ~~~

   ~~~shell
   docker@node2:~$ docker exec -it  c2 ping 192.168.1.123
   PING 192.168.1.123 (192.168.1.123): 56 data bytes
   64 bytes from 192.168.1.123: seq=0 ttl=64 time=1.309 ms
   64 bytes from 192.168.1.123: seq=1 ttl=64 time=0.458 ms
   ~~~

9. container和我的电脑是否互通

   ~~~shell
   docker@node1:~$ docker exec -it c1 ping 192.168.1.101
   PING 192.168.1.101 (192.168.1.101): 56 data bytes
   64 bytes from 192.168.1.101: seq=0 ttl=64 time=1.277 ms
   64 bytes from 192.168.1.101: seq=1 ttl=64 time=0.320 ms
   ~~~

   ~~~shell
   localhost:~ Ben$ ping 192.168.1.123
   PING 192.168.1.123 (192.168.1.123): 56 data bytes
   64 bytes from 192.168.1.123: icmp_seq=0 ttl=64 time=0.452 ms
   64 bytes from 192.168.1.123: icmp_seq=1 ttl=64 time=0.295 ms
   ~~~

   ~~~shell
   docker@node2:~$ docker exec -it  c2 ping 192.168.1.101
   PING 192.168.1.101 (192.168.1.101): 56 data bytes
   64 bytes from 192.168.1.101: seq=0 ttl=64 time=0.656 ms
   64 bytes from 192.168.1.101: seq=1 ttl=64 time=0.419 ms
   ~~~

   ~~~shell
   localhost:~ Ben$ ping 192.168.1.124
   PING 192.168.1.124 (192.168.1.124): 56 data bytes
   64 bytes from 192.168.1.124: icmp_seq=0 ttl=64 time=0.500 ms
   64 bytes from 192.168.1.124: icmp_seq=1 ttl=64 time=0.341 ms
   ~~~

10. container和除它所在host之外的其他host互通，也就是c1和node2，c2和node1

   ~~~shell
   docker@node1:~$ docker exec -it c1 ping 192.168.1.114
   PING 192.168.1.114 (192.168.1.114): 56 data bytes
   64 bytes from 192.168.1.114: seq=0 ttl=64 time=4.136 ms
   64 bytes from 192.168.1.114: seq=1 ttl=64 time=0.388 ms
   ~~~

   ~~~shell
   docker@node2:~$ docker exec -it  c2 ping 192.168.1.113
   PING 192.168.1.113 (192.168.1.113): 56 data bytes
   64 bytes from 192.168.1.113: seq=0 ttl=64 time=4.510 ms
   64 bytes from 192.168.1.113: seq=1 ttl=64 time=0.496 ms
   ~~~

   ~~~shell
   docker@node1:~$ ping 192.168.1.124
   PING 192.168.1.124 (192.168.1.124): 56 data bytes
   64 bytes from 192.168.1.124: seq=0 ttl=64 time=1.375 ms
   64 bytes from 192.168.1.124: seq=1 ttl=64 time=0.509 ms
   64 bytes from 192.168.1.124: seq=2 ttl=64 time=0.531 ms
   ~~~

   ~~~shell
   docker@node2:~$ ping 192.168.1.113
   PING 192.168.1.113 (192.168.1.113): 56 data bytes
   64 bytes from 192.168.1.113: seq=0 ttl=64 time=1.962 ms
   64 bytes from 192.168.1.113: seq=1 ttl=64 time=0.367 ms
   ~~~

11. 不做其他配置的情况下，container和其所在host是不能互通的，即node1和c1之间，node2和c2之间不能ping通，这是macvlan的特性，需要其他办法解决​


##### 总结

macvlan对于那种需要设置container ip的应用场景比较合适，比如外部程序需要直接通过ip访问container，爬虫程序部署到container里是需要设置public ip做一些ip限制的规避，邮件服务器部署到container里需要public ip做spf dkim等验证。另外它也比桥接或者swarm的overlay效率高很多，只是ip的管理比较麻烦，这也是它的一个弊端
