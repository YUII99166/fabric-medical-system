$URL = "http://127.0.0.1:8888/api/v2/register"

Write-Host "开始注册用户..."

# 管理员1
$body1 = @{
    username = "admin1"
    password = "admin123456"
    account_name = "系统管理员1"
    role = "管理员"
    organization = "RegCenterMSP"
    organization_name = "监管中心"
    department = "系统管理部"
    doctor_title = ""
    age = 35
    gender = "男"
} | ConvertTo-Json -Compress

try {
    $response1 = Invoke-RestMethod -Uri $URL -Method Post -Body $body1 -ContentType "application/json; charset=utf-8"
    Write-Host "管理员1注册成功: $($response1.message)"
} catch {
    Write-Host "管理员1注册失败: $($_.Exception.Message)"
}

# 管理员2
$body2 = @{
    username = "admin2"
    password = "admin123456"
    account_name = "系统管理员2"
    role = "管理员"
    organization = "RegCenterMSP"
    organization_name = "监管中心"
    department = "系统管理部"
    doctor_title = ""
    age = 32
    gender = "女"
} | ConvertTo-Json -Compress

try {
    $response2 = Invoke-RestMethod -Uri $URL -Method Post -Body $body2 -ContentType "application/json; charset=utf-8"
    Write-Host "管理员2注册成功: $($response2.message)"
} catch {
    Write-Host "管理员2注册失败: $($_.Exception.Message)"
}

# 医生1
$body3 = @{
    username = "doctor1"
    password = "doctor123456"
    account_name = "张医生"
    role = "医生"
    organization = "TaobaoMSP"
    organization_name = "协和医院"
    department = "心内科"
    doctor_title = "主任医师"
    age = 45
    gender = "男"
} | ConvertTo-Json -Compress

try {
    $response3 = Invoke-RestMethod -Uri $URL -Method Post -Body $body3 -ContentType "application/json; charset=utf-8"
    Write-Host "医生1注册成功: $($response3.message)"
} catch {
    Write-Host "医生1注册失败: $($_.Exception.Message)"
}

# 医生2
$body4 = @{
    username = "doctor2"
    password = "doctor123456"
    account_name = "王医生"
    role = "医生"
    organization = "JDMSP"
    organization_name = "301医院"
    department = "呼吸科"
    doctor_title = "副主任医师"
    age = 38
    gender = "女"
} | ConvertTo-Json -Compress

try {
    $response4 = Invoke-RestMethod -Uri $URL -Method Post -Body $body4 -ContentType "application/json; charset=utf-8"
    Write-Host "医生2注册成功: $($response4.message)"
} catch {
    Write-Host "医生2注册失败: $($_.Exception.Message)"
}

# 病人1
$body5 = @{
    username = "patient1"
    password = "patient123456"
    account_name = "患者张三"
    role = "病人"
    age = 28
    gender = "男"
} | ConvertTo-Json -Compress

try {
    $response5 = Invoke-RestMethod -Uri $URL -Method Post -Body $body5 -ContentType "application/json; charset=utf-8"
    Write-Host "病人1注册成功: $($response5.message)"
} catch {
    Write-Host "病人1注册失败: $($_.Exception.Message)"
}

# 病人2
$body6 = @{
    username = "patient2"
    password = "patient123456"
    account_name = "患者李四"
    role = "病人"
    age = 35
    gender = "女"
} | ConvertTo-Json -Compress

try {
    $response6 = Invoke-RestMethod -Uri $URL -Method Post -Body $body6 -ContentType "application/json; charset=utf-8"
    Write-Host "病人2注册成功: $($response6.message)"
} catch {
    Write-Host "病人2注册失败: $($_.Exception.Message)"
}

# 药店1
$body7 = @{
    username = "pharmacy1"
    password = "pharmacy123456"
    account_name = "同仁堂药店"
    role = "药店"
    age = 0
    gender = ""
} | ConvertTo-Json -Compress

try {
    $response7 = Invoke-RestMethod -Uri $URL -Method Post -Body $body7 -ContentType "application/json; charset=utf-8"
    Write-Host "药店1注册成功: $($response7.message)"
} catch {
    Write-Host "药店1注册失败: $($_.Exception.Message)"
}

# 药店2
$body8 = @{
    username = "pharmacy2"
    password = "pharmacy123456"
    account_name = "百姓大药房"
    role = "药店"
    age = 0
    gender = ""
} | ConvertTo-Json -Compress

try {
    $response8 = Invoke-RestMethod -Uri $URL -Method Post -Body $body8 -ContentType "application/json; charset=utf-8"
    Write-Host "药店2注册成功: $($response8.message)"
} catch {
    Write-Host "药店2注册失败: $($_.Exception.Message)"
}

Write-Host "`n所有用户注册完成！"
