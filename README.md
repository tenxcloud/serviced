# 升级版 `beehive` 提案

在 KubeEdge 中，beehive 作为模块间异步交互的基础框架。
解决的核心问题是，在服务器、客户端间只有一个通信连接时，如何有效的基于 'message'
路由请求到具体的处理上。


## 现存 beehive 实现存在的痛点

模块间调用是'隐式'的，就是说如果不熟悉代码，是不太可能知道模块间调用关系的。
导致的结果有：
1. 有贡献代码意图的开发人员，阅读成本提高；
2. 不容易排错；
3. 模块间交互调用复杂的时候，容易出错；
4. 模块可以 Receive 不应该由它 Receive 的东西。


## 升级版提案的优点

1. 基于接口，有具体类型；
2. 依赖注入，进一步解耦；


### 示例

参见 `example/main.go`


### 提案阶段实现版的局限性（可进一步开发解决）

1. 参数跟返回值无法处理指针类型；
2. 一个接口多个实现类型无法反序列化（error）
