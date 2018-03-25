# mygotools
一些用GO编写的小工具

## Tool List

### mergeFile

包，合并某目录下的文件，按照文件名排序的方式。

### novelDownload

程序，下载笔趣阁类型网站的小说，需事先获取链接，最终将每一章节下载到 output/小说名 目录下，然后调用mergeFile合并文件。先并发获取主页面上每一章节的地址，利用map进行去重操作，然后并发下载，最终合并文件。

使用方法： ```novelDownload http://www.biquge.com.tw/4_4029/```

> bug:   
> 1. windows下中文名会乱码  
> 2. 合并文件用的mergeFile一次性读取文件，耗内存   

### gotunet

tsinghua网络的登录和登出工具。主要是利用do_login.php和rad_user_info.php两个进行get和post请求。登出工具另写了tunetlogout。gotunet的password通过MD5加密。编译后双击即可运行，其中passwd.txt用来存放用户名和密码，用逗号隔开，与编译后的文件放在一起即可运行。
