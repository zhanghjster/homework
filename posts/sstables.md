如果Protocol Buffers是Google的单个数据记录的通用语言，那么Sorted String Table则是存储、处理、交换数据集的的常用载体之一。SSTable是对可存储大量键值对并具有高吞吐量可顺序读写的工作负载的一种抽象。

不幸的是，SSTable已经被业界赋予了超出排序表本身的意义，给这个简单有用的数据结构带来了不必要的混淆。让我们仔细研究一下SSTable的真是内涵以及LevelDB如何使用它

### SSTable: Sorted String Table

想象一下我们要处理大型的数据，输入的大小为G或T。在其上我们要运行多个步骤，由不同的程序执行。换句话说我们正在运行一系列的Map-Reduce作业。由于输入的太大，读写数据占用了大部分运行时间。因此，随机读写不是一个选择，而是应该是流式处理，一旦完成，将结果刷会磁盘，这样我们可以分摊磁盘IO成本。

<img src="http://owo5nif4b.bkt.clouddn.com/xsstable.png.pagespeed.ic.UGWVKAaGWc-2.png" width="400">

正如它的名字一样，SSTable是一个包含任意的、有序的键值对的文件。允许重复的键出现，键值之间不需要任何“填充”，他们是任意的Blob。顺序读取整个文件后会得到有序的索引。此外，如果文件非常大，我们还可以预先创建一个独立的key的偏移索引用于快速访问。这就是SSTable的全部内容，非常简单但也是交换大型有序数据非常有用的方法。

### SSTable and BigTable: Fast random access?

一旦SSTable被写入到磁盘，它就是不可更改的，因为插入和更新都会需要很大的IO来重写文件。话虽如此，但对于静态索引来说这是一个很好的解决方案：要在索引中读取，总是一次磁盘检索，或将整个文件使用mmap映射到内存，随机读取将快速而简单。

然而，随机写入则变得异常困难，除非整个表都在内存中，这样简单的指针操作就可以实现。这也是Google的BigTable要解决的问题：用SSTable支持对PB级数据集的快速读写。

### SSTables and Log Structured Merge Trees

SSTable提供了快速的读取能力，如何才能实现快速的随机写入呢。答案就是把SSTable放在内存中（我们称之为MemTable). 如下图所示。

<img src="http://owo5nif4b.bkt.clouddn.com/xmemtable-sstable.png.pagespeed.ic.VbFiXFpWsL.png" width="400">



1. 在磁盘上的SSTable的索引一直加载在内存
2. 写操作都是针对内存上的MemTable
3. 读操作先检查MemTable然后是SSTable的索引
4. 定时的MemTable会刷新到磁盘上，成为一个新的SSTable
5. 定时的磁盘上的多个SSTable进行合并

写操作由于一直是在内存上，所以会非常快。当MemTable到达一定大小，则会刷写到磁盘上称为一个不可变的SSTable。然而我们会在内存维护一套所有SSTable的索引，这样任何读操作我们都会检查MemTable，然后遍历SSTable的索引来查找数据。这是一套LSM树的实现方式，也是BigTable底层使用的逻辑。

### LSM & SSTables: Updates, Deletes and Maintenance

这一套LSM的架构提供了有趣的特性：不管数据多大，写入的速度会一直很快，而随机的读操作或者在内存中，或者需要一次磁盘查找。

一旦SSTable被写入到磁盘，他就是不可更改的。所以更新或删除操作不能够触及数据。取而代之的手段则是，如果是更新操作则在MemTable中保存最新的数据，如果是删除操作则生成一条标识删除操作的记录。因为我们是顺序的检查索引，所以读取将找到更新的活着删除逇记录，而不会达到旧的值。最后，在磁盘上有成百上千的SSTable也不是一个号的注意，需要定期的将磁盘上的SSTable合并，此过程将更新和删除记录并覆盖旧数据。

### SSTables and LevelDB

SSTable和MemTable结合一套处理约定就会得到一个很好的用于某种类型工作负载的数据库引擎。事实上，BigTable、HBase。Cassandra都在使用这种架构思想。

~~~
原文：https://www.igvita.com/2012/02/06/sstable-and-log-structured-storage-leveldb/
~~~































