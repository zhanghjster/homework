---
layout: post
title: IO 系列3
date: 2018-09-24 9:30:40
tags:
  - LSM Tress
  - Sorted String Tables
---

本文是读后翻译，原文在此[On Disk IO, Part 3: LSM Trees](https://medium.com/databasss/on-disk-io-part-3-lsm-trees-8b2da218496f)

现如今数据库所使用的存储类型有很多种，它们各有优缺点，具体使用哪一种需要做出权衡。本文介绍其中一种经常被使用的类型，LSM(Log Structured Merge tree)。 

### LSM Trees

Log Structured Merge Trees开始于LSM的[开创论文](http://paperhub.s3.amazonaws.com/18e91eb4db2114a06ea614f0384f2784.pdf)，它建议实现类似于B树的磁盘驻留树，不同之处是针对顺序磁盘访问做了优化，节点(Node)可以被完全占用。虽然LSM经常被与B树做对比，但他们最明显的区别是LSM树允许不可变的可合并的文件。

尽管LSM树有各种的不同实现方式，但他们都可能一个共同之处，使用了Sorted String Tables

<!-- more -->

### Sorted String Tables

Sorted String Tables的优点是简单：它的读写和搜索都很简单。 SSTables是从键到值的持久有序不可变映射，其中键和值都是任意字节字符串。它们具有一些很好的属性，例如，随机查询（即按键查找值），顺序扫描（即迭代指定键范围内的所有键/值对）。

<img src="http://owo5nif4b.bkt.clouddn.com/1%2A8zyqPERuEvH2uugprei0DQ.png" width="400">

通常SSTable有两部分：索引和数据块。数据块由一个接一个地连接的键/值对组成，可以进行快速的顺序扫描。索引块包含主键和偏移量，它指向数据块中可以找到实际记录的偏移量。可以使用针对快速搜索优化的格式来实现主索引，例如B树。

由于SStable是不可变的，它对读取很顺序写入进行了优化，没有为修改预留空间，因此插入、更新、删除操作需要重写整个文件。而这则需要LSM树发挥它的特长。

### Anatomy

在LSM树中，所有写入都是针对可变的内存数据结构(B树或SplitList)来执行的。每当数据大小到达一定的阈值或经过一定的时间段，数据会被写入到磁盘，创建一个新的SSTable。这个过程通常称为“刷新”。检索数据可能需要在磁盘上搜索所有的SSTable，检查内存中的表并将他们的内容合并在一起，然后返回结果。

<img src="http://owo5nif4b.bkt.clouddn.com/1%2AXjd5yA7odTaHnyOzjVHwAQ.png" width="400">

在读取过程中，合并操作不可或缺，因为一个数据可以分为好几个部分，比如，插入后跟随删除或更新操作。数据的最终状态需要合并前后的操作来获取。

SSTable中的每个数据单元都有一个时间戳，不管是插入、删除还是更新都会记录当时的时间戳。

### Compaction

随着磁盘驻留表的持续增加，一个key的数据可能存在于多个文件，比如数据的多个版本，被删除操作覆盖的冗余数据。这会会降低读取操作的性能。为了避免这个问题，LSM树会有单独进程读取所有的SSTable并执行类似于检索操作的合并操作。这个过程称为“压缩”。得益于SSTable的结构，这个操作非常高效。记录以顺序方式从多个源读取，并立即可以附加到输出文件，因为所有的输入都是已排序并合并的，生成的文件将具有相同的属性。比较而已，构建一个索引文件则昂贵的多。

<img src="http://owo5nif4b.bkt.clouddn.com/1%2A9QC5iWkltzFE3wpF75GSRw.png" width="400">

### Merge

在讨论“合并”的时候有两点需要注意：复杂度和阴影逻辑。

在复杂度方面，合并SSTables和合并有序集合相同，都具有$O(M)$的内存开销，$M$为要合并的SSTables数量。每一次迭代都是从所有SSTable的当前头部拿去最小的，然后插入到新的集合中。

<img src="http://owo5nif4b.bkt.clouddn.com/1%2AOIVaABxQp11V5TTx9lqz8g.png" width="400">

检索和压缩都使用相同的压缩操作，在压缩时，顺序的SStable读取和顺序的目标SSTable保证了此过程的高效性。

阴影逻辑用来确保更新和删除正确运行：删除操作在LSM树里插入占位符来标明数据时要被删除的，更新操作则是记录数据的值并附以更高的时间戳。当读取的时候，删除的记录将不会被返回，更新的数据将返回时间戳最大的值。

<img src="http://owo5nif4b.bkt.clouddn.com/1%2AcaN0RMzQtLhJAym6Ja32rQ.png" width="400">

#### Summing Up

使用不可变数据结构通常可以简化开发工作。当使用不可变的驻盘结构时，需要偶尔合并表得到更好的空间利用率、并发性和更简单的实现。 LSM树数据库通常是写优化的，因为所有写操作都是针对预写日志和内存驻留表执行的。读操作通常较慢，因为合并过程需要检查磁盘上的多个文件。

由于压缩等维护操作，LSM-Trees可能导致更差的延迟，因为CPU和IO带宽都花费在重新读取和合并表而不仅仅是服务读写。在写操作频繁的情况下，写入和刷新来使IO饱和，此时则需要停止压缩过程。滞后压缩会导致读取速度变慢，增加CPU和IO压力，使事情变得更糟。这都需要关注与权衡。

LSM树还会引起一些写入放大：必须将数据写入预写日志，然后在磁盘上刷新，最终在压缩过程中重新读取和写入数据。不过，可变的B-Tree结构也受到写入放大的影响。

有一些数据库使用了SSTable: RocksDB和Cassandra

