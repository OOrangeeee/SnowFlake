# Snowflake ID Generator 
# 雪花ID生成器

A distributed ID generator implementation based on Twitter's Snowflake algorithm in Go. 适用于分布式系统的雪花ID生成器实现。

## Features 特性
- ✅ Generate 64-bit unique IDs 生成64位唯一ID
- ⏱️ Time-based ID structure 基于时间戳的ID结构
- 🛡️ Sequence rollback protection 序列号回滚保护
- 🧩 Configurable node IDs 可配置节点ID
- 🚀 High concurrency ready 高并发就绪

## Installation 安装

```bash
go get github.com/oorangeeee/snow_flake
```

## Usage 使用

```go
// 单机模式
singleCreator := NewSnowFlakeCreatorForSingle()
singleID := singleCreator.GetId()

// 分布式模式（带数据中心）
// 参数：数据中心ID, 数据中心ID位数, 工作节点ID, 工作节点ID位数，注意：数据中心ID位数+工作节点ID位数<22（64位ID的位数，时间戳占41位，保留1位，序列号至少占1位）
clusterCreator := NewSnowFlakeCreatorForClusterWithDataCenter(3, 5, 7, 5)
clusterID := clusterCreator.GetId()

// 分布式模式（不带数据中心）
// 参数：工作节点ID, 工作节点ID位数，注意：工作节点ID位数<22（64位ID的位数，时间戳占41位，保留1位，序列号至少占1位）
workerOnlyCreator := NewSnowFlakeCreatorForClusterWithoutDataCenter(100, 10)
workerID := workerOnlyCreator.GetId()

```

## Performance 性能

```bash
go test -v -timeout 30s -run TestMaxIDsPerSecond ./...
```

运行此UT即可测试性能，在本机测试结果如下：

```bash
=== RUN   TestMaxIDsPerSecond
    snow_flake_test.go:113: 一秒16个协程16核心，生成ID数量：18684939
--- PASS: TestMaxIDsPerSecond (1.00s)
PASS
ok      github.com/OOrangeeee/SnowFlake/snow_flake      1.481s
```
