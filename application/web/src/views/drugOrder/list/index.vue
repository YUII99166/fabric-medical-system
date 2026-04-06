<template>
  <div class="container">
    <el-alert
      class="info-alert"
      type="success"
    >
      <p>账户ID: {{ account_id }}</p>
      <p>用户名: {{ account_name }}</p>

    </el-alert>
    <div v-if="drugOrderList.length==0" style="text-align: center; padding: 40px;">
      <el-empty description="暂无订单数据">
        <template v-if="roles[0] === 'drugstore' || roles[0] === 'admin'">
          <el-button type="primary" @click="$router.push('/addDrug')">
            <i class="el-icon-plus"></i>
            创建第一个订单
          </el-button>
        </template>
      </el-empty>
    </div>
    <el-row v-loading="loading" :gutter="20">
      <el-col v-for="(val,index) in drugOrderList" :key="index" :span="6" :offset="1">
        <el-card class="drugOrder-card">
          <div slot="header" class="card-header">
            <div class="header-title">
              <i class="el-icon-shopping-bag-2"></i>
              <span>药品订单</span>
            </div>
            <div class="order-id">{{ val.id }}</div>
          </div>

          <div class="card-body">
            <div class="info-row">
              <div class="info-label">
                <i class="el-icon-document"></i>
                <span>病历ID</span>
              </div>
              <div class="info-value">{{ val.prescription }}</div>
            </div>

            <div class="info-row highlight">
              <div class="info-label">
                <i class="el-icon-user"></i>
                <span>病人</span>
              </div>
              <div class="info-value patient-name">{{ getPatientName(val.patient) }}</div>
            </div>

            <div class="info-row secondary">
              <div class="info-label">
                <i class="el-icon-postcard"></i>
                <span>病人ID</span>
              </div>
              <div class="info-value small">{{ val.patient }}</div>
            </div>

            <div class="info-row">
              <div class="info-label">
                <i class="el-icon-medicine-box"></i>
                <span>药品名</span>
              </div>
              <div class="info-value drug-name">{{ val.Name }}</div>
            </div>

            <div class="info-row">
              <div class="info-label">
                <i class="el-icon-goods"></i>
                <span>数量</span>
              </div>
              <div class="info-value amount">{{ val.amount }} 份</div>
            </div>

            <div class="info-row highlight" v-if="val.drug_store">
              <div class="info-label">
                <i class="el-icon-shop"></i>
                <span>药店</span>
              </div>
              <div class="info-value drugstore-name">{{ getDrugstoreName(val.drug_store) }}</div>
            </div>

            <div class="info-row time">
              <div class="info-label">
                <i class="el-icon-time"></i>
                <span>创建时间</span>
              </div>
              <div class="info-value">{{ val.created }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { queryAccountList } from '@/api/accountV2'
import { queryDrugOrderList } from '@/api/drugOrder'

export default {
  name: 'DrugOrder',
  data() {
    return {
      loading: true,
      drugOrderList: [],
      accountList: [],
      accountMap: {} // 账户ID到账户名的映射
    }
  },
  computed: {
    ...mapGetters([
      'account_id',
      'roles',
      'account_name',
    ])
  },
  created() {
    console.log('=== 开始加载药品订单列表 ===')
    console.log('当前角色:', this.roles[0])
    console.log('账户ID:', this.account_id)
    
    // 先加载账户列表
    this.loadAccountList().then(() => {
      // 再加载订单列表
      this.loadDrugOrders()
    })
  },
  methods: {
    // 加载账户列表
    loadAccountList() {
      console.log('=== 加载账户列表 ===')
      return queryAccountList().then(response => {
        console.log('账户列表响应:', response)
        
        let accounts = []
        if (response && response.code === 200 && response.data) {
          accounts = Array.isArray(response.data) ? response.data : []
        } else if (Array.isArray(response)) {
          accounts = response
        } else if (response && typeof response === 'object') {
          accounts = Object.values(response)
        }
        
        this.accountList = accounts
        
        // 创建账户ID到账户名的映射
        this.accountMap = {}
        accounts.forEach(account => {
          if (account && account.account_id && account.account_name) {
            this.accountMap[account.account_id] = account.account_name
          }
        })
        
        console.log('✅ 成功加载账户数量:', accounts.length)
        console.log('账户映射:', this.accountMap)
      }).catch(err => {
        console.error('❌ 加载账户列表失败:', err)
      })
    },
    
    // 加载订单列表
    loadDrugOrders() {
      console.log('=== 加载订单列表 ===')
      
      if (this.roles[0] === 'admin' || this.roles[0] === 'drugstore') {
        queryDrugOrderList().then(response => {
          console.log('药品订单原始响应:', response)
          console.log('响应类型:', typeof response)
          console.log('是否为数组:', Array.isArray(response))
          
          // 正确处理 API 响应格式
          let orders = []
          if (response && response.code === 200) {
            console.log('响应格式: {code: 200, data: ...}')
            if (response.data) {
              orders = Array.isArray(response.data) ? response.data : []
              console.log('从 response.data 提取订单:', orders.length, '条')
            } else {
              console.warn('response.data 为空')
            }
          } else if (Array.isArray(response)) {
            console.log('响应格式: 直接数组')
            orders = response
            console.log('直接使用数组:', orders.length, '条')
          } else if (response !== null && typeof response === 'object') {
            console.log('响应格式: 对象（尝试作为数组）')
            orders = Object.values(response)
            console.log('转换为数组:', orders.length, '条')
          } else {
            console.warn('无法识别的响应格式:', response)
          }
          
          if (!Array.isArray(orders)) {
            console.error('❌ 订单数据不是数组格式:', orders)
            this.drugOrderList = []
            this.loading = false
            this.$message.warning('订单数据格式错误，请联系管理员')
            return
          }
          
          if (orders.length === 0) {
            console.warn('⚠️ 没有订单数据')
          } else {
            console.log('✅ 成功获取', orders.length, '条订单')
            console.log('第一条订单示例:', orders[0])
          }
          
          this.drugOrderList = orders
          this.loading = false
        }).catch(err => {
          console.error('获取订单列表失败:', err)
          this.loading = false
        })
      } else {
        // 病人查看自己的订单
        queryDrugOrderList({ patient: this.account_id }).then(response => {
          console.log('病人订单响应:', response)
          
          let orders = []
          if (response && response.code === 200 && response.data) {
            orders = Array.isArray(response.data) ? response.data : []
          } else if (Array.isArray(response)) {
            orders = response
          } else if (response !== null) {
            orders = []
          }
          
          this.drugOrderList = Array.isArray(orders) ? orders : []
          this.loading = false
        }).catch(err => {
          console.error('获取订单列表失败:', err)
          this.loading = false
        })
      }
    },
    
    // 根据病人ID获取病人名字
    getPatientName(patientId) {
      if (!patientId) {
        return '未知'
      }
      
      const accountName = this.accountMap[patientId]
      if (accountName) {
        // 去掉"病人-"前缀，只显示名字
        return accountName.replace(/^病人-/, '').replace(/病人$/, '')
      }
      
      // 如果找不到，返回ID的前8位
      return patientId.substring(0, 8) + '...'
    },

    // 根据药店ID获取药店名字
    getDrugstoreName(drugstoreId) {
      if (!drugstoreId) {
        return '未指定'
      }
      
      const accountName = this.accountMap[drugstoreId]
      if (accountName) {
        // 去掉"药店-"前缀，只显示名字
        return accountName.replace(/^药店-/, '').replace(/药店$/, '')
      }
      
      // 如果找不到，返回ID的前8位
      return drugstoreId.substring(0, 8) + '...'
    }
  }
}

</script>

<style scoped>
  .container{
    width: 100%;
    min-height: 100%;
    overflow: hidden;
    font-size: 15px;
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

  .drugOrder-card {
    width: 320px;
    margin: 18px;
    border-radius: 16px;
    transition: all 0.3s;
    border: none;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
    overflow: hidden;
  }

  .drugOrder-card:hover {
    transform: translateY(-8px);
    box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
  }

  .drugOrder-card >>> .el-card__header {
    padding: 0;
    border-bottom: none;
  }

  .card-header {
    background: linear-gradient(135deg, #409EFF 0%, #66B1FF 100%);
    padding: 20px;
    color: white;
  }

  .header-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 16px;
    font-weight: 600;
    margin-bottom: 8px;
  }

  .header-title i {
    font-size: 20px;
  }

  .order-id {
    font-size: 13px;
    opacity: 0.95;
    font-family: 'Courier New', monospace;
    word-break: break-all;
  }

  .card-body {
    padding: 20px;
  }

  .info-row {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 16px;
    padding-bottom: 12px;
    border-bottom: 1px solid #f0f0f0;
  }

  .info-row:last-child {
    margin-bottom: 0;
    padding-bottom: 0;
    border-bottom: none;
  }

  .info-row.highlight {
    background: linear-gradient(135deg, #409EFF10 0%, #66B1FF10 100%);
    padding: 10px;
    border-radius: 8px;
    border-bottom: none;
  }

  .info-row.secondary {
    opacity: 0.7;
  }

  .info-row.time {
    margin-top: 8px;
    padding-top: 12px;
    border-top: 1px dashed #e0e0e0;
    border-bottom: none;
  }

  .info-label {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 13px;
    color: #606266;
    font-weight: 500;
    min-width: 80px;
  }

  .info-label i {
    font-size: 16px;
    color: #409EFF;
  }

  .info-value {
    flex: 1;
    text-align: right;
    font-size: 13px;
    color: #303133;
    word-break: break-all;
  }

  .info-value.small {
    font-size: 11px;
    color: #909399;
    font-family: 'Courier New', monospace;
  }

  .info-value.patient-name {
    font-weight: 600;
    font-size: 15px;
    color: #409EFF;
  }

  .info-value.drug-name {
    font-weight: 600;
    color: #409EFF;
  }

  .info-value.drugstore-name {
    font-weight: 600;
    color: #67C23A;
  }

  .info-value.amount {
    font-weight: 600;
    color: #E6A23C;
  }
</style>
