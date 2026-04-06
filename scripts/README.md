# Fabric MIMS 脚本工具

## 测试脚本
- `test-login.sh` - 测试登录功能
- `test-register.sh` - 测试注册功能
- `test-regcenter.sh` / `test-regcenter.bat` - 测试监管中心功能

## 部署脚本
- `本地测试-快速启动.bat` - Windows 本地快速启动
- `本地测试-停止服务.bat` - Windows 本地停止服务
- `服务器快速修复部署.sh` - 服务器快速修复部署
- `快速修复.sh` - 快速修复脚本

## 数据管理
- `清理所有数据.sh` / `清理所有数据.bat` - 清理所有区块链和数据库数据
- `sync-user-to-blockchain.sh` / `sync-user-to-blockchain.bat` - 同步用户到区块链

## 链码管理
- `check-chaincode-status.sh` - 检查链码状态
- `upgrade-chaincode-pharmacy.sh` / `upgrade-chaincode-pharmacy.bat` - 升级药店链码
- `install-chaincode-directly.sh` - 直接安装链码
- `reinstall-chaincode-all-peers.sh` - 在所有节点重新安装链码
- `fix-chaincode-upgrade.sh` - 修复链码升级问题
- `restart-network-with-new-chaincode.sh` - 使用新链码重启网络

## 容器管理
- `start-all-chaincode-containers.sh` - 启动所有链码容器
- `start-endorsing-peers-chaincode.sh` - 启动背书节点链码
- `force-restart-peers.sh` - 强制重启节点

## 使用说明

### 本地开发测试
```bash
# Windows
本地测试-快速启动.bat

# Linux/Mac
./快速修复.sh
```

### 服务器部署
```bash
# 首次部署
参考 ../docs/deployment/服务器部署指南.md

# 快速修复
./服务器快速修复部署.sh
```

### 链码升级
```bash
# 检查当前状态
./check-chaincode-status.sh

# 升级链码
./upgrade-chaincode-pharmacy.sh
```
