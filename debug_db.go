package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 打开数据库
	db, err := sql.Open("sqlite3", "./one-api.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	authCode := "CWuYSBCSbzDykmJY"
	
	fmt.Printf("=== 调试授权码: %s ===\n\n", authCode)
	
	// 1. 检查授权码
	fmt.Println("1. 检查授权码...")
	var id int
	var code, name, group string
	var status int
	var expiredTime int64
	
	err = db.QueryRow("SELECT id, code, name, status, `group`, expired_time FROM auth_codes WHERE code = ?", authCode).
		Scan(&id, &code, &name, &status, &group, &expiredTime)
	
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("❌ 授权码 %s 不存在\n", authCode)
		} else {
			fmt.Printf("❌ 查询授权码失败: %v\n", err)
		}
		return
	}
	
	fmt.Printf("✅ 授权码信息:\n")
	fmt.Printf("  ID: %d\n", id)
	fmt.Printf("  代码: %s\n", code)
	fmt.Printf("  名称: %s\n", name)
	fmt.Printf("  状态: %d (5=激活)\n", status)
	fmt.Printf("  分组: '%s'\n", group)
	fmt.Printf("  过期时间: %d\n\n", expiredTime)
	
	if status != 5 {
		fmt.Printf("❌ 授权码状态不是激活状态(5)，当前状态: %d\n", status)
		return
	}
	
	if group == "" {
		fmt.Println("❌ 授权码没有设置分组")
		return
	}
	
	// 2. 检查所有渠道
	fmt.Println("2. 检查所有渠道...")
	rows, err := db.Query("SELECT id, name, type, status, `group` FROM channels ORDER BY id")
	if err != nil {
		fmt.Printf("❌ 查询渠道失败: %v\n", err)
		return
	}
	defer rows.Close()
	
	var totalChannels, enabledChannels int
	fmt.Println("所有渠道:")
	for rows.Next() {
		var chId, chType, chStatus int
		var chName, chGroup string
		
		err = rows.Scan(&chId, &chName, &chType, &chStatus, &chGroup)
		if err != nil {
			fmt.Printf("❌ 扫描渠道数据失败: %v\n", err)
			continue
		}
		
		totalChannels++
		if chStatus == 1 {
			enabledChannels++
		}
		
		fmt.Printf("  ID:%d, 名称:%s, 类型:%d, 状态:%d, 分组:'%s'\n", 
			chId, chName, chType, chStatus, chGroup)
	}
	
	fmt.Printf("\n总渠道数: %d, 启用渠道数: %d\n\n", totalChannels, enabledChannels)
	
	// 3. 检查分组匹配
	fmt.Println("3. 检查分组匹配...")
	fmt.Printf("授权码分组: '%s'\n", group)
	
	// 使用和代码中相同的查询逻辑
	query := "SELECT id, name, type, status, `group` FROM channels WHERE (',' || `group` || ',') LIKE ? AND status = ?"
	likePattern := "%," + group + ",%"
	
	fmt.Printf("查询条件: %s\n", query)
	fmt.Printf("LIKE 模式: '%s'\n", likePattern)
	
	rows2, err := db.Query(query, likePattern, 1)
	if err != nil {
		fmt.Printf("❌ 查询匹配渠道失败: %v\n", err)
		return
	}
	defer rows2.Close()
	
	var matchedChannels int
	fmt.Println("匹配的启用渠道:")
	for rows2.Next() {
		var chId, chType, chStatus int
		var chName, chGroup string
		
		err = rows2.Scan(&chId, &chName, &chType, &chStatus, &chGroup)
		if err != nil {
			fmt.Printf("❌ 扫描匹配渠道数据失败: %v\n", err)
			continue
		}
		
		matchedChannels++
		fmt.Printf("  ✅ ID:%d, 名称:%s, 类型:%d, 状态:%d, 分组:'%s'\n", 
			chId, chName, chType, chStatus, chGroup)
	}
	
	fmt.Printf("\n匹配的启用渠道数: %d\n", matchedChannels)
	
	if matchedChannels == 0 {
		fmt.Println("\n❌ 没有找到匹配的渠道，可能的原因:")
		fmt.Println("1. 授权码的分组与渠道的分组不匹配")
		fmt.Println("2. 匹配的渠道都被禁用了")
		fmt.Println("3. 分组格式不正确")
		
		// 显示所有不同的分组
		fmt.Println("\n所有渠道的分组:")
		rows3, err := db.Query("SELECT DISTINCT `group` FROM channels WHERE `group` IS NOT NULL AND `group` != ''")
		if err == nil {
			defer rows3.Close()
			for rows3.Next() {
				var distinctGroup string
				rows3.Scan(&distinctGroup)
				fmt.Printf("  '%s'\n", distinctGroup)
			}
		}
	}
}
