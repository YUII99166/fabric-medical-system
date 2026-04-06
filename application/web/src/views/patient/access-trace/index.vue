<template>
  <div class="access-trace-container">
    <el-card class="page-header">
      <h2>
        <i class="el-icon-view"></i>
        隐私溯源
      </h2>
      <p>查看谁访问了您的病历，保护您的隐私安全</p>
      <el-alert
        type="info"
        :closable="false"
        style="margin-top: 10px;"
      >
        <template slot="title">
          <i class="el-icon-info"></i> 功能说明
        </template>
        <div style="font-size: 13px; line-height: 1.8;">
          • <strong>访问记录</strong>：所有非本人的病历访问都会被记录在区块链上，不可篡改<br>
          • <strong>隐私保护</strong>：您可以随时查看谁在什么时间访问了您的病历<br>
          • <strong>安全监控</strong>：如发现异常访问，请及时联系管理员
        </div>
      </el-alert>
    </el-card>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="8">
        <el-card class="stat-card stat-card-blue">
          <div class="stat-content">
            <div class="stat-icon">
              <i class="el-icon-view"></i>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ statistics.total_access || 0 }}</p>
              <p class="stat-label">总访问次数</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="stat-card stat-card-green">
          <div class="stat-content">
            <div class="stat-icon">
              <i class="el-icon-user"></i>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ statistics.unique_accessors || 0 }}</p>
              <p class="stat-label">不同访问者</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="stat-card stat-card-orange">
          <div class="stat-content">
            <div class="stat-icon">
              <i class="el-icon-time"></i>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ statistics.recent_access || 0 }}</p>
              <p class="stat-label">最近7天访问</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 访问日志列表 -->
    <el-card class="logs-card" v-loading="loading">
      <div slot="header" class="card-header">
        <span><i class="el-icon-tickets"></i> 访问记录</span>
        <el-tag type="info" size="small">共 {{ total }} 条记录</el-tag>
      </div>

      <el-timeline v-if="accessLogs.length > 0">
        <el-timeline-item
          v-for="(log, index) in accessLogs"
          :key="index"
          :timestamp="log.accessed_at"
          placement="top"
          :color="getAccessColor(log.access_type)"
        >
          <el-card class="log-card">
            <div class="log-content">
              <div class="log-header">
                <div class="log-title">
                  <el-tag :type="getRoleTagType(log.accessor_role)" size="small">
                    {{ log.accessor_role }}
                  </el-tag>
                  <span class="accessor-name">{{ log.accessor_name }}</span>
                  <el-tag type="warning" size="mini" v-if="log.accessor_organization_name">
                    {{ log.accessor_organization_name }}
                  </el-tag>
                </div>
                <el-tag :type="getAccessTypeTag(log.access_type)" size="small">
                  {{ getAccessTypeText(log.access_type) }}
                </el-tag>
              </div>
              
              <div class="log-body">
                <div class="log-item">
                  <i class="el-icon-document"></i>
                  <span class="log-label">病历编号：</span>
                  <span class="log-value">{{ log.prescription_no }}</span>
                </div>
                <div class="log-item" v-if="log.access_reason">
                  <i class="el-icon-chat-line-square"></i>
                  <span class="log-label">访问原因：</span>
                  <span class="log-value">{{ log.access_reason }}</span>
                </div>
                <div class="log-item" v-if="log.ip_address">
                  <i class="el-icon-location-information"></i>
                  <span class="log-label">IP地址：</span>
                  <span class="log-value">{{ log.ip_address }}</span>
                </div>
              </div>

              <div class="log-footer">
                <el-button 
                  type="text" 
                  size="mini" 
                  @click="viewPrescription(log.prescription_id)"
                >
                  查看病历 <i class="el-icon-arrow-right"></i>
                </el-button>
              </div>
            </div>
          </el-card>
        </el-timeline-item>
      </el-timeline>

      <div v-else class="empty-state">
        <i class="el-icon-document" style="font-size: 64px; color: #C0C4CC;"></i>
        <p style="color: #909399; margin-top: 16px; font-size: 14px;">暂无访问记录</p>
        <p style="color: #C0C4CC; font-size: 12px;">您的病历还没有被其他人访问过</p>
      </div>

      <!-- 分页 -->
      <div class="pagination-container" v-if="total > 0">
        <el-pagination
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="currentPage"
          :page-sizes="[10, 20, 50, 100]"
          :page-size="pageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
        ></el-pagination>
      </div>
    </el-card>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { getMyAccessLogs, getAccessStatistics } from '@/api/accessTrace'

export default {
  name: 'AccessTrace',
  data() {
    return {
      loading: false,
      accessLogs: [],
      statistics: {
        total_access: 0,
        unique_accessors: 0,
        recent_access: 0
      },
      currentPage: 1,
      pageSize: 20,
      total: 0
    }
  },
  computed: {
    ...mapGetters([
      'account_id',
      'account_name'
    ])
  },
  created() {
    this.loadAccessLogs()
    this.loadStatistics()
  },
  methods: {
    async loadAccessLogs() {
      this.loading = true
      console.log('🔍 开始加载访问日志')
      console.log('   当前用户 account_id:', this.account_id)
      console.log('   当前用户 account_name:', this.account_name)
      
      try {
        const params = {
          patient_id: this.account_id,
          page: this.currentPage,
          page_size: this.pageSize
        }
        console.log('   请求参数:', params)
        
        const response = await getMyAccessLogs(params)
        console.log('   API 响应:', response)

        if (response && response.code === 200) {
          this.accessLogs = response.data.logs || []
          this.total = response.data.total || 0
          console.log('   ✅ 加载成功，记录数:', this.accessLogs.length, '总数:', this.total)
        } else {
          console.log('   ❌ 响应失败:', response)
        }
      } catch (error) {
        console.error('❌ 加载访问日志失败:', error)
        this.$message.error('加载访问日志失败')
      } finally {
        this.loading = false
      }
    },

    async loadStatistics() {
      try {
        const response = await getAccessStatistics(this.account_id)
        if (response && response.code === 200) {
          this.statistics = response.data
        }
      } catch (error) {
        console.error('加载统计数据失败:', error)
      }
    },

    handleSizeChange(val) {
      this.pageSize = val
      this.currentPage = 1
      this.loadAccessLogs()
    },

    handleCurrentChange(val) {
      this.currentPage = val
      this.loadAccessLogs()
    },

    getAccessColor(type) {
      const colorMap = {
        'view': '#409EFF',
        'edit': '#E6A23C',
        'download': '#F56C6C'
      }
      return colorMap[type] || '#909399'
    },

    getAccessTypeTag(type) {
      const tagMap = {
        'view': 'primary',
        'edit': 'warning',
        'download': 'danger'
      }
      return tagMap[type] || 'info'
    },

    getAccessTypeText(type) {
      const textMap = {
        'view': '查看',
        'edit': '编辑',
        'download': '下载'
      }
      return textMap[type] || '未知'
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

    viewPrescription(prescriptionId) {
      this.$router.push({
        path: '/prescription/detail',
        query: { id: prescriptionId }
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.access-trace-container {
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

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  transition: all 0.3s;

  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
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
      font-size: 28px;
      color: white;
    }

    .stat-info {
      flex: 1;

      .stat-value {
        font-size: 28px;
        font-weight: bold;
        margin: 0 0 5px 0;
        color: #303133;
      }

      .stat-label {
        font-size: 14px;
        color: #909399;
        margin: 0;
      }
    }
  }

  &.stat-card-blue .stat-icon {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  &.stat-card-green .stat-icon {
    background: linear-gradient(135deg, #67C23A 0%, #85CE61 100%);
  }

  &.stat-card-orange .stat-icon {
    background: linear-gradient(135deg, #E6A23C 0%, #F56C6C 100%);
  }
}

.logs-card {
  border-radius: 12px;
  background: white;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: 600;
    font-size: 16px;

    i {
      margin-right: 8px;
      color: #667eea;
    }
  }
}

.log-card {
  border: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s;

  &:hover {
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
    transform: translateX(5px);
  }

  ::v-deep .el-card__body {
    padding: 15px;
  }
}

.log-content {
  .log-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;

    .log-title {
      display: flex;
      align-items: center;
      gap: 10px;

      .accessor-name {
        font-size: 16px;
        font-weight: 600;
        color: #303133;
      }
    }
  }

  .log-body {
    margin-bottom: 10px;

    .log-item {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 8px;
      font-size: 13px;
      color: #606266;

      i {
        color: #909399;
      }

      .log-label {
        color: #909399;
      }

      .log-value {
        color: #303133;
        font-weight: 500;
      }
    }
  }

  .log-footer {
    border-top: 1px solid #f0f0f0;
    padding-top: 10px;
  }
}

.pagination-container {
  margin-top: 20px;
  text-align: center;
}

::v-deep .el-timeline-item__timestamp {
  font-size: 13px;
  color: #909399;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
  margin-top: 20px;
}
</style>
