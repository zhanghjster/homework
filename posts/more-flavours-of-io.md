---
layout: post
title: IO 系列2
date: 2018-09-23 19:39:44
tags:
  - Memory Mapping
  - Page Cache Optimizations
  - AIO
  - Vectored IO
---

本文是读后翻译，原文在此[On Disk IO, Part 2: More Flavours of IO](https://medium.com/databasss/on-disk-io-part-2-more-flavours-of-io-c945db3edb13)

#### Memory Mapping

内存映射(Memory Mapping)允许你就像文件被完全加载到内存一样访问文件。它简化了文件访问，经常被应用于数据库技术以及应用程序。

<img src="http://owo5nif4b.bkt.clouddn.com/1%2AxvxPz9VhYoJx2AOFPbaSWw.png" width="400">

内存映射将进程虚拟页面直接映射到内核页面缓存，避免了标准IO对用户空间缓存和页面缓存之间的数据复制操作。

使用$mmap$，文件可以以私有或共享的模式映射到内存段。私有映射允许从文件读取，但写入会触发相关页面缓存的写时复制(copy-on-write)，以使原始页面保持原样并保持更改为私有，因此所有更改都不会反映到文件本身上。在共享模式下，文件映射与其他进程共享，以便他们可以查看映射的内存段的更新。此外，更改将传递到底层文件。

除非特殊指定，否则文件内容不会立即加载到内存，而是以惰性方式加载。内存映射所需的空间是保留的，但不会立即分配。第一次读写操作会导致页面错误从而触发相应页面的内存分配。从国MAP_POPULATE可以预先映射区域并强制文件预读取。

内存映射通过页面缓存完成，与标准IO相同，使用按需分页(Damand paging)方法管理这些缓存页面。当第一次访问内存，会发出一个缺页中断(Page Fault)，它会向内核发出请求页面未加载到内存的信号。内核识别出要加载的数据的位置。缺页中断对开发人员来说是透明的，程序会继续进行，就像什么也没发生，不过有时它可能会影响程序性能

也可以使用保护标志将文件映射到内存，比如只读模式。如果一个对映射内存的操作违反了条件，则会触发段错误。

$mmap$是一个非常有用的IO工具：它避免了内存创建缓冲区无关的副本，页避免了触发实际IO操作的系统调用开销。从开发人员角度，对使用$mmapp$的文件发出随机去读看起来就像普通的指针操作，不涉及到lseek调用。

$mmap$存在一些缺点，但对于现代硬件来说不太重要了：

* mmap是内核使用更多的数据结果来管理内存映射
* 内存映射的文件大小限制：大多数时间没和代码应该内存友好。不过对于64位的系统，已经可以映射更大的文件

当然，这并不意味着所有操作都由内存映射文件来完成。它经常被数据库实现者使用。例如，MongoDB的默认存储引擎是mmap支持的，SQLit也广泛使用内存映射

<!-- more -->

#### Page Cache Optimizations

从到目前我们讨论的内容来看，使用标准IO开起来简化了很多东西，同时牺牲了可控性。虽然在内核可以使用内部统计信息更好的预测执行回写和预取页面的时机，但有时候也有必要让内核以有益于程序的方式管理内存。

一种向内核通知你的意图的方式是使用$fadvise$, 使用以下标志，可以向内核指示你的意图，并让他优化页面缓存的使用：

* ADV_SEQUENTIAL，指定按照偏移量从高到低的方式读取文件，因此内核可以确保在读取发生之前提前获取页面
* *FADV_RANDOM*，禁用预读，页面缓存中清除不会很快被访问的页面
* FADV_WILLNEED，通知OS页面将会很快被该进程使用，使得内核有机会提前缓存页面，并且当读取操作发生时，从页面缓存而不是缺页中断提供它
* FADV_DONTNEED， 建议内核可以释放相应的页面

如同$advise$的字面意义，它只是‘建议’，并不意味着内核会按照啊建议去执行。

由于数据库开发人员通常可以预测访问，因此fadvise是一种非常有用的工具。例如，RocksDB使用它来通知内核有关访问模式的信息，具体取决于文件类型（SSTable或HintFile），模式（随机或顺序）和操作（写入或压缩）

另一个有用的调用时$mlock$，它允许强制页面保存在内存中。这意味着，一旦页面内加载到内存，后续操作都是从页面缓存中提供，必须谨慎使用，因为它会耗尽系统资源。

#### AIO

异步IO(Asynchronous IO), 它是一个接口，通过它可以初始化多个IO操作，并注册操作完成时的回调操作。使用异步IO可以帮助程序在IO操作进行时继续运行主线程的任务。

#### Vectored IO

一种可能不太流行的执行IO操作的方法是向量IO（也称为Scatter / Gather）。它在缓冲区向量上运行，并允许每个系统调用使用多个缓冲区从磁盘读取和写入数据。

执行向量读取时，首先将字节从源读取到缓冲区。然后，从第一个缓冲区的长度开始直到第二个缓冲区的长度偏移的源的字节将被读入第二个缓冲区，依此类推，就像源一个接一个地填充缓冲区一样。向量写入以类似的方式工作：缓冲区将被写入，就好像它们在写入之前被连接一样。

这种方法可以通过允许读取较小的块，因此避免为连续块分配大的存储区域，同时减少用来自磁盘的数据填充所有这些缓冲区所需的系统调用量。另一个优点是读取和写入都是原子的：内核阻止其他进程在读取和写入操作期间对同一描述符执行IO，从而保证数据的完整性。

从开发人员的角度来看，如果数据在文件中以某种方式布局。例如，它被拆分为固定大小的头和多个固定大小的块，则可以发出一个单独的调用来填充单独的缓冲区分配给这些部分。

这听起来很有用，但不知何故，只有少数数据库使用向量IO。这可能是因为通用数据库同时处理大量文件，试图保证每个正在运行的操作的活跃性并减少它们的延迟，因此可以按块访问和缓存数据。向量IO对于分析工作负载和/或列式数据库更有用，其中数据连续存储在磁盘上，并且其处理可以在稀疏块中并行完成。其中一个例子是Apache Arrow。

#### Summary

有很多东西可供选择，每一个都有自己的优点和缺点。使用特定工具并不能保证积极的结果，因为它们的具体细节而导致这些IO风格很容易被误解和误用





加载和写入。

### Summary

IO不管是是否使用缓存，了解原理对实际操作中避免误用起到很大的作用。
