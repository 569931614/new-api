package coze_jwt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"one-api/common"
	"sync"
	"time"

	"github.com/coze-dev/coze-go"
)

// Token 缓存结构
type TokenCache struct {
	AccessToken string
	ExpiresAt   time.Time
	mutex       sync.RWMutex
}

var tokenCache = &TokenCache{}

// OAuth 客户端缓存
type OAuthClientCache struct {
	Client *coze.JWTOAuthClient
	Config *coze.OAuthConfig
	mutex  sync.RWMutex
}

var oauthClientCache = &OAuthClientCache{}

// 创建或获取 OAuth 客户端
func getOAuthClient(config *CozeJWTConfig, baseURL string) (*coze.JWTOAuthClient, error) {
	oauthClientCache.mutex.RLock()
	if oauthClientCache.Client != nil && oauthClientCache.Config != nil {
		// 检查配置是否匹配
		if oauthClientCache.Config.ClientID == config.ClientID &&
			oauthClientCache.Config.PublicKeyID == config.PublicKeyID &&
			oauthClientCache.Config.CozeAPIBase == baseURL {
			client := oauthClientCache.Client
			oauthClientCache.mutex.RUnlock()
			return client, nil
		}
	}
	oauthClientCache.mutex.RUnlock()

	// 创建新的 OAuth 配置
	oauthConfig := &coze.OAuthConfig{
		ClientType:  "jwt",
		ClientID:    config.ClientID,
		CozeWWWBase: "https://www.coze.cn", // 固定值
		CozeAPIBase: baseURL,
		PrivateKey:  config.PrivateKey,
		PublicKeyID: config.PublicKeyID,
	}

	common.SysLog(fmt.Sprintf("Creating OAuth config: ClientID=%s, PublicKeyID=%s, BaseURL=%s",
		config.ClientID, config.PublicKeyID, baseURL))

	// 创建 OAuth 客户端
	oauth, err := coze.LoadOAuthAppFromConfig(oauthConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load OAuth config: %w", err)
	}

	jwtClient, ok := oauth.(*coze.JWTOAuthClient)
	if !ok {
		return nil, errors.New("invalid OAuth client type: expected JWT client")
	}

	// 缓存客户端和配置
	oauthClientCache.mutex.Lock()
	oauthClientCache.Client = jwtClient
	oauthClientCache.Config = oauthConfig
	oauthClientCache.mutex.Unlock()

	return jwtClient, nil
}

// 获取 OAuth Access Token (使用官方SDK)
func GetOAuthAccessToken(config *CozeJWTConfig, baseURL string) (string, error) {
	// 检查缓存的 token 是否有效
	tokenCache.mutex.RLock()
	if tokenCache.AccessToken != "" && time.Now().Before(tokenCache.ExpiresAt) {
		token := tokenCache.AccessToken
		tokenCache.mutex.RUnlock()
		return token, nil
	}
	tokenCache.mutex.RUnlock()

	// 获取 OAuth 客户端
	oauthClient, err := getOAuthClient(config, baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to get OAuth client: %w", err)
	}

	// 使用官方SDK获取访问令牌
	ctx := context.Background()
	resp, err := oauthClient.GetAccessToken(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}

	if resp.AccessToken == "" {
		return "", errors.New("no access token received")
	}

	// 缓存 token
	tokenCache.mutex.Lock()
	tokenCache.AccessToken = resp.AccessToken
	// 提前 5 分钟过期以避免边界情况
	tokenCache.ExpiresAt = time.Unix(resp.ExpiresIn-300, 0)
	tokenCache.mutex.Unlock()

	return resp.AccessToken, nil
}

// 从渠道设置中解析配置
func ParseCozeJWTConfigFromSetting(setting map[string]interface{}) (*CozeJWTConfig, error) {
	var config CozeJWTConfig

	// 从设置中提取配置字段
	if clientID, ok := setting["client_id"].(string); ok {
		config.ClientID = clientID
	}
	if publicKeyID, ok := setting["public_key_id"].(string); ok {
		config.PublicKeyID = publicKeyID
	}
	if privateKey, ok := setting["private_key"].(string); ok {
		config.PrivateKey = privateKey
	}
	if spaceID, ok := setting["space_id"].(string); ok {
		config.SpaceID = spaceID
	}
	if defaultBotID, ok := setting["default_bot_id"].(string); ok {
		config.DefaultBotID = defaultBotID
	}

	// 验证必需字段
	if config.ClientID == "" {
		return nil, errors.New("client_id is required in channel setting")
	}
	if config.PublicKeyID == "" {
		return nil, errors.New("public_key_id is required in channel setting")
	}
	if config.PrivateKey == "" {
		return nil, errors.New("private_key is required in channel setting")
	}

	return &config, nil
}

// 解析渠道配置 (保持向后兼容)
func ParseCozeJWTConfig(otherInfo string) (*CozeJWTConfig, error) {
	var config CozeJWTConfig
	err := json.Unmarshal([]byte(otherInfo), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Coze JWT config: %w", err)
	}

	// 验证必需字段
	if config.ClientID == "" {
		return nil, errors.New("client_id is required")
	}
	if config.PublicKeyID == "" {
		return nil, errors.New("public_key_id is required")
	}
	if config.PrivateKey == "" {
		return nil, errors.New("private_key is required")
	}

	return &config, nil
}

// 清理缓存 (用于测试或重置)
func ClearCache() {
	tokenCache.mutex.Lock()
	tokenCache.AccessToken = ""
	tokenCache.ExpiresAt = time.Time{}
	tokenCache.mutex.Unlock()

	oauthClientCache.mutex.Lock()
	oauthClientCache.Client = nil
	oauthClientCache.Config = nil
	oauthClientCache.mutex.Unlock()
}
