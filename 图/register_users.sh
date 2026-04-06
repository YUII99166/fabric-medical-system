#!/bin/bash

URL="http://127.0.0.1:8888/api/v2/register"

echo "开始注册用户..."

# 管理员1
curl -X POST $URL -H "Content-Type: application/json" -d '{"username":"admin1","password":"admin123456","account_name":"系统管理员1","role":"管理员","organization":"RegCenterMSP","organization_name":"监管中心","department":"系统管理部","doctor_title":"","age":35,"gender":"男"}'
echo -e "\n管理员1注册完成"

# 管理员2
curl -X POST $URL -H "Content-Type: application/json" -d '{"username":"admin2","password":"admin123456","account_name":"系统管理员2","role":"管理员","organization":"RegCenterMSP","organization_name":"监管中心","department":"系统管理部","doctor_title":"","age":32,"gender":"女"}'
echo -e "\n管理员2注册完成"

# 医生1
curl -X POST $URL -H "Content-Type: application/json" -d '{"username":"doctor1","password":"doctor123456","account_name":"张医生","role":"医生","organization":"TaobaoMSP","organization_name":"协和医院","department":"心内科","doctor_title":"主任医师","age":45,"gender":"男"}'
echo -e "\n医生1注册完成"

# 医生2
curl -X POST $URL -H "Content-Type: application/json" -d '{"username":"doctor2","password":"doctor123456","account_name":"王医生","role":"医生","organization":"JDMSP","organization_name":"301医院","department":"呼吸科","doctor_title":"副主任医师","age":38,"gender":"女"}'
echo -e "\n医生2注册完成"

# 病人1
curl -X POST $URL -H "Content-Type: application/json" -d '{"username":"patient1","password":"patient123456","account_name":"患者张三","role":"病人","age":28,"gender":"男"}'
echo -e "\n病人1注册完成"

# 病人2
curl -X POST $URL -H "Content-Type: application/json" -d '{"username":"patient2","password":"patient123456","account_name":"患者李四","role":"病人","age":35,"gender":"女"}'
echo -e "\n病人2注册完成"

# 药店1
curl -X POST $URL -H "Content-Type: application/json" -d '{"username":"pharmacy1","password":"pharmacy123456","account_name":"同仁堂药店","role":"药店","age":0,"gender":""}'
echo -e "\n药店1注册完成"

# 药店2
curl -X POST $URL -H "Content-Type: application/json" -d '{"username":"pharmacy2","password":"pharmacy123456","account_name":"百姓大药房","role":"药店","age":0,"gender":""}'
echo -e "\n药店2注册完成"

echo -e "\n所有用户注册完成！"
