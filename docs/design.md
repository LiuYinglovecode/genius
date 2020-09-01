adc-genius设计思路
===

用golang实现k8s相关操作包括helm的解析以及k8s的操作，以微服务的形式提供。


rest web架构参考framework的实现。使用gin。

## 中间件微服务
提供中间件的相关管理，主要功能包括：

- ??兼容Helm repo相关规范?? P2
- 上传中间件 P2
- ??下载中间件?? P2
- 获取中间件基本信息 P1
    - 获取Chart.yaml相关metadata
    - 获取Values.yaml相关内容
    - 获取help说明
- 根据参数渲染成部署文件 P1

实现上参考**chartmusem**
