package unleash

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kromiii/stale-flag-detector/config"
)

func TestGetStaleFlags(t *testing.T) {
	// モックサーバーのセットアップ
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		response := FeatureFlagsResponse{
			Features: []FeatureFlag{
				{Name: "flag1", Type: "release", CreatedAt: time.Now().Add(-41 * 24 * time.Hour), Enabled: true, Stale: false},
				{Name: "flag2", Type: "experiment", CreatedAt: time.Now().Add(-30 * 24 * time.Hour), Enabled: true, Stale: false},
				{Name: "flag3", Type: "operational", CreatedAt: time.Now().Add(-8 * 24 * time.Hour), Enabled: true, Stale: true},
				{Name: "flag4", Type: "kill-switch", CreatedAt: time.Now().Add(-366 * 24 * time.Hour), Enabled: true, Stale: false},
			},
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))
	defer server.Close()

	// テスト用の設定
	cfg := &config.Config{
		ReleaseFlagLifetime:     40,
		ExperimentFlagLifetime:  40,
		OperationalFlagLifetime: 7,
		PermissionFlagLifetime:  -1,
	}

	// クライアントの作成
	client := NewClient(server.URL, "test-token", "test-project", cfg)

	// テストケース
	tests := []struct {
		name           string
		onlyStaleFlags bool
		expected       []string
	}{
		{"全ての古いフラグを取得", false, []string{"flag1", "flag3"}},
		{"Staleフラグのみを取得", true, []string{"flag3"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			flags, err := client.GetStaleFlags(tc.onlyStaleFlags)
			if err != nil {
				t.Fatalf("GetStaleFlags failed: %v", err)
			}

			if len(flags) != len(tc.expected) {
				t.Errorf("Expected %d flags, got %d", len(tc.expected), len(flags))
			}

			for i, flag := range flags {
				if flag != tc.expected[i] {
					t.Errorf("Expected flag %s, got %s", tc.expected[i], flag)
				}
			}
		})
	}
}

func TestGetExpectedLifetime(t *testing.T) {
	cfg := &config.Config{
		ReleaseFlagLifetime:     40,
		ExperimentFlagLifetime:  40,
		OperationalFlagLifetime: 7,
		PermissionFlagLifetime:  -1,
	}
	client := NewClient("", "", "", cfg)

	tests := []struct {
		name     string
		flagType string
		want     time.Duration
	}{
		{"release", "release", 40 * 24 * time.Hour},
		{"experiment", "experiment", 40 * 24 * time.Hour},
		{"operational", "operational", 7 * 24 * time.Hour},
		{"killSwitch", "kill-switch", time.Duration(math.MaxInt64)},
		{"permission", "permission", time.Duration(math.MaxInt64)},
		{"default", "unknown", 30 * 24 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := client.getExpectedLifetime(tt.flagType); got != tt.want {
				t.Errorf("Expected %v for %s, got %v", tt.want, tt.flagType, got)
			}
		})
	}
}
