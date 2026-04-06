@echo off
chcp 65001 >nul
echo ==========================================
echo 清理项目，准备上传到服务器
echo ==========================================

echo.
echo [1/8] 清理前端 node_modules...
if exist "application\web\node_modules" (
    rmdir /s /q "application\web\node_modules"
    echo ✅ 已删除 application\web\node_modules
) else (
    echo ℹ️  node_modules 不存在，跳过
)

echo.
echo [2/8] 清理前端构建产物...
if exist "application\web\dist" rmdir /s /q "application\web\dist"
if exist "application\web\.temp" rmdir /s /q "application\web\.temp"
echo ✅ 已删除前端构建产物

echo.
echo [3/8] 清理后端 vendor...
if exist "application\server\vendor" (
    rmdir /s /q "application\server\vendor"
    echo ✅ 已删除 Go vendor
) else (
    echo ℹ️  vendor 不存在，跳过
)

echo.
echo [4/8] 清理编译产物...
if exist "application\server\server.exe" del /f /q "application\server\server.exe"
if exist "application\server\fabric-server" del /f /q "application\server\fabric-server"
echo ✅ 已删除编译产物

echo.
echo [5/8] 清理 Mac 版本的 Fabric 工具...
if exist "network\hyperledger-fabric-darwin-amd64-1.4.12" (
    rmdir /s /q "network\hyperledger-fabric-darwin-amd64-1.4.12"
    echo ✅ 已删除 Mac 版本（服务器不需要）
) else (
    echo ℹ️  Mac 工具不存在，跳过
)

echo.
echo [6/8] 清理 Fabric 证书和配置（服务器会重新生成）...
if exist "network\crypto-config" rmdir /s /q "network\crypto-config"
if exist "network\config" rmdir /s /q "network\config"
if exist "network\channel-artifacts" rmdir /s /q "network\channel-artifacts"
echo ✅ 已删除 Fabric 证书和配置

echo.
echo [7/8] 清理日志文件...
for /r %%i in (*.log) do del /f /q "%%i"
if exist "logs" rmdir /s /q "logs"
echo ✅ 已删除日志文件

echo.
echo [8/8] 清理临时文件...
for /r %%i in (*.tmp) do del /f /q "%%i"
for /r %%i in (*.bak) do del /f /q "%%i"
for /r %%i in (.DS_Store) do del /f /q "%%i"
echo ✅ 已删除临时文件

echo.
echo ==========================================
echo ✅ 清理完成！
echo ==========================================
echo.
echo 📦 项目体积已大幅减小（预计减少 ~200 MB）
echo.
echo 📋 上传到服务器后的操作步骤：
echo    1. 上传整个项目到服务器
echo    2. 给脚本添加执行权限: chmod +x server-init.sh
echo    3. 运行初始化脚本: ./server-init.sh
echo.
echo 💡 提示：
echo    - 前端依赖会在服务器上自动安装
echo    - Fabric 证书会在启动时自动生成
echo    - 后端会在 Docker 构建时自动编译
echo ==========================================
pause
