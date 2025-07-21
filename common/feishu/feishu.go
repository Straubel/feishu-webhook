package feishu

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	WebhookURL string
	Secret     string
	client     *resty.Client
}

type WebhookRequest struct {
	Timestamp string      `json:"timestamp,omitempty"`
	Sign      string      `json:"sign,omitempty"`
	MsgType   string      `json:"msg_type"`
	Content   interface{} `json:"content"`
}

func NewClient(webhookURL string, secret ...string) *Client {
	client := &Client{
		WebhookURL: webhookURL,
		client:     resty.New(),
	}

	if len(secret) > 0 {
		client.Secret = secret[0]
	}

	return client
}

func (c *Client) SendMessage(message *Message) error {
	if c.Secret != "" {
		return c.sendMessageWithSign(message)
	}
	return c.sendMessageWithoutSign(message)
}

func (c *Client) SendText(text string) error {
	message := NewTextMessage(text)
	return c.SendMessage(message)
}

func (c *Client) SendRichText(title string, content [][]RichTextElement) error {
	message := NewRichTextMessage(title, content)
	return c.SendMessage(message)
}

func (c *Client) SendImage(imageKey string) error {
	message := NewImageMessage(imageKey)
	return c.SendMessage(message)
}

func (c *Client) SendInteractive(config *CardConfig, header *CardHeader, elements []interface{}) error {
	message := NewInteractiveMessage(config, header, elements)
	return c.SendMessage(message)
}

func (c *Client) sendMessageWithSign(message *Message) error {
	timestamp := time.Now().Unix()
	sign, err := GenSign(c.Secret, timestamp)
	if err != nil {
		return fmt.Errorf("generate sign failed: %w", err)
	}

	request := &WebhookRequest{
		Timestamp: fmt.Sprintf("%d", timestamp),
		Sign:      sign,
		MsgType:   string(message.MsgType),
		Content:   message.Content,
	}

	return c.sendRequest(request)
}

func (c *Client) sendMessageWithoutSign(message *Message) error {
	request := &WebhookRequest{
		MsgType: string(message.MsgType),
		Content: message.Content,
	}

	return c.sendRequest(request)
}

func (c *Client) sendRequest(request *WebhookRequest) error {
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post(c.WebhookURL)

	if err != nil {
		return fmt.Errorf("send request failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("request failed with status: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("parse response failed: %w", err)
	}

	if code, ok := result["code"].(float64); ok && code != 0 {
		return fmt.Errorf("feishu webhook error: code=%v, msg=%v", result["code"], result["msg"])
	}

	return nil
}
