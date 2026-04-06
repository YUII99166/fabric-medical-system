<template>
  <div class="container">
    <el-alert
      class="info-alert"
      type="success"
    >
      <p>账户ID: {{ account_id }}</p>
      <p>用户名: {{ account_name }}</p>
    </el-alert>

    <!-- 搜索框 - 仅医生可见 -->
    <el-card v-if="roles[0] === 'doctor' || roles[0] === 'admin'" class="search-card" shadow="hover">
      <div class="search-box">
        <el-input
          v-model="searchKey"
          placeholder="输入患者姓名或用户名搜索，留空显示所有病历"
          clearable
          @clear="loadAllPrescriptions"
          @keyup.enter.native="searchPrescriptions"
        >
          <el-button slot="append" icon="el-icon-search" @click="searchPrescriptions">搜索</el-button>
        </el-input>
      </div>
    </el-card>
    
    <div v-if="prescriptionList.length==0" style="text-align: center; padding: 40px;">
      <el-empty description="暂无病历数据">
        <template v-if="roles[0] === 'doctor' || roles[0] === 'admin'">
          <el-button type="primary" @click="$router.push('/prescription/add')">
            <i class="el-icon-plus"></i>
            创建第一条病历
          </el-button>
        </template>
      </el-empty>
    </div>
    
    <!-- 今天的病历 -->
    <div v-if="hasToday" class="date-group">
      <div class="date-header">
        <i class="el-icon-date"></i>
        <span>今天</span>
        <span class="date-count">({{ groupedPrescriptions.today.length }}条)</span>
      </div>
      <el-row v-loading="loading" :gutter="20">
        <el-col v-for="(val,index) in groupedPrescriptions.today" :key="'today-'+index" :xs="24" :sm="12" :md="8" :lg="6">
          <el-card class="prescription-card" shadow="hover" @click.native="showQuickView(val)">
            <!-- 标签：本院/他院（已授权） -->
            <div v-if="roles[0] === 'doctor'" class="card-tag">
              <el-tag v-if="isOwnHospital(val)" type="success" size="mini">本院</el-tag>
              <el-tag v-else type="info" size="mini">他院（已授权）</el-tag>
            </div>
            
            <div class="card-icon">
              <i class="el-icon-document"></i>
            </div>
            
            <div class="card-title">病历记录</div>
            
            <div class="card-info">
              <div class="info-item">
                <i class="el-icon-user"></i>
                <span>{{ val.patient_name || getAccountName(val.patient) }}</span>
              </div>
              <div class="info-item">
                <i class="el-icon-s-order"></i>
                <span>{{ val.diagnosis }}</span>
              </div>
              <div class="info-item">
                <i class="el-icon-office-building"></i>
                <span>{{ val.hospital_name || val.organization_name }}</span>
              </div>
              <div class="info-item time">
                <i class="el-icon-time"></i>
                <span>{{ val.created }}</span>
              </div>
            </div>
            
            <div class="card-actions">
              <el-button type="text" size="small" @click.native.stop="showQuickView(val)">
                <i class="el-icon-view"></i> 快速预览
              </el-button>
              <el-button type="primary" size="small" @click.native.stop="goToDetail(val)">
                <i class="el-icon-document"></i> 查看详情
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 昨天的病历 -->
    <!-- 昨天的病历 -->
    <div v-if="hasYesterday" class="date-group">
      <div class="date-header">
        <i class="el-icon-date"></i>
        <span>昨天</span>
        <span class="date-count">({{ groupedPrescriptions.yesterday.length }}条)</span>
      </div>
      <el-row v-loading="loading" :gutter="20">
        <el-col v-for="(val,index) in groupedPrescriptions.yesterday" :key="'yesterday-'+index" :xs="24" :sm="12" :md="8" :lg="6">
          <el-card class="prescription-card" shadow="hover" @click.native="showQuickView(val)">
            <!-- 标签：本院/他院（已授权） -->
            <div v-if="roles[0] === 'doctor'" class="card-tag">
              <el-tag v-if="isOwnHospital(val)" type="success" size="mini">本院</el-tag>
              <el-tag v-else type="info" size="mini">他院（已授权）</el-tag>
            </div>
            
            <div class="card-icon">
              <i class="el-icon-document"></i>
            </div>
            
            <div class="card-title">病历记录</div>
            
            <div class="card-info">
              <div class="info-item">
                <i class="el-icon-user"></i>
                <span>{{ val.patient_name || getAccountName(val.patient) }}</span>
              </div>
              <div class="info-item">
                <i class="el-icon-s-order"></i>
                <span>{{ val.diagnosis }}</span>
              </div>
              <div class="info-item">
                <i class="el-icon-office-building"></i>
                <span>{{ val.hospital_name || val.organization_name }}</span>
              </div>
              <div class="info-item time">
                <i class="el-icon-time"></i>
                <span>{{ val.created }}</span>
              </div>
            </div>
            
            <div class="card-actions">
              <el-button type="text" size="small" @click.native.stop="showQuickView(val)">
                <i class="el-icon-view"></i> 快速预览
              </el-button>
              <el-button type="primary" size="small" @click.native.stop="goToDetail(val)">
                <i class="el-icon-document"></i> 查看详情
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 更早的病历 -->
    <div v-if="hasEarlier" class="date-group">
      <div class="date-header">
        <i class="el-icon-date"></i>
        <span>更早</span>
        <span class="date-count">({{ groupedPrescriptions.earlier.length }}条)</span>
      </div>
      <el-row v-loading="loading" :gutter="20">
        <el-col v-for="(val,index) in groupedPrescriptions.earlier" :key="'earlier-'+index" :xs="24" :sm="12" :md="8" :lg="6">
          <el-card class="prescription-card" shadow="hover" @click.native="showQuickView(val)">
            <!-- 标签：本院/他院（已授权） -->
            <div v-if="roles[0] === 'doctor'" class="card-tag">
              <el-tag v-if="isOwnHospital(val)" type="success" size="mini">本院</el-tag>
              <el-tag v-else type="info" size="mini">他院（已授权）</el-tag>
            </div>
            
            <div class="card-icon">
              <i class="el-icon-document"></i>
            </div>
            
            <div class="card-title">病历记录</div>
            
            <div class="card-info">
              <div class="info-item">
                <i class="el-icon-user"></i>
                <span>{{ val.patient_name || getAccountName(val.patient) }}</span>
              </div>
              <div class="info-item">
                <i class="el-icon-s-order"></i>
                <span>{{ val.diagnosis }}</span>
              </div>
              <div class="info-item">
                <i class="el-icon-office-building"></i>
                <span>{{ val.hospital_name || val.organization_name }}</span>
              </div>
              <div class="info-item time">
                <i class="el-icon-time"></i>
                <span>{{ val.created }}</span>
              </div>
            </div>
            
            <div class="card-actions">
              <el-button type="text" size="small" @click.native.stop="showQuickView(val)">
                <i class="el-icon-view"></i> 快速预览
              </el-button>
              <el-button type="primary" size="small" @click.native.stop="goToDetail(val)">
                <i class="el-icon-document"></i> 查看详情
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 详情对话框 -->
    <el-dialog
      title="病历详情"
      :visible.sync="dialogVisible"
      width="700px"
      :close-on-click-modal="false"
    >
      <div v-if="selectedPrescription" class="detail-content">
        <!-- 核心诊疗信息 - 突出显示 -->
        <div class="highlight-section">
          <!-- 主诉 -->
          <div class="info-section">
            <div class="section-header">
              <i class="el-icon-chat-line-square"></i>
              <span>主诉</span>
            </div>
            <div class="section-text">{{ selectedPrescription.chief_complaint || '未填写' }}</div>
          </div>

          <!-- 现病史 -->
          <div class="info-section">
            <div class="section-header">
              <i class="el-icon-document-copy"></i>
              <span>现病史</span>
            </div>
            <div class="section-text">{{ selectedPrescription.present_illness || '未填写' }}</div>
          </div>

          <!-- 体格检查 -->
          <div class="info-section">
            <div class="section-header">
              <i class="el-icon-data-analysis"></i>
              <span>体格检查</span>
            </div>
            <div class="section-text">{{ selectedPrescription.physical_exam || '未填写' }}</div>
          </div>

          <!-- 诊断信息 -->
          <div class="diagnosis-highlight">
            <div class="diagnosis-header">
              <i class="el-icon-document-checked"></i>
              <span>诊断结果</span>
            </div>
            <div class="diagnosis-text">{{ selectedPrescription.diagnosis }}</div>
          </div>

          <!-- 药品信息 -->
          <div class="drug-highlight">
            <div class="drug-header">
              <i class="el-icon-medicine-box"></i>
              <span>处方药品</span>
            </div>
            <div class="drug-list-compact">
              <div v-for="(drug, index) in selectedPrescription.drug" :key="index" class="drug-item-compact">
                <span class="drug-index">{{ index + 1 }}</span>
                <span class="drug-name-compact">{{ drug.Name }}</span>
                <span class="drug-amount-compact">× {{ drug.amount }}</span>
              </div>
            </div>
          </div>

          <!-- 医嘱 -->
          <div class="advice-compact">
            <i class="el-icon-warning-outline"></i>
            <div>
              <div style="font-weight: 600; margin-bottom: 4px;">医嘱</div>
              <div>{{ selectedPrescription.medical_advice || '未填写' }}</div>
            </div>
          </div>

          <!-- 备注 -->
          <div v-if="selectedPrescription.comment" class="comment-compact">
            <i class="el-icon-edit-outline"></i>
            <span>{{ selectedPrescription.comment }}</span>
          </div>
        </div>

        <!-- 基本信息 - 紧凑展示 -->
        <div class="compact-info">
          <div class="info-row">
            <span class="info-icon"><i class="el-icon-user"></i></span>
            <span class="info-text">
              <strong>医生：</strong>{{ selectedPrescription.doctor_name || getAccountName(selectedPrescription.doctor) }}
              <span class="info-detail">（{{ selectedPrescription.doctor_title || '医师' }}）</span>
            </span>
          </div>
          <div class="info-row">
            <span class="info-icon"><i class="el-icon-office-building"></i></span>
            <span class="info-text">
              <strong>医院：</strong>{{ selectedPrescription.organization_name || selectedPrescription.hospital_name || '未设置' }}
              <span class="info-detail" v-if="selectedPrescription.department">- {{ selectedPrescription.department }}</span>
            </span>
          </div>
          <div class="info-row">
            <span class="info-icon"><i class="el-icon-user-solid"></i></span>
            <span class="info-text">
              <strong>患者：</strong>{{ selectedPrescription.patient_name || getAccountName(selectedPrescription.patient) }}
              <span class="info-detail">（ID: {{ selectedPrescription.patient }}）</span>
            </span>
          </div>
          <div class="info-row">
            <span class="info-icon"><i class="el-icon-time"></i></span>
            <span class="info-text">
              <strong>时间：</strong>{{ selectedPrescription.created }}
              <span class="info-detail">（病历ID: {{ selectedPrescription.id }}）</span>
            </span>
          </div>
        </div>

        <!-- 区块链信息 - 可折叠 -->
        <el-collapse v-if="selectedPrescription.tx_id || selectedPrescription.organization_id" class="blockchain-collapse">
          <el-collapse-item title="区块链信息" name="1">
            <div class="blockchain-info">
              <div v-if="selectedPrescription.organization_id" class="blockchain-row">
                <span class="bc-label">组织标识：</span>
                <span class="bc-value">{{ selectedPrescription.organization_id }}</span>
              </div>
              <div v-if="selectedPrescription.creator_msp_id" class="blockchain-row">
                <span class="bc-label">创建者：</span>
                <span class="bc-value">{{ selectedPrescription.creator_msp_id }}</span>
              </div>
              <div v-if="selectedPrescription.tx_id" class="blockchain-row">
                <span class="bc-label">交易ID：</span>
                <span class="bc-value tx-id">{{ selectedPrescription.tx_id }}</span>
              </div>
              <div v-if="selectedPrescription.authorized_orgs && selectedPrescription.authorized_orgs.length > 0" class="blockchain-row">
                <span class="bc-label">授权组织：</span>
                <span class="bc-value">
                  <el-tag v-for="(org, index) in selectedPrescription.authorized_orgs" :key="index" size="mini" style="margin-right: 5px;">
                    {{ org }}
                  </el-tag>
                </span>
              </div>
            </div>
          </el-collapse-item>
        </el-collapse>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="dialogVisible = false">关闭</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { queryPrescriptionList } from '@/api/prescription'
import { getAccountMap } from '@/utils/accountCache'

export default {
  name: 'Prescription',
  data() {
    return {
      loading: true,
      prescriptionList: [],
      accountList: [],
      dialogVisible: false,
      selectedPrescription: null,
      accountMap: {}, // 账户ID到账户名称的映射
      searchKey: '', // 搜索关键字
      allPrescriptions: [] // 保存所有病历，用于搜索
    }
  },
  computed: {
    ...mapGetters([
      'account_id',
      'roles',
      'account_name',
    ]),
    // 按日期分组的病历
    groupedPrescriptions() {
      const groups = {
        today: [],
        yesterday: [],
        earlier: []
      }
      
      // 确保 prescriptionList 是数组
      if (!Array.isArray(this.prescriptionList)) {
        console.warn('prescriptionList 不是数组:', this.prescriptionList)
        return groups
      }
      
      const now = new Date()
      const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
      const yesterday = new Date(today)
      yesterday.setDate(yesterday.getDate() - 1)
      
      this.prescriptionList.forEach(prescription => {
        try {
          // 解析日期，支持多种格式
          let createdDate
          if (prescription.created) {
            // 尝试解析日期字符串
            createdDate = new Date(prescription.created)
            
            // 如果解析失败，尝试其他格式
            if (isNaN(createdDate.getTime())) {
              // 尝试替换格式 "2024-01-01 12:00:00" -> "2024-01-01T12:00:00"
              const dateStr = prescription.created.replace(' ', 'T')
              createdDate = new Date(dateStr)
            }
            
            // 如果还是失败，放到更早分组
            if (isNaN(createdDate.getTime())) {
              console.warn('无法解析日期:', prescription.created)
              groups.earlier.push(prescription)
              return
            }
            
            // 只比较日期部分，忽略时间
            const prescriptionDate = new Date(createdDate.getFullYear(), createdDate.getMonth(), createdDate.getDate())
            
            if (prescriptionDate.getTime() === today.getTime()) {
              groups.today.push(prescription)
            } else if (prescriptionDate.getTime() === yesterday.getTime()) {
              groups.yesterday.push(prescription)
            } else {
              groups.earlier.push(prescription)
            }
          } else {
            // 没有日期，放到更早分组
            groups.earlier.push(prescription)
          }
        } catch (error) {
          console.error('处理病历日期时出错:', error, prescription)
          groups.earlier.push(prescription)
        }
      })
      
      console.log('分组结果:', {
        today: groups.today.length,
        yesterday: groups.yesterday.length,
        earlier: groups.earlier.length
      })
      
      return groups
    },
    // 是否有今天的病历
    hasToday() {
      return this.groupedPrescriptions.today.length > 0
    },
    // 是否有昨天的病历
    hasYesterday() {
      return this.groupedPrescriptions.yesterday.length > 0
    },
    // 是否有更早的病历
    hasEarlier() {
      return this.groupedPrescriptions.earlier.length > 0
    }
  },
  created() {
    // 使用缓存加载账户列表
    getAccountMap().then(accountMap => {
      this.accountMap = accountMap
    }).catch(error => {
      console.error('加载账户映射失败:', error)
    })

    // 加载病历列表
    if (this.roles[0] === 'admin' || this.roles[0] === 'doctor') {
      console.log('=== 开始加载病历列表 ===')
      console.log('当前角色:', this.roles[0])
      console.log('账户ID:', this.account_id)
      
      queryPrescriptionList().then(response => {
        console.log('病历列表原始响应:', response)
        console.log('响应类型:', typeof response)
        console.log('是否为null:', response === null)
        console.log('是否为数组:', Array.isArray(response))
        
        // 正确处理 API 响应格式
        let prescriptions = []
        if (response && response.code === 200) {
          console.log('响应格式: {code: 200, data: ...}')
          if (response.data) {
            prescriptions = Array.isArray(response.data) ? response.data : []
            console.log('从 response.data 提取病历:', prescriptions.length, '条')
          } else {
            console.warn('response.data 为空')
          }
        } else if (Array.isArray(response)) {
          console.log('响应格式: 直接数组')
          prescriptions = response
          console.log('直接使用数组:', prescriptions.length, '条')
        } else if (response !== null && typeof response === 'object') {
          console.log('响应格式: 对象（尝试作为数组）')
          // 尝试将对象转换为数组
          prescriptions = Object.values(response)
          console.log('转换为数组:', prescriptions.length, '条')
        } else {
          console.warn('无法识别的响应格式:', response)
        }
        
        console.log('解析后的病历数据:', prescriptions)
        console.log('病历数据类型:', Array.isArray(prescriptions) ? '数组' : typeof prescriptions)
        console.log('病历数量:', Array.isArray(prescriptions) ? prescriptions.length : 0)
        
        if (!Array.isArray(prescriptions)) {
          console.error('❌ 病历数据不是数组格式:', prescriptions)
          this.prescriptionList = []
          this.allPrescriptions = []
          this.loading = false
          this.$message.warning('病历数据格式错误，请联系管理员')
          return
        }
        
        if (prescriptions.length === 0) {
          console.warn('⚠️ 没有病历数据')
          this.prescriptionList = []
          this.allPrescriptions = []
          this.loading = false
          return
        }
        
        console.log('✅ 成功获取', prescriptions.length, '条病历')
        console.log('第一条病历示例:', prescriptions[0])
        
        if (prescriptions.length > 0) {
          // 如果是医生，需要过滤病历
          if (this.roles[0] === 'doctor') {
            // 获取当前医生的组织信息
            const userInfo = JSON.parse(sessionStorage.getItem('userInfo'))
            console.log('用户信息:', userInfo)
            
            // 从缓存获取医生的组织信息
            getAccountMap().then(accountMap => {
              const doctorAccount = Object.values(accountMap).find(acc => acc.account_id === this.account_id) || 
                                   accountMap[this.account_id]
              console.log('医生账户信息:', doctorAccount)
              
              if (doctorAccount && doctorAccount.organization) {
                const doctorOrg = doctorAccount.organization
                console.log('医生组织:', doctorOrg)
                
                // 过滤病历：只显示本院病历或已授权的他院病历
                const filteredPrescriptions = prescriptions.filter(prescription => {
                  // 本院病历
                  if (prescription.organization_id === doctorOrg) {
                    return true
                  }
                  // 他院病历但已授权
                  if (prescription.authorized_orgs && prescription.authorized_orgs.includes(doctorOrg)) {
                    return true
                  }
                  return false
                })
                
                this.prescriptionList = filteredPrescriptions
                this.allPrescriptions = filteredPrescriptions
                console.log('过滤后的病历数量:', filteredPrescriptions.length)
                
                // 保存组织信息到sessionStorage
                userInfo.organization = doctorOrg
                userInfo.organization_name = doctorAccount.organization_name
                sessionStorage.setItem('userInfo', JSON.stringify(userInfo))
              } else {
                console.warn('未找到医生的组织信息，显示所有病历')
                this.prescriptionList = prescriptions
                this.allPrescriptions = prescriptions
              }
              this.loading = false
            }).catch(err => {
              console.error('获取账户映射失败:', err)
              // 如果获取失败，显示所有病历
              this.prescriptionList = prescriptions
              this.allPrescriptions = prescriptions
              this.loading = false
            })
          } else {
            // 管理员可以看到所有病历
            this.prescriptionList = prescriptions
            this.allPrescriptions = prescriptions
            this.loading = false
          }
          console.log('第一条病历数据:', this.prescriptionList[0])
        } else {
          // 没有病历数据
          this.prescriptionList = []
          this.allPrescriptions = []
          this.loading = false
        }
      }).catch(_ => {
        this.loading = false
      })
    } else {
      queryPrescriptionList({ patient: this.account_id }).then(response => {
        console.log('病历列表响应:', response)
        if (response !== null) {
          this.prescriptionList = response
          this.allPrescriptions = response // 保存所有病历
          console.log('第一条病历数据:', response[0])
        }
        this.loading = false
      }).catch(_ => {
        this.loading = false
      })
    }
  },
  methods: {
    // 快速预览（弹出对话框）
    showQuickView(item) {
      console.log('showQuickView 被调用', item)
      this.selectedPrescription = item
      this.dialogVisible = true
      console.log('dialogVisible 设置为:', this.dialogVisible)
    },
    // 跳转到详情页面
    goToDetail(item) {
      console.log('goToDetail 被调用', item)
      this.$router.push({
        name: 'PrescriptionDetail',
        query: { id: item.id }
      })
    },
    // 获取账户名称
    getAccountName(accountId) {
      return this.accountMap[accountId]?.name || accountId
    },
    // 获取账户用户名
    getAccountUsername(accountId) {
      return this.accountMap[accountId]?.username || '-'
    },
    // 搜索病历
    searchPrescriptions() {
      if (!this.searchKey.trim()) {
        // 如果搜索框为空，显示所有病历
        this.prescriptionList = this.allPrescriptions
        return
      }

      const keyword = this.searchKey.trim().toLowerCase()
      
      // 在本地过滤病历
      this.prescriptionList = this.allPrescriptions.filter(prescription => {
        // 搜索患者姓名
        if (prescription.patient_name && prescription.patient_name.toLowerCase().includes(keyword)) {
          return true
        }
        // 搜索患者ID
        if (prescription.patient && prescription.patient.toLowerCase().includes(keyword)) {
          return true
        }
        // 搜索患者用户名（通过accountMap）
        const patientInfo = this.accountMap[prescription.patient]
        if (patientInfo && patientInfo.username && patientInfo.username.toLowerCase().includes(keyword)) {
          return true
        }
        return false
      })
    },
    // 加载所有病历
    loadAllPrescriptions() {
      this.searchKey = ''
      this.prescriptionList = this.allPrescriptions
    },
    // 判断是否是本院病历
    isOwnHospital(prescription) {
      const userInfo = JSON.parse(sessionStorage.getItem('userInfo'))
      if (userInfo && userInfo.organization) {
        return prescription.organization_id === userInfo.organization
      }
      // 如果没有组织信息，从accountMap中查找
      const doctorAccount = this.accountMap[this.account_id]
      if (doctorAccount && doctorAccount.organization) {
        return prescription.organization_id === doctorAccount.organization
      }
      return false
    }
  }
}
</script>

<style scoped>
  .container {
    width: 100%;
    min-height: 100%;
    padding: 20px;
    background: #f0f2f5;
  }

  .info-alert {
    margin-bottom: 20px;
    border-radius: 8px;
  }

  .info-alert >>> .el-alert__content p {
    margin: 5px 0;
    color: var(--theme-primary, #409EFF);
    font-weight: 500;
  }

  .search-card {
    margin-bottom: 20px;
    border-radius: 8px;
  }

  .search-box {
    display: flex;
    align-items: center;
  }

  .search-box >>> .el-input-group {
    width: 100%;
  }

  .date-group {
    margin-bottom: 30px;
  }

  .date-header {
    display: flex;
    align-items: center;
    padding: 12px 20px;
    background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
    border-radius: 8px;
    margin-bottom: 20px;
    color: white;
    font-size: 16px;
    font-weight: 600;
    box-shadow: 0 2px 8px rgba(79, 172, 254, 0.3);
  }

  .date-header i {
    margin-right: 10px;
    font-size: 18px;
  }

  .date-count {
    margin-left: 10px;
    font-size: 14px;
    opacity: 0.9;
    font-weight: 400;
  }

  .prescription-card {
    margin-bottom: 20px;
    border-radius: 12px;
    transition: all 0.3s;
    cursor: pointer;
    border: 2px solid #EBEEF5;
    overflow: hidden;
    position: relative;
  }

  .card-tag {
    position: absolute;
    top: 10px;
    right: 10px;
    z-index: 1;
  }

  .prescription-card:hover {
    transform: translateY(-8px);
    box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
    border-color: var(--theme-primary, #409EFF);
  }

  .prescription-card >>> .el-card__body {
    padding: 24px;
    text-align: center;
  }

  .card-icon {
    width: 60px;
    height: 60px;
    margin: 0 auto 16px;
    background: var(--theme-light-bg, #ecf5ff);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s;
  }

  .prescription-card:hover .card-icon {
    background: var(--theme-primary, #409EFF);
  }

  .card-icon i {
    font-size: 28px;
    color: var(--theme-primary, #409EFF);
    transition: all 0.3s;
  }

  .prescription-card:hover .card-icon i {
    color: white;
  }

  .card-title {
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    margin-bottom: 16px;
  }

  .card-info {
    margin-bottom: 16px;
  }

  .info-item {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 8px;
    color: #606266;
    font-size: 14px;
  }

  .info-item i {
    margin-right: 6px;
    color: var(--theme-primary, #409EFF);
  }

  .info-item.time {
    font-size: 12px;
    color: #909399;
  }

  .card-actions {
    display: flex;
    gap: 8px;
    padding-top: 12px;
    border-top: 1px solid #EBEEF5;
  }

  .card-actions .el-button {
    flex: 1;
    margin: 0;
  }

  .card-actions .el-button--text {
    color: #606266;
    border: 1px solid #DCDFE6;
    background: #fff;
  }

  .card-actions .el-button--text:hover {
    color: var(--theme-primary, #409EFF);
    border-color: var(--theme-primary, #409EFF);
    background: var(--theme-light-bg, #ecf5ff);
  }

  .card-actions .el-button--primary {
    background: var(--theme-primary, #409EFF);
    border-color: var(--theme-primary, #409EFF);
  }

  .card-actions .el-button i {
    margin-right: 4px;
  }

  /* 详情对话框样式 */
  .detail-content {
    padding: 0;
  }

  /* 病历摘要样式 */
  .summary-section {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 20px;
    border-radius: 12px;
    margin-bottom: 20px;
    color: white;
  }

  .summary-header {
    display: flex;
    align-items: center;
    font-size: 16px;
    font-weight: 600;
    margin-bottom: 12px;
    opacity: 0.95;
  }

  .summary-header i {
    margin-right: 8px;
    font-size: 18px;
  }

  .summary-text {
    background: rgba(255, 255, 255, 0.15);
    padding: 16px;
    border-radius: 8px;
    font-size: 14px;
    line-height: 1.8;
    backdrop-filter: blur(10px);
    text-align: justify;
  }

  .summary-text strong {
    color: #fff;
    font-weight: 600;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  }

  .summary-text .diagnosis-text-inline {
    color: #ffd700;
    font-size: 15px;
    padding: 2px 6px;
    background: rgba(255, 215, 0, 0.2);
    border-radius: 4px;
  }

  /* 核心诊疗信息 - 突出显示 */
  .highlight-section {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 24px;
    border-radius: 12px;
    margin-bottom: 20px;
    color: white;
  }

  .info-section {
    margin-bottom: 20px;
    padding-bottom: 16px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  }

  .section-header {
    display: flex;
    align-items: center;
    font-size: 14px;
    font-weight: 600;
    margin-bottom: 8px;
    opacity: 0.9;
  }

  .section-header i {
    margin-right: 6px;
    font-size: 16px;
  }

  .section-text {
    background: rgba(255, 255, 255, 0.15);
    padding: 12px 16px;
    border-radius: 8px;
    font-size: 14px;
    line-height: 1.6;
    backdrop-filter: blur(10px);
  }

  .diagnosis-highlight {
    margin-bottom: 20px;
  }

  .diagnosis-header {
    display: flex;
    align-items: center;
    font-size: 16px;
    font-weight: 600;
    margin-bottom: 12px;
    opacity: 0.9;
  }

  .diagnosis-header i {
    margin-right: 8px;
    font-size: 18px;
  }

  .diagnosis-text {
    background: rgba(255, 255, 255, 0.2);
    padding: 16px;
    border-radius: 8px;
    font-size: 18px;
    font-weight: 500;
    line-height: 1.6;
    backdrop-filter: blur(10px);
  }

  .drug-highlight {
    margin-bottom: 16px;
  }

  .drug-header {
    display: flex;
    align-items: center;
    font-size: 16px;
    font-weight: 600;
    margin-bottom: 12px;
    opacity: 0.9;
  }

  .drug-header i {
    margin-right: 8px;
    font-size: 18px;
  }

  .drug-list-compact {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
  }

  .drug-item-compact {
    background: rgba(255, 255, 255, 0.2);
    padding: 10px 16px;
    border-radius: 20px;
    display: flex;
    align-items: center;
    gap: 8px;
    backdrop-filter: blur(10px);
    font-size: 14px;
  }

  .drug-index {
    background: rgba(255, 255, 255, 0.3);
    width: 24px;
    height: 24px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    font-size: 12px;
  }

  .drug-name-compact {
    font-weight: 500;
  }

  .drug-amount-compact {
    opacity: 0.9;
    font-size: 13px;
  }

  .advice-compact {
    background: rgba(255, 193, 7, 0.2);
    padding: 12px 16px;
    border-radius: 8px;
    font-size: 13px;
    display: flex;
    align-items: flex-start;
    gap: 10px;
    backdrop-filter: blur(10px);
    margin-bottom: 12px;
    border-left: 3px solid rgba(255, 193, 7, 0.8);
  }

  .advice-compact i {
    font-size: 18px;
    margin-top: 2px;
  }

  .comment-compact {
    background: rgba(255, 255, 255, 0.15);
    padding: 12px 16px;
    border-radius: 8px;
    font-size: 13px;
    display: flex;
    align-items: center;
    gap: 8px;
    backdrop-filter: blur(10px);
  }

  .comment-compact i {
    opacity: 0.8;
  }

  /* 紧凑信息展示 */
  .compact-info {
    background: #f5f7fa;
    padding: 16px;
    border-radius: 8px;
    margin-bottom: 16px;
  }

  .info-row {
    display: flex;
    align-items: center;
    padding: 8px 0;
    font-size: 14px;
    color: #606266;
    border-bottom: 1px solid #e4e7ed;
  }

  .info-row:last-child {
    border-bottom: none;
  }

  .info-icon {
    width: 32px;
    color: var(--theme-primary, #409EFF);
    font-size: 16px;
  }

  .info-text {
    flex: 1;
  }

  .info-text strong {
    color: #303133;
    margin-right: 8px;
  }

  .info-detail {
    color: #909399;
    font-size: 12px;
    margin-left: 8px;
  }

  /* 区块链信息折叠 */
  .blockchain-collapse {
    border: 1px solid #EBEEF5;
    border-radius: 8px;
  }

  .blockchain-collapse >>> .el-collapse-item__header {
    padding: 0 16px;
    font-size: 13px;
    color: #909399;
    background: #fafafa;
  }

  .blockchain-info {
    padding: 12px 16px;
    background: #f9f9f9;
  }

  .blockchain-row {
    display: flex;
    padding: 6px 0;
    font-size: 12px;
  }

  .bc-label {
    width: 80px;
    color: #909399;
    flex-shrink: 0;
  }

  .bc-value {
    flex: 1;
    color: #606266;
    word-break: break-all;
  }

  .bc-value.tx-id {
    font-family: monospace;
    font-size: 11px;
  }

  /* 保留原有样式 */
  .section {
    margin-bottom: 24px;
  }

  .section-title {
    display: flex;
    align-items: center;
    font-size: 16px;
    font-weight: 600;
    color: var(--theme-primary, #409EFF);
    margin-bottom: 16px;
    padding-bottom: 8px;
    border-bottom: 2px solid var(--theme-light-bg, #ecf5ff);
  }

  .section-title i {
    margin-right: 8px;
    font-size: 18px;
  }

  .info-box {
    background: #f5f7fa;
    padding: 12px 16px;
    border-radius: 8px;
    margin-bottom: 12px;
  }

  .info-label {
    font-size: 12px;
    color: #909399;
    margin-bottom: 6px;
  }

  .info-value {
    font-size: 14px;
    color: #303133;
    font-weight: 500;
    word-break: break-all;
  }

  .info-sub {
    font-size: 12px;
    color: #909399;
    margin-top: 4px;
  }

  .diagnosis-box {
    background: #fff9e6;
    border-left: 4px solid #e6a23c;
    padding: 16px;
    border-radius: 4px;
  }

  .diagnosis-label {
    font-size: 13px;
    color: #e6a23c;
    font-weight: 600;
    margin-bottom: 8px;
  }

  .diagnosis-content {
    font-size: 14px;
    color: #303133;
    line-height: 1.8;
  }

  .drug-list {
    background: #f0f9ff;
    padding: 16px;
    border-radius: 8px;
  }

  .drug-item-detail {
    display: flex;
    align-items: center;
    padding: 12px;
    background: white;
    border-radius: 6px;
    margin-bottom: 12px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  }

  .drug-item-detail:last-child {
    margin-bottom: 0;
  }

  .drug-number {
    width: 32px;
    height: 32px;
    background: var(--theme-primary, #409EFF);
    color: white;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    margin-right: 12px;
    flex-shrink: 0;
  }

  .drug-info {
    flex: 1;
  }

  .drug-name {
    font-size: 15px;
    color: #303133;
    font-weight: 600;
    margin-bottom: 4px;
  }

  .drug-amount-text {
    font-size: 13px;
    color: #909399;
  }

  .comment-box {
    background: #f5f7fa;
    padding: 16px;
    border-radius: 8px;
    font-size: 14px;
    color: #606266;
    line-height: 1.8;
    border-left: 4px solid var(--theme-primary, #409EFF);
  }

  >>> .el-dialog__header {
    background: var(--theme-gradient, linear-gradient(135deg, #667eea 0%, #764ba2 100%));
    padding: 20px;
  }

  >>> .el-dialog__title {
    color: white;
    font-weight: 600;
    font-size: 18px;
  }

  >>> .el-dialog__headerbtn .el-dialog__close {
    color: white;
    font-size: 20px;
  }

  >>> .el-dialog__body {
    padding: 24px;
  }

  >>> .el-dialog__footer {
    padding: 16px 24px;
    border-top: 1px solid #EBEEF5;
  }
</style>
