# mygotools
一些用GO编写的小工具

## Tool List

### mergeFile

包，合并某目录下的文件，按照文件名排序的方式。

### novelDownload

程序，下载笔趣阁类型网站的小说，需事先获取链接，最终将每一章节下载到 output/小说名 目录下，然后调用mergeFile合并文件

使用方法： ```novelDownload http://www.biquge.com.tw/4_4029/```

> bug:   
> 1. windows下中文名会乱码  
> 2. 合并文件用的mergeFile一次性读取文件，耗内存   


