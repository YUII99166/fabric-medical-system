<template>
  <div class="app-container">
    <el-card v-loading="loading" class="detail-card">
      <div slot="header" class="card-header">
        <div class="header-left">
          <i class="el-icon-user-solid"></i>
          <span>用户详情</span>
        </div>
        <div class="header-right">
          <el-button size="small" icon="el-icon-back" @click="goBack">返回</el-button>
          <el-button size="small" type="warning" icon="el-icon-edit" @click="editUser">编辑</el-button>
          <el-button size="small" type="danger" icon="el-icon-delete" @click="confirmDelete">删除</el-button>
        </div>
      </div>

      <div v-if="userDetail" class="detail-content">
        <!-- 基本信息 -->
        <div class="info-section">
          <h3 class="section-title">
            <i class="el-icon-user"></i>
            基本信息
          </h3>
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="info-item">
                <span class="label">用户名:</span>
                <span class="value">{{ userDetail.username }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="info-item">
                <span class="label">姓名:</span>
                <span class="value">{{ userDetail.account_name }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="info-item">
                <span class="label">角色:</span>
                <el-tag :type="getRoleTagType(userDetail.role)">{{ userDetail.role }}</el-tag>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="info-item">
                <span class="label">用户ID:</span>
                <span class="value">{{ userDetail.id }}</span>
              </div>
            </el-col>
          </el-row>
        </div>

        <!-- 组织信息 -->
        <div v-if="userDetail.organization || userDetail.department || userDetail.doctor_title" class="info-section">
          <h3 class="section-title">
            <i class="el-icon-office-building"></i>
            组织信息
          </h3>
          <el-row :gutter="20">
            <el-col v-if="userDetail.organization_name" :span="12">
              <div class="info-item">
                <span class="label">所属组织:</span>
                <span class="value">{{ userDetail.organization_name }}</span>
              </div>
            </el-col>
            <el-col v-if="userDetail.organization" :span="12">
              <div class="info-item">
                <span class="label">组织MSP:</span>
                <span class="value">{{ userDetail.organization }}</span>
              </div>
            </el-col>
            <el-col v-if="userDetail.department" :span="12">
              <div class="info-item">
                <span class="label">科室:</span>
                <span class="value">{{ userDetail.department }}</span>
              </div>
            </el-col>
            <el-col v-if="userDetail.doctor_title" :span="12">
              <div class="info-item">
                <span class="label">职称:</span>
                <span class="value">{{ userDetail.doctor_title }}</span>
              </div>
            </el-col>
          </el-row>
        </div>

        <!-- 时间信息 -->
        <div class="info-section">
          <h3 class="section-title">
            <i class="el-icon-time"></i>
            时间信息
          </h3>
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="info-item">
                <span class="label">创建时间:</span>
                <span class="value">{{ formatTime(userDetail.created_at) }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="info-item">
                <span class="label">更新时间:</span>
                <span class="value">{{ formatTime(userDetail.updated_at) }}</span>
              </div>
            </el-col>
          </el-row>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script>
import { getUserDetail, deleteUser } from '@/api/accountV2'
import { mapGetters } from 'vuex'

export default {
  name: 'AccountDetail',
  data() {
    return {
      loading: false,
      userDetail: null
    }
  },
  computed: {
    ...mapGetters(['account_id', 'roles'])
  },
  created() {
    this.loadUserDetail()
  },
  methods: {
    async loadUserDetail() {
      const username = this.$route.query.username
      const id = this.$route.query.id
      
      if (!username && !id) {
        this.$message.error('缺少用户信息')
        this.goBack()
        return
      }
      
      this.loading = true
      
      try {
        // 优先使用ID，否则使用用户名
        const response = await getUserDetail(id || username)
        if (response && response.code === 200) {
          this.userDetail = response.data
        } else {
          this.$message.error(response.msg || '加载用户详情失败')
          this.goBack()
        }
      } catch (error) {
        console.error('加载用户详情失败:', error)
        this.$message.error('加载用户详情失败')
        this.goBack()
      } finally {
        this.loading = false
      }
    },
    
    editUser() {
      if (this.userDetail) {
        this.$router.push({
          path: '/account/edit',
          query: { id: this.userDetail.id }
        })
      }
    },
    
    confirmDelete() {
      if (!this.userDetail) return
      
      this.$confirm(
        `确定要删除用户 "${this.userDetail.account_name}" (${this.userDetail.username}) 吗？此操作不可恢复。`,
        '删除确认',
        {
          confirmButtonText: '确定删除',
          cancelButtonText: '取消',
          type: 'warning'
        }
      ).then(() => {
        this.handleDelete()
      }).catch(() => {
        this.$message.info('已取消删除')
      })
    },
    
    async handleDelete() {
      try {
        const response = await deleteUser(this.userDetail.id)
        if (response && response.code === 200) {
          this.$message.success('删除成功')
          setTimeout(() => {
            this.$router.push('/account/all')
          }, 1000)
        } else {
          this.$message.error(response.msg || '删除失败')
        }
      } catch (error) {
        console.error('删除用户失败:', error)
        if (error.response && error.response.data) {
          this.$message.error(error.response.data.msg || '删除失败')
        } else {
          this.$message.error('删除失败，请重试')
        }
      }
    },
    
    goBack() {
      this.$router.back()
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
    
    formatTime(time) {
      if (!time) return '-'
      return new Date(time).toLocaleString('zh-CN')
    }
  }
}
</script>

<style lang="scss" scoped>
.app-container {
  padding: 20px;
}

.detail-card {
  max-width: 1000px;
  margin: 0 auto;
  border-radius: 12px;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .header-left {
      display: flex;
      align-items: center;
      gap: 10px;
      font-size: 18px;
      font-weight: 600;
      color: #303133;

      i {
        color: #409EFF;
        font-size: 20px;
      }
    }

    .header-right {
      display: flex;
      gap: 10px;
    }
  }
}

.detail-content {
  padding: 20px 0;
}

.info-section {
  margin-bottom: 30px;
  padding: 20px;
  background: #f5f7fa;
  border-radius: 8px;

  &:last-child {
    margin-bottom: 0;
  }

  .section-title {
    margin: 0 0 20px 0;
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    padding-bottom: 10px;
    border-bottom: 2px solid #409EFF;
    display: flex;
    align-items: center;
    gap: 8px;

    i {
      color: #409EFF;
    }
  }

  .info-item {
    padding: 12px 0;
    display: flex;
    align-items: center;

    .label {
      font-weight: 500;
      color: #606266;
      min-width: 100px;
    }

    .value {
      color: #303133;
      flex: 1;
    }
  }
}

@media (max-width: 768px) {
  .card-header {
    flex-direction: column;
    gap: 15px;
    align-items: flex-start !important;

    .header-right {
      width: 100%;
      justify-content: flex-start;
    }
  }
}
</style>
