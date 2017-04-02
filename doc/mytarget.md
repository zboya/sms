* 目标：  
    功能最强大的流媒体服务器

* 支持的协议：  
    rtmp  
    http  
    hls  
    rtp  
    rtsp  
    webrtc  
    dash  

* 功能：  
    第一阶段实现支持协议的推流。
    支持rtmp，rtp推流
    rtp转rtmp
    
* 理念：  
    模块化设计，自定义组合流媒体协议  

* 设计：  
    顶层设计，分模块实现。  
    实现的功能的顺序有上到下。  
    
* 计划：  
    rtmp(2m)  
    http(2m)  
    hls(2m)  
    rtp(1m)  
    rtsp(1m)  
    webrtc(2m) 
    dash(2m)  
    第一阶段实现各个协议的支持。