# pakku 帕克

    pakku 的核心是提供一个对象加载的环境, 对类型为`ipakku.Module`的对象进行加载, 通过对各个`Module`加载节点事件的监听, 提供额外的辅助功能. 如依赖注入、自动完成配置值等操作. 其核心轻量无且三方引用依赖. 使用内置的几个模块, 便可快速搭建一个简单的服务或程序.
    

## 内置的模块

    pakku 默认实现了AppConfig(配置模块)、AppCache(缓存模块)、AppEvent(事件模块)、AppService(NET服务模块)以满足一个基本的服务运行环境. 如需使用默认的接口定义但又需要使用其他方式实现, 可通过重新实现对应`ipakku/Ixxx`的接口后, 重新指定默认调用实例即可(ipakku.Override.RegisterInterfaceImpl).

|  名字 |  可重写接口类  |  描述  |
| ------ | ------ | ------ |
| AppConfig | `ipakku.IConfig` | 使用json格式存储的配置实现, 文件存放在启动目录下`./.conf/{appName}.json`中 |
| AppCache | `ipakku.ICache` | 使用map实现的本地内存缓存, 如需使用其他缓存机制, 如redis需要自己实现 |
| AppEvent | `ipakku.IEvent` | 默认没有实现此接口, 需要自己实现, 如: kafka等 |
| AppService | `-` | 默认实现了http服务和rpc服务, 不可重写, 但可选启用改模块 |


## 特殊标签(tag)

    通过标注在struct的特殊标签值, 来实现一些辅助功能. 

|  TAG |  所属模块  |  作用域  |  格式  |  描述  |
| ------ | ------ | ------ | ------ | ------ |
| `@autowired` | Loader(加载器) | struct成员字段 | `@autowired:"模块名"` | 通过指定该标签, 可实现依赖对象自动注入 |
| `@autoConfig` | AppConfig(配置模块) | struct成员字段 | `@autoConfig:"配置路径"` |  标注当前字段是个配置struct, 可选''配置路径'参数  |
| `@value` | AppConfig(配置模块) | struct成员字段 | `@value:"配置路径"` | 通过'配置路径'查找并自动赋值对应字段, 可选''配置路径'参数 |


## 如何使用

    参考'pakku_demo_test.go'文件示例

```golang
// 实例化一个application, 启用核心模块和网络服务模板并把日志级别设置为DEBUG
app := NewApplication("app-test").EnableCoreModule().EnableNetModule().SetLoggerLevel(logs.DEBUG).BootStart()

// 获取内部的一个模块, 这里使用 AppService 用于开启一个服务
var service ipakku.AppService
if err := app.GetModuleByName("AppService", &service); nil != err {
    logs.Panicln(err)
}

// 设置一个静态页面路径
if err := service.SetStaticDIR("/", "./", nil); nil != err {
    logs.Panicln(err)
}

// 手工注册一个请求路径(可使用Controller接口批量注册)
if err := service.Get("/hello", func(rw http.ResponseWriter, _ *http.Request) {
    rw.Write([]byte("hello!"))
}); nil != err {
    logs.Panicln(err)
}

// 启动服务
service.StartHTTP(ipakku.HTTPServiceConfig{ListenAddr: "127.0.0.1:8080"})
```