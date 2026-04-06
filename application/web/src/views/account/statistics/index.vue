<template>
  <div class="statistics-container">
    <el-card class="page-header">
      <h2>
        <i class="el-icon-data-analysis"></i>
        监管中心 - 数据统计分析
      </h2>
      <p>实时监控医疗联盟链数据，全面掌握系统运行状况</p>
    </el-card>

    <!-- 核心指标卡片 -->
    <el-row :gutter="20" class="stats-cards">
      <el-col :xs="24" :sm="12" :md="6" v-for="(stat, index) in coreStats" :key="index">
        <el-card class="stat-card" :body-style="{ padding: '20px' }">
          <div class="stat-content">
            <div class="stat-icon" :style="{ background: stat.color }">
              <i :class="stat.icon"></i>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ stat.value }}</p>
              <p class="stat-title">{{ stat.title }}</p>
              <p class="stat-desc" v-if="stat.desc">{{ stat.desc }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 用户分析 -->
    <el-row :gutter="20">
      <el-col :xs="24" :md="12">
        <el-card class="chart-card">
          <div slot="header" class="card-header">
            <span><i class="el-icon-user"></i> 用户角色分布</span>
          </div>
          <div class="user-role-stats">
            <div class="role-item" v-for="role in roleStats" :key="role.name">
              <div class="role-info">
                <span class="role-name">{{ role.name }}</span>
                <span class="role-count">{{ role.count }} 人</span>
              </div>
              <el-progress 
                :percentage="role.percentage" 
                :color="role.color"
                :stroke-width="12"
              ></el-progress>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="12">
        <el-card class="chart-card">
          <div slot="header" class="card-header">
            <span><i class="el-icon-office-building"></i> 医院用户分布</span>
          </div>
          <div class="hospital-stats">
            <div class="hospital-item" v-for="(count, hospital) in hospitalStats" :key="hospital">
              <div class="hospital-info">
                <span class="hospital-name">{{ hospital }}</span>
                <span class="hospital-count">{{ count }} 人</span>
              </div>
              <el-progress 
                :percentage="getHospitalPercentage(count)" 
                color="#4facfe"
                :stroke-width="10"
              ></el-progress>
            </div>
            <div v-if="Object.keys(hospitalStats).length === 0" class="empty-data">
              <i class="el-icon-info"></i>
              <p>暂无数据</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 用户列表 -->
    <el-card class="table-card">
      <div slot="header" class="card-header">
        <span><i class="el-icon-s-custom"></i> 用户详细列表</span>
        <div class="header-actions">
          <el-button-group style="margin-right: 10px;">
            <el-button 
              size="small" 
              :type="groupByOrganization ? 'primary' : ''"
              @click="groupByOrganization = true"
              icon="el-icon-s-grid"
            >
              按组织分组
            </el-button>
            <el-button 
              size="small" 
              :type="!groupByOrganization ? 'primary' : ''"
              @click="groupByOrganization = false"
              icon="el-icon-menu"
            >
              列表视图
            </el-button>
          </el-button-group>
          <el-input
            v-model="searchKeyword"
            placeholder="搜索用户名或姓名"
            prefix-icon="el-icon-search"
            size="small"
            style="width: 200px; margin-right: 10px;"
            clearable
          ></el-input>
          <el-select v-model="filterOrganization" placeholder="筛选组织" size="small" style="width: 150px; margin-right: 10px;" clearable>
            <el-option label="全部组织" value=""></el-option>
            <el-option label="协和医院" value="TaobaoMSP"></el-option>
            <el-option label="301医院" value="JDMSP"></el-option>
            <el-option label="温江医疗中心" value="WenjinMSP"></el-option>
            <el-option label="华西医院" value="HuaxiMSP"></el-option>
            <el-option label="监管中心" value="RegCenterMSP"></el-option>
            <el-option label="无组织（病人/药店）" value="none"></el-option>
          </el-select>
          <el-select v-model="filterRole" placeholder="筛选角色" size="small" style="width: 120px;" clearable>
            <el-option label="全部角色" value=""></el-option>
            <el-option label="医生" value="医生"></el-option>
            <el-option label="病人" value="病人"></el-option>
            <el-option label="药店" value="药店"></el-option>
          </el-select>
        </div>
      </div>
      
      <!-- 按组织分组显示 -->
      <div v-if="groupByOrganization">
        <div v-for="(users, orgName) in groupedUsers" :key="orgName" class="org-group">
          <div class="org-group-header">
            <i class="el-icon-office-building"></i>
            <span class="org-name">{{ orgName }}</span>
            <el-tag size="small" type="info">{{ users.length }} 人</el-tag>
          </div>
          <el-table
            :data="users"
            style="width: 100%"
            :default-sort="{prop: 'username', order: 'ascending'}"
          >
            <el-table-column prop="username" label="用户名" width="120"></el-table-column>
            <el-table-column prop="account_name" label="姓名" width="150"></el-table-column>
            <el-table-column prop="role" label="角色" width="100">
              <template slot-scope="scope">
                <el-tag :type="getRoleTagType(scope.row.role)" size="small">
                  {{ scope.row.role }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="department" label="科室" width="120"></el-table-column>
            <el-table-column prop="doctor_title" label="职称" width="120"></el-table-column>
            <el-table-column prop="account_id" label="区块链ID" min-width="150" show-overflow-tooltip></el-table-column>
          </el-table>
        </div>
        <div v-if="Object.keys(groupedUsers).length === 0" class="empty-data">
          <i class="el-icon-info"></i>
          <p>暂无数据</p>
        </div>
      </div>
      
      <!-- 普通列表显示 -->
      <el-table
        v-else
        :data="filteredUserList"
        style="width: 100%"
        :default-sort="{prop: 'created_at', order: 'descending'}"
        v-loading="loading"
      >
        <el-table-column prop="username" label="用户名" width="120"></el-table-column>
        <el-table-column prop="account_name" label="姓名" width="150"></el-table-column>
        <el-table-column prop="role" label="角色" width="100">
          <template slot-scope="scope">
            <el-tag :type="getRoleTagType(scope.row.role)" size="small">
              {{ scope.row.role }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="organization_name" label="所属医院" min-width="150"></el-table-column>
        <el-table-column prop="department" label="科室" width="120"></el-table-column>
        <el-table-column prop="doctor_title" label="职称" width="120"></el-table-column>
        <el-table-column prop="account_id" label="区块链ID" width="150" show-overflow-tooltip></el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          @current-change="handlePageChange"
          :current-page="currentPage"
          :page-size="pageSize"
          layout="total, prev, pager, next"
          :total="filteredUserList.length">
        </el-pagination>
      </div>
    </el-card>

    <!-- 系统信息 -->
    <el-row :gutter="20">
      <el-col :xs="24" :md="12">
        <el-card class="info-card">
          <div slot="header" class="card-header">
            <span><i class="el-icon-connection"></i> 区块链网络信息</span>
          </div>
          <div class="info-list">
            <div class="info-item">
              <span class="info-label">网络名称:</span>
              <span class="info-value">社区医疗管理联盟链</span>
            </div>
            <div class="info-item">
              <span class="info-label">通道名称:</span>
              <span class="info-value">appchannel</span>
            </div>
            <div class="info-item">
              <span class="info-label">链码版本:</span>
              <span class="info-value">v1.0.7</span>
            </div>
            <div class="info-item">
              <span class="info-label">组织数量:</span>
              <span class="info-value">4 个</span>
            </div>
            <div class="info-item">
              <span class="info-label">节点数量:</span>
              <span class="info-value">9 个 (8 peer + 1 orderer)</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="12">
        <el-card class="info-card">
          <div slot="header" class="card-header">
            <span><i class="el-icon-office-building"></i> 联盟组织</span>
          </div>
          <div class="org-list">
            <div class="org-item" v-for="org in organizations" :key="org.name">
              <div class="org-icon" :style="{ background: org.color }">
                <i :class="org.icon"></i>
              </div>
              <div class="org-info">
                <p class="org-name">{{ org.name }}</p>
                <p class="org-desc">{{ org.desc }}</p>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { queryAccountList } from '@/api/accountV2'

export default {
  name: 'Statistics',
  data() {
    return {
      loading: false,
      userList: [],
      searchKeyword: '',
      filterRole: '',
      filterOrganization: '',
      groupByOrganization: true,
      currentPage: 1,
      pageSize: 10,
      organizations: [
        { name: '协和医院', desc: 'TaobaoMSP', icon: 'el-icon-office-building', color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' },
        { name: '301医院', desc: 'JDMSP', icon: 'el-icon-office-building', color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)' },
        { name: '温江医疗中心', desc: 'WenjinMSP', icon: 'el-icon-office-building', color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)' },
        { name: '监管中心', desc: 'RegCenterMSP', icon: 'el-icon-s-data', color: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)' }
      ],
      organizationMap: {
        'TaobaoMSP': '协和医院',
        'JDMSP': '301医院',
        'WenjinMSP': '温江医疗中心',
        'RegCenterMSP': '监管中心'
      }
    }
  },

  computed: {
    coreStats() {
      const doctors = this.userList.filter(u => u.role === '医生')
      const patients = this.userList.filter(u => u.role === '病人')
      const drugstores = this.userList.filter(u => u.role === '药店')
      
      return [
        { 
          title: '总用户数', 
          value: this.userList.length, 
          icon: 'el-icon-user-solid', 
          color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
          desc: '系统注册用户总数'
        },
        { 
          title: '医生数量', 
          value: doctors.length, 
          icon: 'el-icon-user', 
          color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
          desc: '各医院医生总数'
        },
        { 
          title: '病人数量', 
          value: patients.length, 
          icon: 'el-icon-user', 
          color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
          desc: '注册病人总数'
        },
        { 
          title: '药店数量', 
          value: drugstores.length, 
          icon: 'el-icon-shopping-bag-2', 
          color: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
          desc: '合作药店数量'
        }
      ]
    },

    roleStats() {
      const total = this.userList.length || 1
      const roles = [
        { name: '医生', key: '医生', color: '#4facfe' },
        { name: '病人', key: '病人', color: '#43e97b' },
        { name: '药店', key: '药店', color: '#fa709a' },
        { name: '其他', key: '其他', color: '#909399' }
      ]
      
      return roles.map(role => {
        const count = this.userList.filter(u => u.role === role.key).length
        return {
          ...role,
          count,
          percentage: Math.round((count / total) * 100)
        }
      })
    },

    hospitalStats() {
      const stats = {}
      this.userList.forEach(user => {
        if (user.organization_name) {
          stats[user.organization_name] = (stats[user.organization_name] || 0) + 1
        }
      })
      return stats
    },

    filteredUserList() {
      let list = this.userList
      
      if (this.filterOrganization) {
        if (this.filterOrganization === 'none') {
          list = list.filter(u => !u.organization || u.organization === '')
        } else {
          list = list.filter(u => u.organization === this.filterOrganization)
        }
      }
      
      if (this.filterRole) {
        list = list.filter(u => u.role === this.filterRole)
      }
      
      if (this.searchKeyword) {
        const keyword = this.searchKeyword.toLowerCase()
        list = list.filter(u => 
          (u.username && u.username.toLowerCase().includes(keyword)) ||
          (u.account_name && u.account_name.toLowerCase().includes(keyword))
        )
      }
      
      return list
    },
    
    groupedUsers() {
      const groups = {}
      const filteredList = this.filteredUserList
      
      filteredList.forEach(user => {
        let groupName = '无组织（病人/药店）'
        
        if (user.organization && user.organization !== '') {
          groupName = this.organizationMap[user.organization] || user.organization_name || user.organization
        }
        
        if (!groups[groupName]) {
          groups[groupName] = []
        }
        groups[groupName].push(user)
      })
      
      const sortedGroups = {}
      const orgOrder = ['协和医院', '301医院', '温江医疗中心', '监管中心', '无组织（病人/药店）']
      
      orgOrder.forEach(orgName => {
        if (groups[orgName]) {
          sortedGroups[orgName] = groups[orgName]
        }
      })
      
      Object.keys(groups).forEach(orgName => {
        if (!sortedGroups[orgName]) {
          sortedGroups[orgName] = groups[orgName]
        }
      })
      
      return sortedGroups
    }
  },

  mounted() {
    this.loadData()
  },

  methods: {
    async loadData() {
      this.loading = true
      try {
        const response = await queryAccountList()
        const data = response.code === 200 ? response.data : response
        this.userList = Array.isArray(data) ? data : []
        console.log('加载用户数据:', this.userList.length, '条')
      } catch (error) {
        console.error('加载数据失败:', error)
        this.$message.error('加载数据失败')
      } finally {
        this.loading = false
      }
    },

    getHospitalPercentage(count) {
      const total = this.userList.length || 1
      return Math.round((count / total) * 100)
    },

    getRoleTagType(role) {
      const typeMap = {
        '医生': 'primary',
        '病人': 'success',
        '药店': 'warning',
        '管理员': 'danger'
      }
      return typeMap[role] || 'info'
    },

    handlePageChange(page) {
      this.currentPage = page
    }
  }
}
</script>

<style lang="scss" scoped>
.statistics-container {
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

.stats-cards {
  margin-bottom: 20px;
  
  .stat-card {
    border-radius: 12px;
    margin-bottom: 20px;
    border: none;
    transition: all 0.3s;
    
    &:hover {
      transform: translateY(-5px);
      box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
    }
    
    .stat-content {
      display: flex;
      align-items: center;
      gap: 15px;
      
      .stat-icon {
        width: 60px;
        height: 60px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
        
        i {
          font-size: 28px;
          color: white;
        }
      }
      
      .stat-info {
        flex: 1;
        
        .stat-value {
          margin: 0 0 5px 0;
          font-size: 28px;
          font-weight: 600;
          color: #303133;
        }
        
        .stat-title {
          margin: 0 0 3px 0;
          font-size: 14px;
          color: #606266;
        }
        
        .stat-desc {
          margin: 0;
          font-size: 12px;
          color: #909399;
        }
      }
    }
  }
}

.chart-card {
  border-radius: 12px;
  margin-bottom: 20px;
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: 600;
    
    i {
      margin-right: 8px;
      color: #4facfe;
    }
  }
  
  .user-role-stats, .hospital-stats {
    .role-item, .hospital-item {
      margin-bottom: 20px;
      
      &:last-child {
        margin-bottom: 0;
      }
      
      .role-info, .hospital-info {
        display: flex;
        justify-content: space-between;
        margin-bottom: 8px;
        
        .role-name, .hospital-name {
          font-size: 14px;
          color: #606266;
        }
        
        .role-count, .hospital-count {
          font-size: 14px;
          font-weight: 600;
          color: #303133;
        }
      }
    }
    
    .empty-data {
      text-align: center;
      padding: 40px 0;
      color: #909399;
      
      i {
        font-size: 48px;
        margin-bottom: 10px;
      }
      
      p {
        margin: 0;
      }
    }
  }
}

.table-card {
  border-radius: 12px;
  margin-bottom: 20px;
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: 600;
    
    i {
      margin-right: 8px;
      color: #4facfe;
    }
    
    .header-actions {
      display: flex;
      gap: 10px;
    }
  }
  
  .pagination-container {
    margin-top: 20px;
    text-align: right;
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
  
  .info-list {
    .info-item {
      display: flex;
      justify-content: space-between;
      padding: 12px 0;
      border-bottom: 1px solid #f0f0f0;
      
      &:last-child {
        border-bottom: none;
      }
      
      .info-label {
        color: #909399;
        font-size: 14px;
      }
      
      .info-value {
        color: #303133;
        font-size: 14px;
        font-weight: 500;
      }
    }
  }
  
  .org-list {
    .org-item {
      display: flex;
      align-items: center;
      gap: 15px;
      padding: 15px 0;
      border-bottom: 1px solid #f0f0f0;
      
      &:last-child {
        border-bottom: none;
      }
      
      .org-icon {
        width: 50px;
        height: 50px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
        
        i {
          font-size: 24px;
          color: white;
        }
      }
      
      .org-info {
        flex: 1;
        
        .org-name {
          margin: 0 0 5px 0;
          font-size: 15px;
          font-weight: 600;
          color: #303133;
        }
        
        .org-desc {
          margin: 0;
          font-size: 13px;
          color: #909399;
        }
      }
    }
  }
}

.org-group {
  margin-bottom: 30px;
  
  &:last-child {
    margin-bottom: 0;
  }
  
  .org-group-header {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 15px 20px;
    background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
    color: white;
    border-radius: 8px 8px 0 0;
    font-size: 16px;
    font-weight: 600;
    
    i {
      font-size: 20px;
    }
    
    .org-name {
      flex: 1;
    }
    
    .el-tag {
      background: rgba(255, 255, 255, 0.2);
      border: none;
      color: white;
    }
  }
  
  .el-table {
    border-radius: 0 0 8px 8px;
    overflow: hidden;
  }
}

.empty-data {
  text-align: center;
  padding: 60px 20px;
  color: #909399;
  
  i {
    font-size: 64px;
    margin-bottom: 15px;
    opacity: 0.5;
  }
  
  p {
    margin: 0;
    font-size: 16px;
  }
}

.header-actions {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}
</style>
