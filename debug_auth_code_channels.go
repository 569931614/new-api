package main

import (
	"encoding/json"
	"fmt"
	"log"
	"one-api/common"
	"one-api/model"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 初始化数据库连接
	common.SetupGormDB()
	
	// 要调试的授权码
	authCodeParam := "CWuYSBCSbzDykmJY"
	
	fmt.Printf("=== 调试授权码: %s ===\n\n", authCodeParam)
	
	// 1. 检查授权码是否存在
	fmt.Println("1. 检查授权码是否存在...")
	authCode, err := model.GetAuthCodeByCodeForExternal(authCodeParam)
	if err != nil {
		log.Fatalf("授权码不存在: %v", err)
	}
	
	authCodeJSON, _ := json.MarshalIndent(authCode, "", "  ")
	fmt.Printf("授权码信息:\n%s\n\n", authCodeJSON)
	
	// 2. 检查授权码状态
	fmt.Println("2. 检查授权码状态...")
	fmt.Printf("状态: %d (5=激活)\n", authCode.Status)
	fmt.Printf("是否有效: %v\n", authCode.IsValid())
	fmt.Printf("分组: '%s'\n\n", authCode.Group)
	
	if authCode.Status != 5 {
		fmt.Printf("❌ 授权码状态不是激活状态(5)，当前状态: %d\n", authCode.Status)
		return
	}
	
	if authCode.Group == "" {
		fmt.Println("❌ 授权码没有设置分组")
		return
	}
	
	// 3. 解析分组
	fmt.Println("3. 解析授权码分组...")
	groups := strings.Split(strings.Trim(authCode.Group, ","), ",")
	fmt.Printf("原始分组字符串: '%s'\n", authCode.Group)
	fmt.Printf("解析后的分组: %v\n", groups)
	
	var cleanGroups []string
	for _, group := range groups {
		group = strings.TrimSpace(group)
		if group != "" {
			cleanGroups = append(cleanGroups, group)
		}
	}
	fmt.Printf("清理后的分组: %v\n\n", cleanGroups)
	
	if len(cleanGroups) == 0 {
		fmt.Println("❌ 没有有效的分组")
		return
	}
	
	// 4. 查询所有渠道
	fmt.Println("4. 查询所有渠道...")
	var allChannels []*model.Channel
	err = model.DB.Find(&allChannels).Error
	if err != nil {
		log.Fatalf("查询所有渠道失败: %v", err)
	}
	fmt.Printf("数据库中总共有 %d 个渠道\n\n", len(allChannels))
	
	// 5. 查询启用的渠道
	fmt.Println("5. 查询启用的渠道...")
	var enabledChannels []*model.Channel
	err = model.DB.Where("status = ?", 1).Find(&enabledChannels).Error
	if err != nil {
		log.Fatalf("查询启用渠道失败: %v", err)
	}
	fmt.Printf("启用的渠道有 %d 个\n\n", len(enabledChannels))
	
	// 6. 检查每个分组的渠道匹配
	fmt.Println("6. 检查分组匹配...")
	for _, group := range cleanGroups {
		fmt.Printf("检查分组: '%s'\n", group)
		
		// 使用和代码中相同的查询逻辑
		whereClause := "(',' || `group` || ',') LIKE ?"
		args := []interface{}{"%," + group + ",%"}
		
		var matchedChannels []*model.Channel
		err = model.DB.Where(whereClause, args...).Find(&matchedChannels).Error
		if err != nil {
			fmt.Printf("  查询失败: %v\n", err)
			continue
		}
		
		fmt.Printf("  匹配的渠道数量: %d\n", len(matchedChannels))
		for _, ch := range matchedChannels {
			fmt.Printf("    - ID:%d, 名称:%s, 分组:'%s', 状态:%d\n", 
				ch.Id, ch.Name, ch.Group, ch.Status)
		}
		
		// 再加上状态过滤
		var enabledMatchedChannels []*model.Channel
		err = model.DB.Where(whereClause+" AND status = ?", append(args, 1)...).Find(&enabledMatchedChannels).Error
		if err != nil {
			fmt.Printf("  查询启用渠道失败: %v\n", err)
			continue
		}
		
		fmt.Printf("  启用的匹配渠道数量: %d\n", len(enabledMatchedChannels))
		for _, ch := range enabledMatchedChannels {
			fmt.Printf("    - ID:%d, 名称:%s, 分组:'%s', 状态:%d\n", 
				ch.Id, ch.Name, ch.Group, ch.Status)
		}
		fmt.Println()
	}
	
	// 7. 使用完整的查询逻辑
	fmt.Println("7. 使用完整的查询逻辑...")
	channels, err := model.GetChannelsByAuthCodeSimple(authCodeParam)
	if err != nil {
		fmt.Printf("❌ GetChannelsByAuthCodeSimple 失败: %v\n", err)
		return
	}
	
	fmt.Printf("✅ GetChannelsByAuthCodeSimple 返回 %d 个渠道\n", len(channels))
	for _, ch := range channels {
		fmt.Printf("  - ID:%d, 名称:%s, 类型:%d, 分组:'%s', 状态:%d\n", 
			ch.Id, ch.Name, ch.Type, ch.Group, ch.Status)
	}
	
	if len(channels) == 0 {
		fmt.Println("\n❌ 没有找到匹配的渠道，可能的原因:")
		fmt.Println("1. 授权码的分组与渠道的分组不匹配")
		fmt.Println("2. 匹配的渠道都被禁用了")
		fmt.Println("3. 分组格式不正确")
		
		fmt.Println("\n建议检查:")
		fmt.Println("1. 确认渠道的分组格式是否正确（应该用逗号分隔，如: 'group1,group2'）")
		fmt.Println("2. 确认渠道状态是否为启用(1)")
		fmt.Println("3. 确认授权码分组与渠道分组完全匹配")
	}
}
