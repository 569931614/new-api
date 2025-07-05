package model

import (
	"errors"
	"fmt"
	"one-api/common"
	"strings"
	"time"

	"gorm.io/gorm"
)

type AuthCode struct {
	Id          int            `json:"id"`
	Code        string         `json:"code" gorm:"type:varchar(64);uniqueIndex"`
	Name        string         `json:"name" gorm:"index"`
	Description string         `json:"description" gorm:"type:text"`
	Status      int            `json:"status" gorm:"default:1"`               // 1: 启用, 2: 禁用, 3: 已使用, 4: 待激活, 5: 激活
	UserType    int            `json:"user_type" gorm:"default:1"`            // 1: 普通用户, 10: 管理员, 100: 超级管理员
	ExpiredTime int64          `json:"expired_time" gorm:"bigint;default:-1"` // -1 表示永不过期
	CreatedTime int64          `json:"created_time" gorm:"bigint"`
	UsedTime    int64          `json:"used_time" gorm:"bigint;default:0"`
	UsedUserId  int            `json:"used_user_id" gorm:"default:0"`
	IsBot       bool           `json:"is_bot" gorm:"default:false"`             // 是否为机器人账户
	WxAutoXCode string         `json:"wx_auto_x_code" gorm:"type:varchar(255)"` // wxautox码
	MachineCode string         `json:"machine_code" gorm:"type:varchar(255)"`   // 机器码
	Group       string         `json:"group" gorm:"type:varchar(255);index"`    // 分组名称（支持多个分组，用逗号分隔）
	TokenId     int            `json:"token_id" gorm:"default:0;index"`         // 绑定的API密钥ID
	CreatedBy   int            `json:"created_by" gorm:"index"`                 // 创建者ID
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func GetAllAuthCodes(startIdx int, num int) (authCodes []*AuthCode, total int64, err error) {
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取总数
	if err = tx.Model(&AuthCode{}).Count(&total).Error; err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	// 获取分页数据
	err = tx.Order("id desc").Limit(num).Offset(startIdx).Find(&authCodes).Error
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		return nil, 0, err
	}

	return authCodes, total, nil
}

func SearchAuthCodes(keyword string, startIdx int, num int) (authCodes []*AuthCode, total int64, err error) {
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := tx.Model(&AuthCode{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err = query.Count(&total).Error; err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Order("id desc").Limit(num).Offset(startIdx).Find(&authCodes).Error
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		return nil, 0, err
	}

	return authCodes, total, nil
}

func GetAuthCodeById(id int) (*AuthCode, error) {
	if id == 0 {
		return nil, errors.New("id 为空！")
	}
	authCode := AuthCode{Id: id}
	err := DB.First(&authCode, "id = ?", id).Error
	return &authCode, err
}

func GetAuthCodeByCode(code string) (*AuthCode, error) {
	if code == "" {
		return nil, errors.New("授权码为空！")
	}
	var authCode AuthCode
	err := DB.First(&authCode, "code = ? AND (status = 1 OR status = 4 OR status = 5)", code).Error
	return &authCode, err
}

// 根据授权码获取（用于外部接口，不限制状态）
func GetAuthCodeByCodeForExternal(code string) (*AuthCode, error) {
	if code == "" {
		return nil, errors.New("授权码为空！")
	}
	var authCode AuthCode
	err := DB.First(&authCode, "code = ?", code).Error
	return &authCode, err
}

func (authCode *AuthCode) Insert() error {
	authCode.CreatedTime = time.Now().Unix()
	result := DB.Create(authCode)
	return result.Error
}

func (authCode *AuthCode) Update() error {
	return DB.Model(authCode).Updates(authCode).Error
}

func (authCode *AuthCode) Delete() error {
	if authCode.Id == 0 {
		return errors.New("id 为空！")
	}
	return DB.Delete(authCode).Error
}

// 使用授权码
func (authCode *AuthCode) Use(userId int) error {
	if authCode.Status != 1 {
		return errors.New("授权码已被禁用或已使用")
	}

	// 检查是否过期
	if authCode.ExpiredTime != -1 && authCode.ExpiredTime < time.Now().Unix() {
		return errors.New("授权码已过期")
	}

	authCode.Status = 3 // 已使用
	authCode.UsedTime = time.Now().Unix()
	authCode.UsedUserId = userId

	return DB.Model(authCode).Updates(map[string]interface{}{
		"status":       authCode.Status,
		"used_time":    authCode.UsedTime,
		"used_user_id": authCode.UsedUserId,
	}).Error
}

// 绑定机器码并激活
func (authCode *AuthCode) BindMachineCode(machineCode string) error {
	if authCode.Status != 1 && authCode.Status != 4 {
		return errors.New("授权码状态不允许绑定机器码")
	}

	// 检查是否过期
	if authCode.ExpiredTime != -1 && authCode.ExpiredTime < time.Now().Unix() {
		return errors.New("授权码已过期")
	}

	// 检查机器码是否已被其他授权码绑定
	var existingAuthCode AuthCode
	err := DB.Where("machine_code = ? AND machine_code != '' AND id != ?", machineCode, authCode.Id).First(&existingAuthCode).Error
	if err == nil {
		return errors.New("该机器码已被其他授权码绑定")
	}

	authCode.MachineCode = machineCode
	authCode.Status = 5 // 激活状态

	return DB.Model(authCode).Updates(map[string]interface{}{
		"machine_code": authCode.MachineCode,
		"status":       authCode.Status,
	}).Error
}

// 检查授权码是否有效
func (authCode *AuthCode) IsValid() bool {
	if authCode.Status != 1 && authCode.Status != 5 {
		return false
	}

	// 检查是否过期
	if authCode.ExpiredTime != -1 && authCode.ExpiredTime < time.Now().Unix() {
		return false
	}

	return true
}

// 验证授权码和机器码匹配
func (authCode *AuthCode) ValidateWithMachineCode(machineCode string) bool {
	if !authCode.IsValid() {
		return false
	}

	// 如果授权码已绑定机器码，则必须匹配
	if authCode.MachineCode != "" {
		return authCode.MachineCode == machineCode
	}

	// 如果未绑定机器码，则允许通过
	return true
}

// 根据授权码获取可用的渠道列表
func GetChannelsByAuthCode(authCode string, machineCode string) ([]*Channel, error) {
	// 首先验证授权码
	auth, err := GetAuthCodeByCodeForExternal(authCode)
	if err != nil {
		return nil, errors.New("授权码不存在")
	}

	// 验证授权码状态和机器码
	if !auth.ValidateWithMachineCode(machineCode) {
		return nil, errors.New("授权码验证失败")
	}

	// 如果没有设置分组，返回空列表
	if auth.Group == "" {
		return []*Channel{}, nil
	}

	// 解析多个分组（用逗号分隔）
	groups := strings.Split(strings.Trim(auth.Group, ","), ",")
	if len(groups) == 0 {
		return []*Channel{}, nil
	}

	// 根据分组获取渠道列表
	var channels []*Channel

	// 构建查询条件，支持多个分组
	var whereClause string
	var args []interface{}

	for i, group := range groups {
		group = strings.TrimSpace(group)
		if group == "" {
			continue
		}

		if i > 0 {
			whereClause += " OR "
		}

		// 使用LIKE查询支持渠道的多分组匹配
		whereClause += "(',' || `group` || ',') LIKE ?"
		args = append(args, "%,"+group+",%")
	}

	if whereClause == "" {
		return []*Channel{}, nil
	}

	whereClause = "(" + whereClause + ") AND status = ?"
	args = append(args, 1)

	err = DB.Where(whereClause, args...).Find(&channels).Error
	if err != nil {
		return nil, errors.New("获取渠道列表失败")
	}

	return channels, nil
}

// 根据授权码获取可用的渠道列表（简化版，不需要机器码验证）
func GetChannelsByAuthCodeSimple(authCode string) ([]*Channel, error) {
	// 首先验证授权码
	auth, err := GetAuthCodeByCodeForExternal(authCode)
	if err != nil {
		common.SysLog(fmt.Sprintf("授权码查询失败: %s, 错误: %v", authCode, err))
		return nil, errors.New("授权码不存在")
	}

	common.SysLog(fmt.Sprintf("授权码信息: code=%s, status=%d, group='%s'", auth.Code, auth.Status, auth.Group))

	// 验证授权码状态（必须是激活状态）
	if auth.Status != 5 {
		common.SysLog(fmt.Sprintf("授权码状态不是激活状态: %d", auth.Status))
		return nil, errors.New("授权码未激活")
	}

	// 如果没有设置分组，返回所有启用的渠道（修改逻辑）
	if auth.Group == "" {
		common.SysLog("授权码没有设置分组，返回所有启用的渠道")
		var channels []*Channel
		err = DB.Where("status = ?", 1).Find(&channels).Error
		if err != nil {
			common.SysLog(fmt.Sprintf("查询所有启用渠道失败: %v", err))
			return nil, errors.New("获取渠道列表失败")
		}
		common.SysLog(fmt.Sprintf("返回所有启用渠道，共 %d 个", len(channels)))
		return channels, nil
	}

	// 解析多个分组（用逗号分隔）
	groups := strings.Split(strings.Trim(auth.Group, ","), ",")
	common.SysLog(fmt.Sprintf("原始分组: '%s', 解析后: %v", auth.Group, groups))

	var cleanGroups []string
	for _, group := range groups {
		group = strings.TrimSpace(group)
		if group != "" {
			cleanGroups = append(cleanGroups, group)
		}
	}

	if len(cleanGroups) == 0 {
		common.SysLog("没有有效的分组")
		return []*Channel{}, nil
	}

	common.SysLog(fmt.Sprintf("有效分组: %v", cleanGroups))

	// 根据分组获取渠道列表 - 使用更灵活的匹配逻辑
	var channels []*Channel
	var whereConditions []string
	var args []interface{}

	for _, group := range cleanGroups {
		// 支持多种分组格式的匹配
		// 1. 精确匹配: group = 'groupname'
		// 2. 开头匹配: group LIKE 'groupname,%'
		// 3. 中间匹配: group LIKE '%,groupname,%'
		// 4. 结尾匹配: group LIKE '%,groupname'
		condition := "(`group` = ? OR `group` LIKE ? OR `group` LIKE ? OR `group` LIKE ?)"
		whereConditions = append(whereConditions, condition)

		args = append(args, group)           // 精确匹配
		args = append(args, group+",%")      // 开头匹配
		args = append(args, "%,"+group+",%") // 中间匹配
		args = append(args, "%,"+group)      // 结尾匹配
	}

	whereClause := "(" + strings.Join(whereConditions, " OR ") + ") AND status = ?"
	args = append(args, 1)

	common.SysLog(fmt.Sprintf("查询条件: %s", whereClause))
	common.SysLog(fmt.Sprintf("查询参数: %v", args))

	err = DB.Where(whereClause, args...).Find(&channels).Error
	if err != nil {
		common.SysLog(fmt.Sprintf("查询渠道失败: %v", err))
		return nil, errors.New("获取渠道列表失败")
	}

	common.SysLog(fmt.Sprintf("查询到 %d 个匹配的渠道", len(channels)))
	for _, ch := range channels {
		common.SysLog(fmt.Sprintf("  渠道: ID=%d, 名称=%s, 分组='%s', 状态=%d", ch.Id, ch.Name, ch.Group, ch.Status))
	}

	return channels, nil
}

// 调试函数：获取授权码和渠道的详细信息
func DebugAuthCodeChannels(authCode string) map[string]interface{} {
	result := make(map[string]interface{})

	// 1. 获取授权码信息
	auth, err := GetAuthCodeByCodeForExternal(authCode)
	if err != nil {
		result["auth_code_error"] = err.Error()
		return result
	}

	result["auth_code"] = map[string]interface{}{
		"code":   auth.Code,
		"status": auth.Status,
		"group":  auth.Group,
		"name":   auth.Name,
	}

	// 2. 获取所有启用的渠道
	var allChannels []*Channel
	DB.Where("status = ?", 1).Find(&allChannels)

	var channelInfo []map[string]interface{}
	for _, ch := range allChannels {
		channelInfo = append(channelInfo, map[string]interface{}{
			"id":     ch.Id,
			"name":   ch.Name,
			"group":  ch.Group,
			"status": ch.Status,
			"type":   ch.Type,
		})
	}
	result["all_enabled_channels"] = channelInfo

	// 3. 测试分组匹配
	if auth.Group != "" {
		groups := strings.Split(strings.Trim(auth.Group, ","), ",")
		var cleanGroups []string
		for _, group := range groups {
			group = strings.TrimSpace(group)
			if group != "" {
				cleanGroups = append(cleanGroups, group)
			}
		}

		result["parsed_groups"] = cleanGroups

		// 测试每个分组的匹配情况
		var matchResults []map[string]interface{}
		for _, group := range cleanGroups {
			var matchedChannels []*Channel

			// 使用新的匹配逻辑
			condition := "(`group` = ? OR `group` LIKE ? OR `group` LIKE ? OR `group` LIKE ?) AND status = ?"
			args := []interface{}{
				group,               // 精确匹配
				group + ",%",        // 开头匹配
				"%," + group + ",%", // 中间匹配
				"%," + group,        // 结尾匹配
				1,                   // 状态过滤
			}

			DB.Where(condition, args...).Find(&matchedChannels)

			var matched []map[string]interface{}
			for _, ch := range matchedChannels {
				matched = append(matched, map[string]interface{}{
					"id":    ch.Id,
					"name":  ch.Name,
					"group": ch.Group,
				})
			}

			matchResults = append(matchResults, map[string]interface{}{
				"group":            group,
				"matched_channels": matched,
				"count":            len(matchedChannels),
			})
		}
		result["match_results"] = matchResults
	}

	return result
}

func DeleteAuthCodeById(id int) error {
	if id == 0 {
		return errors.New("id 为空！")
	}
	authCode := AuthCode{Id: id}
	return authCode.Delete()
}

// 获取绑定的Token信息
func (authCode *AuthCode) GetBoundToken() (*Token, error) {
	if authCode.TokenId == 0 {
		return nil, nil // 没有绑定Token
	}

	var token Token
	err := DB.First(&token, "id = ? AND status = 1", authCode.TokenId).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// 获取可用的Token列表（用于前端选择）
func GetAvailableTokensForAuthCode(userId int) ([]*Token, error) {
	var tokens []*Token
	var err error

	if userId == 0 {
		// 如果userId为0，获取所有用户的Token（管理员模式）
		err = DB.Where("status = 1").
			Select("id, name, key, status, expired_time, user_id").
			Order("created_time DESC").
			Find(&tokens).Error
	} else {
		// 获取指定用户的Token
		err = DB.Where("user_id = ? AND status = 1", userId).
			Select("id, name, key, status, expired_time, user_id").
			Order("created_time DESC").
			Find(&tokens).Error
	}

	if err != nil {
		return nil, err
	}

	// 清理敏感信息，只保留前8位和后4位
	for _, token := range tokens {
		if len(token.Key) > 12 {
			token.Key = token.Key[:8] + "****" + token.Key[len(token.Key)-4:]
		}

		// 如果是管理员模式，添加用户信息到名称中
		if userId == 0 {
			// 获取用户名
			var user User
			if err := DB.Select("username").First(&user, token.UserId).Error; err == nil {
				token.Name = fmt.Sprintf("[%s] %s", user.Username, token.Name)
			}
		}
	}

	return tokens, nil
}
