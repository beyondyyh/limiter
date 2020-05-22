# 项目名称

Limiter 简易QPS限流策略，如令牌桶等。

## 快速开始

```go
import "github.com/beyondyyh/limiter"

var limit = limiter.NewTokenBucket(&limiter.Config{
    QPS:      100, // 单机qps
    MaxCount: 200,
})

if err := limit.Run(func() {
    // do something here...
}); err != nil {
    log.Print("QPS exceed")
}
```

