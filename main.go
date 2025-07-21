package main

import (
	"fmt"
	"log"

	"github.com/straubel/feishu-webhook/common/feishu"
)

func main() {
	webhookURL := "https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-url"
	secret := "your-secret-key"

	// 示例1: 发送文本消息（带签名）
	sdk := feishu.New(webhookURL, secret)
	err := sdk.SendText("Hello, 这是一条测试消息!")
	if err != nil {
		log.Printf("发送文本消息失败: %v", err)
	}

	// 示例2: 发送文本消息（不带签名）
	sdkNoSign := feishu.New(webhookURL)
	err = sdkNoSign.SendText("Hello, 这是一条不带签名的测试消息!")
	if err != nil {
		log.Printf("发送不带签名的文本消息失败: %v", err)
	}

	// 示例3: 发送富文本消息
	richTextContent := [][]feishu.RichTextElement{
		{
			feishu.CreateRichTextElement("text", "这是一条富文本消息\n"),
		},
		{
			feishu.CreateRichTextElement("text", "包含链接: "),
			feishu.CreateRichTextElement("a", "飞书官网", map[string]string{"href": "https://www.feishu.cn"}),
		},
		{
			feishu.CreateRichTextElement("text", "以及 @用户", map[string]string{"user_id": "user123", "user_name": "张三"}),
		},
	}
	
	err = sdk.SendRichText("富文本标题", richTextContent)
	if err != nil {
		log.Printf("发送富文本消息失败: %v", err)
	}

	// 示例4: 发送图片消息
	err = sdk.SendImage("img_v2_041b28e3-5680-48c2-9af2-497ace79333g")
	if err != nil {
		log.Printf("发送图片消息失败: %v", err)
	}

	// 示例5: 使用便捷函数发送消息
	err = feishu.SendTextMessage(webhookURL, "使用便捷函数发送的消息", secret)
	if err != nil {
		log.Printf("便捷函数发送消息失败: %v", err)
	}

	// 示例6: 发送交互式卡片消息
	header := feishu.CreateCardHeader("卡片标题", "blue")
	config := feishu.CreateCardConfig(true)
	
	// 这里可以添加各种卡片元素
	elements := []interface{}{
		map[string]interface{}{
			"tag": "div",
			"text": map[string]interface{}{
				"content": "这是卡片内容",
				"tag":     "plain_text",
			},
		},
	}
	
	err = sdk.SendInteractive(config, header, elements)
	if err != nil {
		log.Printf("发送交互式卡片失败: %v", err)
	}

	fmt.Println("所有示例执行完成!")
}
