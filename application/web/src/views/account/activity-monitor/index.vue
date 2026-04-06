<template>
  <div class="activity-monitor-container">
    <el-card class="page-header">
      <h2>
        <i class="el-icon-document-copy"></i>
        日志监控
      </h2>
      <p>实时监控联盟中所有组织的业务日志</p>
    </el-card>

    <!-- 筛选工具栏 -->
    <el-card class="filter-card">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-select v-model="filterOrganization" placeholder="筛选组织" clearable style="width: 100%">
            <el-option label="全部组织" value=""></el-option>
            <el-option label="协和医院" value="TaobaoMSP"></el-option>
            <el-option label="301医院" value="JDMSP"></el-option>
            <el-option label="温江社区医疗中心" value="WenjinMSP"></el-option>
            <el-option label="监管中心" value="RegCenterMSP"></el-option>
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-select v-model="filterActivityType" placeholder="活动类型" clearable style="width: 100%">
            <el-option label="全部类型" value=""></el-option>
            <el-option label="创建病历" value="prescription"></el-option>
            <el-option label="补充记录" value="supplement"></el-option>
            <el-option label="授权请求" value="access_request"></el-option>
            <el-option label="药品订单" value="drug_order"></el-option>
            <el-option label="保险报销" value="insurance_claim"></el-option>
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 100%"
            value-format="yyyy-MM-dd"
          ></el-date-picker>
        </el-col>
        <el-col :span="6">
          <el-button type="primary" icon="el-icon-search" @click="loadActivities">查询</el-button>
          <el-button icon="el-icon-refresh" @click="resetFilters">重置</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card stat-card-blue">
          <div class="stat-content">
            <div class="stat-icon">
              <i class="el-icon-document"></i>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ statistics.totalPrescriptions }}</p>
              <p class="stat-label">病历总数</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card stat-card-green">
          <div class="stat-content">
            <div class="stat-icon">
              <i class="el-icon-plus"></i>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ statistics.totalSupplements }}</p>
              <p class="stat-label">补充记录</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card stat-card-orange">
          <div class="stat-content">
            <div class="stat-icon">
              <i class="el-icon-shopping-cart-2"></i>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ statistics.totalOrders }}</p>
              <p class="stat-label">药品订单</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card stat-card-purple">
          <div class="stat-content">
            <div class="stat-icon">
              <i class="el-icon-s-claim"></i>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ statistics.totalClaims }}</p>
              <p class="stat-label">保险报销</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 活动时间线 -->
    <el-card class="activity-card" v-loading="loading">
      <div slot="header" class="card-header">
        <span><i class="el-icon-tickets"></i> 业务日志</span>
        <el-tag type="info" size="small">共 {{ activities.length }} 条记录</el-tag>
      </div>

      <div v-if="activities.length > 0">
        <el-timeline>
          <el-timeline-item
            v-for="(activity, index) in activities"
            :key="index"
            :timestamp="activity.timestamp"
            placement="top"
            :color="getActivityColor(activity.type)"
          >
            <el-card class="timeline-card">
              <div class="timeline-content">
                <div class="timeline-header">
                  <el-tag :type="getActivityTagType(activity.type)" size="small">
                    {{ getActivityTypeText(activity.type) }}
                  </el-tag>
                  <el-tag type="info" size="mini" v-if="activity.organization_name">
                    {{ activity.organization_name }}
                  </el-tag>
                </div>
                <div class="timeline-body">
                  <p class="activity-description">
                    <i :class="activity.icon"></i>
                    {{ activity.content }}
                  </p>
                </div>
                <div class="timeline-footer">
                  <el-button 
                    type="text" 
                    size="mini" 
                    @click="viewDetail(activity)"
                    v-if="activity.resource_id"
                  >
                    查看详情 <i class="el-icon-arrow-right"></i>
                  </el-button>
                </div>
              </div>
            </el-card>
          </el-timeline-item>
        </el-timeline>
      </div>

      <div v-else class="empty-state">
        <i class="el-icon-document" style="font-size: 64px; color: #C0C4CC;"></i>
        <p style="color: #909399; margin-top: 16px; font-size: 14px;">暂无活动记录</p>
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
import { getRecentActivities } from '@/api/accountV2'
import { getActivityStatistics } from '@/api/activity'

export default {
  name: 'ActivityMonitor',
  data() {
    return {
      loading: false,
      activities: [],
      statistics: {
        totalPrescriptions: 0,
        totalSupplements: 0,
        totalOrders: 0,
        totalClaims: 0
      },
      filterOrganization: '',
      filterActivityType: '',
      dateRange: [],
      currentPage: 1,
      pageSize: 20,
      total: 0
    }
  },
  created() {
    this.loadActivities()
    this.loadStatistics()
  },
  methods: {
    async loadActivities() {
      this.loading = true
      try {
        const response = await getRecentActivities()
        
        if (response && response.code === 200) {
          this.activities = response.data || []
          this.total = this.activities.length
        }
      } catch (error) {
        console.error('加载活动失败:', error)
        this.$message.error('加载活动失败')
      } finally {
        this.loading = false
      }
    },

    async loadStatistics() {
      try {
        const response = await getActivityStatistics()
        if (response && response.code === 200) {
          this.statistics = response.data
        }
      } catch (error) {
        console.error('加载统计数据失败:', error)
      }
    },

    resetFilters() {
      this.filterOrganization = ''
      this.filterActivityType = ''
      this.dateRange = []
      this.currentPage = 1
      this.loadActivities()
    },

    handleSizeChange(val) {
      this.pageSize = val
      this.currentPage = 1
      this.loadActivities()
    },

    handleCurrentChange(val) {
      this.currentPage = val
      this.loadActivities()
    },

    getActivityColor(type) {
      const colorMap = {
        'prescription': '#409EFF',
        'supplement': '#67C23A',
        'access_request': '#E6A23C',
        'drug_order': '#F56C6C',
        'insurance_claim': '#909399',
        '病历创建': '#409EFF',
        '用户注册': '#67C23A',
        '补充记录': '#67C23A',
        '授权请求': '#E6A23C',
        '药品订单': '#F56C6C',
        '保险报销': '#909399'
      }
      return colorMap[type] || '#909399'
    },

    getActivityTagType(type) {
      const typeMap = {
        'prescription': 'primary',
        'supplement': 'success',
        'access_request': 'warning',
        'drug_order': 'danger',
        'insurance_claim': 'info',
        '病历创建': 'primary',
        '用户注册': 'success',
        '补充记录': 'success',
        '授权请求': 'warning',
        '药品订单': 'danger',
        '保险报销': 'info'
      }
      return typeMap[type] || 'info'
    },

    getActivityTypeText(type) {
      const textMap = {
        'prescription': '创建病历',
        'supplement': '补充记录',
        'access_request': '授权请求',
        'drug_order': '药品订单',
        'insurance_claim': '保险报销',
        '病历创建': '病历创建',
        '用户注册': '用户注册',
        '补充记录': '补充记录',
        '授权请求': '授权请求',
        '药品订单': '药品订单',
        '保险报销': '保险报销'
      }
      return textMap[type] || type
    },

    getActivityIcon(type) {
      const iconMap = {
        'prescription': 'el-icon-document-add',
        'supplement': 'el-icon-circle-plus-outline',
        'access_request': 'el-icon-key',
        'drug_order': 'el-icon-shopping-cart-full',
        'insurance_claim': 'el-icon-s-finance',
        '病历创建': 'el-icon-document-add',
        '用户注册': 'el-icon-user',
        '补充记录': 'el-icon-circle-plus-outline',
        '授权请求': 'el-icon-key',
        '药品订单': 'el-icon-shopping-cart-full',
        '保险报销': 'el-icon-s-finance'
      }
      return iconMap[type] || 'el-icon-info'
    },

    viewDetail(activity) {
      // 根据活动类型跳转到相应的详情页
      const routeMap = {
        'prescription': '/prescription/detail',
        'supplement': '/prescription/detail',
        'access_request': '/authorization',
        'drug_order': '/drugOrder/list',
        'insurance_claim': '/insurance/list'
      }
      
      const route = routeMap[activity.type]
      if (route) {
        this.$router.push({
          path: route,
          query: { id: activity.resource_id }
        })
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.activity-monitor-container {
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

  &.stat-card-purple .stat-icon {
    background: linear-gradient(135deg, #909399 0%, #C0C4CC 100%);
  }
}

.activity-card {
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

.timeline-card {
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

.timeline-content {
  .timeline-header {
    display: flex;
    gap: 10px;
    margin-bottom: 10px;
  }

  .timeline-body {
    margin-bottom: 10px;

    .activity-description {
      font-size: 15px;
      color: #303133;
      margin: 0 0 10px 0;
      line-height: 1.6;

      i {
        margin-right: 5px;
        color: #667eea;
      }
    }

    .activity-details {
      display: flex;
      flex-wrap: wrap;
      gap: 15px;
      font-size: 13px;
      color: #606266;

      span {
        display: flex;
        align-items: center;
        gap: 5px;

        i {
          color: #909399;
        }
      }
    }
  }

  .timeline-footer {
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
