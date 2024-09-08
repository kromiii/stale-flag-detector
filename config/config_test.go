package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// 正常系のテスト
	t.Run("全ての環境変数が設定されている場合", func(t *testing.T) {
		os.Setenv("UNLEASH_API_ENDPOINT", "https://example.com")
		os.Setenv("UNLEASH_API_TOKEN", "token123")
		os.Setenv("UNLEASH_PROJECT_ID", "project1")
		os.Setenv("OPENAI_API_KEY", "openai123")
		os.Setenv("GITHUB_TOKEN", "github123")
		os.Setenv("GITHUB_OWNER", "owner")
		os.Setenv("GITHUB_REPO", "repo")
		os.Setenv("GITHUB_BASE_URL", "https://api.github.com")

		config, err := Load()
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if config == nil {
			t.Fatal("設定がnilです")
			return
		}
		if config.UnleashAPIEndpoint != "https://example.com" {
			t.Errorf("UnleashAPIEndpointが期待値と異なります: got %v, want %v", config.UnleashAPIEndpoint, "https://example.com")
		}
	})

	// エラー系のテスト
	t.Run("必須の環境変数が欠けている場合", func(t *testing.T) {
		os.Unsetenv("UNLEASH_API_ENDPOINT")
		os.Unsetenv("UNLEASH_API_TOKEN")
		os.Unsetenv("UNLEASH_PROJECT_ID")
		os.Unsetenv("OPENAI_API_KEY")
		os.Unsetenv("GITHUB_TOKEN")
		os.Unsetenv("GITHUB_OWNER")
		os.Unsetenv("GITHUB_REPO")
		os.Unsetenv("GITHUB_BASE_URL")

		_, err := Load()
		if err == nil {
			t.Error("エラーが期待されましたが、nilが返されました")
		}
	})
}
