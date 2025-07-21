# 飞书 Webhook SDK

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue.svg)](https://golang.org/)

一个简单易用的飞书机器人 Webhook SDK，支持多种消息类型和签名验证。

## 功能特性

- ✅ 支持签名和非签名两种发送方式
- ✅ 支持多种消息格式：文本、富文本、图片、交互式卡片
- ✅ 简单易用的API接口
- ✅ 完整的错误处理
- ✅ 支持Go Modules

## 安装

```bash
go mod init your-project
go get github.com/straubel/feishu-webhook
```

## 快速开始

### 基本用法

```go
package main

import (
    "log"
    "github.com/straubel/feishu-webhook/common/feishu"
)

func main() {
    webhookURL := "https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-url"
    secret := "your-secret-key"

    // 创建SDK实例（带签名）
    sdk := feishu.New(webhookURL, secret)
    
    // 发送文本消息
    err := sdk.SendText("Hello, 飞书!")
    if err != nil {
        log.Printf("发送失败: %v", err)
    }
}
```

### 不带签名的使用方式

```go
// 创建SDK实例（不带签名）
sdk := feishu.New(webhookURL)

// 发送文本消息
err := sdk.SendText("Hello, 飞书!")
if err != nil {
    log.Printf("发送失败: %v", err)
}
```

### 发送富文本消息

```go
// 构建富文本内容
richTextContent := [][]feishu.RichTextElement{
    {
        feishu.CreateRichTextElement("text", "这是标题\n"),
    },
    {
        feishu.CreateRichTextElement("text", "包含链接: "),
        feishu.CreateRichTextElement("a", "飞书官网", map[string]string{
            "href": "https://www.feishu.cn",
        }),
    },
    {
        feishu.CreateRichTextElement("text", "以及@用户: "),
        feishu.CreateRichTextElement("at", "张三", map[string]string{
            "user_id": "user123",
        }),
    },
}

err := sdk.SendRichText("富文本标题", richTextContent)
```

### 发送图片消息

```go
// 发送图片（需要先上传图片获取image_key）
err := sdk.SendImage("img_v2_041b28e3-5680-48c2-9af2-497ace79333g")
```

### 发送交互式卡片

```go
header := feishu.CreateCardHeader("卡片标题", "blue")
config := feishu.CreateCardConfig(true)

elements := []interface{}{
    map[string]interface{}{
        "tag": "div",
        "text": map[string]interface{}{
            "content": "这是卡片内容",
            "tag":     "plain_text",
        },
    },
}

err := sdk.SendInteractive(config, header, elements)
```

### 便捷函数

如果只是偶尔发送消息，可以使用便捷函数：

```go
// 发送文本消息（带签名）
err := feishu.SendTextMessage(webhookURL, "Hello", secret)

// 发送文本消息（不带签名）
err := feishu.SendTextMessage(webhookURL, "Hello")

// 发送富文本消息
err := feishu.SendRichTextMessage(webhookURL, "标题", richTextContent, secret)

// 发送图片消息
err := feishu.SendImageMessage(webhookURL, "image_key", secret)
```

## API 文档

### 创建客户端

```go
// 带签名
sdk := feishu.New(webhookURL, secret)

// 不带签名
sdk := feishu.New(webhookURL)
```

### 发送消息方法

- `SendText(text string) error` - 发送文本消息
- `SendRichText(title string, content [][]RichTextElement) error` - 发送富文本消息
- `SendImage(imageKey string) error` - 发送图片消息
- `SendInteractive(config *CardConfig, header *CardHeader, elements []interface{}) error` - 发送交互式卡片
- `SendMessage(message *Message) error` - 发送自定义消息

### 辅助函数

- `CreateRichTextElement(tag, text string, options ...map[string]string) RichTextElement` - 创建富文本元素
- `CreateCardHeader(title, template string) *CardHeader` - 创建卡片头部
- `CreateCardConfig(enableForward bool) *CardConfig` - 创建卡片配置

## 错误处理

所有发送方法都会返回error，建议进行适当的错误处理：

```go
if err := sdk.SendText("Hello"); err != nil {
    log.Printf("发送消息失败: %v", err)
    // 处理错误，比如重试、记录日志等
}
```

## 测试

### 运行测试

```bash
# 运行所有测试
go test ./common/feishu/

# 运行测试并显示覆盖率
go test -cover ./common/feishu/

# 运行基准测试
go test -bench=. ./common/feishu/

# 运行特定测试
go test -run TestSendText ./common/feishu/

# 跳过集成测试（仅运行单元测试）
go test -short ./common/feishu/
```

### 测试覆盖率

```bash
# 生成覆盖率报告
go test -coverprofile=coverage.out ./common/feishu/
go tool cover -html=coverage.out -o coverage.html
```

### 测试分类

- **单元测试**: 测试各个函数和方法的基本功能
- **集成测试**: 测试完整的消息发送流程
- **基准测试**: 测试性能指标
- **错误处理测试**: 测试各种错误场景

## 注意事项

1. 请确保Webhook URL的正确性
2. 如果使用签名验证，请确保secret的正确性
3. 图片消息需要先通过飞书API上传图片获取image_key
4. 富文本和卡片消息的格式请参考飞书官方文档

## 许可证

本项目采用 MIT 许可证。详细信息请查看 [LICENSE](LICENSE) 文件。