# Checklist

## 模块一：gRPC基础
- [ ] protoc环境安装正确
- [ ] buf CLI安装正确，buf.yaml和buf.gen.yaml配置完整
- [ ] `buf generate`能正确生成Go代码
- [ ] proto文件语法正确，能生成Go代码
- [ ] Unary RPC服务端和客户端可运行
- [ ] Server Streaming演示流式返回
- [ ] Client Streaming演示批量上传
- [ ] Bidirectional Streaming演示双向通信

## 模块二：gRPC进阶
- [ ] 日志拦截器正确记录请求
- [ ] 认证拦截器正确拦截未授权请求
- [ ] Recovery拦截器正确恢复panic
- [ ] 三个拦截器通过`grpc.ChainUnaryInterceptor`组合运行，执行顺序符合预期
- [ ] 错误码使用正确，客户端能识别
- [ ] 超时控制正确，能取消请求

## 模块三：微服务实战
- [ ] 用户服务独立运行
- [ ] 订单服务独立运行
- [ ] 订单服务使用`grpc.NewClient`（非废弃的`grpc.Dial`）创建连接
- [ ] gRPC连接全局单例复用，程序退出时正确关闭
- [ ] 订单服务能调用用户服务
- [ ] 完整调用链演示成功

## 代码质量
- [ ] error处理符合Go惯用法：函数返回error，调用方`if err != nil`处理，不用try-catch思维
- [ ] context正确从入口传递到底层调用，没有用`context.Background()`兜底
- [ ] 没有goroutine泄漏（流式RPC的goroutine正确退出）

## 文档完整性
- [ ] 每个示例都有Java对比说明
- [ ] README文档完整
- [ ] proto文件注释清晰
