<template>
  <div class="health-profile-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>健康档案</h2>
      <p>查看您的就医历史和健康统计</p>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-cards">
      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
            <i class="el-icon-document"></i>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ profile.totalPrescriptions }}</div>
            <div class="stat-label">病历总数</div>
          </div>
        </div>
      </el-col>

      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
            <i class="el-icon-s-data"></i>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ profile.totalVisits }}</div>
            <div class="stat-label">就诊次数</div>
          </div>
        </div>
      </el-col>

      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);">
            <i class="el-icon-user"></i>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ profile.authorizedDoctors }}</div>
            <div class="stat-label">授权医生</div>
          </div>
        </div>
      </el-col>

      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);">
            <i class="el-icon-shopping-bag-2"></i>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ profile.totalOrders }}</div>
            <div class="stat-label">药品订单</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 时间筛选 -->
    <div class="time-filter">
      <el-radio-group v-model="timeRange" size="small" @change="handleTimeRangeChange">
        <el-radio-button label="1m">最近一个月</el-radio-button>
        <el-radio-button label="3m">最近三个月</el-radio-button>
        <el-radio-button label="1y">最近一年</el-radio-button>
        <el-radio-button label="all">全部</el-radio-button>
      </el-radio-group>
    </div>

    <!-- 主要内容区域 -->
    <el-row :gutter="20" class="main-content">
      <!-- 左侧：就诊时间线 -->
      <el-col :xs="24" :sm="24" :md="12" :lg="12">
        <el-card class="timeline-card">
          <div slot="header" class="card-header">
            <span><i class="el-icon-time"></i> 就诊时间线</span>
          </div>
          
          <div v-if="profile.timeline && profile.timeline.length > 0" class="timeline-container">
            <el-timeline>
              <el-timeline-item
                v-for="(item, index) in profile.timeline"
                :key="index"
                :timestamp="item.date"
                placement="top"
                :color="getTimelineColor(index)"
              >
                <el-card class="timeline-item-card">
                  <div class="timeline-item-content">
                    <div class="timeline-hospital">
                      <i class="el-icon-office-building"></i>
                      {{ item.hospital }}
                    </div>
                    <div class="timeline-doctor">
                      <i class="el-icon-user-solid"></i>
                      {{ item.doctor }}
                    </div>
                    <div class="timeline-diagnosis">
                      <i class="el-icon-document-checked"></i>
                      诊断：{{ item.diagnosis }}
                    </div>
                    <el-button
                      type="text"
                      size="small"
                      @click="viewPrescriptionDetail(item.prescId)"
                    >
                      查看详情 <i class="el-icon-arrow-right"></i>
                    </el-button>
                  </div>
                </el-card>
              </el-timeline-item>
            </el-timeline>
          </div>
          
          <div v-else class="empty-state">
            <i class="el-icon-document" style="font-size: 48px; color: #dcdfe6;"></i>
            <p style="color: #909399; margin-top: 10px;">暂无就诊记录</p>
          </div>
        </el-card>

        <!-- 转诊记录 -->
        <el-card class="timeline-card" style="margin-top: 20px;">
          <div slot="header" class="card-header">
            <span><i class="el-icon-connection"></i> 转诊记录</span>
          </div>
          
          <div v-if="supplementRecords && supplementRecords.length > 0" class="supplement-records-list">
            <div 
              v-for="(record, index) in supplementRecords" 
              :key="index"
              class="supplement-record-item"
            >
              <div class="record-header">
                <el-tag 
                  :type="getRecordTypeTag(record.recordType)" 
                  size="small"
                >
                  {{ getRecordTypeText(record.recordType) }}
                </el-tag>
                <span class="record-date">{{ record.date }}</span>
              </div>
              <div class="record-content">
                <div class="record-row">
                  <i class="el-icon-office-building"></i>
                  <span>{{ record.hospital }}</span>
                </div>
                <div class="record-row">
                  <i class="el-icon-user-solid"></i>
                  <span>{{ record.doctor }}</span>
                </div>
                <div class="record-row">
                  <i class="el-icon-document-checked"></i>
                  <span>{{ record.diagnosis }}</span>
                </div>
              </div>
            </div>
          </div>
          
          <div v-else class="empty-state">
            <i class="el-icon-connection" style="font-size: 48px; color: #dcdfe6;"></i>
            <p style="color: #909399; margin-top: 10px;">暂无转诊记录</p>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧：统计图表 -->
      <el-col :xs="24" :sm="24" :md="12" :lg="12">
        <!-- 疾病分布 -->
        <el-card class="chart-card">
          <div slot="header" class="card-header">
            <span><i class="el-icon-pie-chart"></i> 疾病分布</span>
          </div>
          
          <div v-if="profile.diseaseStats && profile.diseaseStats.length > 0" class="chart-container">
            <div ref="diseaseChart" class="chart"></div>
          </div>
          
          <div v-else class="empty-state">
            <i class="el-icon-pie-chart" style="font-size: 48px; color: #dcdfe6;"></i>
            <p style="color: #909399; margin-top: 10px;">暂无疾病统计数据</p>
          </div>
        </el-card>

        <!-- 医院分布 -->
        <el-card class="chart-card" style="margin-top: 20px;">
          <div slot="header" class="card-header">
            <span><i class="el-icon-data-line"></i> 就医医院分布</span>
          </div>
          
          <!-- 医院统计列表 -->
          <div v-if="profile.hospitalStats && profile.hospitalStats.length > 0" class="hospital-stats-list">
            <div 
              v-for="(item, index) in profile.hospitalStats" 
              :key="index"
              class="hospital-stat-item"
            >
              <div class="hospital-info">
                <i class="el-icon-office-building"></i>
                <span class="hospital-name">{{ item.hospital }}</span>
              </div>
              <div class="hospital-count">
                <el-tag type="primary" size="small">{{ item.count }} 次</el-tag>
              </div>
            </div>
          </div>
          
          <!-- 图表（可选） -->
          <div v-if="profile.hospitalStats && profile.hospitalStats.length > 0" class="chart-container" style="margin-top: 20px;">
            <div ref="hospitalChart" class="chart"></div>
          </div>
          
          <div v-else class="empty-state">
            <i class="el-icon-data-line" style="font-size: 48px; color: #dcdfe6;"></i>
            <p style="color: #909399; margin-top: 10px;">暂无医院统计数据</p>
          </div>
        </el-card>

        <!-- 药品订单 -->
        <el-card class="chart-card" style="margin-top: 20px;">
          <div slot="header" class="card-header">
            <span><i class="el-icon-shopping-bag-2"></i> 药品订单</span>
          </div>
          
          <div v-if="drugOrders && drugOrders.length > 0" class="drug-orders-list">
            <div 
              v-for="(order, index) in drugOrders" 
              :key="index"
              class="drug-order-item"
            >
              <div class="order-header">
                <span class="order-date">{{ order.created }}</span>
              </div>
              <div class="order-content">
                <div class="order-row">
                  <i class="el-icon-shop"></i>
                  <span>{{ getDrugstoreName(order.drug_store) }}</span>
                </div>
                <div class="order-row">
                  <i class="el-icon-goods"></i>
                  <span>{{ order.Name }}</span>
                </div>
                <div class="order-row">
                  <i class="el-icon-document-copy"></i>
                  <span>数量：{{ order.amount }} 份</span>
                </div>
              </div>
            </div>
          </div>
          
          <div v-else class="empty-state">
            <i class="el-icon-shopping-bag-2" style="font-size: 48px; color: #dcdfe6;"></i>
            <p style="color: #909399; margin-top: 10px;">暂无药品订单</p>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { getHealthProfile } from '@/api/patient'
import { queryDrugOrderList } from '@/api/drugOrder'
import { queryAccountList } from '@/api/accountV2'
import { getAccountMap } from '@/utils/accountCache'
import { mapGetters } from 'vuex'
import * as echarts from 'echarts'

export default {
  name: 'HealthProfile',
  data() {
    return {
      loading: false,
      timeRange: 'all',
      profile: {
        totalPrescriptions: 0,
        totalVisits: 0,
        authorizedDoctors: 0,
        totalOrders: 0,
        timeline: [],
        diseaseStats: [],
        hospitalStats: []
      },
      drugOrders: [],
      supplementRecords: [], // 补充记录（转诊、复诊、急诊）
      accountMap: {},
      diseaseChart: null,
      hospitalChart: null
    }
  },
  computed: {
    ...mapGetters([
      'account_id'
    ])
  },
  mounted() {
    // 使用缓存的账户列表，避免重复查询
    this.loadAccountListFromCache().then(() => {
      // 并行加载数据，提升性能
      Promise.all([
        this.loadHealthProfile(),
        this.loadDrugOrders()
      ]).catch(err => {
        console.error('加载数据失败:', err)
      })
    })
  },
  beforeDestroy() {
    // 销毁图表实例
    if (this.diseaseChart) {
      this.diseaseChart.dispose()
    }
    if (this.hospitalChart) {
      this.hospitalChart.dispose()
    }
  },
  methods: {
    // 从时间线中提取补充记录
    extractSupplementRecords() {
      if (!this.profile.timeline || this.profile.timeline.length === 0) {
        this.supplementRecords = []
        return
      }
      
      // 筛选出带有类型标识的记录（[转诊]、[复诊]、[急诊]、[授权转诊]）
      this.supplementRecords = this.profile.timeline
        .filter(item => {
          return item.diagnosis && (
            item.diagnosis.includes('[转诊]') ||
            item.diagnosis.includes('[复诊]') ||
            item.diagnosis.includes('[急诊]') ||
            item.diagnosis.includes('[授权转诊]')
          )
        })
        .map(item => {
          // 提取记录类型
          let recordType = 'consultation' // 默认转诊
          if (item.diagnosis.includes('[授权转诊]')) {
            recordType = 'authorization' // 授权转诊
          } else if (item.diagnosis.includes('[转诊]')) {
            recordType = 'consultation'
          } else if (item.diagnosis.includes('[复诊]')) {
            recordType = 'followup'
          } else if (item.diagnosis.includes('[急诊]')) {
            recordType = 'emergency'
          }
          
          // 去掉诊断中的类型标识
          const diagnosis = item.diagnosis
            .replace('[授权转诊]', '')
            .replace('[转诊]', '')
            .replace('[复诊]', '')
            .replace('[急诊]', '')
            .trim()
          
          return {
            date: item.date,
            hospital: item.hospital,
            doctor: item.doctor,
            diagnosis: diagnosis,
            recordType: recordType,
            prescId: item.prescId
          }
        })
      
      console.log('📋 提取到补充记录:', this.supplementRecords.length, '条')
    },

    // 获取记录类型标签颜色
    getRecordTypeTag(recordType) {
      const tagMap = {
        'authorization': 'primary',  // 授权转诊 - 蓝色
        'consultation': 'warning',   // 转诊 - 橙色
        'followup': 'success',       // 复诊 - 绿色
        'emergency': 'danger'        // 急诊 - 红色
      }
      return tagMap[recordType] || 'info'
    },

    // 获取记录类型文本
    getRecordTypeText(recordType) {
      const textMap = {
        'authorization': '授权转诊',
        'consultation': '转诊',
        'followup': '复诊',
        'emergency': '急诊'
      }
      return textMap[recordType] || recordType
    },

    // 从缓存加载账户列表（快速）
    async loadAccountListFromCache() {
      try {
        this.accountMap = await getAccountMap()
        console.log('✅ 从缓存加载账户映射:', Object.keys(this.accountMap).length, '个账户')
      } catch (err) {
        console.error('加载账户映射失败:', err)
        // 缓存失败不影响主要功能
      }
    },

    // 加载健康档案数据（带静默重试机制）
    async loadHealthProfile(retryCount = 0) {
      const maxRetries = 3
      
      // 只在第一次加载时显示loading
      if (retryCount === 0) {
        this.loading = true
      }
      
      try {
        const response = await getHealthProfile(this.timeRange)
        
        console.log('📊 健康档案API响应:', response)
        
        if (response.code === 200) {
          this.profile = response.data
          
          console.log('✅ 健康档案数据已加载:')
          console.log('  - 病历总数:', this.profile.totalPrescriptions)
          console.log('  - 就诊次数:', this.profile.totalVisits)
          console.log('  - 时间线记录数:', this.profile.timeline?.length)
          console.log('  - 疾病统计数:', this.profile.diseaseStats?.length)
          console.log('  - 医院统计数:', this.profile.hospitalStats?.length)
          
          // 从时间线中提取补充记录
          this.extractSupplementRecords()
          
          this.$nextTick(() => {
            console.log('🎨 开始初始化图表...')
            this.initCharts()
          })
          this.loading = false
        } else {
          // 如果失败且还有重试次数，静默重试（不显示错误）
          if (retryCount < maxRetries) {
            console.log(`获取健康档案失败，正在静默重试 (${retryCount + 1}/${maxRetries})...`)
            await new Promise(resolve => setTimeout(resolve, 500 * (retryCount + 1)))
            return this.loadHealthProfile(retryCount + 1)
          } else {
            // 所有重试都失败后才显示错误
            this.$message.error(response.message || '获取健康档案失败')
            this.loading = false
          }
        }
      } catch (error) {
        // 只在控制台记录错误，不显示给用户
        if (retryCount === 0) {
          console.error('获取健康档案失败:', error)
        }
        
        // 如果失败且还有重试次数，静默重试（不显示错误）
        if (retryCount < maxRetries) {
          console.log(`获取健康档案失败，正在静默重试 (${retryCount + 1}/${maxRetries})...`)
          await new Promise(resolve => setTimeout(resolve, 500 * (retryCount + 1)))
          return this.loadHealthProfile(retryCount + 1)
        } else {
          // 所有重试都失败后才显示错误
          this.$message.error('获取健康档案失败，请稍后重试')
          this.loading = false
        }
      }
    },

    // 加载药品订单（带静默重试机制）
    async loadDrugOrders(retryCount = 0) {
      const maxRetries = 3
      
      try {
        const response = await queryDrugOrderList({ patient: this.account_id })
        let orders = []
        
        if (response && response.code === 200 && response.data) {
          orders = Array.isArray(response.data) ? response.data : []
        } else if (Array.isArray(response)) {
          orders = response
        } else if (typeof response === 'object') {
          orders = Object.values(response)
        }
        
        // 根据时间范围筛选
        orders = this.filterOrdersByTimeRange(orders)
        
        // 按时间倒序排序
        this.drugOrders = orders.sort((a, b) => {
          return b.created.localeCompare(a.created)
        })
        
        console.log('健康档案-药品订单数量:', this.drugOrders.length)
        
        // 打印订单详情，帮助调试
        if (this.drugOrders.length > 0) {
          console.log('药品订单详情:', this.drugOrders)
          console.log('第一个订单的药店ID:', this.drugOrders[0].drug_store)
          console.log('账户映射表中的药店名:', this.accountMap[this.drugOrders[0].drug_store])
        }
      } catch (err) {
        // 只在控制台记录错误，不显示给用户
        if (retryCount === 0) {
          console.error('加载药品订单失败:', err)
        }
        
        // 如果失败且还有重试次数，静默重试（不显示错误）
        if (retryCount < maxRetries) {
          console.log(`加载药品订单失败，正在静默重试 (${retryCount + 1}/${maxRetries})...`)
          await new Promise(resolve => setTimeout(resolve, 500 * (retryCount + 1)))
          return this.loadDrugOrders(retryCount + 1)
        }
        // 药品订单失败不影响主要功能，所以即使最后失败也不显示错误提示
      }
    },

    // 根据时间范围筛选订单
    filterOrdersByTimeRange(orders) {
      if (this.timeRange === 'all') {
        return orders
      }
      
      const now = new Date()
      let cutoffDate = new Date()
      
      switch (this.timeRange) {
        case '1m':
          cutoffDate.setMonth(now.getMonth() - 1)
          break
        case '3m':
          cutoffDate.setMonth(now.getMonth() - 3)
          break
        case '1y':
          cutoffDate.setFullYear(now.getFullYear() - 1)
          break
        default:
          return orders
      }
      
      return orders.filter(order => {
        if (!order.created) return false
        
        // 解析订单创建时间 (格式: "2026-02-20 12:27:46")
        const orderDate = new Date(order.created.replace(' ', 'T'))
        return orderDate >= cutoffDate
      })
    },

    // 时间范围变化
    handleTimeRangeChange() {
      this.loadHealthProfile()
      this.loadDrugOrders()
    },

    // 根据药店ID获取药店名字
    getDrugstoreName(drugstoreId) {
      if (!drugstoreId) {
        return '未知药店'
      }
      
      const account = this.accountMap[drugstoreId]
      
      if (account) {
        // accountMap 中的值是对象，包含 name, username, role, organization_name 等
        const accountName = account.name || account.account_name || account.username
        
        if (accountName) {
          // 如果账户名包含"药店"，去掉前缀和后缀
          if (/药店/.test(accountName)) {
            return accountName.replace(/^药店-/, '').replace(/-药店$/, '')
          } else {
            // 如果不包含"药店"，直接返回账户名
            return accountName
          }
        }
      }
      
      // 如果找不到账户名，返回"未知药店"
      console.warn('⚠️ 找不到药店账户，ID:', drugstoreId, '账户映射表:', this.accountMap)
      return '未知药店'
    },

    // 获取药品订单颜色
    getDrugOrderColor(index) {
      const colors = ['#43e97b', '#38f9d7', '#4facfe', '#00f2fe', '#f093fb']
      return colors[index % colors.length]
    },

    // 初始化图表
    initCharts() {
      console.log('🎨 initCharts 被调用')
      console.log('  - diseaseStats:', this.profile.diseaseStats)
      console.log('  - hospitalStats:', this.profile.hospitalStats)
      this.initDiseaseChart()
      this.initHospitalChart()
    },

    // 初始化疾病分布图表
    initDiseaseChart() {
      console.log('📊 initDiseaseChart 被调用')
      console.log('  - diseaseStats长度:', this.profile.diseaseStats?.length)
      
      if (!this.profile.diseaseStats || this.profile.diseaseStats.length === 0) {
        console.log('⚠️ 疾病统计数据为空，跳过图表初始化')
        return
      }

      const chartDom = this.$refs.diseaseChart
      console.log('  - chartDom:', chartDom)
      
      if (!chartDom) {
        console.log('❌ 找不到疾病图表DOM元素')
        return
      }

      if (this.diseaseChart) {
        this.diseaseChart.dispose()
      }

      this.diseaseChart = echarts.init(chartDom)
      console.log('  - 疾病图表实例已创建:', this.diseaseChart)

      const option = {
        tooltip: {
          trigger: 'item',
          formatter: '{b}: {c} 次 ({d}%)'
        },
        legend: {
          orient: 'vertical',
          left: 'left',
          top: 'center'
        },
        series: [
          {
            name: '疾病分布',
            type: 'pie',
            radius: ['40%', '70%'],
            avoidLabelOverlap: false,
            itemStyle: {
              borderRadius: 10,
              borderColor: '#fff',
              borderWidth: 2
            },
            label: {
              show: false,
              position: 'center'
            },
            emphasis: {
              label: {
                show: true,
                fontSize: '18',
                fontWeight: 'bold'
              }
            },
            labelLine: {
              show: false
            },
            data: this.profile.diseaseStats.map(item => ({
              name: item.disease,
              value: item.count
            }))
          }
        ],
        color: ['#667eea', '#764ba2', '#f093fb', '#f5576c', '#4facfe', '#00f2fe', '#43e97b', '#38f9d7']
      }

      this.diseaseChart.setOption(option)
      
      console.log('✅ 疾病分布图表已初始化')

      // 响应式
      window.addEventListener('resize', () => {
        if (this.diseaseChart) {
          this.diseaseChart.resize()
        }
      })
    },

    // 初始化医院分布图表
    initHospitalChart() {
      console.log('🏥 initHospitalChart 被调用')
      console.log('  - hospitalStats长度:', this.profile.hospitalStats?.length)
      
      if (!this.profile.hospitalStats || this.profile.hospitalStats.length === 0) {
        console.log('⚠️ 医院统计数据为空，跳过图表初始化')
        return
      }

      const chartDom = this.$refs.hospitalChart
      console.log('  - chartDom:', chartDom)
      
      if (!chartDom) {
        console.log('❌ 找不到医院图表DOM元素')
        return
      }

      if (this.hospitalChart) {
        this.hospitalChart.dispose()
      }

      this.hospitalChart = echarts.init(chartDom)
      console.log('  - 医院图表实例已创建:', this.hospitalChart)

      const option = {
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'shadow'
          }
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          data: this.profile.hospitalStats.map(item => item.hospital),
          axisTick: {
            alignWithLabel: true
          },
          axisLabel: {
            interval: 0,
            rotate: 30
          }
        },
        yAxis: {
          type: 'value',
          name: '就诊次数',
          minInterval: 1
        },
        series: [
          {
            name: '就诊次数',
            type: 'bar',
            barWidth: '60%',
            data: this.profile.hospitalStats.map(item => item.count),
            itemStyle: {
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: '#667eea' },
                { offset: 1, color: '#764ba2' }
              ]),
              borderRadius: [10, 10, 0, 0]
            },
            label: {
              show: true,
              position: 'top',
              formatter: '{c} 次'
            }
          }
        ]
      }

      this.hospitalChart.setOption(option)
      
      console.log('✅ 医院分布图表已初始化')

      // 响应式
      window.addEventListener('resize', () => {
        if (this.hospitalChart) {
          this.hospitalChart.resize()
        }
      })
    },

    // 获取时间线颜色
    getTimelineColor(index) {
      const colors = ['#667eea', '#f093fb', '#4facfe', '#43e97b', '#f5576c']
      return colors[index % colors.length]
    },

    // 查看病历详情
    viewPrescriptionDetail(prescId) {
      this.$router.push({
        path: '/prescription/detail',
        query: { id: prescId }
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.health-profile-container {
  padding: 20px;
  background: #f5f7fa;
  min-height: calc(100vh - 84px);

  .page-header {
    margin-bottom: 30px;
    
    h2 {
      font-size: 28px;
      font-weight: 600;
      color: #303133;
      margin: 0 0 10px 0;
    }
    
    p {
      font-size: 14px;
      color: #909399;
      margin: 0;
    }
  }

  // 统计卡片
  .stats-cards {
    margin-bottom: 30px;

    .stat-card {
      background: white;
      border-radius: 12px;
      padding: 20px;
      display: flex;
      align-items: center;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
      transition: all 0.3s;
      cursor: pointer;

      &:hover {
        transform: translateY(-5px);
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
      }

      .stat-icon {
        width: 60px;
        height: 60px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 15px;

        i {
          font-size: 28px;
          color: white;
        }
      }

      .stat-content {
        flex: 1;

        .stat-value {
          font-size: 28px;
          font-weight: 700;
          color: #303133;
          line-height: 1;
          margin-bottom: 8px;
        }

        .stat-label {
          font-size: 14px;
          color: #909399;
        }
      }
    }
  }

  // 时间筛选
  .time-filter {
    margin-bottom: 20px;
    text-align: center;
  }

  // 主要内容
  .main-content {
    .timeline-card,
    .chart-card {
      margin-bottom: 20px;

      .card-header {
        font-size: 16px;
        font-weight: 600;
        color: #303133;

        i {
          margin-right: 8px;
          color: #667eea;
        }
      }
    }

    // 时间线
    .timeline-container {
      max-height: 600px;
      overflow-y: auto;
      padding-right: 10px;

      &::-webkit-scrollbar {
        width: 6px;
      }

      &::-webkit-scrollbar-thumb {
        background: #dcdfe6;
        border-radius: 3px;
      }

      .timeline-item-card {
        margin-bottom: 10px;

        .timeline-item-content {
          .timeline-hospital,
          .timeline-doctor,
          .timeline-diagnosis,
          .timeline-drugstore,
          .timeline-drug,
          .timeline-amount {
            margin-bottom: 8px;
            font-size: 14px;
            color: #606266;

            i {
              margin-right: 5px;
              color: #909399;
            }
          }

          .timeline-hospital,
          .timeline-drugstore {
            font-weight: 600;
            color: #303133;
          }
        }
      }

      .drug-order-card {
        background: linear-gradient(135deg, #f6f9fc 0%, #ffffff 100%);
        border-left: 3px solid #43e97b;
      }
    }

    // 图表
    .chart-container {
      .chart {
        width: 100%;
        height: 300px;
      }
    }

    // 医院统计列表
    .hospital-stats-list {
      .hospital-stat-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 12px 15px;
        margin-bottom: 10px;
        background: #f9fafc;
        border-radius: 6px;
        border-left: 3px solid #667eea;
        transition: all 0.3s;

        &:hover {
          background: #f0f2f5;
          transform: translateX(5px);
        }

        .hospital-info {
          display: flex;
          align-items: center;
          flex: 1;

          i {
            font-size: 18px;
            color: #667eea;
            margin-right: 10px;
          }

          .hospital-name {
            font-size: 14px;
            font-weight: 600;
            color: #303133;
          }
        }

        .hospital-count {
          .el-tag {
            font-weight: 600;
          }
        }
      }
    }

    // 补充记录列表
    .supplement-records-list {
      .supplement-record-item {
        padding: 15px;
        margin-bottom: 15px;
        background: linear-gradient(135deg, #fff5f5 0%, #ffffff 100%);
        border-radius: 8px;
        border-left: 4px solid #f56c6c;
        transition: all 0.3s;

        &:hover {
          box-shadow: 0 2px 12px rgba(245, 108, 108, 0.2);
          transform: translateY(-2px);
        }

        .record-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 12px;
          padding-bottom: 10px;
          border-bottom: 1px dashed #e4e7ed;

          .el-tag {
            font-weight: 600;
          }

          .record-date {
            font-size: 13px;
            color: #909399;
          }
        }

        .record-content {
          .record-row {
            display: flex;
            align-items: center;
            margin-bottom: 8px;
            font-size: 14px;
            color: #606266;

            i {
              font-size: 16px;
              color: #f56c6c;
              margin-right: 8px;
              width: 20px;
            }

            &:last-child {
              margin-bottom: 0;
            }
          }
        }
      }
    }

    // 药品订单列表
    .drug-orders-list {
      .drug-order-item {
        padding: 15px;
        margin-bottom: 15px;
        background: linear-gradient(135deg, #f0fff4 0%, #ffffff 100%);
        border-radius: 8px;
        border-left: 4px solid #43e97b;
        transition: all 0.3s;

        &:hover {
          box-shadow: 0 2px 12px rgba(67, 233, 123, 0.2);
          transform: translateY(-2px);
        }

        .order-header {
          margin-bottom: 12px;
          padding-bottom: 10px;
          border-bottom: 1px dashed #e4e7ed;

          .order-date {
            font-size: 13px;
            color: #909399;
          }
        }

        .order-content {
          .order-row {
            display: flex;
            align-items: center;
            margin-bottom: 8px;
            font-size: 14px;
            color: #606266;

            i {
              font-size: 16px;
              color: #43e97b;
              margin-right: 8px;
              width: 20px;
            }

            &:last-child {
              margin-bottom: 0;
            }
          }
        }
      }
    }

    // 空状态
    .empty-state {
      text-align: center;
      padding: 60px 20px;
    }
  }
}

// 响应式
@media (max-width: 768px) {
  .health-profile-container {
    padding: 10px;

    .page-header {
      h2 {
        font-size: 22px;
      }
    }

    .stats-cards {
      .stat-card {
        padding: 15px;

        .stat-icon {
          width: 50px;
          height: 50px;

          i {
            font-size: 24px;
          }
        }

        .stat-content {
          .stat-value {
            font-size: 24px;
          }
        }
      }
    }

    .main-content {
      .chart-container {
        .chart {
          height: 250px;
        }
      }
    }
  }
}
</style>
