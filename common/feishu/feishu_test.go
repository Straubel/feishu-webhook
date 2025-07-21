package feishu

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name       string
		webhookURL string
		secret     []string
		wantSecret string
	}{
		{
			name:       "客户端不带签名",
			webhookURL: "https://example.com/webhook",
			secret:     []string{},
			wantSecret: "",
		},
		{
			name:       "客户端带签名",
			webhookURL: "https://example.com/webhook",
			secret:     []string{"test-secret"},
			wantSecret: "test-secret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.webhookURL, tt.secret...)
			if client.WebhookURL != tt.webhookURL {
				t.Errorf("WebhookURL = %v, want %v", client.WebhookURL, tt.webhookURL)
			}
			if client.Secret != tt.wantSecret {
				t.Errorf("Secret = %v, want %v", client.Secret, tt.wantSecret)
			}
			if client.client == nil {
				t.Error("HTTP client should not be nil")
			}
		})
	}
}

func setupMockServer(t *testing.T, responseCode int, responseBody map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", contentType)
		}

		w.WriteHeader(responseCode)
		w.Header().Set("Content-Type", "application/json")

		responseJSON, _ := json.Marshal(responseBody)
		w.Write(responseJSON)
	}))
}

func TestSendText(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		secret       string
		responseCode int
		responseBody map[string]interface{}
		expectError  bool
	}{
		{
			name:         "发送文本消息成功",
			text:         "Hello, World!",
			secret:       "",
			responseCode: 200,
			responseBody: map[string]interface{}{"code": 0, "msg": "success"},
			expectError:  false,
		},
		{
			name:         "发送文本消息失败",
			text:         "Hello, World!",
			secret:       "",
			responseCode: 200,
			responseBody: map[string]interface{}{"code": 1, "msg": "error"},
			expectError:  true,
		},
		{
			name:         "HTTP状态码错误",
			text:         "Hello, World!",
			secret:       "",
			responseCode: 400,
			responseBody: map[string]interface{}{"error": "bad request"},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupMockServer(t, tt.responseCode, tt.responseBody)
			defer server.Close()

			client := NewClient(server.URL, tt.secret)
			err := client.SendText(tt.text)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestSendTextWithSign(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request WebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			return
		}

		if request.Timestamp == "" {
			t.Error("Timestamp should not be empty for signed request")
		}
		if request.Sign == "" {
			t.Error("Sign should not be empty for signed request")
		}
		if request.MsgType != "text" {
			t.Errorf("Expected msg_type text, got %s", request.MsgType)
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code": 0, "msg": "success"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-secret")
	err := client.SendText("Test message with signature")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestSendRichText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request WebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			return
		}

		if request.MsgType != "post" {
			t.Errorf("Expected msg_type post, got %s", request.MsgType)
		}

		content, ok := request.Content.(*RichTextContent)
		if !ok {
			t.Error("Content should be RichTextContent")
			return
		}

		if content.Post == nil || content.Post.ZhCn == nil {
			t.Error("RichText content structure is invalid")
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code": 0, "msg": "success"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)

	richTextContent := [][]RichTextElement{
		{
			CreateRichTextElement("text", "This is a test\n"),
		},
		{
			CreateRichTextElement("text", "Link: "),
			CreateRichTextElement("a", "GitHub", map[string]string{"href": "https://github.com"}),
		},
	}

	err := client.SendRichText("Test Title", richTextContent)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestSendImage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request WebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			return
		}

		if request.MsgType != "image" {
			t.Errorf("Expected msg_type image, got %s", request.MsgType)
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code": 0, "msg": "success"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.SendImage("img_v2_test_image_key")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestSendInteractive(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request WebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			return
		}

		if request.MsgType != "interactive" {
			t.Errorf("Expected msg_type interactive, got %s", request.MsgType)
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code": 0, "msg": "success"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)

	header := CreateCardHeader("Test Card", "blue")
	config := CreateCardConfig(true)
	elements := []interface{}{
		map[string]interface{}{
			"tag": "div",
			"text": map[string]interface{}{
				"content": "Test card content",
				"tag":     "plain_text",
			},
		},
	}

	err := client.SendInteractive(config, header, elements)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestGenSign(t *testing.T) {
	tests := []struct {
		name      string
		secret    string
		timestamp int64
		wantErr   bool
	}{
		{
			name:      "正常签名生成",
			secret:    "test-secret",
			timestamp: time.Now().Unix(),
			wantErr:   false,
		},
		{
			name:      "空密钥签名",
			secret:    "",
			timestamp: time.Now().Unix(),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sign, err := GenSign(tt.secret, tt.timestamp)

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.wantErr && sign == "" {
				t.Error("Sign should not be empty")
			}
		})
	}
}

func TestCreateRichTextElement(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		text     string
		options  map[string]string
		expected RichTextElement
	}{
		{
			name: "文本元素",
			tag:  "text",
			text: "Hello",
			expected: RichTextElement{
				Tag:  "text",
				Text: "Hello",
			},
		},
		{
			name: "链接元素",
			tag:  "a",
			text: "GitHub",
			options: map[string]string{
				"href": "https://github.com",
			},
			expected: RichTextElement{
				Tag:  "a",
				Text: "GitHub",
				Href: "https://github.com",
			},
		},
		{
			name: "@用户元素",
			tag:  "at",
			text: "张三",
			options: map[string]string{
				"user_id":   "user123",
				"user_name": "张三",
			},
			expected: RichTextElement{
				Tag:      "at",
				Text:     "张三",
				UserId:   "user123",
				UserName: "张三",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result RichTextElement
			if tt.options != nil {
				result = CreateRichTextElement(tt.tag, tt.text, tt.options)
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
		})
	}
}

func TestCreateCardHeader(t *testing.T) {
	header := CreateCardHeader("Test Title", "blue")

	if header.Title.Content != "Test Title" {
		t.Errorf("Title content = %v, want %v", header.Title.Content, "Test Title")
	}
	if header.Title.Tag != "plain_text" {
		t.Errorf("Title tag = %v, want %v", header.Title.Tag, "plain_text")
	}
	if header.Template != "blue" {
		t.Errorf("Template = %v, want %v", header.Template, "blue")
	}
}

func TestCreateCardConfig(t *testing.T) {
	config := CreateCardConfig(true)

	if !config.EnableForward {
		t.Error("EnableForward should be true")
	}

	config2 := CreateCardConfig(false)
	if config2.EnableForward {
		t.Error("EnableForward should be false")
	}
}

func TestMessageConstructors(t *testing.T) {
	t.Run("NewTextMessage", func(t *testing.T) {
		msg := NewTextMessage("test")
		if msg.MsgType != MessageTypeText {
			t.Errorf("MsgType = %v, want %v", msg.MsgType, MessageTypeText)
		}

		content, ok := msg.Content.(*TextContent)
		if !ok {
			t.Error("Content should be *TextContent")
		}
		if content.Text != "test" {
			t.Errorf("Text = %v, want %v", content.Text, "test")
		}
	})

	t.Run("NewImageMessage", func(t *testing.T) {
		msg := NewImageMessage("test_key")
		if msg.MsgType != MessageTypeImage {
			t.Errorf("MsgType = %v, want %v", msg.MsgType, MessageTypeImage)
		}

		content, ok := msg.Content.(*ImageContent)
		if !ok {
			t.Error("Content should be *ImageContent")
		}
		if content.ImageKey != "test_key" {
			t.Errorf("ImageKey = %v, want %v", content.ImageKey, "test_key")
		}
	})

	t.Run("NewShareChatMessage", func(t *testing.T) {
		msg := NewShareChatMessage("chat_id")
		if msg.MsgType != MessageTypeShareChat {
			t.Errorf("MsgType = %v, want %v", msg.MsgType, MessageTypeShareChat)
		}

		content, ok := msg.Content.(*ShareChatContent)
		if !ok {
			t.Error("Content should be *ShareChatContent")
		}
		if content.ShareChatId != "chat_id" {
			t.Errorf("ShareChatId = %v, want %v", content.ShareChatId, "chat_id")
		}
	})
}

func TestConvenienceFunctions(t *testing.T) {
	server := setupMockServer(t, 200, map[string]interface{}{"code": 0, "msg": "success"})
	defer server.Close()

	t.Run("SendTextMessage", func(t *testing.T) {
		err := SendTextMessage(server.URL, "test message")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("SendTextMessageWithSecret", func(t *testing.T) {
		err := SendTextMessage(server.URL, "test message", "secret")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("SendImageMessage", func(t *testing.T) {
		err := SendImageMessage(server.URL, "image_key")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func TestErrorHandling(t *testing.T) {
	t.Run("无效的服务器地址", func(t *testing.T) {
		client := NewClient("http://invalid-url-that-does-not-exist.com")
		err := client.SendText("test")
		if err == nil {
			t.Error("Expected error for invalid URL")
		}
		if !strings.Contains(err.Error(), "send request failed") {
			t.Errorf("Error message should contain 'send request failed', got: %v", err)
		}
	})

	t.Run("飞书API错误响应", func(t *testing.T) {
		server := setupMockServer(t, 200, map[string]interface{}{
			"code": 19001,
			"msg":  "param invalid",
		})
		defer server.Close()

		client := NewClient(server.URL)
		err := client.SendText("test")
		if err == nil {
			t.Error("Expected error for API error response")
		}
		if !strings.Contains(err.Error(), "feishu webhook error") {
			t.Errorf("Error message should contain 'feishu webhook error', got: %v", err)
		}
	})
}
