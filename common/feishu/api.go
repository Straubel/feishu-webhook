package feishu

type SDK struct {
	client *Client
}

func New(webhookURL string, secret ...string) *SDK {
	return &SDK{
		client: NewClient(webhookURL, secret...),
	}
}

func (sdk *SDK) SendText(text string) error {
	return sdk.client.SendText(text)
}

func (sdk *SDK) SendRichText(title string, content [][]RichTextElement) error {
	return sdk.client.SendRichText(title, content)
}

func (sdk *SDK) SendImage(imageKey string) error {
	return sdk.client.SendImage(imageKey)
}

func (sdk *SDK) SendInteractive(config *CardConfig, header *CardHeader, elements []interface{}) error {
	return sdk.client.SendInteractive(config, header, elements)
}

func (sdk *SDK) SendMessage(message *Message) error {
	return sdk.client.SendMessage(message)
}

func (sdk *SDK) Client() *Client {
	return sdk.client
}

func SendTextMessage(webhookURL, text string, secret ...string) error {
	client := NewClient(webhookURL, secret...)
	return client.SendText(text)
}

func SendRichTextMessage(webhookURL, title string, content [][]RichTextElement, secret ...string) error {
	client := NewClient(webhookURL, secret...)
	return client.SendRichText(title, content)
}

func SendImageMessage(webhookURL, imageKey string, secret ...string) error {
	client := NewClient(webhookURL, secret...)
	return client.SendImage(imageKey)
}

func CreateRichTextElement(tag, text string, options ...map[string]string) RichTextElement {
	element := RichTextElement{
		Tag:  tag,
		Text: text,
	}
	
	if len(options) > 0 {
		opts := options[0]
		if href, ok := opts["href"]; ok {
			element.Href = href
		}
		if userId, ok := opts["user_id"]; ok {
			element.UserId = userId
		}
		if userName, ok := opts["user_name"]; ok {
			element.UserName = userName
		}
		if imageKey, ok := opts["image_key"]; ok {
			element.ImageKey = imageKey
		}
	}
	
	return element
}

func CreateCardHeader(title, template string) *CardHeader {
	return &CardHeader{
		Title: &CardHeaderTitle{
			Content: title,
			Tag:     "plain_text",
		},
		Template: template,
	}
}

func CreateCardConfig(enableForward bool) *CardConfig {
	return &CardConfig{
		EnableForward: enableForward,
	}
}