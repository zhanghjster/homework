---
layout: post
title: Kafka Consumer基本概念和参数
date: 2018-11-24 11:01:01
tags: 
   - kafka 
   - comsumer
---

kafka的每个topic都被分成一组日志，每个日志被称为一个分区（parition）。生产者向日志尾部写入，而消费者则按照自己的节奏读取日志。kafka通过在一个消费组内的消费者之间分配分区来实现扩展topic的消费，下图显示了一个具有三个分区的topic和有两个消费者的消费组。

<img src="https://www.confluent.io/wp-content/uploads/2016/08/New_Consumer_figure_1.png" width="600">

从版本0.9开始，kafka客户端使用组协调协议进行组管理。对于每个消费组，都会选择一个broker作为协调器来管理组成员的状态。当有旧成员的离开、新成员的加入以及主题元数据更改时调节分区在组成员间的分配，这个过程也称为再平衡(Rebalance)。

<!-- more -->

当一个组首次初始化后，消费者通常会从每个分区中的最早或最新的偏移开始，按照顺序读取日志消息。当消费者处理完已经读的消息后，它需要提交这些消息的偏移量。下图展示了一个消费者当前消费到了偏移量6，它最新提交的偏移量在1。

<img src="https://www.confluent.io/wp-content/uploads/2016/08/New_Consumer_Figure_2.png" width="600">

当一个分区被重新分配给组内的其它消费者时，初始的消费位置会被设置到最后提交的偏移量。如果上图中的消费者突然崩溃，那么接管这个分区的消费者则会从偏移量1开始消费。在这种情况下，它必须重新处理1之后的的数据。

上图中还有另外两个偏移量。“Log End Offset”表示最新一条写入的消息的偏移量，“HighWatermark”表示最新一条被成功复制副本的消息。从消费者的角度，它只能读到“HighWatermark”，这样可以防止消费者读取可能丢失的未复制的消息。

#### 初始化

就像Kafka其他客户端一样，需要使用Properties文件构建Consumer，如下面的示例

~~~java
Properties props = new Properties();
props.put("bootstrap.servers", "localhost:9092");
props.put("group.id", "consumer-tutorial");
props.put("key.deserializer", StringDeserializer.class.getName());
props.put("value.deserializer", StringDeserializer.class.getName());
KafkaConsumer<String, String> consumer = new KafkaConsumer<>(props); 
~~~

在上面代码中可以看到，consumer需要一个初始的broker列表用于去发现kafka集群其他部分，这个列表不需要是集群中所有服务器。

#### Topic订阅

在读取日志之前，Consumer需要订阅要消费的topic。如下面示例

~~~java
consumer.subscribe(Arrays.asList(“foo”, “bar”)); 
~~~

订阅了topic后，broker就会自动协调这个consumer和消费组中的其他成员对分区的分配。另外还可以通过$assgin$接口手动进行分区分配，但无法自动和手动混合使用。

订阅操作不是增量的：用户必须将要订阅的topic列表一次性的包含进去。每次订阅操作都会用新的topic列表替换掉之前订阅的列表。

#### Poll 循环

Consumer需要能并行的获取数据，而这些数据可能来自不同broker上的不同分区。为此它的API调用采用类似Unix的poll或着select风格：一旦订阅了主题，所有协调操作、再平衡、数据获取都是通过poll轮询调用来驱动。这种方式让人可以通过简单而高效的单线程方式实现所有IO操作。

订阅一个topic后，你需要启动事件循环来获取分区数据。看上去有些复杂，但你需要做的就是在一个循环中调用poll，而consumer负责处理其他的事情。每次轮训调用都会从分区返回一组消息。下面示例展示了基本的事件循环。

~~~java
try {
  while (running) {
    ConsumerRecords<String, String> records = consumer.poll(1000);
    for (ConsumerRecord<String, String> record : records)
      System.out.println(record.offset() + ": " + record.value());
  }
} finally {
  consumer.close();
}
~~~

$poll$接口根据‘current position’返回一组消息。当消费组首次创建时，初始读取位置通过重置策略设置，通常为$earliest$或者$latest$ 。一旦consumer开始提交偏移量，之后的每次再平衡都是根据最新一次提交的偏移量来设置初始读取位置。$poll$的参数用于设置consumer阻塞等待‘current position’上的数据的超时时间。当有数据可读时，consumer会尽快返回，否则它会一直等待直至超时。

consumer被设计为单线程运行，没有外部同步机制的情况下在多线程中使用时不安全的。在下面例子中有一个标志用于当应用程序关闭时中断轮询。当此标志从其他线程冲设置为false，则循环将中断。

~~~java
try {   
    while (true) {     
        ConsumerRecords<String, String> records = consumer.poll(Long.MAX_VALUE);
    	for (ConsumerRecord<String, String> record : records)
            System.out.println(record.offset() + “: ” + record.value());
	   	} 
	} catch (WakeupException e) {   
    	// ignore for shutdown 
	} finally {   
    	consumer.close(); 
	}
}
~~~

当不再使用consumer时，需要将其关闭。这样不仅可以清理不用的socket，还可以确保消费者向kafka发出其离开消费组的确切提示。

在上面例子中使用$Long.MAX\_VALUE$作为$poll$超时时间，这意味着consumer无限期的等待，直至有下一个记录。用于触发shutdown的线程可以通过调用$consumer.wakeup()$来中断这个poll，导致它爆出WakeupException异常。这个API是多线程安全的。注意，如果当前没有正在进行的poll，则会在下一次poll调用时引发异常。

下面是放到一起的代码

~~~java
public class ConsumerLoop implements Runnable {
  private final KafkaConsumer<String, String> consumer;
  private final List<String> topics;
  private final int id;

  public ConsumerLoop(int id,
                      String groupId, 
                      List<String> topics) {
    this.id = id;
    this.topics = topics;
    Properties props = new Properties();
    props.put("bootstrap.servers", "localhost:9092");
    props.put(“group.id”, groupId);
    props.put(“key.deserializer”, StringDeserializer.class.getName());
    props.put(“value.deserializer”, StringDeserializer.class.getName());
    this.consumer = new KafkaConsumer<>(props);
  }
 
  @Override
  public void run() {
    try {
      consumer.subscribe(topics);

      while (true) {
        ConsumerRecords<String, String> records = consumer.poll(Long.MAX_VALUE);
        for (ConsumerRecord<String, String> record : records) {
          Map<String, Object> data = new HashMap<>();
          data.put("partition", record.partition());
          data.put("offset", record.offset());
          data.put("value", record.value());
          System.out.println(this.id + ": " + data);
        }
      }
    } catch (WakeupException e) {
      // ignore for shutdown 
    } finally {
      consumer.close();
    }
  }

  public void shutdown() {
    consumer.wakeup();
  }
}
~~~

下面是运行上面的consumer loop的驱动程序

~~~java
public static void main(String[] args) { 
  int numConsumers = 3;
  String groupId = "consumer-tutorial-group"
  List<String> topics = Arrays.asList("consumer-tutorial");
  ExecutorService executor = Executors.newFixedThreadPool(numConsumers);

  final List<ConsumerLoop> consumers = new ArrayList<>();
  for (int i = 0; i < numConsumers; i++) {
    ConsumerLoop consumer = new ConsumerLoop(i, groupId, topics);
    consumers.add(consumer);
    executor.submit(consumer);
  }

  Runtime.getRuntime().addShutdownHook(new Thread() {
    @Override
    public void run() {
      for (ConsumerLoop consumer : consumers) {
        consumer.shutdown();
      } 
      executor.shutdown();
      try {
        executor.awaitTermination(5000, TimeUnit.MILLISECONDS);
      } catch (InterruptedException e) {
        e.printStackTrace;
      }
    }
  });
}
~~~

#### Consumer超时

当consumer分配了它所订阅的topic的一系列partition后，这些分区在组内是被锁定的，只要组成员保持不变，分区则会一直被固定consumer消费。如果消费者因机器或程序故障而消失，这个锁会被释放，进而触发再平衡，分区会分配给其他的成员。

Kafka组协调协议通过心跳机制来实现对组成员健康的探查。每次再平衡之后，组成员需要向组协调器发送定期心跳，只要协调器收到心跳则假定成员健康。每次收到心跳时，协调器会启动或重置计时器，如果在计时周期内没有收到心跳，则会认为成员已经死亡，进而会向组内其他成员发出信号，表明他们应该重新加入以便进行分区重新分配。计时器的时长被称为会话超时时间。客户端通过$session.timeout.ms$进行设置。

~~~java
props.put(“session.timeout.ms”, “60000”);
~~~

会话超时可确保由于机器故障、程序崩溃、网络故障导致consumer被隔离时消费组的锁被释放。然而consumer持续发送心跳不代表程序是依然健康的。

consumer 的poll循环则避免了这个问题。所有的网络IO都在调用poll时执行，而不是使用单独的后台线程。这意味着只有在调用poll时才会将心跳发送到协调器。如果程序停止轮询（无论是因为处理代码异常还是下游系统崩溃），那么将不会发送心跳，进而会话超时。

这种方式的唯一问题是如果consumer用太长的时间去处理读取的消息，则也会触发超时。这则需要根据实际情况，适当的调长超时时间。但超时时间又不宜太长，因为太长会导致协调器需要等待更长的时间才能检测到consumer异常。

#### 偏移量更新

当consumer组首次创建时，初始的偏移量根据$auto.offset.reset$设置。consumer开始处理消息，它需要适当的提交偏移量。每次再平衡后，都会从最新提交的偏移量开始读取消息。如果consumer在提交偏移量之前崩溃，则其他consumer可能会重复读取一部分消息。

在默认情况下$enale.auto.commit$是被设置为$true$的，表示consumer会自动提交偏移量。参数$auto.commit.interval.ms$用于设置自动提交偏移量的时间间隔。如果心跳一样，自动偏移量提交也是通过poll来实现。当poll时，如果发现上一次提交的时间到现在超过的设置的时长，则提交偏移量。

要想使用consumer的commit API，则需要通过设置$enable.auto.commit$为$false$来关闭自动提交。

~~~java
props.put(“enable.auto.commit”, “false”);
~~~

Commit API很容易使用，但他在poll循环中不同的集成方式会产生不同的提交策略。比如面示例

~~~java
try {
  while (running) {
    ConsumerRecords<String, String> records = consumer.poll(1000);
    for (ConsumerRecord<String, String> record : records)
      System.out.println(record.offset() + ": " + record.value());

    try {
      consumer.commitSync();
    } catch (CommitFailedException e) {
      // application specific failure handling
    }
  }
} finally {
  consumer.close();
}
~~~

在处理完poll返回的所有消息后，使用commitSync提交偏移量，它会一直阻塞到服务器返回提交成功或者error异常。在这里主要需要考虑的error是由于数据处理时间过长导致会话超时，此时协调器会将consumer踢出消费组并返回CommitFailedException，需要程序对这种异常进行适当处理。

通常情况下我们需要确保偏移量在消息被处理完后被提交。如果消费者在偏移量提交前崩溃会导致部分消息会被重复处理。这种确保最后提交的偏移量永不会超过‘current position’的策略被称为‘at last once’。

与上面策略相反，如果确保偏移量一直超过‘current position’的策略称为‘at most once’。如下图所示

<img src="https://www.confluent.io/wp-content/uploads/2016/08/New_Consumer_figure_3.png" width="600">

当cosumer崩溃时，‘current position’和 ‘last committed offset’之间的数据将“丢失”。但可以确保没有消息会被重复处理。下面代码实现了这个策略。

~~~java
try {
  while (running) {
  ConsumerRecords<String, String> records = consumer.poll(1000);

  try {
    consumer.commitSync();
    for (ConsumerRecord<String, String> record : records)
      System.out.println(record.offset() + ": " + record.value());
    } catch (CommitFailedException e) {
      // application specific failure handling
    }
  }
} finally {
  consumer.close();
}
~~~

使用Commit API还可以更精细的控制偏移量提交的时机，在极端情况下，没处理一个消息都提交一次偏移量。如下面示例所示

~~~java
try {
  while (running) {
    ConsumerRecords<String, String> records = consumer.poll(1000);

    try {
      for (ConsumerRecord<String, String> record : records) {
        System.out.println(record.offset() + ": " + record.value());
        consumer.commitSync(Collections.singletonMap(record.partition(), new OffsetAndMetadata(record.offset() + 1)));
      }
    } catch (CommitFailedException e) {
      // application specific failure handling
    }
  }
} finally {
  consumer.close();
} 
~~~

在上面例子中，在commitSync调用时显示的设置了要其提交的偏移量和分区。当没有这两个参数时，则会提交独处数据的最后一个偏移量+1。

很显然，每处理一个消息都提交一次偏移量会降低consumer的吞吐量。与之相比，一个更合理的办法是当处理完一个分区的数据后在提交这个分区的偏移量，如下例所示

~~~java
try {
  while (running) {
    ConsumerRecords<String, String> records = consumer.poll(Long.MAX_VALUE);
    for (TopicPartition partition : records.partitions()) {
      List<ConsumerRecord<String, String>> partitionRecords = records.records(partition);
      for (ConsumerRecord<String, String> record : partitionRecords)
        System.out.println(record.offset() + ": " + record.value());

      long lastoffset = partitionRecords.get(partitionRecords.size() - 1).offset();
      consumer.commitSync(Collections.singletonMap(partition, new OffsetAndMetadata(lastoffset + 1)));
    }
  }
} finally {
  consumer.close();
} 
~~~

到目前为止，所有示例都调用的是同步提交偏移量接口。还有一种异步方式，commitAsync。使用异步方式可以获得更搞得吞吐量，因为处理程序不用等待commit操作返回再进行下一步操作。但代价就是对commit失败的处理会延迟。如下例所示

~~~java
try {
  while (running) {
    ConsumerRecords<String, String> records = consumer.poll(Long.MAX_VALUE);
    for (TopicPartition partition : records.partitions()) {
      List<ConsumerRecord<String, String>> partitionRecords = records.records(partition);
      for (ConsumerRecord<String, String> record : partitionRecords)
        System.out.println(record.offset() + ": " + record.value());

      long lastoffset = partitionRecords.get(partitionRecords.size() - 1).offset();
      consumer.commitSync(Collections.singletonMap(partition, new OffsetAndMetadata(lastoffset + 1)));
    }
  }
} finally {
  consumer.close();
} 
~~~

#### 总结

1. consumer的核心操作时poll循环
2. 心跳检测、自动偏移量提交都是有poll触发。数据处理时间会影响poll调用之间的时间间隔，进而影响这超时或偏移量自动提交
3. 手动偏移量提交分‘at least one’、‘at most once‘两种策略，取决于提交偏移量的时机是在poll之前，数据处理之后。
4. 更精细的偏移量提交控制可以在数据处理之间，但会影响吞吐量
5. 使用异步偏移量提交可以增加吞吐量但失败处理会出现延迟

以上几种情况需要根据实际使用情况作出权衡取舍。
