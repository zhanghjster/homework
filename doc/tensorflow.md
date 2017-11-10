##### 计算图（计算模型）

Tensorflow中的tensor是张量的意思，可以理解为多维数组。flow是流的意思，直观的表达了张量之间通过计算相互转化的过程。Tensorflow是通过计算图的形式来表达计算的编程系统。每个计算都是计算图上的一个节点，节点之间的变描述了计算之间的依赖关系。

在tensorflow程序中，系统会自动会维护一个默认的计算图

```python
import tensorflow as tf
tf.get_default_grtaph()
```

可以通过tf.Graph来生成新的计算图，不同计算图上的张量和运算都不会共享

```python
g = tf.Graph()
with g.as_default():
  ......
```

#####  张量（数据模型）

从功能角度，张量可以理解为一个多维数组。但张量在TF中实现不是一个数组而是对一个运算结果的引用。张量中并没有保存数字，保存的是如何得到这个数字的计算过程

```python
import tensoflow as tf

a = tf.constant([1.0, 2.0], name="a")
b = tf.constant([3.0, 4.0], name="b")

# result为一个张量，并不是具体的加法结果
result = a + b
print(result)

# 输出
# Tensor("add:0", shape=(2,), dtype=float32)
```

张量有名称(name),维度(shape)和类型(type)三个属性

* name，不仅仅是一个张良的唯一标识，同样给出这个张量是如何计算出来的。命名规则"node:src_output", node为节点名称,src_output表示来自节点的第几个输出
* shape，表示张量的维度信息
* type表示数据类型，如果张量之间的数据类型不合符会报错

##### 会话（运行模型） 

会话用于管理TF程序运行时的所有资源，当计算完成后，需要关闭会话来回收资源。会话的使用一般有两种模式

第一种是明确调用会话生成、运行、关闭

```python
sess = tf.Session()
sess.run(...)
sess.close()
```

第二种通过上下文管理器

```python
with tf.Session() as sess:
	sess.run(...)
```

TF不会生成默认的会话而需要手工指定。指定默认session后就可以通过tf.Tensor.eval()来计算张量的值

```python
sess = tf.Session()
with sess.as_default():
	print(result.eval())
```

第三种，在交互式环境下直接构建默认会话，他会将生成的会话注册为默认会话

```python
sess = tf.InteractiveSession()
print(result.eval())
```

##### 神经元

一个神经元有多个输入和一个输出，每个神经元的出入既可以是其他神经元的输出也可以是整个神经网络的输入。神经网络的结构将不同神经元连接的结果。神经元在神经网络中也可以成为节点。

神经网络向前传播需要三部分信息

1. 神经网络的输入
2. 神经网络的连接
3. 每个神经元中的参数

##### 神经网络参数与TF变量

神经网络中的参数是实现分类或回归问题的重要组成部分。使用tf.Variable保存和更新神经网络中的参数



