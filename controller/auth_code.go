package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"one-api/common"
	"one-api/model"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllAuthCodes(c *gin.Context) {
	p, _ := strconv.Atoi(c.Query("p"))
	size, _ := strconv.Atoi(c.Query("page_size"))
	if p < 1 {
		p = 1
	}
	if size <= 0 {
		size = common.ItemsPerPage
	} else if size > 100 {
		size = 100
	}

	authCodes, total, err := model.GetAllAuthCodes((p-1)*size, size)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"items":     authCodes,
			"total":     total,
			"page":      p,
			"page_size": size,
		},
	})
}

func SearchAuthCodes(c *gin.Context) {
	keyword := c.Query("keyword")
	p, _ := strconv.Atoi(c.Query("p"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if p < 1 {
		p = 1
	}
	if pageSize <= 0 {
		pageSize = common.ItemsPerPage
	}

	authCodes, total, err := model.SearchAuthCodes(keyword, (p-1)*pageSize, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"items":     authCodes,
			"total":     total,
			"page":      p,
			"page_size": pageSize,
		},
	})
}

// 获取可用的Token列表（用于授权码绑定）
func GetAvailableTokens(c *gin.Context) {
	userId := c.GetInt("id")

	tokens, err := model.GetAvailableTokensForAuthCode(userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    tokens,
	})
}

func GetAuthCode(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	authCode, err := model.GetAuthCodeById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    authCode,
	})
}

func AddAuthCode(c *gin.Context) {
	var authCode model.AuthCode
	err := json.NewDecoder(c.Request.Body).Decode(&authCode)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}

	// 验证必填字段
	if authCode.Code == "" || authCode.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "授权码和名称不能为空",
		})
		return
	}

	// 验证用户类型
	if !common.IsValidateRole(authCode.UserType) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的用户类型",
		})
		return
	}

	// 设置创建者
	authCode.CreatedBy = c.GetInt("id")

	// 清理授权码（去除空格）
	authCode.Code = strings.TrimSpace(authCode.Code)

	if err := authCode.Insert(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
}

func UpdateAuthCode(c *gin.Context) {
	statusOnly := c.Query("status_only")
	var authCode model.AuthCode
	err := json.NewDecoder(c.Request.Body).Decode(&authCode)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}

	// 获取原始授权码
	originAuthCode, err := model.GetAuthCodeById(authCode.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if statusOnly != "" {
		// 只更新状态
		originAuthCode.Status = authCode.Status
	} else {
		// 验证必填字段
		if authCode.Code == "" || authCode.Name == "" {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "授权码和名称不能为空",
			})
			return
		}

		// 验证用户类型
		if !common.IsValidateRole(authCode.UserType) {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无效的用户类型",
			})
			return
		}

		// 更新所有字段
		originAuthCode.Code = strings.TrimSpace(authCode.Code)
		originAuthCode.Name = authCode.Name
		originAuthCode.Description = authCode.Description
		originAuthCode.Status = authCode.Status
		originAuthCode.UserType = authCode.UserType
		originAuthCode.ExpiredTime = authCode.ExpiredTime
		originAuthCode.IsBot = authCode.IsBot
		originAuthCode.WxAutoXCode = authCode.WxAutoXCode
		originAuthCode.Group = authCode.Group
		originAuthCode.TokenId = authCode.TokenId
		// 只有在未激活状态或管理员操作时才允许修改机器码
		if originAuthCode.Status != 5 || authCode.MachineCode == "" {
			originAuthCode.MachineCode = authCode.MachineCode
		}
	}

	if err := originAuthCode.Update(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    originAuthCode,
	})
}

func DeleteAuthCode(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	err = model.DeleteAuthCodeById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
}

// 批量生成授权码
func BatchCreateAuthCodes(c *gin.Context) {
	var req struct {
		Count       int    `json:"count"`
		Name        string `json:"name"`
		Description string `json:"description"`
		UserType    int    `json:"user_type"`
		ExpiredTime int64  `json:"expired_time"`
		IsBot       bool   `json:"is_bot"`
		WxAutoXCode string `json:"wx_auto_x_code"`
		MachineCode string `json:"machine_code"`
		Group       string `json:"group"`
		TokenId     int    `json:"token_id"`
	}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}

	if req.Count <= 0 || req.Count > 100 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "生成数量必须在1-100之间",
		})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "名称不能为空",
		})
		return
	}

	if !common.IsValidateRole(req.UserType) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的用户类型",
		})
		return
	}

	createdBy := c.GetInt("id")
	var authCodes []model.AuthCode

	for i := 0; i < req.Count; i++ {
		authCode := model.AuthCode{
			Code:        common.GetRandomString(16), // 生成16位随机字符串
			Name:        req.Name + "_" + strconv.Itoa(i+1),
			Description: req.Description,
			UserType:    req.UserType,
			ExpiredTime: req.ExpiredTime,
			IsBot:       req.IsBot,
			WxAutoXCode: req.WxAutoXCode,
			MachineCode: req.MachineCode,
			Group:       req.Group,
			TokenId:     req.TokenId,
			CreatedBy:   createdBy,
			Status:      1,
		}

		if err := authCode.Insert(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "批量创建失败: " + err.Error(),
			})
			return
		}

		authCodes = append(authCodes, authCode)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "批量创建成功",
		"data":    authCodes,
	})
}

// 外部接口：绑定机器码
func BindMachineCode(c *gin.Context) {
	var req struct {
		AuthCode    string `json:"auth_code" binding:"required"`
		MachineCode string `json:"machine_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 获取授权码
	authCode, err := model.GetAuthCodeByCodeForExternal(req.AuthCode)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "授权码不存在",
		})
		return
	}

	// 绑定机器码
	if err := authCode.BindMachineCode(req.MachineCode); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "机器码绑定成功",
	})
}

// 生成验证挑战值
func generateChallenge(authCode string, machineCode string, timestamp int64) string {
	// 使用HMAC-SHA256生成挑战值，包含授权码、机器码和时间戳
	data := fmt.Sprintf("%s:%s:%d", authCode, machineCode, timestamp)
	h := hmac.New(sha256.New, []byte(common.CryptoSecret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// 验证挑战响应
func verifyChallenge(authCode string, machineCode string, timestamp int64, challenge string, response string) bool {
	// 重新生成挑战值
	expectedChallenge := generateChallenge(authCode, machineCode, timestamp)
	if expectedChallenge != challenge {
		return false
	}

	// 验证响应（客户端应该返回挑战值的SHA256哈希）
	h := sha256.New()
	h.Write([]byte(challenge))
	expectedResponse := hex.EncodeToString(h.Sum(nil))

	return expectedResponse == response
}

// 外部接口：验证授权码
func ValidateAuthCode(c *gin.Context) {
	var req struct {
		AuthCode    string `json:"auth_code" binding:"required"`
		MachineCode string `json:"machine_code" binding:"required"`
		Challenge   string `json:"challenge,omitempty"`
		Response    string `json:"response,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 获取授权码
	authCode, err := model.GetAuthCodeByCodeForExternal(req.AuthCode)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "授权码不存在",
		})
		return
	}

	// 检查授权码是否有效
	if !authCode.IsValid() {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "授权码无效或已过期",
		})
		return
	}

	// 验证机器码匹配
	if !authCode.ValidateWithMachineCode(req.MachineCode) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "机器码不匹配",
		})
		return
	}

	// 如果没有提供挑战响应，则生成新的挑战
	if req.Challenge == "" || req.Response == "" {
		timestamp := time.Now().Unix()
		challenge := generateChallenge(req.AuthCode, req.MachineCode, timestamp)

		c.JSON(http.StatusOK, gin.H{
			"success":    true,
			"message":    "请完成验证挑战",
			"challenge":  challenge,
			"timestamp":  timestamp,
			"expires_in": 300, // 5分钟有效期
		})
		return
	}

	// 验证挑战响应
	timestamp := time.Now().Unix()
	// 允许5分钟的时间窗口
	for i := int64(0); i < 300; i++ {
		if verifyChallenge(req.AuthCode, req.MachineCode, timestamp-i, req.Challenge, req.Response) {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "授权码验证成功",
				"data": gin.H{
					"user_type":      authCode.UserType,
					"is_bot":         authCode.IsBot,
					"wx_auto_x_code": authCode.WxAutoXCode,
					"expired_time":   authCode.ExpiredTime,
				},
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"message": "验证挑战失败",
	})
}

// 外部接口：根据授权码获取渠道列表
func GetChannelsByAuthCode(c *gin.Context) {
	// 从URL参数获取授权码
	authCodeParam := c.Query("auth_code")
	if authCodeParam == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "授权码参数不能为空",
		})
		return
	}

	// 从URL参数获取业务类型（可选）
	businessTypeParam := c.Query("business_type")
	var businessTypeFilter int = 0 // 0表示不过滤，返回所有类型
	if businessTypeParam != "" {
		if businessType, err := strconv.Atoi(businessTypeParam); err == nil {
			if businessType >= 1 && businessType <= 3 { // 1:对话, 2:应用, 3:工作流
				businessTypeFilter = businessType
			}
		}
	}

	// 获取授权码
	authCode, err := model.GetAuthCodeByCodeForExternal(authCodeParam)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "授权码不存在",
		})
		return
	}

	// 检查授权码是否有效
	if !authCode.IsValid() {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "授权码无效或已过期",
		})
		return
	}

	// 检查授权码状态（必须是激活状态才能获取渠道）
	if authCode.Status != 5 { // 5表示激活状态
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "授权码未激活",
		})
		return
	}

	// 获取渠道列表（简化版，不需要机器码验证）
	channels, err := model.GetChannelsByAuthCodeSimple(authCodeParam)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 构建返回的渠道信息（只返回必要信息，不暴露敏感数据）
	var channelList []gin.H
	for _, channel := range channels {
		// 如果指定了业务类型过滤，则只返回匹配的渠道
		if businessTypeFilter > 0 && channel.GetBusinessType() != businessTypeFilter {
			continue
		}

		channelInfo := gin.H{
			"id":            channel.Id,
			"name":          channel.Name,
			"type":          channel.Type,
			"business_type": channel.GetBusinessType(),
			"status":        channel.Status,
			"models":        channel.GetModels(),
			"group":         channel.Group,
			"priority":      channel.Priority,
			"weight":        channel.Weight,
		}
		channelList = append(channelList, channelInfo)
	}

	// 解析授权码的分组信息
	var authGroups []string
	if authCode.Group != "" {
		groups := strings.Split(strings.Trim(authCode.Group, ","), ",")
		for _, group := range groups {
			group = strings.TrimSpace(group)
			if group != "" {
				authGroups = append(authGroups, group)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取渠道列表成功",
		"data": gin.H{
			"auth_groups": authGroups,
			"channels":    channelList,
			"total":       len(channelList),
		},
	})
}

// 外部接口：根据授权码获取绑定的API密钥
func GetApiKeyByAuthCode(c *gin.Context) {
	// 从URL参数获取授权码
	authCodeParam := c.Query("auth_code")
	if authCodeParam == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "授权码参数不能为空",
		})
		return
	}

	// 获取授权码
	authCode, err := model.GetAuthCodeByCodeForExternal(authCodeParam)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "授权码不存在",
		})
		return
	}

	// 检查授权码是否有效
	if !authCode.IsValid() {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "授权码无效或已过期",
		})
		return
	}

	// 检查授权码状态（必须是激活状态才能获取API密钥）
	if authCode.Status != 5 { // 5表示激活状态
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "授权码未激活",
		})
		return
	}

	// 检查是否绑定了API密钥
	if authCode.TokenId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "授权码未绑定API密钥",
		})
		return
	}

	// 获取绑定的Token信息
	token, err := authCode.GetBoundToken()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "获取API密钥失败: " + err.Error(),
		})
		return
	}

	if token == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "绑定的API密钥不存在或已被禁用",
		})
		return
	}

	// 返回API密钥信息（不暴露敏感信息）
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取API密钥成功",
		"data": gin.H{
			"token_id":        token.Id,
			"token_name":      token.Name,
			"api_key":         token.Key,
			"status":          token.Status,
			"expired_time":    token.ExpiredTime,
			"remain_quota":    token.RemainQuota,
			"unlimited_quota": token.UnlimitedQuota,
			"group":           token.Group,
			"auth_code_info": gin.H{
				"code":      authCode.Code,
				"name":      authCode.Name,
				"user_type": authCode.UserType,
				"is_bot":    authCode.IsBot,
				"group":     authCode.Group,
			},
		},
	})
}
