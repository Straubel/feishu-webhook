package feishu

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkGenSign(b *testing.B) {
	secret := "test-secret-key-for-benchmark"
	timestamp := int64(1640995200)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GenSign(secret, timestamp)
		if err != nil {
			b.Fatalf("GenSign failed: %v", err)
		}
	}
}

func BenchmarkNewTextMessage(b *testing.B) {
	text := "This is a benchmark test message for performance testing"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg := NewTextMessage(text)
		if msg == nil {
			b.Fatal("NewTextMessage returned nil")
		}
	}
}

func BenchmarkNewRichTextMessage(b *testing.B) {
	title := "Benchmark Test Title"
	content := [][]RichTextElement{
		{
			CreateRichTextElement("text", "Benchmark test content "),
			CreateRichTextElement("a", "link", map[string]string{"href": "https://example.com"}),
		},
		{
			CreateRichTextElement("text", "Second line content"),
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg := NewRichTextMessage(title, content)
		if msg == nil {
			b.Fatal("NewRichTextMessage returned nil")
		}
	}
}

func BenchmarkCreateRichTextElement(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		element := CreateRichTextElement("a", "Benchmark Link", map[string]string{
			"href":      "https://example.com",
			"user_id":   "user123",
			"user_name": "Test User",
		})
		if element.Tag == "" {
			b.Fatal("CreateRichTextElement returned empty tag")
		}
	}
}

func BenchmarkSendTextMessage(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code": 0, "msg": "success"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := client.SendText("Benchmark test message")
		if err != nil {
			b.Fatalf("SendText failed: %v", err)
		}
	}
}

func BenchmarkSendTextMessageWithSign(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code": 0, "msg": "success"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "benchmark-secret")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := client.SendText("Benchmark test message with signature")
		if err != nil {
			b.Fatalf("SendText with sign failed: %v", err)
		}
	}
}

func BenchmarkSendRichTextMessage(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code": 0, "msg": "success"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	content := [][]RichTextElement{
		{
			CreateRichTextElement("text", "Benchmark rich text "),
			CreateRichTextElement("a", "link", map[string]string{"href": "https://example.com"}),
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := client.SendRichText("Benchmark Title", content)
		if err != nil {
			b.Fatalf("SendRichText failed: %v", err)
		}
	}
}

func BenchmarkSDKCreation(b *testing.B) {
	webhookURL := "https://example.com/webhook"
	secret := "test-secret"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sdk := New(webhookURL, secret)
		if sdk == nil {
			b.Fatal("New returned nil")
		}
	}
}

func BenchmarkConvenienceFunction(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code": 0, "msg": "success"}`))
	}))
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := SendTextMessage(server.URL, "Benchmark convenience function")
		if err != nil {
			b.Fatalf("SendTextMessage convenience function failed: %v", err)
		}
	}
}
