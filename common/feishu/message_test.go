package feishu

import (
	"encoding/json"
	"testing"
)

func TestMessageTypes(t *testing.T) {
	tests := []struct {
		name     string
		msgType  MessageType
		expected string
	}{
		{"文本消息类型", MessageTypeText, "text"},
		{"富文本消息类型", MessageTypeRichText, "post"},
		{"交互式消息类型", MessageTypeInteractive, "interactive"},
		{"分享群聊类型", MessageTypeShareChat, "share_chat"},
		{"图片消息类型", MessageTypeImage, "image"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.msgType) != tt.expected {
				t.Errorf("MessageType = %v, want %v", tt.msgType, tt.expected)
			}
		})
	}
}

func TestNewTextMessage(t *testing.T) {
	text := "Hello, World!"
	msg := NewTextMessage(text)

	if msg.MsgType != MessageTypeText {
		t.Errorf("MsgType = %v, want %v", msg.MsgType, MessageTypeText)
	}

	content, ok := msg.Content.(*TextContent)
	if !ok {
		t.Fatalf("Content should be *TextContent, got %T", msg.Content)
	}

	if content.Text != text {
		t.Errorf("Text = %v, want %v", content.Text, text)
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if unmarshaled["msg_type"] != "text" {
		t.Errorf("JSON msg_type = %v, want %v", unmarshaled["msg_type"], "text")
	}
}

func TestNewRichTextMessage(t *testing.T) {
	title := "Test Title"
	content := [][]RichTextElement{
		{
			{Tag: "text", Text: "Hello "},
			{Tag: "a", Text: "World", Href: "https://example.com"},
		},
		{
			{Tag: "text", Text: "New line"},
		},
	}

	msg := NewRichTextMessage(title, content)

	if msg.MsgType != MessageTypeRichText {
		t.Errorf("MsgType = %v, want %v", msg.MsgType, MessageTypeRichText)
	}

	richContent, ok := msg.Content.(*RichTextContent)
	if !ok {
		t.Fatalf("Content should be *RichTextContent, got %T", msg.Content)
	}

	if richContent.Post == nil {
		t.Fatal("Post should not be nil")
	}

	if richContent.Post.ZhCn == nil {
		t.Fatal("ZhCn should not be nil")
	}

	if richContent.Post.ZhCn.Title != title {
		t.Errorf("Title = %v, want %v", richContent.Post.ZhCn.Title, title)
	}

	if len(richContent.Post.ZhCn.Content) != 2 {
		t.Errorf("Content length = %v, want %v", len(richContent.Post.ZhCn.Content), 2)
	}

	if len(richContent.Post.ZhCn.Content[0]) != 2 {
		t.Errorf("First line elements = %v, want %v", len(richContent.Post.ZhCn.Content[0]), 2)
	}

	if richContent.Post.ZhCn.Content[0][0].Tag != "text" {
		t.Errorf("First element tag = %v, want %v", richContent.Post.ZhCn.Content[0][0].Tag, "text")
	}
	if richContent.Post.ZhCn.Content[0][0].Text != "Hello " {
		t.Errorf("First element text = %v, want %v", richContent.Post.ZhCn.Content[0][0].Text, "Hello ")
	}

	if richContent.Post.ZhCn.Content[0][1].Tag != "a" {
		t.Errorf("Second element tag = %v, want %v", richContent.Post.ZhCn.Content[0][1].Tag, "a")
	}
	if richContent.Post.ZhCn.Content[0][1].Href != "https://example.com" {
		t.Errorf("Second element href = %v, want %v", richContent.Post.ZhCn.Content[0][1].Href, "https://example.com")
	}
}

func TestNewInteractiveMessage(t *testing.T) {
	config := &CardConfig{EnableForward: true}
	header := &CardHeader{
		Title: &CardHeaderTitle{
			Content: "Test Card",
			Tag:     "plain_text",
		},
		Template: "blue",
	}
	elements := []interface{}{
		map[string]interface{}{
			"tag": "div",
			"text": map[string]interface{}{
				"content": "Card content",
				"tag":     "plain_text",
			},
		},
	}

	msg := NewInteractiveMessage(config, header, elements)

	if msg.MsgType != MessageTypeInteractive {
		t.Errorf("MsgType = %v, want %v", msg.MsgType, MessageTypeInteractive)
	}

	interactiveContent, ok := msg.Content.(*InteractiveContent)
	if !ok {
		t.Fatalf("Content should be *InteractiveContent, got %T", msg.Content)
	}

	if interactiveContent.Config == nil {
		t.Fatal("Config should not be nil")
	}

	if !interactiveContent.Config.EnableForward {
		t.Error("EnableForward should be true")
	}

	if interactiveContent.Header == nil {
		t.Fatal("Header should not be nil")
	}

	if interactiveContent.Header.Title.Content != "Test Card" {
		t.Errorf("Header title = %v, want %v", interactiveContent.Header.Title.Content, "Test Card")
	}

	if len(interactiveContent.Elements) != 1 {
		t.Errorf("Elements length = %v, want %v", len(interactiveContent.Elements), 1)
	}
}

func TestNewImageMessage(t *testing.T) {
	imageKey := "img_v2_test_key"
	msg := NewImageMessage(imageKey)

	if msg.MsgType != MessageTypeImage {
		t.Errorf("MsgType = %v, want %v", msg.MsgType, MessageTypeImage)
	}

	content, ok := msg.Content.(*ImageContent)
	if !ok {
		t.Fatalf("Content should be *ImageContent, got %T", msg.Content)
	}

	if content.ImageKey != imageKey {
		t.Errorf("ImageKey = %v, want %v", content.ImageKey, imageKey)
	}
}

func TestNewShareChatMessage(t *testing.T) {
	shareChatId := "oc_test_chat_id"
	msg := NewShareChatMessage(shareChatId)

	if msg.MsgType != MessageTypeShareChat {
		t.Errorf("MsgType = %v, want %v", msg.MsgType, MessageTypeShareChat)
	}

	content, ok := msg.Content.(*ShareChatContent)
	if !ok {
		t.Fatalf("Content should be *ShareChatContent, got %T", msg.Content)
	}

	if content.ShareChatId != shareChatId {
		t.Errorf("ShareChatId = %v, want %v", content.ShareChatId, shareChatId)
	}
}

func TestRichTextElementSerialization(t *testing.T) {
	element := RichTextElement{
		Tag:      "a",
		Text:     "Link Text",
		Href:     "https://example.com",
		UserId:   "user123",
		UserName: "Test User",
		ImageKey: "img_key",
	}

	jsonData, err := json.Marshal(element)
	if err != nil {
		t.Fatalf("Failed to marshal RichTextElement: %v", err)
	}

	var unmarshaled RichTextElement
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal RichTextElement: %v", err)
	}

	if unmarshaled.Tag != element.Tag {
		t.Errorf("Tag = %v, want %v", unmarshaled.Tag, element.Tag)
	}
	if unmarshaled.Text != element.Text {
		t.Errorf("Text = %v, want %v", unmarshaled.Text, element.Text)
	}
	if unmarshaled.Href != element.Href {
		t.Errorf("Href = %v, want %v", unmarshaled.Href, element.Href)
	}
	if unmarshaled.UserId != element.UserId {
		t.Errorf("UserId = %v, want %v", unmarshaled.UserId, element.UserId)
	}
	if unmarshaled.UserName != element.UserName {
		t.Errorf("UserName = %v, want %v", unmarshaled.UserName, element.UserName)
	}
	if unmarshaled.ImageKey != element.ImageKey {
		t.Errorf("ImageKey = %v, want %v", unmarshaled.ImageKey, element.ImageKey)
	}
}

func TestMessageSerialization(t *testing.T) {
	msg := NewTextMessage("Test message")

	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	var unmarshaled Message
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	if unmarshaled.MsgType != msg.MsgType {
		t.Errorf("MsgType = %v, want %v", unmarshaled.MsgType, msg.MsgType)
	}

	var expected map[string]interface{}
	var actual map[string]interface{}

	expectedJSON, _ := json.Marshal(msg.Content)
	actualJSON, _ := json.Marshal(unmarshaled.Content)

	json.Unmarshal(expectedJSON, &expected)
	json.Unmarshal(actualJSON, &actual)

	expectedText, _ := expected["text"].(string)
	actualText, _ := actual["text"].(string)

	if expectedText != actualText {
		t.Errorf("Content text = %v, want %v", actualText, expectedText)
	}
}

func TestComplexRichTextMessage(t *testing.T) {
	content := [][]RichTextElement{
		{
			CreateRichTextElement("text", "普通文本 "),
			CreateRichTextElement("a", "链接", map[string]string{"href": "https://example.com"}),
			CreateRichTextElement("text", " 和 "),
			CreateRichTextElement("at", "@用户", map[string]string{
				"user_id":   "user123",
				"user_name": "张三",
			}),
		},
		{
			CreateRichTextElement("text", "第二行文本"),
		},
		{
			CreateRichTextElement("text", "第三行包含图片 "),
			CreateRichTextElement("img", "", map[string]string{"image_key": "img_key_123"}),
		},
	}

	msg := NewRichTextMessage("复杂富文本测试", content)

	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal complex rich text message: %v", err)
	}

	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if unmarshaled["msg_type"] != "post" {
		t.Errorf("msg_type = %v, want post", unmarshaled["msg_type"])
	}

	contentMap, ok := unmarshaled["content"].(map[string]interface{})
	if !ok {
		t.Fatal("content should be a map")
	}

	postMap, ok := contentMap["post"].(map[string]interface{})
	if !ok {
		t.Fatal("post should be a map")
	}

	zhCnMap, ok := postMap["zh_cn"].(map[string]interface{})
	if !ok {
		t.Fatal("zh_cn should be a map")
	}

	if zhCnMap["title"] != "复杂富文本测试" {
		t.Errorf("title = %v, want 复杂富文本测试", zhCnMap["title"])
	}

	contentArray, ok := zhCnMap["content"].([]interface{})
	if !ok {
		t.Fatal("content should be an array")
	}

	if len(contentArray) != 3 {
		t.Errorf("content array length = %v, want 3", len(contentArray))
	}
}
