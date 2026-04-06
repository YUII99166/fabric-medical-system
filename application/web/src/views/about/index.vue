<template>
  <div class="about-container">
    <el-card class="about-card">
      <div slot="header" class="card-header">
        <i class="el-icon-info"></i>
        <span>关于系统</span>
      </div>

      <!-- 项目简介 -->
      <div class="section">
        <h2>基于区块链的社区医疗管理系统</h2>
        <p class="intro">这是一个基于 Hyperledger Fabric 区块链的社区医疗管理系统，支持多角色用户管理、病历管理、跨院诊疗、授权管理和药品订单管理。</p>
      </div>

      <!-- 技术栈 -->
      <div class="section">
        <h3><i class="el-icon-cpu"></i> 技术栈</h3>
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12">
            <div class="tech-box">
              <h4>后端技术</h4>
              <ul>
                <li>Go 1.16+</li>
                <li>Gin Web Framework</li>
                <li>MySQL 数据库</li>
                <li>Hyperledger Fabric 1.4.12</li>
              </ul>
            </div>
          </el-col>
          <el-col :xs="24" :sm="12">
            <div class="tech-box">
              <h4>前端技术</h4>
              <ul>
                <li>Vue.js 2.x</li>
                <li>Element UI</li>
                <li>Vuex</li>
                <li>Vue Router</li>
              </ul>
            </div>
          </el-col>
        </el-row>
      </div>

      <!-- 系统架构 -->
      <div class="section">
        <h3><i class="el-icon-connection"></i> 系统架构</h3>
        <div class="architecture">
          <el-tag type="success" size="large">前端 (Vue.js)</el-tag>
          <i class="el-icon-right"></i>
          <el-tag type="primary" size="large">后端 (Go)</el-tag>
          <i class="el-icon-right"></i>
          <el-tag type="warning" size="large">MySQL 数据库</el-tag>
          <div class="blockchain-layer">
            <i class="el-icon-bottom"></i>
            <el-tag type="info" size="large">Hyperledger Fabric 区块链网络（可选）</el-tag>
          </div>
        </div>
      </div>

      <!-- 用户角色 -->
      <div class="section">
        <h3><i class="el-icon-user-solid"></i> 用户角色</h3>
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12" :md="8" :lg="4" v-for="role in roles" :key="role.name">
            <div class="role-card">
              <div class="role-icon">
                <i :class="role.icon"></i>
              </div>
              <h4>{{ role.name }}</h4>
              <p>{{ role.desc }}</p>
            </div>
          </el-col>
        </el-row>
      </div>

      <!-- 主要功能 -->
      <div class="section">
        <h3><i class="el-icon-menu"></i> 主要功能</h3>
        <el-collapse accordion>
          <el-collapse-item v-for="feature in features" :key="feature.title" :title="feature.title" :name="feature.title">
            <ul>
              <li v-for="item in feature.items" :key="item">{{ item }}</li>
            </ul>
          </el-collapse-item>
        </el-collapse>
      </div>

      <!-- 许可证 -->
      <div class="section footer-section">
        <el-divider></el-divider>
        <p class="license">
          <i class="el-icon-document"></i>
          许可证：MIT License
        </p>
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'About',
  computed: {
    primaryColor() {
      return getComputedStyle(document.documentElement).getPropertyValue('--theme-primary').trim() || '#409EFF'
    }
  },
  data() {
    return {
      roles: [
        { name: '管理员', desc: '系统管理，所有功能权限', icon: 'el-icon-s-custom' },
        { name: '医生', desc: '病历管理（查看、新增）、跨院诊疗', icon: 'el-icon-user' },
        { name: '病人', desc: '查看病历、授权管理', icon: 'el-icon-user' },
        { name: '药店', desc: '药品订单管理', icon: 'el-icon-shopping-bag-2' }
      ],
      features: [
        {
          title: '账户管理',
          items: ['用户注册', '用户登录', '角色权限控制']
        },
        {
          title: '病历管理',
          items: ['医生创建病历', '病人查看个人病历', '管理员查看所有病历']
        },
        {
          title: '跨院诊疗',
          items: ['医生申请访问授权', '病人批准授权', '补充诊疗记录']
        },
        {
          title: '药品订单',
          items: ['药店创建订单', '病人查看个人订单', '订单状态管理']
        }
      ]
    }
  }
}
</script>

<style lang="scss" scoped>
.about-container {
  padding: 20px;
  background: #f0f2f5;
  min-height: calc(100vh - 60px);
}

.about-card {
  max-width: 1200px;
  margin: 0 auto;
  border-radius: 12px;

  .card-header {
    font-size: 20px;
    font-weight: 600;
    color: #303133;

    i {
      margin-right: 8px;
      color: var(--theme-primary, #409EFF);
    }
  }
}

.section {
  margin-bottom: 30px;

  h2 {
    font-size: 24px;
    color: #303133;
    margin-bottom: 10px;
  }

  h3 {
    font-size: 18px;
    color: #303133;
    margin-bottom: 15px;
    padding-bottom: 10px;
    border-bottom: 2px solid #EBEEF5;

    i {
      margin-right: 8px;
      color: var(--theme-primary, #409EFF);
    }
  }

  h4 {
    font-size: 16px;
    color: #606266;
    margin-bottom: 10px;
  }

  .intro {
    font-size: 16px;
    line-height: 1.8;
    color: #606266;
  }
}

.tech-box {
  background: var(--theme-light-bg, #ecf5ff);
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 15px;
  border: 1px solid var(--theme-primary, #409EFF);
  border-left: 4px solid var(--theme-primary, #409EFF);

  h4 {
    margin-top: 0;
    color: var(--theme-primary, #409EFF);
  }

  ul {
    list-style: none;
    padding: 0;
    margin: 0;

    li {
      padding: 8px 0;
      color: #606266;
      border-bottom: 1px solid #e4e7ed;

      &:last-child {
        border-bottom: none;
      }

      &:before {
        content: "✓";
        color: var(--theme-primary, #409EFF);
        font-weight: bold;
        margin-right: 10px;
      }
    }
  }
}

.architecture {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 15px;
  padding: 30px;
  background: var(--theme-light-bg, #ecf5ff);
  border-radius: 8px;
  border: 1px solid var(--theme-primary, #409EFF);

  .el-tag {
    padding: 12px 20px;
    font-size: 14px;
  }

  i {
    font-size: 20px;
    color: var(--theme-primary, #409EFF);
  }

  .blockchain-layer {
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
    margin-top: 10px;
  }
}

.role-card {
  text-align: center;
  padding: 20px;
  border: 2px solid var(--theme-primary, #409EFF);
  border-radius: 12px;
  margin-bottom: 15px;
  transition: all 0.3s;
  background: var(--theme-light-bg, #ecf5ff);

  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
    border-color: var(--theme-primary, #409EFF);
  }

  .role-icon {
    width: 60px;
    height: 60px;
    border-radius: 50%;
    background: var(--theme-gradient, linear-gradient(135deg, #667eea 0%, #764ba2 100%));
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 15px;

    i {
      font-size: 28px;
      color: white;
    }
  }

  h4 {
    margin: 10px 0;
    font-size: 16px;
    color: #303133;
  }

  p {
    margin: 0;
    font-size: 13px;
    color: #909399;
    line-height: 1.5;
  }
}

::v-deep .el-collapse {
  border: none;

  .el-collapse-item__header {
    font-size: 15px;
    font-weight: 500;
    background: var(--theme-light-bg, #ecf5ff);
    padding: 0 15px;
    border-radius: 4px;
    margin-bottom: 5px;
    border: 1px solid var(--theme-primary, #409EFF);
    color: var(--theme-primary, #409EFF);

    &:hover {
      background: var(--theme-primary, #409EFF);
      color: white;
    }
  }

  .el-collapse-item__content {
    padding: 15px 20px;
  }

  ul {
    list-style: none;
    padding: 0;
    margin: 0;

    li {
      padding: 8px 0;
      color: #606266;

      &:before {
        content: "•";
        color: var(--theme-primary, #409EFF);
        font-weight: bold;
        margin-right: 10px;
      }
    }
  }
}

.footer-section {
  margin-bottom: 0;

  .license {
    text-align: center;
    color: #909399;
    font-size: 14px;

    i {
      margin-right: 5px;
    }
  }
}

@media (max-width: 768px) {
  .architecture {
    flex-direction: column;

    i.el-icon-right {
      transform: rotate(90deg);
    }
  }
}
</style>
