<template>
  <div class="architecture-container">
    <!-- 页面标题 -->
    <el-card class="page-header">
      <h2>
        <i class="el-icon-connection"></i>
        联盟链架构概览
      </h2>
      <p>社区医疗管理联盟链系统架构与组织关系</p>
    </el-card>

    <!-- 系统概述 -->
    <el-row :gutter="20">
      <el-col :xs="24" :md="12">
        <el-card class="info-card">
          <div slot="header" class="card-header">
            <i class="el-icon-info"></i>
            <span>系统概述</span>
          </div>
          <div class="info-content">
            <div class="info-item">
              <span class="label">系统名称：</span>
              <span class="value">社区医疗管理联盟链</span>
            </div>
            <div class="info-item">
              <span class="label">区块链平台：</span>
              <span class="value">Hyperledger Fabric 1.4.12</span>
            </div>
            <div class="info-item">
              <span class="label">链码版本：</span>
              <span class="value">v1.0.7</span>
            </div>
            <div class="info-item">
              <span class="label">通道名称：</span>
              <span class="value">appchannel</span>
            </div>
            <div class="info-item">
              <span class="label">组织数量：</span>
              <span class="value">4 个</span>
            </div>
            <div class="info-item">
              <span class="label">节点总数：</span>
              <span class="value">9 个（8 peer + 1 orderer）</span>
            </div>
            <div class="info-item">
              <span class="label">背书策略：</span>
              <span class="value">OutOf(3,4)（至少3个组织背书）</span>
            </div>
            <div class="info-item">
              <span class="label">认证组织：</span>
              <span class="value">8 个（5 个已认证，3 个审核中）</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="12">
        <el-card class="info-card">
          <div slot="header" class="card-header">
            <i class="el-icon-data-line"></i>
            <span>系统统计</span>
          </div>
          <div class="stats-grid">
            <div class="stat-item" v-for="stat in systemStats" :key="stat.label">
              <div class="stat-icon" :style="{ background: stat.color }">
                <i :class="stat.icon"></i>
              </div>
              <div class="stat-info">
                <p class="stat-value">{{ stat.value }}</p>
                <p class="stat-label">{{ stat.label }}</p>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 联盟链架构图 -->
    <el-card class="architecture-card">
      <div slot="header" class="card-header">
        <i class="el-icon-share"></i>
        <span>联盟链架构图</span>
      </div>
      
      <div class="architecture-diagram">
        <!-- 应用层 -->
        <div class="layer app-layer">
          <h3>应用层</h3>
          <div class="app-roles">
            <div class="role-box" v-for="role in roles" :key="role.name">
              <i :class="role.icon"></i>
              <p>{{ role.name }}</p>
              <span>{{ role.count }}</span>
            </div>
          </div>
        </div>

        <!-- 箭头：应用层 -> 背书节点 -->
        <div class="layer-arrow">
          <i class="el-icon-bottom"></i>
          <span>1. 提交交易提案</span>
        </div>

        <!-- 背书节点层 -->
        <div class="layer endorser-layer">
          <h3>背书节点层（Endorsing Peers）</h3>
          <div class="endorser-info">
            <el-tag type="success" size="small">背书策略: OutOf(3,4) (至少3个组织背书)</el-tag>
          </div>
          <div class="organizations">
            <div class="organization" v-for="org in organizations" :key="org.msp">
              <div class="org-header" :style="{ background: org.color }">
                <i :class="org.icon"></i>
                <h4>{{ org.name }}</h4>
                <span class="msp-id">{{ org.msp }}</span>
              </div>
              <div class="org-peers">
                <div class="peer-node endorser" v-for="peer in org.peers" :key="peer.name">
                  <i class="el-icon-circle-check"></i>
                  <p>{{ peer.name }}</p>
                  <span class="port">{{ peer.port }}</span>
                  <span class="badge">背书节点</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 箭头：背书节点 -> 排序服务 -->
        <div class="layer-arrow">
          <i class="el-icon-bottom"></i>
          <span>2. 返回背书结果 → 3. 提交已背书交易</span>
        </div>

        <!-- 排序服务层 -->
        <div class="layer orderer-layer">
          <h3>排序服务层（Ordering Service）</h3>
          <div class="node-container">
            <div class="node orderer-node">
              <i class="el-icon-s-grid"></i>
              <p>Orderer</p>
              <span>orderer.qq.com</span>
              <span class="port">7050</span>
              <span class="badge">排序节点</span>
            </div>
          </div>
          <div class="orderer-desc">
            <el-tag type="warning" size="small">负责交易排序和区块生成</el-tag>
          </div>
        </div>

        <!-- 箭头：排序服务 -> 通道 -->
        <div class="layer-arrow">
          <i class="el-icon-bottom"></i>
          <span>4. 生成区块并分发</span>
        </div>

        <!-- 通道层 -->
        <div class="layer channel-layer">
          <div class="channel-box">
            <i class="el-icon-connection"></i>
            <span>Channel: appchannel</span>
          </div>
        </div>

        <!-- 箭头：通道 -> 提交节点 -->
        <div class="layer-arrow">
          <i class="el-icon-bottom"></i>
          <span>5. 区块广播到所有节点</span>
        </div>

        <!-- 提交节点层（Committing Peers）-->
        <div class="layer committer-layer">
          <h3>提交节点层（Committing Peers）</h3>
          <div class="committer-info">
            <el-tag type="info" size="small">所有 Peer 节点都是提交节点</el-tag>
          </div>
          <div class="organizations">
            <div class="organization" v-for="org in organizations" :key="'commit-' + org.msp">
              <div class="org-header" :style="{ background: org.color }">
                <i :class="org.icon"></i>
                <h4>{{ org.name }}</h4>
                <span class="msp-id">{{ org.msp }}</span>
              </div>
              <div class="org-peers">
                <div class="peer-node committer" v-for="peer in org.peers" :key="'commit-' + peer.name">
                  <i class="el-icon-document-checked"></i>
                  <p>{{ peer.name }}</p>
                  <span class="port">{{ peer.port }}</span>
                  <span class="badge">提交节点</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 箭头：提交节点 -> 链码 -->
        <div class="layer-arrow">
          <i class="el-icon-bottom"></i>
          <span>6. 验证并提交到账本</span>
        </div>

        <!-- 链码和账本层 -->
        <div class="layer chaincode-layer">
          <h3>链码与账本层</h3>
          <div class="chaincode-ledger-container">
            <div class="chaincode-box">
              <i class="el-icon-document-copy"></i>
              <div class="chaincode-info">
                <p class="chaincode-name">Chaincode</p>
                <span class="chaincode-version">fabric-mims v1.0.7</span>
                <span class="chaincode-lang">Golang</span>
              </div>
            </div>
            <div class="ledger-box">
              <i class="el-icon-coin"></i>
              <div class="ledger-info">
                <p class="ledger-name">Ledger</p>
                <span>世界状态 (World State)</span>
                <span>区块链 (Blockchain)</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 数据流向图 -->
    <el-card class="dataflow-card">
      <div slot="header" class="card-header">
        <i class="el-icon-sort"></i>
        <span>数据流向示意</span>
      </div>
      
      <div class="dataflow-diagram">
        <div class="flow-step">
          <div class="step-box user-step">
            <i class="el-icon-user"></i>
            <p>用户操作</p>
            <span>医生/病人/药店</span>
          </div>
          <i class="el-icon-right flow-arrow"></i>
        </div>

        <div class="flow-step">
          <div class="step-box frontend-step">
            <i class="el-icon-monitor"></i>
            <p>前端应用</p>
            <span>Vue.js + Element UI</span>
          </div>
          <i class="el-icon-right flow-arrow"></i>
        </div>

        <div class="flow-step">
          <div class="step-box backend-step">
            <i class="el-icon-cpu"></i>
            <p>后端服务</p>
            <span>Go + Gin + MySQL</span>
          </div>
          <i class="el-icon-right flow-arrow"></i>
        </div>

        <div class="flow-step">
          <div class="step-box sdk-step">
            <i class="el-icon-connection"></i>
            <p>Fabric SDK</p>
            <span>Go SDK</span>
          </div>
          <i class="el-icon-right flow-arrow"></i>
        </div>

        <div class="flow-step">
          <div class="step-box blockchain-step">
            <i class="el-icon-coin"></i>
            <p>区块链网络</p>
            <span>Hyperledger Fabric</span>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 组织详情 -->
    <el-card class="organizations-card">
      <div slot="header" class="card-header">
        <i class="el-icon-office-building"></i>
        <span>组织详细信息</span>
      </div>
      
      <el-table :data="organizationsDetail" style="width: 100%">
        <el-table-column prop="name" label="组织名称" width="180"></el-table-column>
        <el-table-column prop="msp" label="MSP ID" width="150"></el-table-column>
        <el-table-column prop="domain" label="域名" width="180"></el-table-column>
        <el-table-column prop="peerCount" label="节点数量" width="100"></el-table-column>
        <el-table-column prop="ports" label="端口" min-width="150"></el-table-column>
        <el-table-column prop="role" label="角色" width="120"></el-table-column>
        <el-table-column label="证书状态" width="120" align="center">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.status === 'pending'" type="warning" size="small">
              <i class="el-icon-loading"></i>
              审核中
            </el-tag>
            <el-tag v-else-if="scope.row.certified" type="success" size="small">
              <i class="el-icon-circle-check"></i>
              已认证
            </el-tag>
            <el-tag v-else type="danger" size="small">
              <i class="el-icon-circle-close"></i>
              未认证
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script>
import { queryAccountList } from '@/api/accountV2'

export default {
  name: 'Architecture',
  data() {
    return {
      userCount: 0,
      organizations: [
        {
          name: '协和医院',
          msp: 'TaobaoMSP',
          icon: 'el-icon-office-building',
          color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
          peers: [
            { name: 'peer0.taobao.com', port: '7051' },
            { name: 'peer1.taobao.com', port: '17051' }
          ]
        },
        {
          name: '301医院',
          msp: 'JDMSP',
          icon: 'el-icon-office-building',
          color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
          peers: [
            { name: 'peer0.jd.com', port: '27051' },
            { name: 'peer1.jd.com', port: '37051' }
          ]
        },
        {
          name: '温江医疗中心',
          msp: 'WenjinMSP',
          icon: 'el-icon-office-building',
          color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
          peers: [
            { name: 'peer0.wenjin.com', port: '47051' },
            { name: 'peer1.wenjin.com', port: '57051' }
          ]
        },
        {
          name: '监管中心',
          msp: 'RegCenterMSP',
          icon: 'el-icon-s-data',
          color: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
          peers: [
            { name: 'peer0.regcenter.com', port: '8051' },
            { name: 'peer1.regcenter.com', port: '9051' }
          ]
        }
      ],
      organizationsDetail: [
        { name: '协和医院', msp: 'TaobaoMSP', domain: 'taobao.com', peerCount: 2, ports: '7051, 17051', role: '医疗机构', certified: true },
        { name: '301医院', msp: 'JDMSP', domain: 'jd.com', peerCount: 2, ports: '27051, 37051', role: '医疗机构', certified: true },
        { name: '温江医疗中心', msp: 'WenjinMSP', domain: 'wenjin.com', peerCount: 2, ports: '47051, 57051', role: '医疗机构', certified: true },
        { name: '监管中心', msp: 'RegCenterMSP', domain: 'regcenter.com', peerCount: 2, ports: '8051, 9051', role: '监管机构', certified: true },
        { name: '排序服务', msp: 'OrdererMSP', domain: 'qq.com', peerCount: 1, ports: '7050', role: '排序节点', certified: true },
        { name: '温江人寿', msp: 'InsuranceMSP', domain: 'insurance.com', peerCount: 2, ports: '待分配', role: '保险机构', certified: false, status: 'pending' },
        { name: '成都中医药大学第二附属医院', msp: 'CDUTCMMSP', domain: 'cdutcm.edu.cn', peerCount: 2, ports: '待分配', role: '医疗机构', certified: false, status: 'pending' },
        { name: '四川省药品监督管理局', msp: 'DrugRegMSP', domain: 'scda.gov.cn', peerCount: 2, ports: '待分配', role: '药品监管机构', certified: false, status: 'pending' }
      ],
      roles: [
        { name: '医生', icon: 'el-icon-user', count: 0 },
        { name: '病人', icon: 'el-icon-user', count: 0 },
        { name: '药店', icon: 'el-icon-shopping-bag-2', count: 0 },
        { name: '监管中心', icon: 'el-icon-s-custom', count: 0 }
      ]
    }
  },
  computed: {
    systemStats() {
      return [
        { label: '总用户数', value: this.userCount, icon: 'el-icon-user-solid', color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' },
        { label: '已认证组织', value: '5', icon: 'el-icon-circle-check', color: 'linear-gradient(135deg, #67C23A 0%, #85ce61 100%)' },
        { label: '待审核组织', value: '3', icon: 'el-icon-loading', color: 'linear-gradient(135deg, #ffa726 0%, #fb8c00 100%)' },
        { label: 'Peer节点', value: '8', icon: 'el-icon-monitor', color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)' },
        { label: '排序节点', value: '1', icon: 'el-icon-s-grid', color: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)' },
        { label: '通道数量', value: '1', icon: 'el-icon-connection', color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)' }
      ]
    }
  },
  mounted() {
    this.loadUserData()
  },
  methods: {
    async loadUserData() {
      try {
        const response = await queryAccountList()
        const userList = response.code === 200 ? response.data : response
        this.userCount = Array.isArray(userList) ? userList.length : 0
        
        if (Array.isArray(userList)) {
          this.roles[0].count = userList.filter(u => u.role === '医生').length
          this.roles[1].count = userList.filter(u => u.role === '病人').length
          this.roles[2].count = userList.filter(u => u.role === '药店').length
          this.roles[3].count = userList.filter(u => u.organization === 'RegCenterMSP').length
        }
      } catch (error) {
        console.error('加载用户数据失败:', error)
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.architecture-container {
  padding: 20px;
  background: #f5f7fa;
  min-height: calc(100vh - 60px);
}

.page-header {
  margin-bottom: 20px;
  border-radius: 12px;
  
  h2 {
    margin: 0 0 10px 0;
    font-size: 24px;
    color: #303133;
    
    i {
      color: #4facfe;
      margin-right: 10px;
    }
  }
  
  p {
    margin: 0;
    color: #909399;
    font-size: 14px;
  }
}

.info-card {
  border-radius: 12px;
  margin-bottom: 20px;
  
  .card-header {
    font-weight: 600;
    
    i {
      margin-right: 8px;
      color: #4facfe;
    }
  }
  
  .info-content {
    .info-item {
      display: flex;
      justify-content: space-between;
      padding: 12px 0;
      border-bottom: 1px solid #f0f0f0;
      
      &:last-child {
        border-bottom: none;
      }
      
      .label {
        color: #909399;
        font-size: 14px;
      }
      
      .value {
        color: #303133;
        font-size: 14px;
        font-weight: 500;
      }
    }
  }
  
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
    
    .stat-item {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 12px;
      background: #f9fafc;
      border-radius: 8px;
      transition: all 0.3s;
      
      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
      }
      
      .stat-icon {
        width: 45px;
        height: 45px;
        border-radius: 10px;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
        
        i {
          font-size: 22px;
          color: white;
        }
      }
      
      .stat-info {
        flex: 1;
        min-width: 0;
        
        .stat-value {
          margin: 0 0 3px 0;
          font-size: 20px;
          font-weight: 600;
          color: #303133;
        }
        
        .stat-label {
          margin: 0;
          font-size: 12px;
          color: #909399;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }
      }
    }
  }
}

.architecture-card, .dataflow-card, .organizations-card {
  border-radius: 12px;
  margin-bottom: 20px;
  
  .card-header {
    font-weight: 600;
    
    i {
      margin-right: 8px;
      color: #4facfe;
    }
  }
}

.architecture-diagram {
  .layer {
    margin-bottom: 20px;
    
    h3 {
      text-align: center;
      font-size: 16px;
      color: #606266;
      margin-bottom: 15px;
      font-weight: 600;
    }
  }

  .layer-arrow {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin: 15px 0;
    
    i {
      font-size: 28px;
      color: #4facfe;
      animation: bounce 2s infinite;
    }
    
    span {
      margin-top: 5px;
      font-size: 13px;
      color: #909399;
      font-weight: 500;
    }
  }

  @keyframes bounce {
    0%, 100% {
      transform: translateY(0);
    }
    50% {
      transform: translateY(5px);
    }
  }

  .app-layer {
    .app-roles {
      display: flex;
      justify-content: center;
      gap: 20px;
      flex-wrap: wrap;
      
      .role-box {
        background: white;
        border: 2px solid #4facfe;
        padding: 20px 30px;
        border-radius: 12px;
        text-align: center;
        min-width: 120px;
        transition: all 0.3s;
        
        &:hover {
          transform: translateY(-5px);
          box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
        }
        
        i {
          font-size: 32px;
          color: #4facfe;
          margin-bottom: 10px;
        }
        
        p {
          margin: 5px 0;
          font-size: 15px;
          color: #303133;
          font-weight: 600;
        }
        
        span {
          font-size: 20px;
          color: #4facfe;
          font-weight: 600;
        }
      }
    }
  }

  .endorser-layer, .committer-layer {
    background: #f9fafc;
    padding: 20px;
    border-radius: 12px;
    border: 2px solid #e4e7ed;

    .endorser-info, .committer-info {
      text-align: center;
      margin-bottom: 15px;
    }

    .organizations {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
      gap: 20px;
      
      .organization {
        border-radius: 12px;
        overflow: hidden;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        
        .org-header {
          padding: 15px;
          color: white;
          text-align: center;
          
          i {
            font-size: 28px;
            margin-bottom: 8px;
          }
          
          h4 {
            margin: 5px 0;
            font-size: 16px;
          }
          
          .msp-id {
            font-size: 12px;
            opacity: 0.9;
          }
        }
        
        .org-peers {
          background: white;
          padding: 15px;
          
          .peer-node {
            background: #f9fafc;
            padding: 12px;
            border-radius: 8px;
            margin-bottom: 10px;
            text-align: center;
            position: relative;
            
            &:last-child {
              margin-bottom: 0;
            }

            &.endorser {
              border: 2px solid #67c23a;
              background: #f0f9ff;
              
              i {
                color: #67c23a;
              }
            }

            &.committer {
              border: 2px solid #409eff;
              background: #ecf5ff;
              
              i {
                color: #409eff;
              }
            }
            
            i {
              font-size: 24px;
              margin-bottom: 5px;
            }
            
            p {
              margin: 5px 0;
              font-size: 13px;
              color: #303133;
              font-weight: 500;
            }
            
            .port {
              display: block;
              font-size: 12px;
              color: #909399;
              margin-top: 3px;
            }

            .badge {
              display: inline-block;
              margin-top: 5px;
              padding: 2px 8px;
              font-size: 11px;
              border-radius: 10px;
              background: rgba(0, 0, 0, 0.1);
              color: #303133;
            }
          }
        }
      }
    }
  }
  
  .orderer-layer {
    .node-container {
      display: flex;
      justify-content: center;
      
      .orderer-node {
        background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
        color: white;
        padding: 20px 30px;
        border-radius: 12px;
        text-align: center;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        
        i {
          font-size: 32px;
          margin-bottom: 10px;
        }
        
        p {
          margin: 5px 0;
          font-size: 16px;
          font-weight: 600;
        }
        
        span {
          display: block;
          font-size: 13px;
          opacity: 0.9;
          
          &.port {
            margin-top: 5px;
            font-size: 12px;
          }

          &.badge {
            display: inline-block;
            margin-top: 8px;
            padding: 3px 10px;
            font-size: 11px;
            border-radius: 10px;
            background: rgba(255, 255, 255, 0.3);
          }
        }
      }
    }

    .orderer-desc {
      text-align: center;
      margin-top: 10px;
    }
  }
  
  .channel-layer {
    display: flex;
    justify-content: center;
    margin: 20px 0;
    
    .channel-box {
      background: #ecf5ff;
      border: 2px dashed #4facfe;
      padding: 15px 40px;
      border-radius: 8px;
      display: flex;
      align-items: center;
      gap: 10px;
      
      i {
        font-size: 20px;
        color: #4facfe;
      }
      
      span {
        font-size: 15px;
        font-weight: 600;
        color: #4facfe;
      }
    }
  }
  
  .chaincode-layer {
    .chaincode-ledger-container {
      display: flex;
      justify-content: center;
      gap: 30px;
      flex-wrap: wrap;

      .chaincode-box, .ledger-box {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
        padding: 20px 40px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        gap: 15px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        min-width: 250px;
        
        i {
          font-size: 32px;
        }
        
        .chaincode-info, .ledger-info {
          .chaincode-name, .ledger-name {
            margin: 0 0 5px 0;
            font-size: 18px;
            font-weight: 600;
          }
          
          span {
            display: block;
            margin-top: 3px;
            font-size: 13px;
            opacity: 0.9;
          }
        }
      }

      .ledger-box {
        background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
      }
    }
  }
}

.dataflow-diagram {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 15px;
  padding: 20px;
  
  .flow-step {
    display: flex;
    align-items: center;
    gap: 15px;
    
    .step-box {
      padding: 20px 30px;
      border-radius: 12px;
      text-align: center;
      min-width: 150px;
      color: white;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
      
      i {
        font-size: 32px;
        margin-bottom: 10px;
      }
      
      p {
        margin: 5px 0;
        font-size: 15px;
        font-weight: 600;
      }
      
      span {
        font-size: 12px;
        opacity: 0.9;
      }
      
      &.user-step {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      }
      
      &.frontend-step {
        background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
      }
      
      &.backend-step {
        background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
      }
      
      &.sdk-step {
        background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
      }
      
      &.blockchain-step {
        background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
      }
    }
    
    .flow-arrow {
      font-size: 24px;
      color: #909399;
    }
  }
}

@media (max-width: 768px) {
  .dataflow-diagram {
    flex-direction: column;
    
    .flow-step {
      flex-direction: column;
      
      .flow-arrow {
        transform: rotate(90deg);
      }
    }
  }
  
  .endorser-layer .organizations,
  .committer-layer .organizations {
    grid-template-columns: 1fr;
  }

  .chaincode-layer .chaincode-ledger-container {
    flex-direction: column;
    align-items: center;
  }
}
</style>
