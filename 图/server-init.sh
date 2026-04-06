#!/bin/bash

echo "=========================================="
echo "服务器端项目初始化脚本"
echo "域名: http://www.lyt9999.cloud"
echo "=========================================="

# 检查必要的工具
echo -e "\n[检查] 验证必要工具..."
command -v docker >/dev/null 2>&1 || { echo "❌ Docker 未安装"; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "❌ Docker Compose 未安装"; exit 1; }
command -v node >/dev/null 2>&1 || { echo "❌ Node.js 未安装"; exit 1; }
command -v npm >/dev/null 2>&1 || { echo "❌ npm 未安装"; exit 1; }
echo "✅ 所有必要工具已安装"

# 1. 安装前端依赖
echo -e "\n[1/5] 安装前端依赖..."
cd application/web
if [ ! -d "node_modules" ]; then
    echo "正在安装 npm 包，这可能需要几分钟..."
    npm install --registry=https://registry.npmmirror.com
    if [ $? -ne 0 ]; then
        echo "❌ 前端依赖安装失败"
        exit 1
    fi
    echo "✅ 前端依赖安装完成"
else
    echo "✅ node_modules 已存在，跳过安装"
fi
cd ../..

# 2. 启动区块链网络
echo -e "\n[2/5] 启动 Hyperledger Fabric 区块链网络..."
cd network
chmod +x *.sh
./start.sh
if [ $? -ne 0 ]; then
    echo "❌ 区块链网络启动失败"
    exit 1
fi
echo "✅ 区块链网络启动成功"
cd ..

# 3. 构建后端 Docker 镜像
echo -e "\n[3/5] 构建后端服务..."
cd application/server
chmod +x build.sh
./build.sh
if [ $? -ne 0 ]; then
    echo "❌ 后端镜像构建失败"
    exit 1
fi
echo "✅ 后端镜像构建完成"
cd ..

# 4. 构建前端 Docker 镜像
echo -e "\n[4/5] 构建前端服务..."
cd web
if [ -f build.sh ]; then
    chmod +x build.sh
    ./build.sh
else
    docker build -t fabric-mims.web:latest .
fi
if [ $? -ne 0 ]; then
    echo "❌ 前端镜像构建失败"
    exit 1
fi
echo "✅ 前端镜像构建完成"
cd ..

# 5. 启动应用服务
echo -e "\n[5/5] 启动应用服务..."
docker-compose up -d
if [ $? -ne 0 ]; then
    echo "❌ 应用服务启动失败"
    exit 1
fi
echo "✅ 应用服务启动成功"
cd ..

# 等待服务启动
echo -e "\n等待服务完全启动..."
sleep 10

# 显示服务状态
echo -e "\n=========================================="
echo "✅ 部署完成！"
echo "=========================================="
echo ""
echo "📊 服务状态："
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
echo ""
echo "🌐 访问地址："
echo "   前端: http://www.lyt9999.cloud:8000"
echo "   后端: http://www.lyt9999.cloud:8888"
echo ""
echo "📝 常用命令："
echo "   查看后端日志: docker logs -f fabric-mims.server"
echo "   查看前端日志: docker logs -f fabric-mims.web"
echo "   重启服务: cd application && docker-compose restart"
echo "   停止服务: cd application && docker-compose stop"
echo "=========================================="
