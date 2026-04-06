-- ============================================
-- 修复用户表中缺失的 organization_name
-- ============================================

USE fabric;

-- 更新所有 organization 为 TaobaoMSP 但 organization_name 为空的用户
UPDATE users 
SET organization_name = '协和医院' 
WHERE organization = 'TaobaoMSP' AND (organization_name IS NULL OR organization_name = '');

-- 更新所有 organization 为 JDMSP 但 organization_name 为空的用户
UPDATE users 
SET organization_name = '301医院' 
WHERE organization = 'JDMSP' AND (organization_name IS NULL OR organization_name = '');

-- 更新所有 organization 为 WenjinMSP 但 organization_name 为空的用户
UPDATE users 
SET organization_name = '温江社区医疗中心' 
WHERE organization = 'WenjinMSP' AND (organization_name IS NULL OR organization_name = '');

-- 更新所有 organization 为 RegCenterMSP 但 organization_name 为空的用户
UPDATE users 
SET organization_name = '监管中心' 
WHERE organization = 'RegCenterMSP' AND (organization_name IS NULL OR organization_name = '');

-- 查看更新结果
SELECT id, username, account_name, role, organization, organization_name 
FROM users 
WHERE organization IS NOT NULL AND organization != ''
ORDER BY id;
