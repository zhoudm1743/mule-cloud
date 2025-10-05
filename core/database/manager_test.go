package database

import (
	"context"
	"testing"
	
	tenantCtx "mule-cloud/core/context"
)

// TestDatabaseManager 测试数据库管理器
func TestDatabaseManager(t *testing.T) {
	// 注意：这个测试需要Mock，这里只是示例
	
	t.Run("GetTenantDatabaseName", func(t *testing.T) {
		tests := []struct {
			tenantID string
			want     string
		}{
			{"abc123", "mule_abc123"},
			{"test", "mule_test"},
			{"68dda6cd04ba0d6c8dda4b7a", "mule_68dda6cd04ba0d6c8dda4b7a"},
		}
		
		for _, tt := range tests {
			got := GetTenantDatabaseName(tt.tenantID)
			if got != tt.want {
				t.Errorf("GetTenantDatabaseName(%s) = %s, want %s", tt.tenantID, got, tt.want)
			}
		}
	})
}

// TestContextTenantID 测试Context租户ID传递
func TestContextTenantID(t *testing.T) {
	ctx := context.Background()
	
	t.Run("WithTenantID and GetTenantID", func(t *testing.T) {
		testTenantID := "test-tenant-123"
		ctx = tenantCtx.WithTenantID(ctx, testTenantID)
		
		got := tenantCtx.GetTenantID(ctx)
		if got != testTenantID {
			t.Errorf("GetTenantID() = %s, want %s", got, testTenantID)
		}
	})
	
	t.Run("GetTenantID from empty context", func(t *testing.T) {
		emptyCtx := context.Background()
		got := tenantCtx.GetTenantID(emptyCtx)
		if got != "" {
			t.Errorf("GetTenantID() = %s, want empty string", got)
		}
	})
	
	t.Run("WithUserInfo", func(t *testing.T) {
		ctx = tenantCtx.WithUserInfo(ctx, "tenant1", "user1", "admin", []string{"admin", "user"})
		
		tenantID := tenantCtx.GetTenantID(ctx)
		userID := tenantCtx.GetUserID(ctx)
		username := tenantCtx.GetUsername(ctx)
		roles := tenantCtx.GetRoles(ctx)
		
		if tenantID != "tenant1" {
			t.Errorf("GetTenantID() = %s, want tenant1", tenantID)
		}
		if userID != "user1" {
			t.Errorf("GetUserID() = %s, want user1", userID)
		}
		if username != "admin" {
			t.Errorf("GetUsername() = %s, want admin", username)
		}
		if len(roles) != 2 {
			t.Errorf("GetRoles() length = %d, want 2", len(roles))
		}
	})
}

