#### 入门

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

```python
v = tf.Variable()
```

在使用前需要初始化

```python
sess.run(v.initializer)
```

如果有多个变量可以使用下面方法一次性全部初始化

```python
sess.run(tf.global_varivables_initializer())
```

所有变量都会自动加入到tf.VARIABLES, 如果一个变量声明时候trainable为True，则自动加入到tf.TRAINABLE_VARIABLES, 分别使用下面方法获取

```python
tf.all_variables()
tf.trainable_variables()
```

监督学习的重要思想是在已知答案的标注数据集上模型给出的预测结果要尽量接近真实答案，神经网络的重要优化算法中，最常用的是反向传播(backpropagation)，用于优化神经网络的参数取值

为了在每轮训练迭代中更改测试数据，tf提供了placeholder来动态的获取输入数据，避免大量的常量带来的开销问题。

##### 完整举例

```python
# coding:utf8
import tensorflow as tf
from numpy.random import RandomState
import os
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '2'

batch_size = 8

w1 = tf.Variable(tf.random_normal([2, 3], stddev=1, seed=1))
w2 = tf.Variable(tf.random_normal([3, 1], stddev=1, seed=1))

x = tf.placeholder(tf.float32, shape=(None, 2), name='x-input')
y_ = tf.placeholder(tf.float32, shape=(None, 1), name='y-input')

a = tf.matmul(x, w1)
y = tf.matmul(a, w2)

# 损失函数
cross_entropy = -tf.reduce_mean(
   y_ * tf.log(tf.clip_by_value(y, 1e-10, 1.0))
)

# 定义反向传播算法
train_step = tf.train.AdamOptimizer(0.001).minimize(cross_entropy)

# 通过随机数生成训练集
rdm = RandomState(1)

dataset_size = 128

X = rdm.rand(dataset_size, 2)
Y = [[int(x1+x2 < 1)] for (x1, x2) in X]

with tf.Session() as sess:
    # 初始化所有变量
    sess.run(tf.global_variables_initializer())

    # 训练前的w1 w2
    print sess.run(w1)
    print sess.run(w2)

    STEPS = 5000
    for i in range(STEPS):
        start = (i*batch_size) % dataset_size
        end = min(start+batch_size, dataset_size)

        sess.run(train_step, feed_dict={x: X[start:end], y_: Y[start:end]})
        if i % 1000 == 0:
            total_cross_entropy = sess.run(
                cross_entropy, feed_dict={x: X, y_: Y}
            )
            print("After %d training step(s), cross entropy on all data is %g" % (i, total_cross_entropy))

    # 训练后的w1 w2
    print sess.run(w1)
    print sess.run(w2)

```

训练神经网络的过程：

1. 定义神经网络结构和向前传播输出结果
2. 定义损失函数以及选择反向传播优化的算法
3. 生成会话并且在训练数据上反复进行反向传播优化算法

#### 深层神经网络

##### 深度学习与深层神经网络

深度学习是一类通过多层非线性变化对高复杂性数据建模算法的合集

##### 线性模型的局限性

只通过线性变化，任意层的全连接神经网络和单层的神经网络模型的表达能力没有区别，线性模型解决的问题是有线的，他不能解决非线性问题

##### 激活函数实现去线性化

激活函数将神经元的输出做非线性变化，使整个神经网络模型不再是线性，比如

```python
a = tf.nn.relu(tf.matmul(x, w1) + biases1
```

##### 损失函数

神经网络解决分类问题最常用的办法是设置n个输出节点，n为类别的个数，对每个样例，神经网络可以得到一个n维数组作为输出结果，用交叉熵判断输出向量和期望的向量之间距离

softmax层将神经网络的输出变成一个概率分布，也就是一个样例为不同类别的概率分别是多大

$$\\ softmax(y)_i = y^{'}=\frac{e^{(yi)}}{\sum_{j=1}{^n}e^{yj}}$$

这样就可以通过交叉熵计算来预测的概率分布和真实答案的概率分布之间的距离了

$$\\H(p,q) = - \sum_xp(x)logq(x)​$$

p为代表正确答案，q代表预测答案

