---
title: ELK Metricbeat体验
date: 2018-05-27 08:34:50
tags:
    - docker 
    - elk 
    - metricbeat
    - kibana
---

1. 安装docker, https://store.docker.com/editions/community/docker-ce-desktop-mac

2. 下载 elk 的 docker-compose 配置文件 

   ~~~
   git clone https://github.com/deviantony/docker-elk.git
   ~~~

3. 启动elk三个container

   ~~~
   cd docker-elk
   docker-compose up 
   ~~~

<!-- more -->

   第一次启动会很慢，因为他要先下载elk的docker镜像

   当开始输出日志后，说明启动了，日志类似下面的内容

   ~~~
   ......
   kibana_1         | {"type":"log","@timestamp":"2018-05-26T23:40:28Z","tags":["warning","elasticsearch","admin"],"pid":1,"message":"Unable to revive connection: http://elasticsearch:9200/"}
   kibana_1         | {"type":"log","@timestamp":"2018-05-26T23:40:28Z","tags":["warning","elasticsearch","admin"],"pid":1,"message":"No living connections"}
   elasticsearch_1  | [2018-05-26T23:40:28,518][INFO ][o.e.t.TransportService   ] [NDeKMVX] publish_address {172.24.0.2:9300}, bound_addresses {0.0.0.0:9300}
   elasticsearch_1  | [2018-05-26T23:40:28,719][INFO ][o.e.h.n.Netty4HttpServerTransport] [NDeKMVX] publish_address {172.24.0.2:9200}, bound_addresses {0.0.0.0:9200}
   elasticsearch_1  | [2018-05-26T23:40:28,719][INFO ][o.e.n.Node               ] [NDeKMVX] started
   elasticsearch_1  | [2018-05-26T23:40:28,784][INFO ][o.e.g.GatewayService     ] [NDeKMVX] recovered [0] indices into cluster_state
   kibana_1         | {"type":"log","@timestamp":"2018-05-26T23:40:31Z","tags":["status","plugin:elasticsearch@6.2.3","info"],"pid":1,"state":"green","message":"Status changed from red to green - Ready","prevState":"red","prevMsg":"Unable to connect to Elasticsearch at http://elasticsearch:9200."}
   logstash_1       | Sending Logstash's logs to /usr/share/logstash/logs which is now configured via log4j2.properties
   ......

   ~~~

   ​

   在另外的terminal里看一下启动是否完成,

   ~~~
   localhost:~ ben$ docker ps -a 
   CONTAINER ID        IMAGE                                       COMMAND                  CREATED             STATUS              PORTS                                                  NAMES
   780529b3fdb5        docker-elk_logstash                         "/usr/local/bin/dock…"   2 minutes ago       Up 2 minutes        5044/tcp, 0.0.0.0:5000->5000/tcp, 9600/tcp             docker-elk_logstash_1
   055ea02dfdc7        docker-elk_kibana                           "/bin/bash /usr/loca…"   2 minutes ago       Up 2 minutes        0.0.0.0:5601->5601/tcp                                 docker-elk_kibana_1
   82670cb37a90        docker-elk_elasticsearch                    "/usr/local/bin/dock…"   2 minutes ago       Up 2 minutes        0.0.0.0:9200->9200/tcp, 0.0.0.0:9300->9300/tcp         docker-elk_elasticsearch_1
   ~~~

   这样说明已经都启动了

4. 浏览器打开 http://localhost:5601 kibana

5. 下载metricbeat https://www.elastic.co/downloads/beats/metricbeat

   metricbeat是一个轻量的数据收集工具，可以收集System、Nginx、Mysql等诸多系统的统计信息，不同的系统都有各自的统计指标以及kiban的图表配置

   可以在本机运行metricbeat来监控System的CPU、Memoty、Network这些指标来体验一下kibana的每种图表的用法

   生产环境也可以用来收集监控数据，但可能业务需求不同会有特殊的指标要收集，可以通过改它的代码来定制

6. 运行metricbeat

   ~~~
   tar zxvf metricbeat-6.2.4-darwin-x86_64.tar
   cd metricbeat-6.2.4-darwin-x86_64

   ~~~

   第一次运行会慢一些，他会到es里初始化一些索引，运行后他就开始收集本机的System数据到es里了，就是前面启动的el的container

7. 配置kibana的dashboard

   ~~~
   ./metricbeat setup --dashboards
   ~~~

   在新的terminal里运行这个命令会通过kibana的接口直接创建一些dashboard

8. kibana上的discover下查看，应该有数据出现了，索引的名字是  ’metricbeat -*‘

   <img src="http://owo5nif4b.bkt.clouddn.com/QQ20180527-080607.png" width="400">

   ​

9. kibana上的visualize下查看，应该有各种dashbaord了，搜出"system"的一些dashboard可以看到本机的一些统计了

   <img src="http://owo5nif4b.bkt.clouddn.com/QQ20180527-080754.png" width="400">
