# sms
用golang实现的流媒体服务器,来自开源，现在回馈给开源社区。

## 特性
* 模块化，易扩展
* 高性能，易部署，跨平台
* 支持rtmp推拉流
* 支持http-flv分发
* 支持hls分发
* 支持http接口请求来推拉流
* 支持rtmp流录制成flv文件

## 源码编译

```
git clone https://github.com/sheepbao/sms.git  
cd sms  
./build.sh  
```

## 交叉编译
```
./build.sh linux arm // or ./build.sh os arch
```

## 运行
```
./bin/sms
```