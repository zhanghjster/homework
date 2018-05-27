---
title: kafka connect
date: 2018-05-20 12:59:17
tags:
    - kafka
    - kafka-connect
    - CDC
---

在一条数据的生命周期中，它通常要经历提取、存储、加载三套系统。数据提取系统接收数据并将其保存到存储系统，数据加载系统从则从存储系统中的加载数据进行业务处理。而在众多场景中，当数据提取系统接收并保存新的数据时，往往要加载系统做出实时的更新。

比如，一个评论系统，当保存评论内容后会要求缓存、搜索引擎、垃圾检测等使用到评论的系统做相应的更新

<!-- more -->

如果是简单的单体式应用，评论系统会在保存完评论后调用更新缓存、重建搜索引擎索引、垃圾检测的代码，显然这种强耦合办法日积月累的话会让项目臃肿得一团乱麻。

那么改进一下，把缓存、搜索、垃圾检测拆分出来做成单独服务，评论系统通过rpc或api调用方式触发他们的逻辑来进行解耦。但看上去还不够优雅，因为评论系统还是要维护一个要调用的外界服务的列表。

再改进一下，加入数据总线，评论系统将新的评论放到数据库后发布一条消息到总线，所有应用到品论数据的系统均监听总线上的数据然后进行自己的业务逻辑处理。这样看，系统间的耦合已经很低了，评论系统只需要维护一个发送消息的逻辑就可以了。

那么，有没有办法评论系统连消息都不用发呢？答案是看情况。

如果使用kafka作为消息总线，这个则是有可能的。kafka提供了一套框架用于从数据源中拉取数据并转化为kafka数据流的框架 [kafka connect](https://www.confluent.io/product/connectors/),  它支持诸如MySQL、ActiveMQ、ES、MongoDB等等多种数据源，并且可以根据需求定制connector。

对于前面的评论系统，如果将数据保存到MySQL中。使用kafka connect监听评论数据库的更新操作(INSERT、UPDATE、DELETE), 将这些操作转化成一条消息发布到数据总线上，所有使用到评论的系统监听到这些操作后所处相应的实时响应。这样一来，评论系统就只与数据库打交道，逻辑更加简单了

Kafka官方的MySQL Connector只能监听到数据的INSERT和UPDATE操作，如果要将监听到DELETE，需要采取一些技巧，方法是DELETE操作不是直接从数据库中删除，而是通过UPDATE一个字段为‘deleted’状态来触发connector拉取数据。之所以这样做，是因为这个connector使用的JDBC方式拉取数据，它是监听不到一个被真正delete掉的数据的。

一个更好的办法是采用CDC(Change Data Capture)技术。[Debezium](http://debezium.io) 就是其中之一，它监听的是MySQL的binlog，可以对数据库的DELETE，包括数据库的表的操作都能够获取到，速度上比JDBC的更快。

它的[工作原理](http://debezium.io/docs/connectors/mysql/), 大概分为如下两步骤：

1. 同步镜像数据

   通过[一致性读取](https://dev.mysql.com/doc/refman/5.6/en/innodb-consistent-read.html)的方式读取到当前时刻MySQL的快照数据，包括binlog的位置、数据库表的结构、所有表里的所有数据

2. 监听binlog

   同步完镜像数据后会从上一步的binlog位置开始读取binlog数据
