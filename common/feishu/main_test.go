package feishu

import (
	"os"
	"testing"
)

var (
	testWebhookURL = "https://open.feishu.cn/open-apis/bot/v2/hook/test-webhook-url"
	testSecret     = "test-secret-key"
)

func TestMain(m *testing.M) {
	// 从环境变量读取测试配置
	if url := os.Getenv("FEISHU_WEBHOOK_URL"); url != "" {
		testWebhookURL = url
	}

	if secret := os.Getenv("FEISHU_SECRET"); secret != "" {
		testSecret = secret
	}

	// 运行测试
	code := m.Run()

	// 清理工作（如果需要）

	os.Exit(code)
}

func skipIfNoCredentials(t *testing.T) {
	if testWebhookURL == "https://open.feishu.cn/open-apis/bot/v2/hook/test-webhook-url" {
		t.Skip("跳过真实API测试 - 请设置 FEISHU_WEBHOOK_URL 环境变量")
	}
}

func TestRealAPI(t *testing.T) {
	skipIfNoCredentials(t)

	if testing.Short() {
		t.Skip("跳过真实API测试")
	}

	t.Run("真实API文本消息测试", func(t *testing.T) {
		sdk := New(testWebhookURL, testSecret)
		err := sdk.SendText("SDK单元测试消息 - 请忽略")
		if err != nil {
			t.Errorf("真实API测试失败: %v", err)
		}
	})
}
