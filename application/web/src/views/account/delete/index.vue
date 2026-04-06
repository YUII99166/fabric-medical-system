<template>
  <div class="container">
    <el-card class="page-header">
      <h2>
        <i class="el-icon-delete"></i>
        删除用户
      </h2>
      <p>停用或恢复系统用户账户</p>
      <el-alert
        type="warning"
        :closable="false"
        style="margin-top: 10px;"
      >
        <template slot="title">
          <i class="el-icon-warning"></i> 重要说明
        </template>
        <div style="font-size: 13px; line-height: 1.8;">
          • <strong>停用用户</strong>：用户将无法登录系统，但所有数据保留在区块链和数据库中，可随时恢复<br>
          • <strong>恢复用户</strong>：已停用的用户可以重新启用，恢复正常使用<br>
          • <strong>数据安全</strong>：用户的病历、订单等数据不会丢失，符合医疗行业合规要求<br>
          • <strong>搜索功能</strong>：可以按用户名、姓名、组织或科室进行搜索
        </div>
      </el-alert>
    </el-card>

    <!-- 标签页切换 -->
    <el-card class="tabs-card">
      <el-tabs v-model="activeTab" @tab-click="handleTabClick">
        <el-tab-pane label="正常用户" name="active">
          <span slot="label">
            <i class="el-icon-user"></i> 正常用户 ({{ filteredActiveUsers.length }})
          </span>
        </el-tab-pane>
        <el-tab-pane label="已停用用户" name="disabled">
          <span slot="label">
            <i class="el-icon-remove-outline"></i> 已停用用户 ({{ filteredDisabledUsers.length }})
          </span>
        </el-tab-pane>
      </el-tabs>

      <!-- 搜索框 -->
      <div style="margin-top: 15px;">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索用户名、姓名、组织或科室"
          prefix-icon="el-icon-search"
          clearable
          style="max-width: 400px;"
        ></el-input>
      </div>
    </el-card>

    <!-- 用户列表 -->
    <el-card class="list-card" v-loading="loading">
      <div v-if="currentUsers.length === 0" style="text-align: center; padding: 40px;">
        <el-empty :description="searchKeyword ? '没有找到符合条件的用户' : '暂无用户数据'">
          <span style="color: #909399;">{{ searchKeyword ? '请尝试其他搜索关键词' : '当前分类下没有用户' }}</span>
        </el-empty>
      </div>
      
      <el-table
        v-else
        :data="currentUsers"
        style="width: 100%"
        :row-class-name="tableRowClassName"
      >
        <el-table-column prop="username" label="用户名" width="150"></el-table-column>
        <el-table-column prop="account_name" label="姓名" width="150"></el-table-column>
        <el-table-column prop="role" label="角色" width="100">
          <template slot-scope="scope">
            <el-tag :type="getRoleTagType(scope.row.role)" size="small">{{ scope.row.role }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="组织" width="200">
          <template slot-scope="scope">
            {{ getOrganizationDisplayName(scope.row.organization, scope.row.organization_name) }}
          </template>
        </el-table-column>
        <el-table-column prop="department" label="科室" width="150"></el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180"></el-table-column>
        <el-table-column label="当前状态" width="100">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.status === 1" type="success" size="small">正常</el-tag>
            <el-tag v-else type="danger" size="small">已停用</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="150">
          <template slot-scope="scope">
            <el-button
              v-if="scope.row.status === 1"
              type="danger"
              size="mini"
              icon="el-icon-delete"
              @click="confirmDelete(scope.row)"
            >停用</el-button>
            <el-button
              v-else
              type="success"
              size="mini"
              icon="el-icon-refresh-right"
              @click="confirmRestore(scope.row)"
            >恢复</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { queryAccountListFromDB, getUserDetail, deleteUser, restoreUser } from '@/api/accountV2'

export default {
  name: 'DeleteAccount',
  data() {
    return {
      loading: true,
      activeTab: 'active',
      activeUsers: [],
      disabledUsers: [],
      searchKeyword: ''
    }
  },
  computed: {
    ...mapGetters([
      'account_id',
      'roles',
      'account_name',
    ]),
    filteredActiveUsers() {
      return this.filterUsers(this.activeUsers)
    },
    filteredDisabledUsers() {
      return this.filterUsers(this.disabledUsers)
    },
    currentUsers() {
      return this.activeTab === 'active' ? this.filteredActiveUsers : this.filteredDisabledUsers
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

    filterUsers(users) {
      if (!this.searchKeyword) {
        return users
      }
      
      const keyword = this.searchKeyword.toLowerCase()
      return users.filter(u => 
        (u.username && u.username.toLowerCase().includes(keyword)) ||
        (u.account_name && u.account_name.toLowerCase().includes(keyword)) ||
        (u.organization_name && u.organization_name.toLowerCase().includes(keyword)) ||
        (u.department && u.department.toLowerCase().includes(keyword))
      )
    },

    async loadData() {
      this.loading = true
      try {
        // 加载正常用户
        const activeResponse = await queryAccountListFromDB(1)
        if (activeResponse && activeResponse.code === 200) {
          this.activeUsers = activeResponse.data || []
        }

        // 加载已停用用户
        const disabledResponse = await queryAccountListFromDB(0)
        if (disabledResponse && disabledResponse.code === 200) {
          this.disabledUsers = disabledResponse.data || []
        }

        console.log('✅ 加载用户列表成功')
      } catch (error) {
        console.error('❌ 加载用户列表失败:', error)
        this.$message.error('加载用户列表失败')
      } finally {
        this.loading = false
      }
    },

    handleTabClick(tab) {
      console.log('切换标签:', tab.name)
    },

    tableRowClassName({ row }) {
      return row.status === 0 ? 'disabled-row' : ''
    },

    async confirmDelete(user) {
      try {
        const response = await getUserDetail(user.username)
        if (response && response.code === 200) {
          const userDetail = response.data

          // 检查是否为管理员
          if (userDetail.role === '管理员') {
            this.$confirm(
              `"${userDetail.account_name}" 是管理员账户，停用后可能影响系统管理。确定要停用吗？`,
              '停用管理员确认',
              {
                confirmButtonText: '确定停用',
                cancelButtonText: '取消',
                type: 'error'
              }
            ).then(() => {
              this.handleDelete(userDetail.id)
            }).catch(() => {
              this.$message.info('已取消操作')
            })
          } else {
            this.$confirm(
              `确定要停用用户 "${userDetail.account_name}" (${userDetail.username}) 吗？<br><br>
              <span style="color: #E6A23C;">停用后该用户将无法登录系统，但数据会保留，可以随时恢复。</span>`,
              '停用确认',
              {
                confirmButtonText: '确定停用',
                cancelButtonText: '取消',
                type: 'warning',
                dangerouslyUseHTMLString: true
              }
            ).then(() => {
              this.handleDelete(userDetail.id)
            }).catch(() => {
              this.$message.info('已取消操作')
            })
          }
        } else {
          this.$message.error('获取用户信息失败')
        }
      } catch (error) {
        console.error('获取用户信息失败:', error)
        this.$message.error('获取用户信息失败')
      }
    },

    async handleDelete(userId) {
      try {
        const response = await deleteUser(userId)
        if (response && response.code === 200) {
          this.$message.success('用户已停用')
          this.loadData()
        } else {
          this.$message.error(response.msg || '停用失败')
        }
      } catch (error) {
        console.error('停用用户失败:', error)
        if (error.response && error.response.data) {
          this.$message.error(error.response.data.msg || '停用失败')
        } else {
          this.$message.error('停用失败，请重试')
        }
      }
    },

    async confirmRestore(user) {
      try {
        const response = await getUserDetail(user.username)
        if (response && response.code === 200) {
          const userDetail = response.data

          this.$confirm(
            `确定要恢复用户 "${userDetail.account_name}" (${userDetail.username}) 吗？<br><br>
            <span style="color: #67C23A;">恢复后该用户可以正常登录使用系统。</span>`,
            '恢复确认',
            {
              confirmButtonText: '确定恢复',
              cancelButtonText: '取消',
              type: 'success',
              dangerouslyUseHTMLString: true
            }
          ).then(() => {
            this.handleRestore(userDetail.id)
          }).catch(() => {
            this.$message.info('已取消操作')
          })
        } else {
          this.$message.error('获取用户信息失败')
        }
      } catch (error) {
        console.error('获取用户信息失败:', error)
        this.$message.error('获取用户信息失败')
      }
    },

    async handleRestore(userId) {
      try {
        const response = await restoreUser(userId)
        if (response && response.code === 200) {
          this.$message.success('用户已恢复')
          this.loadData()
        } else {
          this.$message.error(response.msg || '恢复失败')
        }
      } catch (error) {
        console.error('恢复用户失败:', error)
        if (error.response && error.response.data) {
          this.$message.error(error.response.data.msg || '恢复失败')
        } else {
          this.$message.error('恢复失败，请重试')
        }
      }
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
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
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
      color: #f5576c;
      margin-right: 10px;
    }
  }

  p {
    margin: 0;
    color: #909399;
    font-size: 14px;
  }
}

.tabs-card {
  margin-bottom: 20px;
  border-radius: 12px;
  background: white;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);

  ::v-deep .el-input__inner {
    border-radius: 20px;
    border: 2px solid #DCDFE6;
    transition: all 0.3s;

    &:focus {
      border-color: #f5576c;
      box-shadow: 0 0 0 2px rgba(245, 87, 108, 0.1);
    }
  }
}

.list-card {
  border-radius: 12px;
  background: white;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

::v-deep .disabled-row {
  background: #f5f7fa;
  opacity: 0.7;
}

::v-deep .el-tabs__item {
  font-size: 15px;
  font-weight: 500;
}

::v-deep .el-tabs__item.is-active {
  color: #f5576c;
}

::v-deep .el-tabs__active-bar {
  background-color: #f5576c;
}
</style>
