<template>
  <div class="container">
    <el-card class="page-header">
      <h2>
        <i class="el-icon-user-solid"></i>
        用户管理
      </h2>
      <p>查看和管理系统中的所有用户账户</p>
      <el-alert
        type="info"
        :closable="false"
        style="margin-top: 10px;"
      >
        <template slot="title">
          <i class="el-icon-info"></i> 功能说明
        </template>
        <div style="font-size: 13px; line-height: 1.8;">
          • <strong>停用用户</strong>：点击"停用"按钮，用户将无法登录，但数据保留，可随时恢复<br>
          • <strong>恢复用户</strong>：在"已停用"筛选中找到用户，点击"恢复"按钮即可重新启用<br>
          • <strong>状态筛选</strong>：使用状态筛选器查看正常用户、已停用用户或全部用户
        </div>
      </el-alert>
    </el-card>

    <!-- 筛选工具栏 -->
    <el-card class="filter-card">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索用户名或姓名"
            prefix-icon="el-icon-search"
            clearable
          ></el-input>
        </el-col>
        <el-col :span="6">
          <el-select v-model="filterOrganization" placeholder="筛选组织" clearable style="width: 100%">
            <el-option label="全部组织" value=""></el-option>
            <el-option label="协和医院" value="协和医院"></el-option>
            <el-option label="301医院" value="301医院"></el-option>
            <el-option label="温江医疗中心" value="温江医疗中心"></el-option>
            <el-option label="监管中心" value="监管中心"></el-option>
            <el-option label="无组织" value="none"></el-option>
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-select v-model="filterRole" placeholder="筛选角色" clearable style="width: 100%">
            <el-option label="全部角色" value=""></el-option>
            <el-option label="医生" value="医生"></el-option>
            <el-option label="病人" value="病人"></el-option>
            <el-option label="药店" value="药店"></el-option>
            <el-option label="管理员" value="管理员"></el-option>
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-button type="primary" icon="el-icon-refresh" @click="loadData">刷新</el-button>
          <el-tag type="info" style="margin-left: 10px;">共 {{ filteredAccountList.length }} 个</el-tag>
        </el-col>
      </el-row>
    </el-card>

    <div v-if="filteredAccountList.length === 0 && !loading" style="text-align: center; padding: 40px;">
      <el-card style="background: white; border-radius: 12px;">
        <el-empty description="暂无数据">
          <span style="color: #909399;">没有找到符合条件的用户账户</span>
        </el-empty>
      </el-card>
    </div>

    <div v-loading="loading" class="account-grid-container">
      <el-card 
        v-for="(val, index) in filteredAccountList" 
        :key="index" 
        class="account-card" 
        :class="{ 'disabled-card': val.status === 0 }"
      >
          <div slot="header" class="card-header">
            <!-- 使用自定义图标或 Element UI 图标 -->
            <img v-if="val.role === '药店'" src="/image/药店-copy.png" class="role-icon-img" alt="药店" />
            <i v-else :class="getRoleIcon(val.role)" class="role-icon"></i>
            <span class="account-id">{{ val.account_id }}</span>
            <!-- 状态标签 -->
            <el-tag v-if="val.status === 0" type="danger" size="mini" class="status-tag">已停用</el-tag>
          </div>

          <div class="account-info">
            <!-- 固定显示6行，确保所有卡片高度一致 -->
            <div class="info-row">
              <div class="info-item">
                <i class="el-icon-user"></i>
                <span class="label">用户名:</span>
                <span class="value">{{ val.username }}</span>
              </div>
            </div>

            <div class="info-row">
              <div class="info-item">
                <i class="el-icon-postcard"></i>
                <span class="label">姓名:</span>
                <span class="value">{{ val.account_name }}</span>
              </div>
            </div>

            <div class="info-row">
              <div class="info-item">
                <i class="el-icon-s-custom"></i>
                <span class="label">角色:</span>
                <el-tag :type="getRoleTagType(val.role)" size="small">{{ val.role }}</el-tag>
              </div>
            </div>

            <div class="info-row">
              <div class="info-item">
                <i class="el-icon-office-building"></i>
                <span class="label">组织:</span>
                <span class="value">{{ getOrganizationDisplayName(val.organization, val.organization_name) }}</span>
              </div>
            </div>

            <div class="info-row">
              <div class="info-item">
                <i class="el-icon-folder-opened"></i>
                <span class="label">科室:</span>
                <span class="value">{{ val.department || '-' }}</span>
              </div>
            </div>

            <div class="info-row">
              <div class="info-item">
                <i class="el-icon-medal"></i>
                <span class="label">职称:</span>
                <span class="value">{{ val.doctor_title || '-' }}</span>
              </div>
            </div>
          </div>

        <div class="card-actions">
          <el-button type="primary" size="mini" icon="el-icon-view" @click="viewUser(val)">查看详情</el-button>
          <el-button type="warning" size="mini" icon="el-icon-edit" @click="editUser(val)">编辑信息</el-button>
        </div>
        </el-card>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { queryAccountListFromDB, getUserDetail } from '@/api/accountV2'

export default {
  name: 'Account',
  data() {
    return {
      loading: true,
      accountList: [],
      searchKeyword: '',
      filterOrganization: '',
      filterRole: ''
    }
  },
  computed: {
    ...mapGetters([
      'account_id',
      'roles',
      'account_name',
    ]),
    
    // 过滤后的账户列表
    filteredAccountList() {
      let list = this.accountList
      
      // 组织筛选 - 使用 organization_name 而不是 organization
      if (this.filterOrganization) {
        if (this.filterOrganization === 'none') {
          list = list.filter(u => !u.organization_name || u.organization_name === '')
        } else {
          list = list.filter(u => u.organization_name === this.filterOrganization)
        }
      }
      
      // 角色筛选
      if (this.filterRole) {
        list = list.filter(u => u.role === this.filterRole)
      }
      
      // 关键词搜索
      if (this.searchKeyword) {
        const keyword = this.searchKeyword.toLowerCase()
        list = list.filter(u => 
          (u.username && u.username.toLowerCase().includes(keyword)) ||
          (u.account_name && u.account_name.toLowerCase().includes(keyword))
        )
      }
      
      return list
    }
  },
  created() {
    this.loadData()
  },
  methods: {
    // 组织MSPID到名称的映射
    getOrganizationDisplayName(org, orgName) {
      // 如果有organization_name，直接使用
      if (orgName && orgName !== '') {
        return orgName
      }
      
      // 如果没有，根据organization MSPID映射
      const orgMap = {
        'TaobaoMSP': '协和医院',
        'JDMSP': '301医院',
        'WenjinMSP': '温江社区医疗中心',
        'RegCenterMSP': '监管中心'
      }
      
      return orgMap[org] || '-'
    },

    loadData() {
      console.log('=== 用户管理页面加载（从数据库） ===')
      this.loading = true
      
      // 从数据库查询用户列表，只查询正常用户
      queryAccountListFromDB(1).then(response => {
          console.log('API 响应:', response)
          
          // 正确处理响应格式
          let accountList = []
          if (response && response.code === 200 && response.data) {
            accountList = response.data
          } else if (Array.isArray(response)) {
            // 兼容旧格式
            accountList = response
          }
          
          this.accountList = accountList
          console.log('✅ 从数据库加载账户列表成功，数量:', this.accountList.length)
          this.loading = false
        }).catch(error => {
          console.error('❌ 加载账户列表失败:', error)
          this.$message.error('加载账户列表失败')
          this.loading = false
        })
    },
    
    viewUser(user) {
      // 通过用户名跳转到详情页
      this.$router.push({
        path: '/account/detail',
        query: { username: user.username }
      })
    },
    
    editUser(user) {
      // 通过用户名跳转到编辑页
      this.$router.push({
        path: '/account/edit',
        query: { username: user.username }
      })
    },
    
    getRoleIcon(role) {
      const iconMap = {
        '医生': 'el-icon-user',
        '病人': 'el-icon-user',
        '药店': 'el-icon-shop',  // 改为商店图标
        '管理员': 'el-icon-s-custom'
      }
      return iconMap[role] || 'el-icon-user'
    },
    
    getRoleTagType(role) {
      const typeMap = {
        '医生': 'primary',
        '病人': 'success',
        '药店': 'warning',
        '管理员': 'danger'
      }
      return typeMap[role] || 'info'
    }
  }
}
</script>

<style lang="scss" scoped>
.container {
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: calc(100vh - 60px);
}

.page-header {
  margin-bottom: 20px;
  border-radius: 12px;
  background: white;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  
  h2 {
    margin: 0 0 10px 0;
    font-size: 24px;
    color: #303133;
    
    i {
      color: #667eea;
      margin-right: 10px;
    }
  }
  
  p {
    margin: 0;
    color: #909399;
    font-size: 14px;
  }
}

.filter-card {
  margin-bottom: 20px;
  border-radius: 12px;
  background: white;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.account-grid-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-top: 20px;
  align-items: start;
}

.account-card {
  width: 100%;
  border-radius: 12px;
  transition: all 0.3s;
  border: 2px solid #EBEEF5;
  background: white;
  display: flex;
  flex-direction: column;
  
  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 8px 20px rgba(102, 126, 234, 0.3);
    border-color: #667eea;
  }
  
  &.disabled-card {
    opacity: 0.7;
    background: #f5f7fa;
    
    &:hover {
      border-color: #e4e7ed;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    }
    
    ::v-deep .el-card__header {
      background: #e4e7ed;
    }
  }
  
  ::v-deep .el-card__body {
    padding: 15px;
    display: flex;
    flex-direction: column;
    flex: 1;
  }
  
  .card-header {
    display: flex;
    align-items: center;
    gap: 10px;
    font-weight: 600;
    color: #303133;
    position: relative;
    
    .role-icon {
      font-size: 18px;
      color: #667eea;
      flex-shrink: 0;
    }
    
    .role-icon-img {
      width: 20px;
      height: 20px;
      flex-shrink: 0;
      object-fit: contain;
    }
    
    .account-id {
      flex: 1;
      font-size: 13px;
      color: #667eea;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    
    .status-tag {
      position: absolute;
      right: 0;
      top: 50%;
      transform: translateY(-50%);
    }
  }
  
  ::v-deep .el-card__header {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-bottom: none;
    padding: 12px 15px;
    
    .card-header {
      color: white;
      
      .role-icon {
        color: white;
      }
      
      .account-id {
        color: rgba(255, 255, 255, 0.9);
      }
    }
  }
  
  .account-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    
    .info-row {
      min-height: 36px;
      display: flex;
      align-items: center;
      padding: 6px 0;
      border-bottom: 1px solid #f0f0f0;
      
      &:last-child {
        border-bottom: none;
      }
      
      .info-item {
        width: 100%;
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 13px;
        
        i {
          color: #667eea;
          font-size: 14px;
          width: 16px;
          flex-shrink: 0;
        }
        
        .label {
          color: #909399;
          min-width: 50px;
          flex-shrink: 0;
        }
        
        .value {
          color: #303133;
          font-weight: 500;
          flex: 1;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }
    }
  }
  
  .card-actions {
    margin-top: 15px;
    padding-top: 15px;
    border-top: 1px solid #f0f0f0;
    display: flex;
    justify-content: space-between;
    gap: 5px;
    
    .el-button {
      flex: 1;
      padding: 7px 10px;
      font-size: 12px;
    }
  }
}

@media (max-width: 768px) {
  .filter-card {
    ::v-deep .el-col {
      margin-bottom: 10px;
    }
  }
}
</style>
