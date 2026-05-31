# 基于区块链的社区医疗档案管理系统

> Hyperledger Fabric 医疗联盟链

## 系统架构

```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   Taobao    │  │     JD      │  │   Wenjin    │  │  RegCenter  │
│  (药房组织)  │  │ (互联网医疗) │  │ (温江社区)   │  │  (监管中心)  │
│ peer0/peer1 │  │ peer0/peer1 │  │ peer0/peer1 │  │ peer0/peer1 │
└──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘
       └────────────────┴───────┬────────┴───────────────┘
                               │
                     ┌─────────▼─────────┐
                     │  Orderer (QQ.com)  │
                     │    共识排序节点     │
                     └───────────────────┘
```

- **4 个联盟组织**: 淘宝（药房）、京东（互联网医疗）、温江（社区医疗中心）、RegCenter（监管中心）
- **1 个排序节点**: QQ.com
- **通道**: appchannel

## 技术栈

| 组件 | 技术 |
|---|---|
| 区块链网络 | Hyperledger Fabric 1.4.12 |
| 智能合约(链码) | Go |
| 后端服务 | Go + Gin 框架 |
| 前端 | Vue 2 + Element UI + ECharts |
| 数据库 | MySQL |
| SDK | Fabric-SDK-Go |

## 功能模块

- **用户管理**: 注册、登录、信息管理、软删除
- **病历管理**: 创建、查询、跨机构访问
- **授权管理**: 跨院授权、访问控制
- **药品订单**: 订单创建、状态流转
- **保险报销**: 保险覆盖信息管理
- **联盟活动监控**: 联盟内链上活动追踪
- **审计日志**: 病历访问日志、操作审计
- **统计报表**: ECharts 数据可视化

## 快速启动

### 1. 启动区块链网络

```bash
cd network
./start.sh
```

### 2. 安装链码

```bash
./install-all-peers.sh
```

### 3. 启动后端

```bash
cd application/server
go run main.go
```

### 4. 启动前端

```bash
cd application/web
npm install
npm run dev
```

## 目录结构

```
├── application/            # 应用层
│   ├── server/             # Go 后端服务
│   │   ├── api/            # API 接口
│   │   ├── blockchain/     # Fabric SDK 集成
│   │   ├── db/             # 数据库初始化与操作
│   │   ├── model/          # 数据模型
│   │   ├── pkg/            # 工具包
│   │   ├── routers/        # 路由定义
│   │   ├── service/        # 业务逻辑
│   │   ├── main.go         # 入口
│   │   └── config.yaml     # Fabric SDK 配置
│   └── web/                # Vue 前端
│       └── src/
│           ├── api/        # 前端 API 封装
│           ├── views/      # 页面组件
│           ├── router/     # 前端路由
│           └── store/      # 状态管理
├── chaincode/              # Fabric 链码（智能合约）
│   ├── api/                # 链码接口实现
│   ├── model/              # 链码数据模型
│   ├── pkg/utils/          # 工具函数
│   └── chaincode.go        # 链码入口
├── network/                # Fabric 网络配置
│   ├── configtx.yaml       # 通道配置
│   ├── crypto-config.yaml  # 证书配置
│   ├── docker-compose.yaml # 容器编排
│   ├── start.sh            # 网络启动脚本
│   └── explorer/           # 区块链浏览器
├── scripts/                # 部署脚本
└── tests/                  # 测试
    └── jmeter/             # JMeter 性能测试
```

## API 接口（v2）

| 接口 | 说明 |
|---|---|
| `/api/v2/register` | 用户注册 |
| `/api/v2/loginWithPassword` | 密码登录 |
| `/api/v2/getUserInfo` | 获取用户信息 |
| `/api/v2/createPrescription` | 创建病历 |
| `/api/v2/queryPrescription` | 查询病历列表 |
| `/api/v2/createDrugOrder` | 创建药品订单 |
| `/api/v2/queryDrugOrderList` | 查询药品订单 |
| `/api/v2/createInsuranceCover` | 创建保险报销信息 |
| `/api/v2/queryAccessTrace` | 查询病历访问轨迹 |

完整 API 列表见 `application/server/routers/router.go`。
