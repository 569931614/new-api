-- 调试授权码 CWuYSBCSbzDykmJY 获取渠道为空的问题

-- 1. 检查授权码是否存在
SELECT 'Step 1: 检查授权码' as step;
SELECT id, code, name, status, group_name as 'group', expired_time, created_time 
FROM auth_codes 
WHERE code = 'CWuYSBCSbzDykmJY';

-- 2. 检查所有渠道
SELECT 'Step 2: 检查所有渠道' as step;
SELECT id, name, type, status, group_name as 'group', priority, weight 
FROM channels 
ORDER BY id;

-- 3. 检查启用的渠道
SELECT 'Step 3: 检查启用的渠道' as step;
SELECT id, name, type, status, group_name as 'group', priority, weight 
FROM channels 
WHERE status = 1
ORDER BY id;

-- 4. 检查分组匹配（假设授权码分组是 'default'）
SELECT 'Step 4: 检查分组匹配' as step;
SELECT id, name, type, status, group_name as 'group', priority, weight 
FROM channels 
WHERE (',' || group_name || ',') LIKE '%,default,%' AND status = 1;

-- 5. 检查所有可能的分组匹配
SELECT 'Step 5: 检查所有可能的分组' as step;
SELECT DISTINCT group_name as 'group' FROM channels WHERE group_name IS NOT NULL AND group_name != '';
