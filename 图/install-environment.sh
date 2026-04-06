#!/bin/bash

echo "=========================================="
echo "服务器环境一键安装脚本"
echo "适用于: CentOS 7/8, Ubuntu 18.04+, Debian 10+"
echo "=========================================="

# 检测操作系统
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$ID
    VERSION=$VERSION_ID
else
    echo "❌ 无法检测操作系统类型"
    exit 1
fi

echo "检测到操作系统: $OS $VERSION"

# 1. 安装 Docker
echo -e "\n[1/5] 安装 Docker..."
if command -v docker &> /dev/null; then
    echo "✅ Docker 已安装: $(docker --version)"
else
    echo "正在安装 Docker..."
    curl -fsSL https://get.docker.com | bash -s docker --mirror Aliyun
    
    # 启动 Docker
    systemctl start docker
    systemctl enable docker
    
    if command -v docker &> /dev/null; then
        echo "✅ Docker 安装成功: $(docker --version)"
    else
        echo "❌ Docker 安装失败"
        exit 1
    fi
fi

# 2. 安装 Docker Compose
echo -e "\n[2/5] 安装 Docker Compose..."
if command -v docker-compose &> /dev/null; then
    echo "✅ Docker Compose 已安装: $(docker-compose --version)"
else
    echo "正在安装 Docker Compose..."
    
    # 下载 Docker Compose
    COMPOSE_VERSION="v2.24.0"
    curl -L "https://github.com/docker/compose/releases/download/${COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    
    # 如果 GitHub 下载失败，使用国内镜像
    if [ $? -ne 0 ]; then
        echo "GitHub 下载失败，尝试使用国内镜像..."
        curl -L "https://get.daocloud.io/docker/compose/releases/download/${COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    fi
    
    chmod +x /usr/local/bin/docker-compose
    
    if command -v docker-compose &> /dev/null; then
        echo "✅ Docker Compose 安装成功: $(docker-compose --version)"
    else
        echo "❌ Docker Compose 安装失败"
        exit 1
    fi
fi

# 3. 安装 Node.js 和 npm
echo -e "\n[3/5] 安装 Node.js 和 npm..."
if command -v node &> /dev/null; then
    echo "✅ Node.js 已安装: $(node --version)"
    echo "✅ npm 已安装: $(npm --version)"
else
    echo "正在安装 Node.js..."
    
    if [[ "$OS" == "centos" ]] || [[ "$OS" == "rhel" ]] || [[ "$OS" == "fedora" ]]; then
        # CentOS/RHEL/Fedora
        curl -fsSL https://rpm.nodesource.com/setup_18.x | bash -
        yum install -y nodejs
    elif [[ "$OS" == "ubuntu" ]] || [[ "$OS" == "debian" ]]; then
        # Ubuntu/Debian
        curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
        apt-get install -y nodejs
    else
        echo "❌ 不支持的操作系统: $OS"
        exit 1
    fi
    
    if command -v node &> /dev/null; then
        echo "✅ Node.js 安装成功: $(node --version)"
        echo "✅ npm 安装成功: $(npm --version)"
    else
        echo "❌ Node.js 安装失败"
        exit 1
    fi
fi

# 4. 配置 npm 镜像源
echo -e "\n[4/5] 配置 npm 镜像源..."
npm config set registry https://registry.npmmirror.com
echo "✅ npm 镜像源已设置为: $(npm config get registry)"

# 5. 安装 Git（如果没有）
echo -e "\n[5/5] 检查 Git..."
if command -v git &> /dev/null; then
    echo "✅ Git 已安装: $(git --version)"
else
    echo "正在安装 Git..."
    if [[ "$OS" == "centos" ]] || [[ "$OS" == "rhel" ]] || [[ "$OS" == "fedora" ]]; then
        yum install -y git
    elif [[ "$OS" == "ubuntu" ]] || [[ "$OS" == "debian" ]]; then
        apt-get update
        apt-get install -y git
    fi
    echo "✅ Git 安装成功: $(git --version)"
fi

# 显示安装结果
echo -e "\n=========================================="
echo "✅ 环境安装完成！"
echo "=========================================="
echo ""
echo "📦 已安装的工具："
echo "   Docker: $(docker --version)"
echo "   Docker Compose: $(docker-compose --version)"
echo "   Node.js: $(node --version)"
echo "   npm: $(npm --version)"
echo "   Git: $(git --version)"
echo ""
echo "🚀 下一步操作："
echo "   1. 上传项目到服务器"
echo "   2. 运行: cd /root/fabric-mims-main && chmod +x server-init.sh && ./server-init.sh"
echo "=========================================="
