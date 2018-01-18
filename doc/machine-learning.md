#### 线性模型

##### 基本形式

通过学得属性的线性组合来进行预测的函数

$$\\ f(x) = w_1x_1+w_2x_2+….+w_dx_d + b$$

用向量表示

$$\\ f(x) = w^Tx + b$$

学得 $w = (w_1;w_2;…;w_d)$ 后，模型就得以确定

##### 线性回归

给定数据集 $D=\{(x_1, y_1), (x_2, y_2),…,(x_d, y_d)\}$ , 其中 $x_i = (x_{i1};x_{i2};…x_{id})$ , $y_i \in \mathbb{R}$ ，试图学得一个线性模型尽可能的预测实值输出标记

$$\\ f(x_i) = wx_i + b$$ 

使得

$$\\ f(x_i) \approx y_i$$ 

如何确定$w, b$ 呢？通过实现最小化均方差

$$\\ (w^*, b^*) = \underset{(w,b)}{\arg\min}\sum_{i=1}^m(f(x_i) - y_i)^2 \\ = \underset{(w,b)}{\arg\min}\sum_{i=1}^m(y_i - wx_i - b)^2 $$

均方差对应了常用的欧几里得距离，在线性回归中，使用最小二乘法试图找到一个直线，使所有样本到直线上的欧氏距离之和最小

$$\\ E(w, b) = \sum_{i=1}^m(y_i - wx_i - b)^2$$

分别对$w$ 和 $b$ 求导

$$\\ \frac{\partial E_{(w,b)}}{\partial w} = 2(w\sum_{i=1}^mx_i^2 - \sum_{i=1}^m(y_i-b)x_i )$$

$$\\ \frac{\partial E_{(w, b)}}{\partial b} = 2(wb - \sum_{i=1}^m(y_i - wx_i))$$

上面两式为0可得w b最优解

$$\\ w = \frac{\sum_{i=1}^my_i(x_i - \frac{1}{m}\sum_{i=1}^mx_i)}{\sum_{i=1}^mx_i^2 - \frac{1}{m}(\sum_{i=1}^mx_i)^2}$$

$$\\ b = \frac{1}{m}\sum_{i=1}^m(y_i - wx_i) $$

更一般的形式

$$\\ f(x_i) = w^Tx_i + b $$

使得

$$\\ f(x_i) \approx y_i$$

称为多元线性回归

