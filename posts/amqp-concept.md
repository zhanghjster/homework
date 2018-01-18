####  概念

1. 代理服务器 Message Broker 
2. 队列，Queue 
3. 路由， Route 
4. 路由键，Routing key 
5. 绑定件，Binding key
6. 交换器，Exchange 
7. 直接交换器，Direct Exchange
8. 默认交换器，Default Exchange 



#### Fanout 交换器

最简单也是最快的交换器，直接将消息广播给所有绑定到它的所有队列或交换器，提供了典型的pub/sub机制

<img src="http://owo5nif4b.bkt.clouddn.com/QQ20180112-130224@2x.png" width="400">

#### 直接交换器

直接交换器使用'路由键'来路由消息。路由键由发布者在发布消息时设置，通常是一个由点号链接的多个字符串组成，比如 "booking.new"

队列通过'绑定键'与直接交换器绑定，直接交换器会将消息发送给所有'绑定键'与'路由键'相同的队列

<img src="http://owo5nif4b.bkt.clouddn.com/QQ20180112-133507.png" width="400">

#### 默认交换器

‘默认交换器‘是由’代理服务器‘预定义的一种没有名称的‘直接交换器’。每个创建的’队列‘都会以它的名称为'路由键'自动绑定到'默认交换器'。

假设一个队列被定义为以'search-index'为名称。一个发给'默认交换器'的消息，如果它的'路由健'是'search-indexing'，则这个消息就会路由给‘search-index’这个队列。这个特性使得看上去我们可以直接发送消息给指定名称的队列，但实际上底层还是经过了交换器。



