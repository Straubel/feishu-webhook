package feishu

import (
	"testing"
	"time"
)

func TestGenSignDetailed(t *testing.T) {
	tests := []struct {
		name      string
		secret    string
		timestamp int64
		wantErr   bool
	}{
		{
			name:      "标准签名测试",
			secret:    "SECc32bb7cea4f15d6f55f3190529259423343da36045b20de6ad7d3e4fb03ec51d14b",
			timestamp: 1499827200,
			wantErr:   false,
		},
		{
			name:      "当前时间戳测试",
			secret:    "test-secret-key",
			timestamp: time.Now().Unix(),
			wantErr:   false,
		},
		{
			name:      "空密钥测试",
			secret:    "",
			timestamp: time.Now().Unix(),
			wantErr:   false,
		},
		{
			name:      "特殊字符密钥测试",
			secret:    "secret!@#$%^&*()_+{}|:<>?[]\\;'\",./-=",
			timestamp: time.Now().Unix(),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sign1, err1 := GenSign(tt.secret, tt.timestamp)
			if (err1 != nil) != tt.wantErr {
				t.Errorf("GenSign() error = %v, wantErr %v", err1, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if sign1 == "" {
					t.Error("GenSign() should not return empty string")
				}

				sign2, err2 := GenSign(tt.secret, tt.timestamp)
				if err2 != nil {
					t.Errorf("Second GenSign() call failed: %v", err2)
				}

				if sign1 != sign2 {
					t.Error("GenSign() should be deterministic - same input should produce same output")
				}
			}
		})
	}
}

func TestSignConsistency(t *testing.T) {
	secret := "test-secret"
	timestamp := int64(1640995200)

	expectedSign, err := GenSign(secret, timestamp)
	if err != nil {
		t.Fatalf("Failed to generate expected sign: %v", err)
	}

	for i := 0; i < 10; i++ {
		sign, err := GenSign(secret, timestamp)
		if err != nil {
			t.Errorf("Iteration %d: GenSign failed: %v", i, err)
		}
		if sign != expectedSign {
			t.Errorf("Iteration %d: Sign mismatch. Expected %s, got %s", i, expectedSign, sign)
		}
	}
}

func TestSignatureValidation(t *testing.T) {
	secret := "test-secret"
	timestamp1 := int64(1640995200)
	timestamp2 := int64(1640995201)

	sign1, err1 := GenSign(secret, timestamp1)
	if err1 != nil {
		t.Fatalf("Failed to generate sign1: %v", err1)
	}

	sign2, err2 := GenSign(secret, timestamp2)
	if err2 != nil {
		t.Fatalf("Failed to generate sign2: %v", err2)
	}

	if sign1 == sign2 {
		t.Error("Different timestamps should produce different signatures")
	}

	secret2 := "different-secret"
	sign3, err3 := GenSign(secret2, timestamp1)
	if err3 != nil {
		t.Fatalf("Failed to generate sign3: %v", err3)
	}

	if sign1 == sign3 {
		t.Error("Different secrets should produce different signatures")
	}
}
