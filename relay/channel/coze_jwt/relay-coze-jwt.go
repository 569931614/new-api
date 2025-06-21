package coze_jwt

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"one-api/common"
	"one-api/dto"
	relaycommon "one-api/relay/common"
	"one-api/relay/helper"
	"one-api/service"
	"strings"

	"github.com/gin-gonic/gin"
)

// 转换智能体聊天请求
func ConvertCozeChatRequest(c *gin.Context, request dto.GeneralOpenAIRequest, config *CozeJWTConfig) *CozeChatRequest {
	var messages []CozeEnterMessage
	// 将 request的messages的role为user的content转换为CozeMessage
	for _, message := range request.Messages {
		if message.Role == "user" {
			messages = append(messages, CozeEnterMessage{
				Role:    "user",
				Content: message.Content,
				// TODO: support more content type
				ContentType: "text",
			})
		}
	}
	user := request.User
	if user == "" {
		user = helper.GetResponseID(c)
	}

	// 获取 bot_id，优先使用配置中的默认值
	botId := config.DefaultBotID
	if botId == "" {
		botId = c.GetString("bot_id")
	}

	cozeRequest := &CozeChatRequest{
		BotId:              botId,
		UserId:             user,
		AdditionalMessages: messages,
		Stream:             request.Stream,
	}
	return cozeRequest
}

// 转换工作流请求
func ConvertCozeWorkflowRequest(c *gin.Context, request dto.GeneralOpenAIRequest, config *CozeJWTConfig) *CozeWorkflowRequest {
	// 从模型名称中提取工作流ID
	workflowId := strings.TrimPrefix(request.Model, "workflow:")

	// 将消息转换为参数
	var parameters map[string]interface{}
	if len(request.Messages) > 0 {
		// 简单处理：将最后一条用户消息作为输入参数
		for i := len(request.Messages) - 1; i >= 0; i-- {
			if request.Messages[i].Role == "user" {
				parameters = map[string]interface{}{
					"input": request.Messages[i].Content,
				}
				break
			}
		}
	}

	parametersJson, _ := json.Marshal(parameters)

	// 检查是否为异步工作流请求
	// 可以通过模型名称的特殊标记来判断，例如 "workflow-async:id"
	isAsync := strings.HasPrefix(request.Model, "workflow-async:")
	if isAsync {
		workflowId = strings.TrimPrefix(request.Model, "workflow-async:")
	}

	cozeRequest := &CozeWorkflowRequest{
		WorkflowId: workflowId,
		Parameters: parametersJson,
		Stream:     request.Stream,
		IsAsync:    isAsync || !request.Stream, // 非流式请求默认使用异步模式
	}
	return cozeRequest
}

// 智能体聊天响应处理
func CozeChatHandler(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (*dto.OpenAIErrorWithStatusCode, *dto.Usage) {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return service.OpenAIErrorWrapper(err, "read_response_body_failed", http.StatusInternalServerError), nil
	}
	err = resp.Body.Close()
	if err != nil {
		return service.OpenAIErrorWrapperLocal(err, "close_response_body_failed", http.StatusInternalServerError), nil
	}
	// convert coze response to openai response
	var response dto.TextResponse
	var cozeResponse CozeChatDetailResponse
	response.Model = info.UpstreamModelName
	err = json.Unmarshal(responseBody, &cozeResponse)
	if err != nil {
		return service.OpenAIErrorWrapper(err, "unmarshal_response_body_failed", http.StatusInternalServerError), nil
	}
	if cozeResponse.Code != 0 {
		return service.OpenAIErrorWrapper(errors.New(cozeResponse.Msg), fmt.Sprintf("%d", cozeResponse.Code), http.StatusInternalServerError), nil
	}
	// 从上下文获取 usage
	var usage dto.Usage
	usage.PromptTokens = c.GetInt("coze_input_count")
	usage.CompletionTokens = c.GetInt("coze_output_count")
	usage.TotalTokens = c.GetInt("coze_token_count")
	response.Usage = usage
	response.Id = helper.GetResponseID(c)

	var responseContent json.RawMessage
	for _, data := range cozeResponse.Data {
		if data.Type == "answer" {
			responseContent = data.Content
			response.Created = data.CreatedAt
		}
	}
	// 添加 response.Choices
	response.Choices = []dto.OpenAITextResponseChoice{
		{
			Index:        0,
			Message:      dto.Message{Role: "assistant", Content: responseContent},
			FinishReason: "stop",
		},
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return service.OpenAIErrorWrapper(err, "marshal_response_body_failed", http.StatusInternalServerError), nil
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(resp.StatusCode)
	_, _ = c.Writer.Write(jsonResponse)

	return nil, &usage
}

// 工作流响应处理
func CozeWorkflowHandler(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (*dto.OpenAIErrorWithStatusCode, *dto.Usage) {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return service.OpenAIErrorWrapper(err, "read_response_body_failed", http.StatusInternalServerError), nil
	}
	err = resp.Body.Close()
	if err != nil {
		return service.OpenAIErrorWrapperLocal(err, "close_response_body_failed", http.StatusInternalServerError), nil
	}

	// 解析工作流响应
	var workflowResponse CozeWorkflowStatusResponse
	err = json.Unmarshal(responseBody, &workflowResponse)
	if err != nil {
		return service.OpenAIErrorWrapper(err, "unmarshal_response_body_failed", http.StatusInternalServerError), nil
	}
	if workflowResponse.Code != 0 {
		return service.OpenAIErrorWrapper(errors.New(workflowResponse.Msg), fmt.Sprintf("%d", workflowResponse.Code), http.StatusInternalServerError), nil
	}

	// 转换为OpenAI格式响应
	var response dto.TextResponse
	response.Model = info.UpstreamModelName
	response.Id = helper.GetResponseID(c)
	response.Created = workflowResponse.Data.CreatedAt

	// 从上下文获取 usage
	var usage dto.Usage
	usage.PromptTokens = c.GetInt("coze_input_count")
	usage.CompletionTokens = c.GetInt("coze_output_count")
	usage.TotalTokens = c.GetInt("coze_token_count")
	response.Usage = usage

	// 处理工作流输出
	var responseContent string
	if workflowResponse.Data.Result != nil {
		responseContent = string(workflowResponse.Data.Result)
	} else if len(workflowResponse.Data.Outputs) > 0 {
		// 如果有多个输出，合并为一个响应
		outputs := make([]string, 0, len(workflowResponse.Data.Outputs))
		for _, output := range workflowResponse.Data.Outputs {
			outputs = append(outputs, string(output.Output))
		}
		responseContent = strings.Join(outputs, "\n")
	} else {
		responseContent = "工作流执行完成"
	}

	// 添加 response.Choices
	response.Choices = []dto.OpenAITextResponseChoice{
		{
			Index:        0,
			Message:      dto.Message{Role: "assistant", Content: json.RawMessage(`"` + responseContent + `"`)},
			FinishReason: "stop",
		},
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return service.OpenAIErrorWrapper(err, "marshal_response_body_failed", http.StatusInternalServerError), nil
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(resp.StatusCode)
	_, _ = c.Writer.Write(jsonResponse)

	return nil, &usage
}

// 流式响应处理
func CozeChatStreamHandler(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (*dto.OpenAIErrorWithStatusCode, *dto.Usage) {
	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(bufio.ScanLines)
	helper.SetEventStreamHeaders(c)
	id := helper.GetResponseID(c)
	var responseText string

	var currentEvent string
	var currentData string
	var usage dto.Usage

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			if currentEvent != "" && currentData != "" {
				// handle last event
				handleCozeEvent(c, currentEvent, currentData, &responseText, &usage, id, info)
				currentEvent = ""
				currentData = ""
			}
			continue
		}

		if strings.HasPrefix(line, "event:") {
			currentEvent = strings.TrimSpace(line[6:])
			continue
		}

		if strings.HasPrefix(line, "data:") {
			currentData = strings.TrimSpace(line[5:])
			continue
		}
	}

	// Last event
	if currentEvent != "" && currentData != "" {
		handleCozeEvent(c, currentEvent, currentData, &responseText, &usage, id, info)
	}

	if err := scanner.Err(); err != nil {
		return service.OpenAIErrorWrapper(err, "stream_scanner_error", http.StatusInternalServerError), nil
	}
	helper.Done(c)

	if usage.TotalTokens == 0 {
		usage.PromptTokens = info.PromptTokens
		usage.CompletionTokens, _ = service.CountTextToken("gpt-3.5-turbo", responseText)
		usage.TotalTokens = usage.PromptTokens + usage.CompletionTokens
	}

	return nil, &usage
}

// 处理Coze事件
func handleCozeEvent(c *gin.Context, event string, data string, responseText *string, usage *dto.Usage, id string, info *relaycommon.RelayInfo) {
	switch event {
	case "conversation.chat.completed":
		// 将 data 解析为 CozeChatResponseData
		var chatData CozeChatResponseData
		err := json.Unmarshal([]byte(data), &chatData)
		if err != nil {
			common.SysError("error_unmarshalling_stream_response: " + err.Error())
			return
		}

		usage.PromptTokens = chatData.Usage.InputCount
		usage.CompletionTokens = chatData.Usage.OutputCount
		usage.TotalTokens = chatData.Usage.TokenCount

		finishReason := "stop"
		stopResponse := helper.GenerateStopResponse(id, common.GetTimestamp(), info.UpstreamModelName, finishReason)
		helper.ObjectData(c, stopResponse)

	case "conversation.message.delta":
		// 将 data 解析为 CozeChatV3MessageDetail
		var messageData CozeChatV3MessageDetail
		err := json.Unmarshal([]byte(data), &messageData)
		if err != nil {
			common.SysError("error_unmarshalling_stream_response: " + err.Error())
			return
		}

		var content string
		err = json.Unmarshal(messageData.Content, &content)
		if err != nil {
			common.SysError("error_unmarshalling_stream_response: " + err.Error())
			return
		}

		*responseText += content

		openaiResponse := dto.ChatCompletionsStreamResponse{
			Id:      id,
			Object:  "chat.completion.chunk",
			Created: common.GetTimestamp(),
			Model:   info.UpstreamModelName,
		}

		choice := dto.ChatCompletionsStreamResponseChoice{
			Index: 0,
		}
		choice.Delta.SetContentString(content)
		openaiResponse.Choices = append(openaiResponse.Choices, choice)

		helper.ObjectData(c, openaiResponse)

	case "error":
		var errorData CozeError
		err := json.Unmarshal([]byte(data), &errorData)
		if err != nil {
			common.SysError("error_unmarshalling_stream_response: " + err.Error())
			return
		}

		common.SysError(fmt.Sprintf("stream event error: %d %s", errorData.Code, errorData.Message))
	}
}
