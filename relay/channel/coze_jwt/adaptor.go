package coze_jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"one-api/dto"
	"one-api/relay/channel"
	"one-api/relay/common"
	"one-api/service"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Adaptor struct {
}

// ConvertAudioRequest implements channel.Adaptor.
func (a *Adaptor) ConvertAudioRequest(c *gin.Context, info *common.RelayInfo, request dto.AudioRequest) (io.Reader, error) {
	return nil, errors.New("not implemented")
}

// ConvertClaudeRequest implements channel.Adaptor.
func (a *Adaptor) ConvertClaudeRequest(c *gin.Context, info *common.RelayInfo, request *dto.ClaudeRequest) (any, error) {
	return nil, errors.New("not implemented")
}

// ConvertEmbeddingRequest implements channel.Adaptor.
func (a *Adaptor) ConvertEmbeddingRequest(c *gin.Context, info *common.RelayInfo, request dto.EmbeddingRequest) (any, error) {
	return nil, errors.New("not implemented")
}

// ConvertImageRequest implements channel.Adaptor.
func (a *Adaptor) ConvertImageRequest(c *gin.Context, info *common.RelayInfo, request dto.ImageRequest) (any, error) {
	return nil, errors.New("not implemented")
}

// ConvertOpenAIRequest implements channel.Adaptor.
func (a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *common.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}

	// 从gin上下文获取渠道设置
	channelSetting := c.GetStringMap("channel_setting")

	// 解析配置
	config, err := ParseCozeJWTConfigFromSetting(channelSetting)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// 检查是否为工作流请求 (通过模型名称判断)
	if strings.HasPrefix(request.Model, "workflow:") {
		return ConvertCozeWorkflowRequest(c, *request, config), nil
	}

	// 默认为智能体请求
	return ConvertCozeChatRequest(c, *request, config), nil
}

// ConvertOpenAIResponsesRequest implements channel.Adaptor.
func (a *Adaptor) ConvertOpenAIResponsesRequest(c *gin.Context, info *common.RelayInfo, request dto.OpenAIResponsesRequest) (any, error) {
	return nil, errors.New("not implemented")
}

// ConvertRerankRequest implements channel.Adaptor.
func (a *Adaptor) ConvertRerankRequest(c *gin.Context, relayMode int, request dto.RerankRequest) (any, error) {
	return nil, errors.New("not implemented")
}

// DoRequest implements channel.Adaptor.
func (a *Adaptor) DoRequest(c *gin.Context, info *common.RelayInfo, requestBody io.Reader) (any, error) {
	if info.IsStream {
		return channel.DoApiRequest(a, c, info, requestBody)
	}

	// 检查是否为工作流请求
	modelName := info.UpstreamModelName
	if strings.HasPrefix(modelName, "workflow:") {
		return a.doWorkflowRequest(c, info, requestBody)
	}

	// 智能体请求处理 (与原 Coze 相同)
	return a.doChatRequest(c, info, requestBody)
}

// DoResponse implements channel.Adaptor.
func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *common.RelayInfo) (usage any, err *dto.OpenAIErrorWithStatusCode) {
	// 检查是否为工作流请求
	modelName := info.UpstreamModelName
	if strings.HasPrefix(modelName, "workflow:") || strings.HasPrefix(modelName, "workflow-async:") {
		err, usage = CozeWorkflowHandler(c, resp, info)
		return
	}

	// 智能体请求处理
	if info.IsStream {
		err, usage = CozeChatStreamHandler(c, resp, info)
	} else {
		err, usage = CozeChatHandler(c, resp, info)
	}
	return
}

// GetChannelName implements channel.Adaptor.
func (a *Adaptor) GetChannelName() string {
	return ChannelName
}

// GetModelList implements channel.Adaptor.
func (a *Adaptor) GetModelList() []string {
	return ModelList
}

// GetRequestURL implements channel.Adaptor.
func (a *Adaptor) GetRequestURL(info *common.RelayInfo) (string, error) {
	// 根据模型名称判断请求类型
	modelName := info.UpstreamModelName
	if strings.HasPrefix(modelName, "workflow:") {
		return fmt.Sprintf("%s%s", info.BaseUrl, WorkflowRunEndpoint), nil
	}

	// 默认为智能体聊天
	return fmt.Sprintf("%s/v3/chat", info.BaseUrl), nil
}

// Init implements channel.Adaptor.
func (a *Adaptor) Init(info *common.RelayInfo) {
}

// SetupRequestHeader implements channel.Adaptor.
func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *common.RelayInfo) error {
	channel.SetupApiRequestHeader(info, c, req)

	// 从gin上下文获取渠道设置
	channelSetting := c.GetStringMap("channel_setting")

	// 解析配置并获取 OAuth token
	config, err := ParseCozeJWTConfigFromSetting(channelSetting)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	accessToken, err := GetOAuthAccessToken(config, info.BaseUrl)
	if err != nil {
		return fmt.Errorf("failed to get OAuth access token: %w", err)
	}

	req.Set("Authorization", "Bearer "+accessToken)
	return nil
}

// 智能体聊天请求处理
func (a *Adaptor) doChatRequest(c *gin.Context, info *common.RelayInfo, requestBody io.Reader) (any, error) {
	// 首先发送创建消息请求，成功后再发送获取消息请求
	resp, err := channel.DoApiRequest(a, c, info, requestBody)
	if err != nil {
		return nil, err
	}

	// 解析 resp
	var cozeResponse CozeChatResponse
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respBody, &cozeResponse)
	if cozeResponse.Code != 0 {
		return nil, errors.New(cozeResponse.Msg)
	}
	c.Set("coze_conversation_id", cozeResponse.Data.ConversationId)
	c.Set("coze_chat_id", cozeResponse.Data.Id)

	// 轮询检查消息是否完成
	for {
		err, isComplete := a.checkIfChatComplete(c, info)
		if err != nil {
			return nil, err
		} else {
			if isComplete {
				break
			}
		}
		time.Sleep(time.Second * 1)
	}

	// 发送获取消息请求
	return a.getChatDetail(c, info)
}

// 工作流请求处理
func (a *Adaptor) doWorkflowRequest(c *gin.Context, info *common.RelayInfo, requestBody io.Reader) (any, error) {
	// 发送工作流执行请求
	resp, err := channel.DoApiRequest(a, c, info, requestBody)
	if err != nil {
		return nil, err
	}

	// 解析工作流响应
	var workflowResponse CozeWorkflowResponse
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read workflow response body failed: %w", err)
	}
	err = json.Unmarshal(respBody, &workflowResponse)
	if err != nil {
		return nil, fmt.Errorf("unmarshal workflow response failed: %w", err)
	}
	if workflowResponse.Code != 0 {
		return nil, fmt.Errorf("workflow request failed: %s", workflowResponse.Msg)
	}

	// 检查是否为异步工作流
	if workflowResponse.Data.Status == WorkflowStatusRunning {
		// 异步工作流，需要轮询状态
		return a.pollWorkflowStatus(c, info, workflowResponse.Data.ExecuteId)
	}

	// 同步工作流，直接返回结果
	return resp, nil
}

// 检查聊天是否完成
func (a *Adaptor) checkIfChatComplete(c *gin.Context, info *common.RelayInfo) (error, bool) {
	requestURL := fmt.Sprintf("%s/v3/chat/retrieve", info.BaseUrl)
	requestURL = requestURL + "?conversation_id=" + c.GetString("coze_conversation_id") + "&chat_id=" + c.GetString("coze_chat_id")

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return err, false
	}
	err = a.SetupRequestHeader(c, &req.Header, info)
	if err != nil {
		return err, false
	}

	resp, err := a.doRequest(req, info)
	if err != nil {
		return err, false
	}
	if resp == nil {
		return fmt.Errorf("resp is nil"), false
	}
	defer resp.Body.Close()

	var cozeResponse CozeChatResponse
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body failed: %w", err), false
	}
	err = json.Unmarshal(responseBody, &cozeResponse)
	if err != nil {
		return fmt.Errorf("unmarshal response body failed: %w", err), false
	}
	if cozeResponse.Data.Status == "completed" {
		c.Set("coze_token_count", cozeResponse.Data.Usage.TokenCount)
		c.Set("coze_output_count", cozeResponse.Data.Usage.OutputCount)
		c.Set("coze_input_count", cozeResponse.Data.Usage.InputCount)
		return nil, true
	} else if cozeResponse.Data.Status == "failed" || cozeResponse.Data.Status == "canceled" || cozeResponse.Data.Status == "requires_action" {
		return fmt.Errorf("chat status: %s", cozeResponse.Data.Status), false
	} else {
		return nil, false
	}
}

// 获取聊天详情
func (a *Adaptor) getChatDetail(c *gin.Context, info *common.RelayInfo) (*http.Response, error) {
	requestURL := fmt.Sprintf("%s/v3/chat/message/list", info.BaseUrl)
	requestURL = requestURL + "?conversation_id=" + c.GetString("coze_conversation_id") + "&chat_id=" + c.GetString("coze_chat_id")

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("new request failed: %w", err)
	}
	err = a.SetupRequestHeader(c, &req.Header, info)
	if err != nil {
		return nil, fmt.Errorf("setup request header failed: %w", err)
	}
	resp, err := a.doRequest(req, info)
	if err != nil {
		return nil, fmt.Errorf("do request failed: %w", err)
	}
	return resp, nil
}

// 轮询工作流状态
func (a *Adaptor) pollWorkflowStatus(c *gin.Context, info *common.RelayInfo, executeId string) (*http.Response, error) {
	for attempt := 0; attempt < MaxWorkflowPollingAttempts; attempt++ {
		// 构建状态查询URL
		statusURL := fmt.Sprintf("%s%s?execute_id=%s", info.BaseUrl, WorkflowStatusEndpoint, executeId)

		req, err := http.NewRequest("GET", statusURL, nil)
		if err != nil {
			return nil, fmt.Errorf("create status request failed: %w", err)
		}

		err = a.SetupRequestHeader(c, &req.Header, info)
		if err != nil {
			return nil, fmt.Errorf("setup status request header failed: %w", err)
		}

		resp, err := a.doRequest(req, info)
		if err != nil {
			return nil, fmt.Errorf("do status request failed: %w", err)
		}

		// 解析状态响应
		var statusResponse CozeWorkflowStatusResponse
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("read status response body failed: %w", err)
		}
		resp.Body.Close()

		err = json.Unmarshal(responseBody, &statusResponse)
		if err != nil {
			return nil, fmt.Errorf("unmarshal status response failed: %w", err)
		}

		if statusResponse.Code != 0 {
			return nil, fmt.Errorf("status query failed: %s", statusResponse.Msg)
		}

		// 检查工作流状态
		switch statusResponse.Data.Status {
		case WorkflowStatusCompleted:
			// 工作流完成，设置使用量信息并返回结果
			if statusResponse.Data.Usage != nil {
				c.Set("coze_token_count", statusResponse.Data.Usage.TokenCount)
				c.Set("coze_output_count", statusResponse.Data.Usage.OutputCount)
				c.Set("coze_input_count", statusResponse.Data.Usage.InputCount)
			}
			// 重新构建响应
			return a.buildWorkflowResponse(statusResponse), nil

		case WorkflowStatusFailed, WorkflowStatusCanceled:
			// 工作流失败或取消
			errorMsg := fmt.Sprintf("workflow %s", statusResponse.Data.Status)
			if statusResponse.Data.Error != nil {
				errorMsg = fmt.Sprintf("workflow failed: %s", statusResponse.Data.Error.Message)
			}
			return nil, fmt.Errorf(errorMsg)

		case WorkflowStatusRunning:
			// 继续轮询
			time.Sleep(time.Duration(WorkflowPollingInterval) * time.Second)
			continue

		default:
			return nil, fmt.Errorf("unknown workflow status: %s", statusResponse.Data.Status)
		}
	}

	return nil, fmt.Errorf("workflow polling timeout after %d attempts", MaxWorkflowPollingAttempts)
}

// 构建工作流响应
func (a *Adaptor) buildWorkflowResponse(statusResponse CozeWorkflowStatusResponse) *http.Response {
	// 将状态响应转换为标准HTTP响应
	responseBody, _ := json.Marshal(statusResponse)
	resp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(string(responseBody))),
	}
	resp.Header.Set("Content-Type", "application/json")
	return resp
}

// HTTP 请求执行
func (a *Adaptor) doRequest(req *http.Request, info *common.RelayInfo) (*http.Response, error) {
	var client *http.Client
	var err error
	if proxyURL, ok := info.ChannelSetting["proxy"]; ok {
		client, err = service.NewProxyHttpClient(proxyURL.(string))
		if err != nil {
			return nil, fmt.Errorf("new proxy http client failed: %w", err)
		}
	} else {
		client = service.GetHttpClient()
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do failed: %w", err)
	}
	return resp, nil
}
