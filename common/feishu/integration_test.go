package feishu

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIntegrationComplete(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	var lastRequest *WebhookRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req WebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			w.WriteHeader(400)
			return
		}
		lastRequest = &req

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code": 0, "msg": "success"}`))
	}))
	defer server.Close()

	t.Run("完整的文本消息流程", func(t *testing.T) {
		sdk := New(server.URL, "integration-test-secret")

		err := sdk.SendText("集成测试文本消息")
		if err != nil {
			t.Fatalf("发送文本消息失败: %v", err)
		}

		if lastRequest == nil {
			t.Fatal("没有收到请求")
		}

		if lastRequest.MsgType != "text" {
			t.Errorf("消息类型错误: %v, 期望: text", lastRequest.MsgType)
		}

		if lastRequest.Timestamp == "" {
			t.Error("时间戳不应为空")
		}

		if lastRequest.Sign == "" {
			t.Error("签名不应为空")
		}

		content, ok := lastRequest.Content.(map[string]interface{})
		if !ok {
			t.Fatal("内容应为map类型")
		}

		if content["text"] != "集成测试文本消息" {
			t.Errorf("文本内容错误: %v", content["text"])
		}
	})

	t.Run("完整的富文本消息流程", func(t *testing.T) {
		sdk := New(server.URL)

		content := [][]RichTextElement{
			{
				CreateRichTextElement("text", "这是一条富文本消息\n"),
			},
			{
				CreateRichTextElement("text", "包含链接: "),
				CreateRichTextElement("a", "飞书官网", map[string]string{
					"href": "https://www.feishu.cn",
				}),
			},
			{
				CreateRichTextElement("text", "还有@用户: "),
				CreateRichTextElement("at", "@张三", map[string]string{
					"user_id":   "ou_123456",
					"user_name": "张三",
				}),
			},
		}

		err := sdk.SendRichText("集成测试富文本", content)
		if err != nil {
			t.Fatalf("发送富文本消息失败: %v", err)
		}

		if lastRequest.MsgType != "post" {
			t.Errorf("消息类型错误: %v, 期望: post", lastRequest.MsgType)
		}

		if lastRequest.Timestamp != "" {
			t.Error("无签名模式下时间戳应为空")
		}

		if lastRequest.Sign != "" {
			t.Error("无签名模式下签名应为空")
		}
	})

	t.Run("完整的图片消息流程", func(t *testing.T) {
		sdk := New(server.URL, "test-secret")

		err := sdk.SendImage("img_v2_integration_test_key")
		if err != nil {
			t.Fatalf("发送图片消息失败: %v", err)
		}

		if lastRequest.MsgType != "image" {
			t.Errorf("消息类型错误: %v, 期望: image", lastRequest.MsgType)
		}

		content, ok := lastRequest.Content.(map[string]interface{})
		if !ok {
			t.Fatal("内容应为map类型")
		}

		if content["image_key"] != "img_v2_integration_test_key" {
			t.Errorf("图片键错误: %v", content["image_key"])
		}
	})

	t.Run("完整的交互式卡片流程", func(t *testing.T) {
		sdk := New(server.URL)

		header := CreateCardHeader("集成测试卡片", "blue")
		config := CreateCardConfig(true)
		elements := []interface{}{
			map[string]interface{}{
				"tag": "div",
				"text": map[string]interface{}{
					"content": "这是卡片内容",
					"tag":     "plain_text",
				},
			},
			map[string]interface{}{
				"tag": "action",
				"actions": []interface{}{
					map[string]interface{}{
						"tag": "button",
						"text": map[string]interface{}{
							"content": "点击按钮",
							"tag":     "plain_text",
						},
						"value": map[string]interface{}{
							"key": "button_click",
						},
						"type": "primary",
					},
				},
			},
		}

		err := sdk.SendInteractive(config, header, elements)
		if err != nil {
			t.Fatalf("发送交互式卡片失败: %v", err)
		}

		if lastRequest.MsgType != "interactive" {
			t.Errorf("消息类型错误: %v, 期望: interactive", lastRequest.MsgType)
		}

		content, ok := lastRequest.Content.(map[string]interface{})
		if !ok {
			t.Fatal("内容应为map类型")
		}

		if content["config"] == nil {
			t.Error("配置不应为空")
		}

		if content["header"] == nil {
			t.Error("头部不应为空")
		}

		if content["elements"] == nil {
			t.Error("元素不应为空")
		}
	})
}

func TestErrorHandlingIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	t.Run("服务器错误响应", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
			w.Write([]byte("Internal Server Error"))
		}))
		defer server.Close()

		sdk := New(server.URL)
		err := sdk.SendText("测试错误处理")

		if err == nil {
			t.Error("应该返回错误")
		}

		if err.Error() == "" {
			t.Error("错误信息不应为空")
		}
	})

	t.Run("飞书API错误", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"code": 19001, "msg": "param invalid"}`))
		}))
		defer server.Close()

		sdk := New(server.URL)
		err := sdk.SendText("测试API错误")

		if err == nil {
			t.Error("应该返回错误")
		}

		expectedMsg := "feishu webhook error"
		if err.Error()[:len(expectedMsg)] != expectedMsg {
			t.Errorf("错误信息应包含 '%s', 实际: %v", expectedMsg, err.Error())
		}
	})

	t.Run("无效JSON响应", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`invalid json`))
		}))
		defer server.Close()

		sdk := New(server.URL)
		err := sdk.SendText("测试无效JSON")

		if err == nil {
			t.Error("应该返回错误")
		}
	})
}

func TestRealWorldScenarios(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code": 0, "msg": "success"}`))
	}))
	defer server.Close()

	t.Run("报警通知场景", func(t *testing.T) {
		sdk := New(server.URL, "alert-secret")

		content := [][]RichTextElement{
			{
				CreateRichTextElement("text", "🚨 系统报警通知\n"),
			},
			{
				CreateRichTextElement("text", "服务器: "),
				CreateRichTextElement("text", "web-server-01", map[string]string{}),
			},
			{
				CreateRichTextElement("text", "状态: "),
				CreateRichTextElement("text", "CPU使用率过高 (95%)", map[string]string{}),
			},
			{
				CreateRichTextElement("text", "时间: "),
				CreateRichTextElement("text", time.Now().Format("2006-01-02 15:04:05"), map[string]string{}),
			},
			{
				CreateRichTextElement("text", "处理人: "),
				CreateRichTextElement("at", "@运维小组", map[string]string{
					"user_id": "ou_ops_team",
				}),
			},
		}

		err := sdk.SendRichText("系统报警", content)
		if err != nil {
			t.Fatalf("发送报警通知失败: %v", err)
		}
	})

	t.Run("部署通知场景", func(t *testing.T) {
		sdk := New(server.URL)

		header := CreateCardHeader("部署完成通知", "green")
		config := CreateCardConfig(false)
		elements := []interface{}{
			map[string]interface{}{
				"tag": "div",
				"text": map[string]interface{}{
					"content": "✅ 应用部署成功",
					"tag":     "plain_text",
				},
			},
			map[string]interface{}{
				"tag": "div",
				"fields": []interface{}{
					map[string]interface{}{
						"is_short": true,
						"text": map[string]interface{}{
							"content": "**项目名称:**\nmyapp-v2.1.0",
							"tag":     "lark_md",
						},
					},
					map[string]interface{}{
						"is_short": true,
						"text": map[string]interface{}{
							"content": "**环境:**\nproduction",
							"tag":     "lark_md",
						},
					},
				},
			},
			map[string]interface{}{
				"tag": "action",
				"actions": []interface{}{
					map[string]interface{}{
						"tag": "button",
						"text": map[string]interface{}{
							"content": "查看详情",
							"tag":     "plain_text",
						},
						"url":  "https://deploy.example.com/logs/12345",
						"type": "default",
					},
				},
			},
		}

		err := sdk.SendInteractive(config, header, elements)
		if err != nil {
			t.Fatalf("发送部署通知失败: %v", err)
		}
	})

	t.Run("日程提醒场景", func(t *testing.T) {
		err := SendTextMessage(server.URL, "📅 会议提醒：项目评审会议将在15分钟后开始，请相关同事准时参加。", "meeting-secret")
		if err != nil {
			t.Fatalf("发送会议提醒失败: %v", err)
		}
	})
}

func TestConcurrentRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 模拟网络延迟
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code": 0, "msg": "success"}`))
	}))
	defer server.Close()

	sdk := New(server.URL)

	t.Run("并发发送消息", func(t *testing.T) {
		concurrency := 10
		done := make(chan error, concurrency)

		for i := 0; i < concurrency; i++ {
			go func(id int) {
				err := sdk.SendText("并发测试消息 " + string(rune('A'+id)))
				done <- err
			}(i)
		}

		for i := 0; i < concurrency; i++ {
			if err := <-done; err != nil {
				t.Errorf("并发请求 %d 失败: %v", i, err)
			}
		}
	})
}
