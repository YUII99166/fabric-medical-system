@echo off
REM 隐私溯源功能一键部署脚本 (Windows)
REM 用途：部署访问日志记录功能到区块链

echo ==========================================
echo   隐私溯源功能部署脚本
echo ==========================================
echo.

REM 获取脚本所在目录的父目录（项目根目录）
set SCRIPT_DIR=%~dp0
set PROJECT_ROOT=%SCRIPT_DIR%..

echo 项目根目录: %PROJECT_ROOT%
echo.

REM 步骤1: 打包 chaincode
echo 步骤 1/3: 打包 chaincode...
cd /d "%PROJECT_ROOT%\chaincode"
echo    当前目录: %CD%
echo    执行: set GO111MODULE=on ^& go mod vendor
set GO111MODULE=on
go mod vendor
if errorlevel 1 (
    echo    ❌ Chaincode 打包失败
    pause
    exit /b 1
)
echo    ✅ Chaincode 打包完成
echo.

REM 步骤2: 部署 chaincode
echo 步骤 2/3: 部署 chaincode 到区块链网络...
cd /d "%PROJECT_ROOT%\network"
echo    当前目录: %CD%
echo    执行: network.sh deployCC -ccn medical -ccp ../chaincode -ccl go
call network.sh deployCC -ccn medical -ccp ../chaincode -ccl go
if errorlevel 1 (
    echo    ❌ Chaincode 部署失败
    pause
    exit /b 1
)
echo    ✅ Chaincode 部署完成
echo.

REM 步骤3: 提示重启后端
echo 步骤 3/3: 需要重启后端服务器
echo    请在新终端执行以下命令：
echo    cd %PROJECT_ROOT%\application\server
echo    go run main.go
echo.

echo ==========================================
echo   ✅ 部署完成！
echo ==========================================
echo.
echo 测试步骤：
echo 1. 重启后端服务器（见上方命令）
echo 2. 用管理员账号访问病人病历
echo 3. 用病人账号查看'隐私溯源'页面
echo.

pause
