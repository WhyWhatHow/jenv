package env

import (
	"testing"
)

func TestQueryUserEnvironmentVariable(t *testing.T) {
	// 测试查询用户级别的环境变量
	key := "PATH"
	value, err := QueryUserEnvironmentVariable(key)
	if err != nil {
		t.Fatalf("QueryUserEnvironmentVariable failed: %v", err)
	}
	if value == "" {
		t.Errorf("Expected non-empty value for key %s, got empty string", key)
	}
}

func TestQuerySystemEnvironmentVariable(t *testing.T) {
	// 测试查询系统级别的环境变量
	key := "PATH"
	value, err := QuerySystemEnvironmentVariable(key)
	if err != nil {
		t.Fatalf("QuerySystemEnvironmentVariable failed: %v", err)
	}
	if value == "" {
		t.Errorf("Expected non-empty value for key %s, got empty string", key)
	}
}
