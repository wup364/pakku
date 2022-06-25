# pakku 帕克

    pakku 的核心是提供一个对象加载的环境, 在加载对象的同时提供额外的服务, 如依赖注入等操作. 其核心轻量无且三方引用依赖. 使用内置的几个模块, 便可快速搭建一个简单的服务或程序.

## 内置的接口

    pakku 在加载实现了`ipakku.Module`、`ipakku.Controller`、`ipakku.Router`接口的对象以及RPC注册对象时, 会对内部的成员变量(包括私有)进行扫描, 对包含`@autowired:"接口名"`标签的字段进行自动赋值, 实现了模块间的解耦和依赖注入.

* `Module` 即模块, 通常为有一定功能单元的对象, 在`Controller`、`rpc服务`、`其他Module`中被注入和调用.
* `Router` 即http路由定义模块, 通过定义接口的url、地址列表实现对url的路由.
* `Controller` 即http服务定义模块, 通过定义接口的url、地址列表、过滤器等实现对url的路由.
* `RpcService` 即rpc服务定义模块, 默认使用自带的rpc服务进行注册.

## 内置的模块

    pakku 默认实现了AppConfig(配置模块)、AppCache(缓存模块)、AppEvent(事件模块)、AppService(NET服务模块)以满足一个基本的服务运行环境. 如需使用默认的接口定义但又需要使用其他方式实现, 可通过重新实现对应`ipakku/Ixxx`的接口后, 重新指定默认调用实例即可(ipakku.Override.RegisterInterfaceImpl).

|  名字 |  可重写接口类  |  描述  |
| ------ | ------ | ------ |
| AppConfig | `ipakku.IConfig` | 使用json格式存储的配置实现, 文件存放在启动目录下`./.conf/{appName}.json`中 |
| AppCache | `ipakku.ICache` | 使用map实现的本地内存缓存, 如需使用其他缓存机制, 如redis需要自己实现 |
| AppEvent | `ipakku.IEvent` | 默认没有实现此接口, 需要自己实现, 如: kafka等 |
| AppService | `-` | 默认实现了http服务和rpc服务, 不可重写 |

## 如何使用

    参考'pakku_test.go'文件示例
