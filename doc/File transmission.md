# 关于文件传输

## 1.分块缓存
将文件按64kb进行分块，并将最近最常用的块存放在缓存中。
> 为什么不以文件为单位进行缓存？文件太大了，太吃内存。
## 2.随机读取
参考：
* https://simpsonlab.github.io/2015/05/19/io-performance/
* https://blog.cloudera.com/apache-hbase-i-o-hfile/

HBase实现原理：顺序文件+索引。顺序文件只能进行append操作，索引也是一个顺序文件，用于记录每N条记录的偏移。
## 3.断点续传