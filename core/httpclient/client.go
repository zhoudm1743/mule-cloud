package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/consul/api"
)

// Context keys for passing authentication info
type contextKey string

const (
	// ContextKeyToken JWT token key in context
	ContextKeyToken contextKey = "jwt_token"
	// ContextKeyUserID User ID key in context
	ContextKeyUserID contextKey = "user_id"
	// ContextKeyTenantID Tenant ID key in context
	ContextKeyTenantID contextKey = "tenant_id"
)

// ServiceClient HTTP服务客户端（用于服务间调用）
type ServiceClient struct {
	consulClient *api.Client
	httpClient   *http.Client
	// defaultToken 默认的服务间调用 token（可选）
	defaultToken string
}

// NewServiceClient 创建服务客户端
func NewServiceClient(consulAddr string) (*ServiceClient, error) {
	config := api.DefaultConfig()
	config.Address = consulAddr

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("连接Consul失败: %v", err)
	}

	return &ServiceClient{
		consulClient: client,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

// SetDefaultToken 设置默认的服务间调用 token
func (c *ServiceClient) SetDefaultToken(token string) {
	c.defaultToken = token
}

// WithToken 将 JWT token 添加到 context 中
func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, ContextKeyToken, token)
}

// WithUserID 将用户 ID 添加到 context 中
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}

// WithTenantID 将租户 ID 添加到 context 中
func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, ContextKeyTenantID, tenantID)
}

// GetTokenFromContext 从 context 中获取 token
func GetTokenFromContext(ctx context.Context) string {
	if token, ok := ctx.Value(ContextKeyToken).(string); ok {
		return token
	}
	return ""
}

// GetServiceAddress 从Consul获取服务地址
func (c *ServiceClient) GetServiceAddress(serviceName string) (string, error) {
	services, _, err := c.consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", fmt.Errorf("查询服务失败: %v", err)
	}

	if len(services) == 0 {
		return "", fmt.Errorf("未找到服务: %s", serviceName)
	}

	// 简单负载均衡：返回第一个健康实例
	service := services[0].Service
	return fmt.Sprintf("http://%s:%d", service.Address, service.Port), nil
}

// Get 发送GET请求
func (c *ServiceClient) Get(ctx context.Context, serviceName, path string, headers map[string]string) ([]byte, error) {
	return c.doRequest(ctx, "GET", serviceName, path, nil, headers)
}

// Post 发送POST请求
func (c *ServiceClient) Post(ctx context.Context, serviceName, path string, body interface{}, headers map[string]string) ([]byte, error) {
	return c.doRequest(ctx, "POST", serviceName, path, body, headers)
}

// Put 发送PUT请求
func (c *ServiceClient) Put(ctx context.Context, serviceName, path string, body interface{}, headers map[string]string) ([]byte, error) {
	return c.doRequest(ctx, "PUT", serviceName, path, body, headers)
}

// Delete 发送DELETE请求
func (c *ServiceClient) Delete(ctx context.Context, serviceName, path string, headers map[string]string) ([]byte, error) {
	return c.doRequest(ctx, "DELETE", serviceName, path, nil, headers)
}

// doRequest 执行HTTP请求
func (c *ServiceClient) doRequest(ctx context.Context, method, serviceName, path string, body interface{}, headers map[string]string) ([]byte, error) {
	// 1. 从Consul获取服务地址
	baseURL, err := c.GetServiceAddress(serviceName)
	if err != nil {
		return nil, err
	}

	// 2. 构建完整URL
	url := baseURL + path

	// 3. 序列化请求体
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %v", err)
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	// 4. 创建请求
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 5. 设置Header
	req.Header.Set("Content-Type", "application/json")

	// 5.1 自动添加认证信息
	// 优先级：context 中的 token > 手动传递的 headers > 默认 token
	token := GetTokenFromContext(ctx)
	if token == "" && c.defaultToken != "" {
		token = c.defaultToken
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// 5.2 传递用户信息（用于审计和权限检查）
	if userID, ok := ctx.Value(ContextKeyUserID).(string); ok && userID != "" {
		req.Header.Set("X-User-ID", userID)
	}
	if tenantID, ok := ctx.Value(ContextKeyTenantID).(string); ok && tenantID != "" {
		req.Header.Set("X-Tenant-ID", tenantID)
	}

	// 5.3 设置自定义 headers（可以覆盖默认值）
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 6. 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 7. 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 8. 检查状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// CallService 通用服务调用方法（自动解析响应）
func (c *ServiceClient) CallService(ctx context.Context, method, serviceName, path string, reqBody, respData interface{}, headers map[string]string) error {
	var respBody []byte
	var err error

	switch method {
	case "GET":
		respBody, err = c.Get(ctx, serviceName, path, headers)
	case "POST":
		respBody, err = c.Post(ctx, serviceName, path, reqBody, headers)
	case "PUT":
		respBody, err = c.Put(ctx, serviceName, path, reqBody, headers)
	case "DELETE":
		respBody, err = c.Delete(ctx, serviceName, path, headers)
	default:
		return fmt.Errorf("不支持的HTTP方法: %s", method)
	}

	if err != nil {
		return err
	}

	// 解析响应
	if respData != nil {
		if err := json.Unmarshal(respBody, respData); err != nil {
			return fmt.Errorf("解析响应失败: %v", err)
		}
	}

	return nil
}
