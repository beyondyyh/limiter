# 项目名称

Limiter 简易QPS限流策略，如令牌桶等。

## 快速开始

```go
import "gitlab.weibo.cn/gdp/limiter"

var limit = limiter.NewTokenBucket(&limiter.Config{
    QPS:      100, // 单机qps
    MaxCount: 200,
})

if err := limit.Run(func() {
    // do something...
}); err != nil {
    ctx.AddNotice("QPS.exceed", 1)
}
```

## 测试

如何执行自动化测试？

## 如何贡献

提交PR -> fork分支 -> 开发代码 -> 提CR -> 管理员merge

## 讨论

QQ讨论群：XXXX
