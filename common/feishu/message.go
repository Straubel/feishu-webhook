package feishu

type MessageType string

const (
	MessageTypeText        MessageType = "text"
	MessageTypeRichText    MessageType = "post"
	MessageTypeInteractive MessageType = "interactive"
	MessageTypeShareChat   MessageType = "share_chat"
	MessageTypeImage       MessageType = "image"
)

type Message struct {
	MsgType MessageType `json:"msg_type"`
	Content interface{} `json:"content"`
}

type TextContent struct {
	Text string `json:"text"`
}

type RichTextContent struct {
	Post *Post `json:"post"`
}

type Post struct {
	ZhCn *PostContent `json:"zh_cn,omitempty"`
	EnUs *PostContent `json:"en_us,omitempty"`
}

type PostContent struct {
	Title   string     `json:"title"`
	Content [][]RichTextElement `json:"content"`
}

type RichTextElement struct {
	Tag      string `json:"tag"`
	Text     string `json:"text,omitempty"`
	Href     string `json:"href,omitempty"`
	UserId   string `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
	ImageKey string `json:"image_key,omitempty"`
}

type InteractiveContent struct {
	Config   *CardConfig   `json:"config,omitempty"`
	Elements []interface{} `json:"elements"`
	Header   *CardHeader   `json:"header,omitempty"`
}

type CardConfig struct {
	EnableForward bool `json:"enable_forward"`
}

type CardHeader struct {
	Title    *CardHeaderTitle `json:"title"`
	Template string           `json:"template,omitempty"`
}

type CardHeaderTitle struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type ImageContent struct {
	ImageKey string `json:"image_key"`
}

type ShareChatContent struct {
	ShareChatId string `json:"share_chat_id"`
}

func NewTextMessage(text string) *Message {
	return &Message{
		MsgType: MessageTypeText,
		Content: &TextContent{
			Text: text,
		},
	}
}

func NewRichTextMessage(title string, content [][]RichTextElement) *Message {
	return &Message{
		MsgType: MessageTypeRichText,
		Content: &RichTextContent{
			Post: &Post{
				ZhCn: &PostContent{
					Title:   title,
					Content: content,
				},
			},
		},
	}
}

func NewInteractiveMessage(config *CardConfig, header *CardHeader, elements []interface{}) *Message {
	return &Message{
		MsgType: MessageTypeInteractive,
		Content: &InteractiveContent{
			Config:   config,
			Header:   header,
			Elements: elements,
		},
	}
}

func NewImageMessage(imageKey string) *Message {
	return &Message{
		MsgType: MessageTypeImage,
		Content: &ImageContent{
			ImageKey: imageKey,
		},
	}
}

func NewShareChatMessage(shareChatId string) *Message {
	return &Message{
		MsgType: MessageTypeShareChat,
		Content: &ShareChatContent{
			ShareChatId: shareChatId,
		},
	}
}