<template>
  <div class="container">
    <el-alert class="info-alert" type="success">
      <p>账户ID: {{ account_id }}</p>
      <p>用户名: {{ account_name }}</p>
    </el-alert>
    
    <div v-if="drugOrderList.length === 0 && !loading" style="text-align: center; padding: 40px;">
      <el-empty description="暂无药品订单">
        <el-button type="primary" @click="$router.push('/prescription/mine')">
          <i class="el-icon-document"></i>
          查看我的病历
        </el-button>
      </el-empty>
    </div>
    
    <el-row v-loading="loading" :gutter="20">
      <el-col v-for="(val, index) in drugOrderList" :key="index" :xs="24" :sm="12" :md="8" :lg="6">
        <el-card class="drugOrder-card" shadow="hover">
          <div slot="header" class="clearfix">
            药品订单ID:
            <span class="highlight-text">{{ val.id }}</span>
          </div>

          <div class="item">
            <el-tag>病历ID: </el-tag>
            <span style="margin-left: 5px;">{{ val.prescription }}</span>
          </div>
          <div class="item">
            <el-tag type="success">药店: </el-tag>
            <span style="margin-left: 5px; font-weight: 600;">{{ getDrugstoreName(val.drug_store) }}</span>
          </div>
          <div class="item">
            <el-tag type="warning">药品名: </el-tag>
            <span style="margin-left: 5px; font-weight: 600;">{{ val.Name }}</span>
          </div>
          <div class="item">
            <el-tag type="warning">药品数量: </el-tag>
            <span style="margin-left: 5px;">{{ val.amount }} 份</span>
          </div>
          <div class="item">
            <el-tag type="info">创建时间: </el-tag>
            <span style="margin-left: 5px;">{{ val.created }}</span>
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
    console.log('=== 病人端：开始加载我的药品订单 ===')
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
      
      if (this.roles[0] === 'admin') {
        // 管理员查看所有订单
        queryDrugOrderList().then(response => {
          console.log('管理员订单响应:', response)
          
          let orders = []
          if (response && response.code === 200 && response.data) {
            orders = Array.isArray(response.data) ? response.data : []
          } else if (Array.isArray(response)) {
            orders = response
          } else if (response && typeof response === 'object') {
            orders = Object.values(response)
          }
          
          this.drugOrderList = Array.isArray(orders) ? orders : []
          console.log('✅ 管理员订单数量:', this.drugOrderList.length)
          this.loading = false
        }).catch(err => {
          console.error('❌ 获取订单列表失败:', err)
          this.loading = false
        })
      } else {
        // 病人查看自己的订单
        queryDrugOrderList({ patient: this.account_id }).then(response => {
          console.log('病人订单原始响应:', response)
          console.log('响应类型:', typeof response)
          console.log('是否为数组:', Array.isArray(response))
          
          let orders = []
          
          if (!response) {
            console.warn('⚠️ 响应为空')
          } else if (response.code === 200) {
            console.log('✓ 响应格式: {code: 200, data: ...}')
            if (response.data) {
              if (Array.isArray(response.data)) {
                orders = response.data
              } else if (typeof response.data === 'object') {
                orders = Object.values(response.data)
              }
            }
          } else if (Array.isArray(response)) {
            console.log('✓ 响应格式: 直接数组')
            orders = response
          } else if (typeof response === 'object') {
            console.log('✓ 响应格式: 对象')
            orders = Object.values(response)
          }
          
          this.drugOrderList = Array.isArray(orders) ? orders : []
          console.log('✅ 病人订单数量:', this.drugOrderList.length)
          
          if (this.drugOrderList.length > 0) {
            console.log('订单列表:', this.drugOrderList)
          } else {
            console.warn('⚠️ 没有找到订单')
          }
          
          this.loading = false
        }).catch(err => {
          console.error('❌ 获取订单列表失败:', err)
          this.loading = false
        })
      }
    },
    
    // 根据药店ID获取药店名字
    getDrugstoreName(drugstoreId) {
      if (!drugstoreId) {
        return '未知药店'
      }
      
      const accountName = this.accountMap[drugstoreId]
      
      if (accountName) {
        // 如果账户名包含"药店"，去掉前缀和后缀
        if (/药店/.test(accountName)) {
          return accountName.replace(/^药店-/, '').replace(/-药店$/, '')
        } else {
          // 如果不包含"药店"，直接返回账户名
          return accountName
        }
      }
      
      // 如果找不到账户名，返回"未知药店"
      console.warn('⚠️ 找不到药店账户，ID:', drugstoreId, '账户映射表:', this.accountMap)
      return '未知药店'
    }
  }
}

</script>

<style scoped>
  .container {
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

  .tag {
    float: left;
  }

  .item {
    font-size: 14px;
    margin-bottom: 18px;
    color: #606266;
  }

  .clearfix:before,
  .clearfix:after {
    display: table;
  }
  
  .clearfix:after {
    clear: both;
  }

  .highlight-text {
    color: var(--theme-primary, #409EFF);
    font-weight: 600;
  }

  .drugOrder-card {
    width: 280px;
    height: 330px;
    margin: 18px;
    border-radius: 12px;
    transition: all 0.3s;
    border: 1px solid #EBEEF5;
  }

  .drugOrder-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
    border-color: var(--theme-primary, #409EFF);
  }

  .drugOrder-card >>> .el-card__header {
    background: var(--theme-light-bg, #ecf5ff);
    border-bottom: 2px solid var(--theme-primary, #409EFF);
    font-weight: 600;
  }

  /* 响应式设计 */
  @media (max-width: 768px) {
    .container {
      padding: 10px;
    }
    
    .drugOrder-card {
      width: 100%;
      margin: 10px 0;
    }
  }
</style>
