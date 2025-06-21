package coze_jwt

import "encoding/json"

// JWT Claims 结构
type JWTClaims struct {
	Iss string `json:"iss"` // client_id
	Aud string `json:"aud"` // https://api.coze.cn/api/permission/oauth2/token
	Iat int64  `json:"iat"` // 签发时间
	Exp int64  `json:"exp"` // 过期时间
	Jti string `json:"jti"` // JWT ID (随机字符串)
}

// OAuth Token 请求结构
type OAuthTokenRequest struct {
	GrantType string `json:"grant_type"`
	Assertion string `json:"assertion"`
}

// OAuth Token 响应结构
type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Error       string `json:"error,omitempty"`
	ErrorDesc   string `json:"error_description,omitempty"`
}

// Coze JWT 配置结构
type CozeJWTConfig struct {
	ClientID     string `json:"client_id"`
	PublicKeyID  string `json:"public_key_id"`
	PrivateKey   string `json:"private_key"`
	SpaceID      string `json:"space_id"`
	DefaultBotID string `json:"default_bot_id"`
}

// 继承原有的 Coze 数据结构
type CozeError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CozeEnterMessage struct {
	Role        string          `json:"role"`
	Type        string          `json:"type,omitempty"`
	Content     any             `json:"content,omitempty"`
	MetaData    json.RawMessage `json:"meta_data,omitempty"`
	ContentType string          `json:"content_type,omitempty"`
}

type CozeChatRequest struct {
	BotId              string             `json:"bot_id"`
	UserId             string             `json:"user_id"`
	AdditionalMessages []CozeEnterMessage `json:"additional_messages,omitempty"`
	Stream             bool               `json:"stream,omitempty"`
	CustomVariables    json.RawMessage    `json:"custom_variables,omitempty"`
	AutoSaveHistory    bool               `json:"auto_save_history,omitempty"`
	MetaData           json.RawMessage    `json:"meta_data,omitempty"`
	ExtraParams        json.RawMessage    `json:"extra_params,omitempty"`
	ShortcutCommand    json.RawMessage    `json:"shortcut_command,omitempty"`
	Parameters         json.RawMessage    `json:"parameters,omitempty"`
}

// 工作流请求结构
type CozeWorkflowRequest struct {
	WorkflowId      string          `json:"workflow_id"`
	Parameters      json.RawMessage `json:"parameters,omitempty"`
	Stream          bool            `json:"stream,omitempty"`
	CustomVariables json.RawMessage `json:"custom_variables,omitempty"`
	IsAsync         bool            `json:"is_async,omitempty"`
}

// 工作流响应结构
type CozeWorkflowResponse struct {
	Code int                      `json:"code"`
	Msg  string                   `json:"msg"`
	Data CozeWorkflowResponseData `json:"data"`
}

type CozeWorkflowResponseData struct {
	ExecuteId string               `json:"execute_id"`
	Status    string               `json:"status"` // running, completed, failed, canceled
	Result    json.RawMessage      `json:"result,omitempty"`
	Error     *CozeError           `json:"error,omitempty"`
	CreatedAt int64                `json:"created_at"`
	UpdatedAt int64                `json:"updated_at"`
	Usage     *CozeWorkflowUsage   `json:"usage,omitempty"`
	Outputs   []CozeWorkflowOutput `json:"outputs,omitempty"`
}

type CozeWorkflowUsage struct {
	TokenCount  int `json:"token_count"`
	OutputCount int `json:"output_count"`
	InputCount  int `json:"input_count"`
}

type CozeWorkflowOutput struct {
	NodeId   string          `json:"node_id"`
	NodeType string          `json:"node_type"`
	Output   json.RawMessage `json:"output"`
}

// 工作流状态查询响应
type CozeWorkflowStatusResponse struct {
	Code int                      `json:"code"`
	Msg  string                   `json:"msg"`
	Data CozeWorkflowResponseData `json:"data"`
}

type CozeChatResponse struct {
	Code int                  `json:"code"`
	Msg  string               `json:"msg"`
	Data CozeChatResponseData `json:"data"`
}

type CozeChatResponseData struct {
	Id             string        `json:"id"`
	ConversationId string        `json:"conversation_id"`
	BotId          string        `json:"bot_id"`
	CreatedAt      int64         `json:"created_at"`
	LastError      CozeError     `json:"last_error"`
	Status         string        `json:"status"`
	Usage          CozeChatUsage `json:"usage"`
}

type CozeChatUsage struct {
	TokenCount  int `json:"token_count"`
	OutputCount int `json:"output_count"`
	InputCount  int `json:"input_count"`
}

type CozeChatDetailResponse struct {
	Data   []CozeChatV3MessageDetail `json:"data"`
	Code   int                       `json:"code"`
	Msg    string                    `json:"msg"`
	Detail CozeResponseDetail        `json:"detail"`
}

type CozeChatV3MessageDetail struct {
	Id               string          `json:"id"`
	Role             string          `json:"role"`
	Type             string          `json:"type"`
	BotId            string          `json:"bot_id"`
	ChatId           string          `json:"chat_id"`
	Content          json.RawMessage `json:"content"`
	MetaData         json.RawMessage `json:"meta_data"`
	CreatedAt        int64           `json:"created_at"`
	SectionId        string          `json:"section_id"`
	UpdatedAt        int64           `json:"updated_at"`
	ContentType      string          `json:"content_type"`
	ConversationId   string          `json:"conversation_id"`
	ReasoningContent string          `json:"reasoning_content"`
}

type CozeResponseDetail struct {
	Logid string `json:"logid"`
}
