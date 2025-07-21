package feishu

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSDK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code": 0, "msg": "success"}`))
	}))
	defer server.Close()

	t.Run("创建SDK实例", func(t *testing.T) {
		sdk := New(server.URL)
		if sdk == nil {
			t.Error("SDK should not be nil")
		}
		if sdk.client == nil {
			t.Error("SDK client should not be nil")
		}
		if sdk.client.WebhookURL != server.URL {
			t.Errorf("WebhookURL = %v, want %v", sdk.client.WebhookURL, server.URL)
		}
	})

	t.Run("创建带密钥的SDK实例", func(t *testing.T) {
		secret := "test-secret"
		sdk := New(server.URL, secret)
		if sdk.client.Secret != secret {
			t.Errorf("Secret = %v, want %v", sdk.client.Secret, secret)
		}
	})

	t.Run("SDK发送文本消息", func(t *testing.T) {
		sdk := New(server.URL)
		err := sdk.SendText("Test message")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("SDK发送富文本消息", func(t *testing.T) {
		sdk := New(server.URL)
		content := [][]RichTextElement{
			{CreateRichTextElement("text", "Test content")},
		}
		err := sdk.SendRichText("Test Title", content)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("SDK发送图片消息", func(t *testing.T) {
		sdk := New(server.URL)
		err := sdk.SendImage("test_image_key")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("SDK发送交互式消息", func(t *testing.T) {
		sdk := New(server.URL)
		config := CreateCardConfig(true)
		header := CreateCardHeader("Test Card", "blue")
		elements := []interface{}{
			map[string]interface{}{
				"tag": "div",
				"text": map[string]interface{}{
					"content": "Test content",
					"tag":     "plain_text",
				},
			},
		}
		err := sdk.SendInteractive(config, header, elements)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("获取客户端实例", func(t *testing.T) {
		sdk := New(server.URL)
		client := sdk.Client()
		if client != sdk.client {
			t.Error("Client() should return the same client instance")
		}
	})
}

func TestConvenienceFunctionsDetailed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code": 0, "msg": "success"}`))
	}))
	defer server.Close()

	t.Run("SendTextMessage不带密钥", func(t *testing.T) {
		err := SendTextMessage(server.URL, "Test message")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("SendTextMessage带密钥", func(t *testing.T) {
		err := SendTextMessage(server.URL, "Test message", "secret")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("SendRichTextMessage不带密钥", func(t *testing.T) {
		content := [][]RichTextElement{
			{CreateRichTextElement("text", "Test")},
		}
		err := SendRichTextMessage(server.URL, "Title", content)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("SendRichTextMessage带密钥", func(t *testing.T) {
		content := [][]RichTextElement{
			{CreateRichTextElement("text", "Test")},
		}
		err := SendRichTextMessage(server.URL, "Title", content, "secret")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("SendImageMessage不带密钥", func(t *testing.T) {
		err := SendImageMessage(server.URL, "image_key")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("SendImageMessage带密钥", func(t *testing.T) {
		err := SendImageMessage(server.URL, "image_key", "secret")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func TestCreateRichTextElementExtended(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		text     string
		options  []map[string]string
		expected RichTextElement
	}{
		{
			name: "无选项的文本元素",
			tag:  "text",
			text: "Hello",
			expected: RichTextElement{
				Tag:  "text",
				Text: "Hello",
			},
		},
		{
			name: "带href的链接元素",
			tag:  "a",
			text: "Click here",
			options: []map[string]string{
				{"href": "https://example.com"},
			},
			expected: RichTextElement{
				Tag:  "a",
				Text: "Click here",
				Href: "https://example.com",
			},
		},
		{
			name: "带用户信息的@元素",
			tag:  "at",
			text: "@用户",
			options: []map[string]string{
				{
					"user_id":   "ou_123456",
					"user_name": "张三",
				},
			},
			expected: RichTextElement{
				Tag:      "at",
				Text:     "@用户",
				UserId:   "ou_123456",
				UserName: "张三",
			},
		},
		{
			name: "带图片键的图片元素",
			tag:  "img",
			text: "",
			options: []map[string]string{
				{"image_key": "img_v2_test_key"},
			},
			expected: RichTextElement{
				Tag:      "img",
				Text:     "",
				ImageKey: "img_v2_test_key",
			},
		},
		{
			name: "包含所有属性的元素",
			tag:  "a",
			text: "复杂链接",
			options: []map[string]string{
				{
					"href":      "https://example.com",
					"user_id":   "user123",
					"user_name": "测试用户",
					"image_key": "img_key",
				},
			},
			expected: RichTextElement{
				Tag:      "a",
				Text:     "复杂链接",
				Href:     "https://example.com",
				UserId:   "user123",
				UserName: "测试用户",
				ImageKey: "img_key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result RichTextElement
			if len(tt.options) > 0 {
				result = CreateRichTextElement(tt.tag, tt.text, tt.options[0])
			} else {
				result = CreateRichTextElement(tt.tag, tt.text)
			}

			if result.Tag != tt.expected.Tag {
				t.Errorf("Tag = %v, want %v", result.Tag, tt.expected.Tag)
			}
			if result.Text != tt.expected.Text {
				t.Errorf("Text = %v, want %v", result.Text, tt.expected.Text)
			}
			if result.Href != tt.expected.Href {
				t.Errorf("Href = %v, want %v", result.Href, tt.expected.Href)
			}
			if result.UserId != tt.expected.UserId {
				t.Errorf("UserId = %v, want %v", result.UserId, tt.expected.UserId)
			}
			if result.UserName != tt.expected.UserName {
				t.Errorf("UserName = %v, want %v", result.UserName, tt.expected.UserName)
			}
			if result.ImageKey != tt.expected.ImageKey {
				t.Errorf("ImageKey = %v, want %v", result.ImageKey, tt.expected.ImageKey)
			}
		})
	}
}

func TestCreateCardHeaderExtended(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		template string
		expected *CardHeader
	}{
		{
			name:     "蓝色模板",
			title:    "蓝色卡片",
			template: "blue",
			expected: &CardHeader{
				Title: &CardHeaderTitle{
					Content: "蓝色卡片",
					Tag:     "plain_text",
				},
				Template: "blue",
			},
		},
		{
			name:     "红色模板",
			title:    "警告卡片",
			template: "red",
			expected: &CardHeader{
				Title: &CardHeaderTitle{
					Content: "警告卡片",
					Tag:     "plain_text",
				},
				Template: "red",
			},
		},
		{
			name:     "空模板",
			title:    "默认卡片",
			template: "",
			expected: &CardHeader{
				Title: &CardHeaderTitle{
					Content: "默认卡片",
					Tag:     "plain_text",
				},
				Template: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateCardHeader(tt.title, tt.template)

			if result.Title.Content != tt.expected.Title.Content {
				t.Errorf("Title.Content = %v, want %v", result.Title.Content, tt.expected.Title.Content)
			}
			if result.Title.Tag != tt.expected.Title.Tag {
				t.Errorf("Title.Tag = %v, want %v", result.Title.Tag, tt.expected.Title.Tag)
			}
			if result.Template != tt.expected.Template {
				t.Errorf("Template = %v, want %v", result.Template, tt.expected.Template)
			}
		})
	}
}

func TestCreateCardConfigExtended(t *testing.T) {
	tests := []struct {
		name          string
		enableForward bool
		expected      *CardConfig
	}{
		{
			name:          "启用转发",
			enableForward: true,
			expected:      &CardConfig{EnableForward: true},
		},
		{
			name:          "禁用转发",
			enableForward: false,
			expected:      &CardConfig{EnableForward: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateCardConfig(tt.enableForward)

			if result.EnableForward != tt.expected.EnableForward {
				t.Errorf("EnableForward = %v, want %v", result.EnableForward, tt.expected.EnableForward)
			}
		})
	}
}
