package coze_jwt

var ModelList = []string{
	"moonshot-v1-8k",
	"moonshot-v1-32k",
	"moonshot-v1-128k",
	"Baichuan4",
	"abab6.5s-chat-pro",
	"glm-4-0520",
	"qwen-max",
	"deepseek-r1",
	"deepseek-v3",
	"deepseek-r1-distill-qwen-32b",
	"deepseek-r1-distill-qwen-7b",
	"step-1v-8k",
	"step-1.5v-mini",
	"Doubao-pro-32k",
	"Doubao-pro-256k",
	"Doubao-lite-128k",
	"Doubao-lite-32k",
	"Doubao-vision-lite-32k",
	"Doubao-vision-pro-32k",
	"Doubao-1.5-pro-vision-32k",
	"Doubao-1.5-lite-32k",
	"Doubao-1.5-pro-32k",
	"Doubao-1.5-thinking-pro",
	"Doubao-1.5-pro-256k",
}

var ChannelName = "coze_jwt"

// OAuth JWT 相关常量
const (
	// JWT 算法
	JWTAlgorithm = "RS256"

	// Token 有效期 (1小时)
	TokenExpirationTime = 3600

	// OAuth 端点
	OAuthTokenEndpoint = "/api/permission/oauth2/token"

	// 授权类型
	GrantType = "urn:ietf:params:oauth:grant-type:jwt-bearer"
)

// 工作流相关常量
const (
	// 工作流执行端点
	WorkflowRunEndpoint = "/v1/workflow/run"

	// 工作流状态查询端点
	WorkflowStatusEndpoint = "/v1/workflow/run/status"

	// 工作流状态
	WorkflowStatusRunning   = "running"
	WorkflowStatusCompleted = "completed"
	WorkflowStatusFailed    = "failed"
	WorkflowStatusCanceled  = "canceled"

	// 轮询间隔 (秒)
	WorkflowPollingInterval = 2

	// 最大轮询次数 (总计约5分钟)
	MaxWorkflowPollingAttempts = 150
)
