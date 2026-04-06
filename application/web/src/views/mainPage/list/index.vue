<template>
  <div class="dashboard-container">
    <!-- 欢迎信息 -->
    <div class="welcome-section">
      <el-card class="welcome-card">
        <div class="welcome-content">
          <div class="welcome-text">
            <div class="welcome-quote-top">
              {{ roleQuote }}
            </div>
            <div class="welcome-header">
              <h2>欢迎回来，{{ displayAccountName }}！</h2>
              <div class="status-info">
                <span class="online-status">
                  <span class="status-dot-pulse"></span>
                  在线
                </span>
                <span class="security-info">
                  <i class="el-icon-time"></i>
                  {{ lastLoginTime }}
                </span>
                <span class="security-info">
                  <i class="el-icon-monitor"></i>
                  {{ loginDevice }}
                </span>
                <span class="security-info">
                  <i class="el-icon-location-information"></i>
                  {{ loginIP }}
                </span>
              </div>
            </div>
            <p class="role-tag">
              <i :class="roleIcon"></i>
              {{ roleText }}
            </p>
            
            <!-- 用户详细信息 -->
            <div class="user-details">
              <div class="detail-item" v-if="account_id">
                <i class="el-icon-link"></i>
                <span class="detail-label">区块链ID:</span>
                <span class="detail-value">{{ displayAccountId }}</span>
              </div>
              <div class="detail-item" v-if="userDetails && userDetails.age">
                <i class="el-icon-user"></i>
                <span class="detail-label">年龄:</span>
                <span class="detail-value">{{ userDetails.age }} 岁</span>
              </div>
              <div class="detail-item" v-if="userDetails && userDetails.gender">
                <i class="el-icon-male"></i>
                <span class="detail-label">性别:</span>
                <span class="detail-value">{{ userDetails.gender }}</span>
              </div>
              <div class="detail-item" v-if="userDetails && userDetails.organization_name && roles[0] !== 'patient'">
                <i class="el-icon-office-building"></i>
                <span class="detail-label">所属组织:</span>
                <span class="detail-value">{{ userDetails.organization_name }}</span>
              </div>
              <div class="detail-item" v-if="userDetails && userDetails.department && roles[0] !== 'patient'">
                <i class="el-icon-folder-opened"></i>
                <span class="detail-label">科室:</span>
                <span class="detail-value">{{ userDetails.department }}</span>
              </div>
              <div class="detail-item" v-if="userDetails && userDetails.doctor_title && roles[0] !== 'patient'">
                <i class="el-icon-medal"></i>
                <span class="detail-label">职称:</span>
                <span class="detail-value">{{ userDetails.doctor_title }}</span>
              </div>
            </div>
          </div>
          <div class="welcome-avatar">
            <img :src="avatarUrl" alt="avatar">
          </div>
        </div>
      </el-card>
    </div>

    <!-- 快速统计 -->
    <div class="stats-section" v-if="showStats">
      <h3 class="section-title">
        <i class="el-icon-data-line"></i>
        数据概览
      </h3>
      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :md="6" v-for="(stat, index) in enhancedStats" :key="stat.title">
          <el-card class="stat-card" shadow="hover" :style="{ '--card-gradient': stat.gradient }">
            <div class="stat-content">
              <div class="stat-left">
                <div class="stat-icon">
                  <i :class="stat.icon"></i>
                </div>
                <div class="stat-info">
                  <p class="stat-value">{{ stat.value }}</p>
                  <p class="stat-title">{{ stat.title }}</p>
                  <p class="stat-subtitle" v-if="stat.subtitle">{{ stat.subtitle }}</p>
                </div>
              </div>
              <div class="stat-chart">
                <div :ref="'chart' + index" class="mini-chart"></div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 最近活动时间线 -->
    <div class="activity-section" v-if="showActivity">
      <div class="section-header">
        <h3 class="section-title">
          <i class="el-icon-time"></i>
          最近活动
        </h3>
        <el-button 
          type="text" 
          size="small" 
          @click="toggleActivityExpand"
          class="expand-btn"
        >
          {{ activityExpanded ? '收起' : '查看更多' }}
          <i :class="activityExpanded ? 'el-icon-arrow-up' : 'el-icon-arrow-down'"></i>
        </el-button>
      </div>
      <el-card class="activity-card">
        <el-timeline>
          <el-timeline-item
            v-for="(activity, index) in displayedActivities"
            :key="index"
            :color="activity.color"
            placement="top"
          >
            <div class="activity-item">
              <div class="activity-content">
                <i :class="activity.icon" :style="{ color: activity.color }"></i>
                <span class="activity-type">{{ activity.type }}</span>
                <span class="activity-text">{{ activity.content }}</span>
              </div>
              <div class="activity-time">{{ activity.timestamp }}</div>
            </div>
          </el-timeline-item>
          <el-timeline-item v-if="recentActivities.length === 0" color="#909399">
            <p style="color: #909399;">暂无最近活动</p>
          </el-timeline-item>
        </el-timeline>
      </el-card>
    </div>

    <!-- 功能模块 -->
    <div class="modules-section">
      <h3 class="section-title">
        <i class="el-icon-menu"></i>
        快捷功能
      </h3>
      <el-row :gutter="20">
        <el-col 
          v-for="module in availableModules" 
          :key="module.path"
          :xs="24" 
          :sm="12" 
          :md="8" 
          :lg="8"
          :xl="6"
        >
          <el-card 
            class="module-card" 
            :body-style="{ padding: '0px' }"
            shadow="hover"
            @click.native="navigateTo(module.path)"
          >
            <div class="module-content">
              <div class="module-icon" :style="{ background: module.color }">
                <i :class="module.icon"></i>
              </div>
              <div class="module-info">
                <h4>{{ module.title }}</h4>
                <p>{{ module.description }}</p>
              </div>
              <div class="module-arrow">
                <i class="el-icon-arrow-right"></i>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { queryAccountList, getStatistics, getRecentActivities } from '@/api/accountV2'
import * as echarts from 'echarts'

export default {
  name: 'Dashboard',
  computed: {
    ...mapGetters([
      'account_id',
      'roles',
      'account_name'
    ]),
    
    displayAccountName() {
      if (this.roles[0] === 'admin' && this.account_name) {
        return this.account_name.replace(/^医生-/, '管理员-')
      }
      return this.account_name
    },
    
    displayAccountId() {
      return this.account_id
    },
    
    roleText() {
      const roleMap = {
        'admin': '系统管理员',
        'doctor': '医生',
        'patient': '病人',
        'drugstore': '药店'
      }
      return roleMap[this.roles[0]] || '用户'
    },

    roleIcon() {
      const iconMap = {
        'admin': 'el-icon-s-custom',
        'doctor': 'el-icon-user',
        'patient': 'el-icon-user',
        'drugstore': 'el-icon-shopping-bag-2'
      }
      return iconMap[this.roles[0]] || 'el-icon-user'
    },

    avatarUrl() {
      const avatarMap = {
        'admin': '/image/系统管理员.png',
        'doctor': '/image/男医生.png',
        'patient': '/image/病人-男.png',
        'drugstore': '/image/药店-copy.png'
      }
      return avatarMap[this.roles[0]] || '/image/病人-男.png'
    },

    availableModules() {
      const role = this.roles[0]
      
      const allModules = {
        admin: [
          { path: '/account/architecture', title: '架构概览', description: '查看联盟链架构图', icon: 'el-icon-connection', color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)' },
          { path: '/account/statistics', title: '数据统计', description: '查看系统数据统计', icon: 'el-icon-data-line', color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' },
          { path: '/account/activity-monitor', title: '日志监控', description: '监控联盟中所有组织的业务日志', icon: 'el-icon-document-copy', color: 'linear-gradient(135deg, #30cfd0 0%, #330867 100%)' },
          { path: '/account/all', title: '用户管理', description: '查看、编辑、停用和恢复用户', icon: 'el-icon-user-solid', color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)' },
          { path: '/account/add', title: '新增账户', description: '创建新的用户账户', icon: 'el-icon-plus', color: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)' },
          { path: '/account/delete', title: '删除用户', description: '停用或恢复用户账户', icon: 'el-icon-delete', color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)' },
          { path: '/prescription/all', title: '病历管理', description: '查看和管理所有病历', icon: 'el-icon-document', color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)' },
          { path: '/drug/all', title: '药品订单', description: '查看所有药品订单', icon: 'el-icon-shopping-cart-2', color: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)' },
          { path: '/about', title: '关于系统', description: '查看系统介绍', icon: 'el-icon-info', color: 'linear-gradient(135deg, #30cfd0 0%, #330867 100%)' }
        ],
        doctor: [
          { path: '/prescription/all', title: '所有病历', description: '查看所有患者病历', icon: 'el-icon-document', color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' },
          { path: '/prescription/patient-search', title: '患者病历查询', description: '搜索患者病历记录', icon: 'el-icon-search', color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)' },
          { path: '/prescription/add', title: '新增病历', description: '为患者创建病历', icon: 'el-icon-edit', color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)' },
          { path: '/authorization', title: '授权管理', description: '查看授权申请记录', icon: 'el-icon-key', color: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)' },
          { path: '/about', title: '关于系统', description: '查看系统介绍', icon: 'el-icon-info', color: 'linear-gradient(135deg, #30cfd0 0%, #330867 100%)' }
        ],
        patient: [
          { path: '/prescription/health-profile', title: '健康档案', description: '查看健康统计和就医历史', icon: 'el-icon-data-analysis', color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' },
          { path: '/prescription/mine', title: '我的病历', description: '查看个人病历记录', icon: 'el-icon-document', color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)' },
          { path: '/prescription/access-trace', title: '隐私溯源', description: '查看病历访问记录', icon: 'el-icon-view', color: 'linear-gradient(135deg, #30cfd0 0%, #330867 100%)' },
          { path: '/authorization', title: '授权管理', description: '管理病历授权申请', icon: 'el-icon-key', color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)' },
          { path: '/drug/mine', title: '我的订单', description: '查看药品订单', icon: 'el-icon-shopping-cart-2', color: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)' },
          { path: '/about', title: '关于系统', description: '查看系统介绍', icon: 'el-icon-info', color: 'linear-gradient(135deg, #30cfd0 0%, #330867 100%)' }
        ],
        drugstore: [
          { path: '/drug/all', title: '所有订单', description: '查看所有药品订单', icon: 'el-icon-shopping-cart-2', color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)' },
          { path: '/addDrug', title: '新增订单', description: '创建新的药品订单', icon: 'el-icon-plus', color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)' },
          { path: '/about', title: '关于系统', description: '查看系统介绍', icon: 'el-icon-info', color: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)' }
        ]
      }
      return allModules[role] || []
    },

    showStats() {
      return this.roles[0] === 'admin'
    },

    showActivity() {
      return this.roles[0] === 'admin'
    },

    displayedActivities() {
      // 如果展开，显示所有活动；否则只显示前5条
      return this.activityExpanded 
        ? this.recentActivities 
        : this.recentActivities.slice(0, 5)
    },

    roleQuote() {
      const quoteMap = {
        'admin': '春风若有怜花意，我当悬壶济世人',
        'doctor': '春风若有怜花意，我当悬壶济世人',
        'patient': '春风十里不如你，愿君康健乐无忧',
        'drugstore': '春风若有怜花意，良药济世暖人心'
      }
      return quoteMap[this.roles[0]] || '春风若有怜花意，我当悬壶济世人'
    },

    enhancedStats() {
      return [
        { 
          title: '总用户数', 
          value: this.statsData.userCount || '0', 
          subtitle: `今日新增 ${this.statsData.todayNewUsers || 0}`,
          icon: 'el-icon-user-solid', 
          color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
          gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
          chartData: this.statsData.userTrendData && this.statsData.userTrendData.length > 0 
            ? this.statsData.userTrendData 
            : [5, 8, 9, 10, 11, 12],
          chartColor: ['#667eea', '#764ba2']
        },
        { 
          title: '病历记录', 
          value: this.statsData.prescriptionCount || '0', 
          subtitle: `今日新增 ${this.statsData.todayNewPrescriptions || 0}`,
          icon: 'el-icon-document', 
          color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
          gradient: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
          chartData: this.statsData.prescriptionTrendData && this.statsData.prescriptionTrendData.length > 0
            ? this.statsData.prescriptionTrendData 
            : [1, 1, 2, 2, 3, 3],
          chartColor: ['#f093fb', '#f5576c']
        },
        { 
          title: '药品订单', 
          value: this.statsData.drugOrderCount || '0', 
          subtitle: `今日新增 ${this.statsData.todayNewOrders || 0}`,
          icon: 'el-icon-shopping-cart-2', 
          color: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
          gradient: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
          chartData: this.statsData.orderTrendData && this.statsData.orderTrendData.length > 0
            ? this.statsData.orderTrendData 
            : [0, 0, 0, 0, 0, 0],
          chartColor: ['#fa709a', '#fee140']
        },
        { 
          title: '医疗组织', 
          value: this.statsData.organizationCount || '4', 
          subtitle: '联盟链节点',
          icon: 'el-icon-office-building', 
          color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
          gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
          chartData: [3, 3, 4, 4, 4, 4],
          chartColor: ['#4facfe', '#00f2fe']
        }
      ]
    }
  },

  data() {
    return {
      statsData: {
        userCount: 0,
        organizationCount: 4,
        doctorCount: 0,
        patientCount: 0,
        prescriptionCount: 0,
        drugOrderCount: 0,
        todayNewUsers: 0,
        todayNewPrescriptions: 0,
        todayNewOrders: 0,
        userTrend: 0,
        prescriptionTrend: 0,
        orderTrend: 0,
        userTrendData: [],
        prescriptionTrendData: [],
        orderTrendData: []
      },
      charts: [],
      userDetails: {
        organization_name: '',
        department: '',
        doctor_title: '',
        age: 0,
        gender: ''
      },
      recentActivities: [],
      lastLoginTime: '',
      loginDevice: '',
      loginIP: '',
      activityExpanded: false
    }
  },

  mounted() {
    console.log('=== 主页加载 ===')
    console.log('account_id:', this.account_id)
    console.log('account_name:', this.account_name)
    console.log('roles:', this.roles)
    
    // 初始化登录信息
    this.lastLoginTime = this.getLastLoginTime()
    this.loginDevice = this.getLoginDevice()
    this.loginIP = this.getLoginIP()
    
    try {
      const userInfo = sessionStorage.getItem('userInfo')
      if (userInfo) {
        const user = JSON.parse(userInfo)
        this.userDetails = {
          organization_name: user.organization_name || '',
          department: user.department || '',
          doctor_title: user.doctor_title || '',
          age: user.age || 0,
          gender: user.gender || ''
        }
        console.log('✅ 从 sessionStorage 快速加载用户信息:', this.userDetails)
      }
    } catch (e) {
      console.error('❌ 从 sessionStorage 读取失败:', e)
    }
    
    if (this.roles[0] === 'admin') {
      this.loadStats()
      this.loadRecentActivities()
    }
    
    this.loadUserDetails()
  },

  updated() {
    // 当数据更新后，重新初始化图表
    if (this.roles[0] === 'admin' && this.statsData.userCount > 0 && this.charts.length === 0) {
      this.$nextTick(() => {
        this.initCharts()
      })
    }
  },

  beforeDestroy() {
    // 销毁图表实例
    this.charts.forEach(chart => {
      if (chart) {
        chart.dispose()
      }
    })
  },

  methods: {
    navigateTo(path) {
      this.$router.push(path)
    },

    getLastLoginTime() {
      // 尝试从 sessionStorage 获取上次登录时间
      const lastLogin = sessionStorage.getItem('lastLoginTime')
      if (lastLogin) {
        return lastLogin
      }
      
      // 如果没有，使用当前时间并保存
      const now = new Date()
      const timeStr = now.getFullYear() + '-' + 
                      String(now.getMonth() + 1).padStart(2, '0') + '-' + 
                      String(now.getDate()).padStart(2, '0') + ' ' + 
                      String(now.getHours()).padStart(2, '0') + ':' + 
                      String(now.getMinutes()).padStart(2, '0')
      sessionStorage.setItem('lastLoginTime', timeStr)
      return timeStr
    },

    getLoginDevice() {
      // 尝试从 sessionStorage 获取登录设备
      const device = sessionStorage.getItem('loginDevice')
      if (device) {
        return device
      }
      
      // 检测设备类型
      const ua = navigator.userAgent
      let deviceType = 'Unknown'
      
      if (/Windows/i.test(ua)) {
        deviceType = 'Windows PC'
      } else if (/Macintosh|MacIntel|MacPPC|Mac68K/i.test(ua)) {
        deviceType = 'Mac'
      } else if (/Linux/i.test(ua)) {
        deviceType = 'Linux PC'
      } else if (/Android/i.test(ua)) {
        deviceType = 'Android'
      } else if (/iPhone|iPad|iPod/i.test(ua)) {
        deviceType = 'iOS'
      }
      
      // 检测浏览器
      let browser = ''
      if (/Edge/i.test(ua)) {
        browser = 'Edge'
      } else if (/Chrome/i.test(ua)) {
        browser = 'Chrome'
      } else if (/Firefox/i.test(ua)) {
        browser = 'Firefox'
      } else if (/Safari/i.test(ua)) {
        browser = 'Safari'
      }
      
      const deviceInfo = browser ? `${deviceType} · ${browser}` : deviceType
      sessionStorage.setItem('loginDevice', deviceInfo)
      return deviceInfo
    },

    getLoginIP() {
      // 尝试从 sessionStorage 获取IP
      const ip = sessionStorage.getItem('loginIP')
      if (ip) {
        return ip
      }
      
      // 默认显示本地IP（实际项目中应该从后端获取）
      const defaultIP = '127.0.0.1'
      sessionStorage.setItem('loginIP', defaultIP)
      
      // 异步获取真实IP（可选）
      this.fetchRealIP()
      
      return defaultIP
    },

    async fetchRealIP() {
      try {
        // 这里可以调用后端API获取真实IP
        // 或者使用第三方服务（注意隐私和安全）
        // const response = await fetch('https://api.ipify.org?format=json')
        // const data = await response.json()
        // this.loginIP = data.ip
        // sessionStorage.setItem('loginIP', data.ip)
      } catch (e) {
        console.log('获取IP失败，使用默认值')
      }
    },

    async loadStats() {
      try {
        // 调用新的统计 API
        const statsRes = await getStatistics()
        if (statsRes.code === 200 && statsRes.data) {
          const data = statsRes.data
          this.statsData.userCount = data.userCount || 0
          this.statsData.prescriptionCount = data.prescriptionCount || 0
          this.statsData.drugOrderCount = data.drugOrderCount || 0
          this.statsData.organizationCount = data.organizationCount || 4
          this.statsData.todayNewUsers = data.todayNewUsers || 0
          this.statsData.todayNewPrescriptions = data.todayNewPrescriptions || 0
          this.statsData.todayNewOrders = data.todayNewOrders || 0
          this.statsData.userTrendData = data.userTrendData || []
          this.statsData.prescriptionTrendData = data.prescriptionTrendData || []
          this.statsData.orderTrendData = data.orderTrendData || []
          
          console.log('✅ 统计数据加载成功:', this.statsData)
        }
      } catch (error) {
        console.error('❌ 加载统计数据失败:', error)
        // 降级方案：使用区块链账户列表
        try {
          const userRes = await queryAccountList()
          const userList = userRes.code === 200 ? userRes.data : userRes
          this.statsData.userCount = Array.isArray(userList) ? userList.length : 0
          this.statsData.organizationCount = 4
        } catch (e) {
          console.error('❌ 降级方案也失败:', e)
        }
      }
      
      // 等待 DOM 更新后初始化图表
      this.$nextTick(() => {
        this.initCharts()
      })
    },

    initCharts() {
      // 清除旧图表
      this.charts.forEach(chart => {
        if (chart) {
          try {
            chart.dispose()
          } catch (e) {
            console.error('销毁图表失败:', e)
          }
        }
      })
      this.charts = []

      // 延迟确保 DOM 完全渲染
      setTimeout(() => {
        // 为每个统计卡片创建图表
        for (let index = 0; index < this.enhancedStats.length; index++) {
          const stat = this.enhancedStats[index]
          const refName = 'chart' + index
          const chartRefs = this.$refs[refName]
          
          let chartDom = null
          if (Array.isArray(chartRefs) && chartRefs.length > 0) {
            chartDom = chartRefs[0]
          } else if (chartRefs) {
            chartDom = chartRefs
          }
          
          if (chartDom) {
            try {
              // 先销毁可能存在的实例
              const existingChart = echarts.getInstanceByDom(chartDom)
              if (existingChart) {
                existingChart.dispose()
              }
              
              const chart = echarts.init(chartDom)
              
              const option = {
                grid: {
                  left: 5,
                  right: 5,
                  top: 10,
                  bottom: 10
                },
                xAxis: {
                  type: 'category',
                  show: false,
                  boundaryGap: false,
                  data: ['1', '2', '3', '4', '5', '6']
                },
                yAxis: {
                  type: 'value',
                  show: false
                },
                series: [{
                  data: stat.chartData,
                  type: 'line',
                  smooth: true,
                  symbol: 'none',
                  lineStyle: {
                    width: 3,
                    color: {
                      type: 'linear',
                      x: 0,
                      y: 0,
                      x2: 1,
                      y2: 0,
                      colorStops: [
                        { offset: 0, color: stat.chartColor[0] },
                        { offset: 1, color: stat.chartColor[1] }
                      ]
                    },
                    shadowColor: stat.chartColor[0] + '40',
                    shadowBlur: 10,
                    shadowOffsetY: 3
                  },
                  areaStyle: {
                    color: {
                      type: 'linear',
                      x: 0,
                      y: 0,
                      x2: 0,
                      y2: 1,
                      colorStops: [
                        { offset: 0, color: stat.chartColor[0] + '80' },
                        { offset: 0.5, color: stat.chartColor[1] + '40' },
                        { offset: 1, color: stat.chartColor[1] + '05' }
                      ]
                    }
                  },
                  emphasis: {
                    lineStyle: {
                      width: 4
                    }
                  }
                }]
              }
              
              chart.setOption(option)
              
              // 延迟调用 resize 确保渲染
              setTimeout(() => {
                chart.resize()
              }, 50)
              
              this.charts.push(chart)
            } catch (e) {
              console.error(`图表 ${index} 初始化失败:`, e)
            }
          }
        }
      }, 100)
    },

    async loadRecentActivities() {
      try {
        const activitiesRes = await getRecentActivities()
        if (activitiesRes.code === 200 && activitiesRes.data) {
          this.recentActivities = activitiesRes.data
          console.log('✅ 最近活动加载成功:', this.recentActivities.length, '条')
        }
      } catch (error) {
        console.error('❌ 加载最近活动失败:', error)
        this.recentActivities = []
      }
    },

    async loadUserDetails() {
      try {
        const response = await queryAccountList()
        const accountList = response.code === 200 ? response.data : response
        
        if (accountList && Array.isArray(accountList)) {
          const currentUser = accountList.find(acc => acc.account_id === this.account_id)
          if (currentUser) {
            this.userDetails = {
              organization_name: currentUser.organization_name || '',
              department: currentUser.department || '',
              doctor_title: currentUser.doctor_title || '',
              age: currentUser.age || 0,
              gender: currentUser.gender || ''
            }
            console.log('✅ 用户详细信息加载成功:', this.userDetails)
          } else {
            console.warn('⚠️ 未找到当前用户的区块链信息，account_id:', this.account_id)
            const userInfo = sessionStorage.getItem('userInfo')
            if (userInfo) {
              const user = JSON.parse(userInfo)
              this.userDetails = {
                organization_name: user.organization_name || '',
                department: user.department || '',
                doctor_title: user.doctor_title || '',
                age: user.age || 0,
                gender: user.gender || ''
              }
              console.log('✅ 从 sessionStorage 获取用户信息:', this.userDetails)
            }
          }
        }
      } catch (error) {
        console.error('❌ 加载用户详细信息失败:', error)
        try {
          const userInfo = sessionStorage.getItem('userInfo')
          if (userInfo) {
            const user = JSON.parse(userInfo)
            this.userDetails = {
              organization_name: user.organization_name || '',
              department: user.department || '',
              doctor_title: user.doctor_title || '',
              age: user.age || 0,
              gender: user.gender || ''
            }
            console.log('✅ 从 sessionStorage 获取用户信息（fallback）:', this.userDetails)
          }
        } catch (e) {
          console.error('❌ 从 sessionStorage 获取用户信息失败:', e)
        }
      }
    },

    toggleActivityExpand() {
      this.activityExpanded = !this.activityExpanded
    }
  }
}
</script>

<style lang="scss" scoped>
.dashboard-container {
  padding: 20px;
  background: linear-gradient(to bottom, #f5f7fa 0%, #e8eef5 100%);
  min-height: calc(100vh - 60px);
}

.welcome-section {
  margin-bottom: 20px;
  position: relative;

  .welcome-card {
    border-radius: 20px;
    background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
    border: none;
    box-shadow: 0 8px 32px rgba(168, 237, 234, 0.3);
    transition: all 0.3s ease;
    
    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 12px 48px rgba(168, 237, 234, 0.4);
    }
    
    ::v-deep .el-card__body {
      padding: 32px;
    }

    .welcome-content {
      display: flex;
      justify-content: space-between;
      align-items: center;
      color: rgba(255, 255, 255, 0.98);

      .welcome-text {
        flex: 1;
        
        .welcome-quote-top {
          font-size: 15px;
          color: rgba(255, 255, 255, 0.75);
          font-style: normal;
          font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', sans-serif;
          font-weight: 300;
          letter-spacing: 0.5px;
          margin-bottom: 18px;
          line-height: 1.5;
        }
        
        .welcome-header {
          display: flex;
          align-items: flex-start;
          justify-content: space-between;
          gap: 15px;
          margin-bottom: 16px;
          flex-wrap: wrap;
        }
        
        h2 {
          margin: 0;
          font-size: 32px;
          font-weight: 700;
          flex: 1;
          min-width: 200px;
          color: rgba(255, 255, 255, 1);
          text-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
          letter-spacing: 0.5px;
        }

        .status-info {
          display: flex;
          flex-direction: column;
          align-items: flex-end;
          gap: 6px;
        }

        .online-status {
          display: inline-flex;
          align-items: center;
          gap: 8px;
          background: rgba(255, 255, 255, 0.95);
          padding: 8px 16px;
          border-radius: 20px;
          font-size: 14px;
          font-weight: 600;
          color: #67C23A;
          box-shadow: 0 4px 12px rgba(103, 194, 58, 0.25);

          .status-dot-pulse {
            width: 10px;
            height: 10px;
            border-radius: 50%;
            background: #67C23A;
            box-shadow: 0 0 0 0 rgba(103, 194, 58, 0.7);
            animation: pulse-ring 2s infinite;
          }
        }

        .last-login {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          font-size: 12px;
          color: rgba(255, 255, 255, 0.95);
          background: rgba(255, 255, 255, 0.25);
          padding: 5px 12px;
          border-radius: 12px;
          backdrop-filter: blur(10px);

          i {
            font-size: 13px;
            color: rgba(255, 255, 255, 0.9);
          }
        }

        .security-info {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          font-size: 11px;
          color: rgba(255, 255, 255, 0.9);
          background: rgba(255, 255, 255, 0.2);
          padding: 4px 10px;
          border-radius: 10px;
          backdrop-filter: blur(10px);
          white-space: nowrap;

          i {
            font-size: 12px;
            color: rgba(255, 255, 255, 0.85);
          }
        }

        .role-tag {
          font-size: 16px;
          margin-bottom: 16px;
          color: rgba(255, 255, 255, 0.95);
          font-weight: 500;
          
          i {
            margin-right: 8px;
            color: rgba(255, 255, 255, 0.9);
          }
        }

        .user-details {
          margin-top: 15px;
          display: flex;
          flex-direction: column;
          gap: 8px;

          .detail-item {
            display: flex;
            align-items: center;
            font-size: 13px;
            color: rgba(255, 255, 255, 0.9);
            font-weight: 400;

            i {
              margin-right: 8px;
              font-size: 15px;
              color: rgba(255, 255, 255, 0.85);
            }

            .detail-label {
              font-weight: 400;
              margin-right: 8px;
              color: rgba(255, 255, 255, 0.85);
            }

            .detail-value {
              font-weight: 500;
              background: rgba(255, 255, 255, 0.2);
              padding: 3px 10px;
              border-radius: 6px;
              color: rgba(255, 255, 255, 0.98);
              backdrop-filter: blur(5px);
            }
          }
        }
      }

      .welcome-avatar {
        img {
          width: 90px;
          height: 90px;
          border-radius: 20px;
          border: 3px solid rgba(255, 255, 255, 0.4);
          object-fit: contain;
          object-position: center;
          background: rgba(255, 255, 255, 0.9);
          box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
          transition: all 0.3s ease;
          
          &:hover {
            transform: scale(1.05);
            box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
          }
        }
      }
    }
  }
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  
  i {
    margin-right: 8px;
    color: #4facfe;
    font-size: 22px;
  }
}

.stats-section {
  margin-bottom: 30px;

  .stat-card {
    margin-bottom: 20px;
    border-radius: 16px;
    transition: all 0.3s;
    border: none;
    background: white;
    overflow: hidden;
    position: relative;

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      height: 4px;
      background: var(--card-gradient);
    }

    &:hover {
      transform: translateY(-5px);
      box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
    }

    ::v-deep .el-card__body {
      padding: 20px;
    }

    .stat-content {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 16px;

      .stat-left {
        display: flex;
        align-items: center;
        gap: 15px;
        flex: 1;
        min-width: 0;

        .stat-icon {
          width: 56px;
          height: 56px;
          border-radius: 14px;
          display: flex;
          align-items: center;
          justify-content: center;
          flex-shrink: 0;
          position: relative;
          overflow: hidden;

          &::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: var(--card-gradient);
            opacity: 0.15;
          }

          i {
            font-size: 28px;
            background: var(--card-gradient);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            position: relative;
            z-index: 1;
          }
        }

        .stat-info {
          flex: 1;
          min-width: 0;

          .stat-value {
            margin: 0 0 4px 0;
            font-size: 28px;
            font-weight: 700;
            background: var(--card-gradient);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            line-height: 1;
          }

          .stat-title {
            margin: 0 0 3px 0;
            font-size: 14px;
            color: #606266;
            font-weight: 600;
          }

          .stat-subtitle {
            margin: 0;
            font-size: 12px;
            color: #909399;
          }
        }
      }

      .stat-chart {
        width: 100px;
        height: 70px;
        flex-shrink: 0;
        position: relative;
        border-radius: 8px;
        overflow: hidden;

        .mini-chart {
          width: 100%;
          height: 100%;
          position: absolute;
          top: 0;
          left: 0;
        }
      }
    }
  }
}

.activity-section {
  margin-bottom: 30px;

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    .expand-btn {
      color: #409EFF;
      font-size: 14px;
      padding: 8px 16px;
      
      &:hover {
        background: #ecf5ff;
      }

      i {
        margin-left: 4px;
        transition: transform 0.3s;
      }
    }
  }

  .activity-card {
    border-radius: 12px;
    border: none;

    .activity-item {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      gap: 16px;

      .activity-content {
        display: flex;
        align-items: center;
        gap: 10px;
        font-size: 14px;
        flex: 1;
        min-width: 0;

        i {
          font-size: 18px;
          flex-shrink: 0;
        }

        .activity-type {
          font-weight: 600;
          color: #303133;
          flex-shrink: 0;
        }

        .activity-text {
          color: #606266;
          word-break: break-word;
        }
      }

      .activity-time {
        font-size: 13px;
        color: #909399;
        white-space: nowrap;
        flex-shrink: 0;
        padding: 4px 12px;
        background: #f5f7fa;
        border-radius: 8px;
        font-family: 'Courier New', monospace;
      }
    }
  }
}

.modules-section {
  margin-bottom: 30px;

  .module-card {
    margin-bottom: 20px;
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.3s;
    border: none;

    &:hover {
      transform: translateY(-5px);
      box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
    }

    .module-content {
      padding: 20px;
      display: flex;
      align-items: center;
      gap: 15px;

      .module-icon {
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

      .module-info {
        flex: 1;

        h4 {
          margin: 0 0 5px 0;
          font-size: 16px;
          font-weight: 600;
          color: #303133;
        }

        p {
          margin: 0;
          font-size: 13px;
          color: #909399;
        }
      }

      .module-arrow {
        i {
          font-size: 18px;
          color: #C0C4CC;
          transition: all 0.3s;
        }
      }
    }

    &:hover .module-arrow i {
      color: #4facfe;
      transform: translateX(5px);
    }
  }
}

@keyframes pulse-ring {
  0% {
    box-shadow: 0 0 0 0 rgba(103, 194, 58, 0.7);
  }
  50% {
    box-shadow: 0 0 0 8px rgba(103, 194, 58, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(103, 194, 58, 0);
  }
}

@keyframes pulse-dot {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.7;
    transform: scale(1.2);
  }
}

@media (max-width: 768px) {
  .welcome-content {
    flex-direction: column;
    text-align: center;

    .welcome-avatar {
      margin-top: 20px;
    }
  }
}
</style>
