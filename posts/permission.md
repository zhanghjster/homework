---
title: 权限
date: 2017-11-30 23:58:48
tags:
   - 基础
---

通过本文总结一下如何识别和设置文件和目录的权限，以及如何更改他们的所有权、群组等信息

#### 文件权限

##### 查看

通过 ls命令查看，以下面结果为例

```
drwxr-xr-x  2 Ben  staff    64B 11 30 22:41 bar
-rw-r--r--  1 Ben  staff     0B 11 30 22:41 bar.txt
lrwxr-xr-x  1 Ben  staff     7B 11 30 22:42 bar_link.txt -> bar.txt
drwxr-xr-x  2 Ben  staff    64B 11 30 22:41 foo
-rw-r--r--  1 Ben  staff     0B 11 30 22:41 foo.txt
```
<!-- more -->
第一列为表示权限，分为四部分，其中

- 第一个字符为用于区分“目录”、“文件”、链接
  - ’-‘ 表示文件
  - ’d'表示目录
  - ‘l'表示链接
- 第2-4字符用于表示文件所有者的权限分别用'r''w''x'表示读、写、执行的权限
- 第5-7字符表示群组所拥有的权限，定义同上
- 第8–10字符表示其他用户拥有的权限，定义同上

##### 修改

- 增加权限，使用chmod命令结合‘+’以及'r'(读)、'w'(写)、'x'(执行)组合来增加权限

  比如  ‘chmod +rw’ 增加读写权限

- 减少权限， 使用chmod命令结合‘-’以及'r'(读)、'w'(写)、'x'(执行)组合来减少权限

  比如‘chmod -rw’去掉执行权限

```
localhost:test Ben$ ll foo.txt 
-rw-r--r--  1 Ben  staff     0B 11 30 22:41 foo.txt
localhost:test Ben$ chmod +x foo.txt 
localhost:test Ben$ ll foo.txt 
-rwxr-xr-x  1 Ben  staff     0B 11 30 22:41 foo.txt
```

上面演示了如何增加和减少可执行权限, 这个办法只更改了文件所有者的权限， 如果要更改’组’或者‘其他用户’的权限‘则需要结合'g'或者‘o'来更改， 如下

```
localhost:test Ben$ chmod g+w foo.txt 
localhost:test Ben$ ll foo.txt 
-rwxrwxr-x  1 Ben  staff     0B 11 30 22:41 foo.txt
localhost:test Ben$ chmod g-w foo.txt
localhost:test Ben$ ll foo.txt 
-rwxr-xr-x  1 Ben  staff     0B 11 30 22:41 foo.txt
```

上面增加和减少foo.txt的组的可执行权限

```
localhost:test Ben$ chmod o+w foo.txt 
localhost:test Ben$ ll foo.txt 
-rwxr-xrwx  1 Ben  staff     0B 11 30 22:41 foo.txt
localhost:test Ben$ chmod o-w foo.txt 
localhost:test Ben$ ll foo.txt 
-rwxr-xr-x  1 Ben  staff     0B 11 30 22:41 foo.txt
```

上面演示了增加和减少foo.txt的其他用户的写权限

##### 数字模式权限

上面都是以 'w' 'r' 'x' 定义权限，还可以用数字表示权限

- 0 - 没有权限
- 1 - 执行权限
- 2 - 写权限
- 4 - 读权限

可以通过数的相加来组合权限，比如 6 表示可读可写， 全部组合如下

- 0 = ---
- 1 = --x
- 2 = -w-
- 3 = -wx
- 4 = r--
- 5 = r-x
- 6 = rw-
- 7 = rwx

#### 目录权限

##### 查看

通过ls命令查看，比如下面结果

```
drwxr-xr-x  2 Ben  staff  64 11 30 22:41 foo
```

第一列结果为权限部分，同文件一样也分为四块，每块定义与文件也相同，但'w','r','x'的定义和文件的有所不同

- 'w', 用户可以重命名目录。在有’x'的权限下可以删除目录，还可以在增加、删除、更新、重命名目录里的文件,

  ```
  localhost:test Ben$ ll
  total 0
  d-w-------  2 Ben  staff    64B 11 30 23:45 foo
  localhost:test Ben$ rm -fr foo/
  rm: foo/: Permission denied
  localhost:test Ben$ mv foo/ bar
  localhost:test Ben$ ll 
  total 0
  d-w-------  2 Ben  staff    64B 11 30 23:45 bar
  ```

  只有'w'权限可以mv但不能删除目录

- 'x', 用户可以进入目录，但需要有'w'权限才能在目录里增加、删除、更新文件或子目录

  ```
  localhost:test Ben$ ll 
  total 0
  d--x--x--x  2 Ben  staff    64B 11 30 23:45 bar
  localhost:test Ben$ cd bar/
  localhost:bar Ben$ touch a
  touch: a: Permission denied
  ```

  只有'x'权限，只能进入目录，不能增加、删除、更改文件

  ```
  localhost:test Ben$ ll
  total 0
  d-wx--x--x  2 Ben  staff    64B 11 30 23:45 bar
  localhost:test Ben$ touch bar/a
  localhost:test Ben$ rm bar/a
  ```

  有'w'和'x'权限可以增加更改目录里的文件

- 'r', 用户可以查看目录的信息，比如子目录和文件

##### 更改

操作如同文件的

#### 总结

很基础的内容，需要系统总结才能记得牢
