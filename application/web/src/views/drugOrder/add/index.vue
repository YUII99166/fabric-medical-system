<template>
  <div class="app-container">
    <!-- 患者搜索卡片 -->
    <el-card class="search-card" shadow="hover">
      <div slot="header" class="card-header">
        <i class="el-icon-search"></i>
        <span>患者查询</span>
      </div>
      <div class="search-box">
        <el-input
          v-model="searchKey"
          placeholder="请输入患者姓名或用户名，例如：①号病人"
          clearable
          @keyup.enter.native="searchPatient"
        >
          <el-button 
            slot="append" 
            icon="el-icon-search" 
            :loading="searching"
            @click="searchPatient"
          >
            {{ searching ? '搜索中...' : '搜索患者' }}
          </el-button>
        </el-input>
        <div class="search-actions">
          <div class="search-hint">
            <i class="el-icon-info"></i>
            <span>提示：请输入患者姓名，如"①号病人"或"①号"</span>
          </div>
          <el-button 
            type="text" 
            size="small" 
            icon="el-icon-view"
            @click="showAllPatients"
          >
            显示所有患者
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 搜索结果提示 -->
    <div v-if="searched && filteredAccountList.length === 0" class="no-result">
      <el-alert title="未找到匹配的患者" type="warning" :closable="false" />
    </div>

    <!-- 订单创建表单 -->
    <el-card v-if="filteredAccountList.length > 0 || !searched" class="form-card" shadow="hover">
      <div slot="header" class="card-header">
        <i class="el-icon-document-add"></i>
        <span>创建药品订单</span>
      </div>
      
      <el-form ref="ruleForm" v-loading="loading" :model="ruleForm" :rules="rules" label-width="100px">
        <el-form-item label="病人" prop="patient">
          <el-select 
            v-model="ruleForm.patient" 
            placeholder="请先搜索患者，然后选择" 
            filterable
            @change="selectGetPatient"
          >
            <el-option
              v-for="item in filteredAccountList"
              :key="item.account_id"
              :label="item.account_name"
              :value="item.account_id"
            >
              <span style="float: left">{{ item.account_name }}</span>
              <span style="float: right; color: #8492a6; font-size: 13px">{{ item.account_id }}</span>
            </el-option>
          </el-select>
          <div v-if="!searched" class="field-hint">
            <i class="el-icon-warning"></i>
            <span>请先使用上方搜索框查询患者</span>
          </div>
        </el-form-item>

      <el-form-item label="病历" prop="prescription">
        <el-select 
          v-model="ruleForm.prescription" 
          placeholder="请先选择病人" 
          :disabled="!ruleForm.patient" 
          @change="selectGetPrescription"
        >
          <el-option
            v-for="item in prescriptionList"
            :key="item.id"
            :label="`${item.patient_name || '患者'} - ${item.diagnosis}`"
            :value="item.id"
          >
            <div style="display: flex; justify-content: space-between; align-items: center;">
              <span style="font-weight: 600; color: #303133;">
                <i class="el-icon-user" style="margin-right: 4px; color: #409EFF;"></i>
                {{ item.patient_name || '患者' }}
              </span>
              <span style="color: #606266; margin: 0 10px;">|</span>
              <span style="flex: 1; color: #606266;">
                <i class="el-icon-document-checked" style="margin-right: 4px; color: #67C23A;"></i>
                {{ item.diagnosis }}
              </span>
              <span style="color: #909399; font-size: 12px; margin-left: 10px;">
                <i class="el-icon-time" style="margin-right: 4px;"></i>
                {{ item.created }}
              </span>
            </div>
          </el-option>
        </el-select>
        <div v-if="ruleForm.patient && prescriptionList.length === 0" class="field-hint warning">
          <i class="el-icon-warning"></i>
          <span>该患者暂无病历记录</span>
        </div>
      </el-form-item>

      <el-form-item label="药品名" prop="drug_name">
        <el-input v-model="ruleForm.drug_name" placeholder="请输入药品名称" style="width: 197px" />
      </el-form-item>
      
      <el-form-item label="药品数量" prop="drug_amount">
        <el-input-number 
          v-model="ruleForm.drug_amount" 
          :precision="0" 
          :step="1" 
          :min="1" 
          style="width: 197px" 
        />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="submitForm('ruleForm')">立即创建</el-button>
        <el-button @click="resetForm('ruleForm')">重置</el-button>
      </el-form-item>
    </el-form>
    </el-card>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { queryAccountList } from '@/api/accountV2'
import { queryPrescriptionList } from '@/api/prescription'
import { createDrugOrder } from '@/api/drugOrder'

export default {
  name: 'AddDrugOrder',
  data() {
    return {
      searchKey: '',
      searched: false,
      searching: false,
      ruleForm: {
        patient: '',
        prescription: '',
        drug_name: '',
        drug_amount: 1,
        drug_store: '' // 将在 created 中设置为当前登录的药店ID
      },
      accountList: [],
      filteredAccountList: [],
      prescriptionList: [],
      rules: {
        patient: [
          { required: true, message: '请选择病人', trigger: 'change' }
        ],
        drug_name: [
          { required: true, message: '请输入药品名称', trigger: 'blur' }
        ],
        drug_amount: [
          { required: true, message: '请输入药品数量', trigger: 'blur' }
        ]
      },
      loading: false
    }
  },
  computed: {
    ...mapGetters([
      'account_id'
    ])
  },
  created() {
    console.log('=== 药店新增订单页面初始化 ===')
    console.log('当前账户ID:', this.account_id)
    
    // 设置药店ID为当前登录的账户ID
    this.ruleForm.drug_store = this.account_id
    console.log('✓ 药店ID已设置:', this.ruleForm.drug_store)
    
    // 加载所有病人账户（用于搜索）
    this.loadAllPatients()
  },
  methods: {
    // 加载所有病人账户
    loadAllPatients() {
      console.log('=== 开始加载所有病人账户 ===')
      queryAccountList().then(response => {
        console.log('账户列表原始响应:', response)
        console.log('响应类型:', typeof response)
        console.log('是否为数组:', Array.isArray(response))
        
        // 处理不同的响应格式
        let accounts = []
        
        if (!response) {
          console.warn('⚠ 响应为空')
          this.$message.warning('未获取到账户列表')
          return
        }
        
        if (Array.isArray(response)) {
          // 格式1: 直接返回数组 [...]
          console.log('✓ 响应格式: 直接数组')
          accounts = response
        } else if (typeof response === 'object') {
          if (response.code === 200) {
            // 格式2: {code: 200, msg: "成功", data: [...]}
            console.log('✓ 响应格式: {code: 200, data: ...}')
            if (Array.isArray(response.data)) {
              accounts = response.data
            } else if (response.data && typeof response.data === 'object') {
              // 格式3: data 是对象，转换为数组
              console.log('✓ 响应格式: data 是对象，转换为数组')
              accounts = Object.values(response.data)
            } else {
              console.warn('⚠ data 字段不是数组:', response.data)
            }
          } else if (response.data && Array.isArray(response.data)) {
            // 格式4: {data: [...]}（没有 code 字段）
            console.log('✓ 响应格式: {data: ...}')
            accounts = response.data
          } else {
            console.warn('⚠ 未知的响应格式')
            console.log('响应内容:', JSON.stringify(response, null, 2))
          }
        }
        
        console.log('提取到的账户数组:', accounts)
        console.log('账户数组长度:', accounts.length)
        
        if (!Array.isArray(accounts)) {
          console.error('✗ 无法将响应转换为数组')
          this.$message.error('账户列表格式错误')
          return
        }
        
        // 过滤出病人账户
        this.accountList = accounts.filter(item => {
          if (!item || !item.account_name) {
            return false
          }
          const accountName = item.account_name
          // 支持多种病人账户格式：
          // 1. 以"病人"结尾：①号病人、②号病人
          // 2. 以"病人-"开头：病人-张三、病人-李富贵
          // 3. 包含"病人"：病人李四
          return /病人/.test(accountName)
        })
        
        console.log('✓ 成功加载病人账户数量:', this.accountList.length)
        
        if (this.accountList.length === 0) {
          console.warn('⚠ 没有找到病人账户')
          this.$message.info('系统中暂无病人账户')
        } else {
          console.log('病人账户列表:', this.accountList.map(a => ({
            name: a.account_name,
            id: a.account_id
          })))
        }
      }).catch(err => {
        console.error('✗ 加载账户列表失败:', err)
        this.$message.error('加载患者列表失败：' + err.message)
      })
    },

    // 搜索患者
    searchPatient() {
      if (!this.searchKey.trim()) {
        this.$message.warning('请输入患者姓名或用户名')
        return
      }

      console.log('=== 开始搜索患者 ===')
      console.log('搜索关键词:', this.searchKey)
      console.log('可用账户总数:', this.accountList.length)
      console.log('可用账户列表:', this.accountList.map(a => a.account_name))

      this.searching = true
      this.searched = true
      this.filteredAccountList = []

      // 在本地账户列表中搜索
      const keyword = this.searchKey.trim().toLowerCase()
      console.log('处理后的关键词:', keyword)
      
      this.filteredAccountList = this.accountList.filter(item => {
        if (!item) {
          return false
        }
        
        const accountName = (item.account_name || '').toLowerCase()
        const accountId = (item.account_id || '').toLowerCase()
        
        console.log('检查账户:', {
          account_name: accountName,
          account_id: accountId,
          keyword: keyword
        })
        
        // 支持多种搜索方式
        const matchByName = accountName.includes(keyword)
        const matchById = accountId.includes(keyword)
        
        // 去掉"病人-"前缀后再匹配
        const nameWithoutPrefix = accountName.replace(/^病人-/, '')
        const matchByNameWithoutPrefix = nameWithoutPrefix.includes(keyword)
        
        // 去掉"病人"后缀后再匹配（如：①号病人 -> ①号）
        const nameWithoutSuffix = accountName.replace(/病人$/, '')
        const matchByNameWithoutSuffix = nameWithoutSuffix.includes(keyword)
        
        const isMatch = matchByName || matchById || matchByNameWithoutPrefix || matchByNameWithoutSuffix
        
        if (isMatch) {
          console.log('✓ 匹配成功:', item.account_name)
        }
        
        return isMatch
      })

      this.searching = false

      console.log('搜索结果数量:', this.filteredAccountList.length)
      console.log('搜索结果:', this.filteredAccountList.map(a => a.account_name))
      
      if (this.filteredAccountList.length === 0) {
        const availablePatients = this.accountList.map(a => a.account_name).join('、')
        this.$message({
          type: 'info',
          message: `未找到匹配的患者。当前系统中的患者：${availablePatients}`,
          duration: 5000
        })
        console.log('提示：当前系统中的患者:', availablePatients)
      } else {
        this.$message.success(`找到 ${this.filteredAccountList.length} 位患者`)
      }
    },

    // 显示所有患者
    showAllPatients() {
      console.log('=== 显示所有患者 ===')
      this.searched = true
      this.filteredAccountList = [...this.accountList]
      this.searchKey = ''
      
      if (this.filteredAccountList.length > 0) {
        this.$message.success(`共有 ${this.filteredAccountList.length} 位患者`)
      } else {
        this.$message.info('系统中暂无患者')
      }
    },

    // 选择患者后加载病历
    selectGetPatient(account_id) {
      console.log('=== 选择患者 ===')
      console.log('患者ID:', account_id)
      
      this.ruleForm.patient = account_id
      this.prescriptionList = []
      this.ruleForm.prescription = ''

      // 加载该患者的病历列表
      queryPrescriptionList({ patient: account_id }).then(response => {
        console.log('病历列表原始响应:', response)
        
        // 处理不同的响应格式
        let prescriptions = []
        
        if (response && typeof response === 'object') {
          if (response.code === 200 && response.data) {
            // 格式: {code: 200, data: [...]}
            prescriptions = Array.isArray(response.data) ? response.data : []
          } else if (Array.isArray(response)) {
            // 格式: [...]
            prescriptions = response
          } else if (response.data && Array.isArray(response.data)) {
            // 格式: {data: [...]}
            prescriptions = response.data
          }
        }
        
        this.prescriptionList = prescriptions
        console.log('✓ 成功加载病历数量:', this.prescriptionList.length)
        
        if (this.prescriptionList.length === 0) {
          this.$message.info('该患者暂无病历记录')
        }
      }).catch(err => {
        console.error('✗ 加载病历列表失败:', err)
        this.$message.error('加载病历列表失败')
      })
    },

    // 选择病历
    selectGetPrescription(prescription) {
      console.log('选择病历ID:', prescription)
      this.ruleForm.prescription = prescription
    },

    // 提交表单
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.$confirm('是否立即创建药品订单?', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'success'
          }).then(() => {
            this.loading = true
            
            console.log('=== 创建药品订单 ===')
            console.log('订单数据:', {
              patient: this.ruleForm.patient,
              prescription: this.ruleForm.prescription,
              drug_name: this.ruleForm.drug_name,
              drug_amount: this.ruleForm.drug_amount.toString(),
              drug_store: this.ruleForm.drug_store
            })
            
            createDrugOrder({
              patient: this.ruleForm.patient,
              prescription: this.ruleForm.prescription,
              drug_name: this.ruleForm.drug_name,
              drug_amount: this.ruleForm.drug_amount.toString(),
              drug_store: this.ruleForm.drug_store // 使用当前登录的药店ID
            }).then(response => {
              this.loading = false
              console.log('创建订单响应:', response)
              
              if (response !== null) {
                this.$message({
                  type: 'success',
                  message: '创建成功!'
                })
                // 重置表单
                this.resetForm(formName)
                // 跳转到订单列表
                this.$router.push('/drug/list')
              } else {
                this.$message({
                  type: 'error',
                  message: '创建失败!'
                })
              }
            }).catch(err => {
              this.loading = false
              console.error('✗ 创建订单失败:', err)
              this.$message.error('创建失败：' + err.message)
            })
          }).catch(() => {
            this.loading = false
            this.$message({
              type: 'info',
              message: '已取消创建'
            })
          })
        } else {
          this.$message.warning('请完整填写表单')
          return false
        }
      })
    },

    // 重置表单
    resetForm(formName) {
      this.$refs[formName].resetFields()
      this.searchKey = ''
      this.searched = false
      this.filteredAccountList = []
      this.prescriptionList = []
    }
  }
}
</script>

<style scoped>
.app-container {
  padding: 20px;
}

/* 搜索卡片样式 */
.search-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: bold;
  color: #303133;
}

.card-header i {
  margin-right: 8px;
  color: #409EFF;
  font-size: 18px;
}

.search-box {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.search-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-hint {
  font-size: 13px;
  color: #909399;
  display: flex;
  align-items: center;
  gap: 5px;
  padding-left: 5px;
  flex: 1;
}

.search-hint i {
  color: #409EFF;
}

.search-box >>> .el-input-group__append {
  background-color: #409EFF;
  border-color: #409EFF;
  transition: all 0.3s;
}

.search-box >>> .el-input-group__append:hover {
  background-color: #66b1ff;
  border-color: #66b1ff;
}

.search-box >>> .el-input-group__append .el-button {
  color: white;
  font-weight: 500;
}

.search-box >>> .el-input-group__append .el-button:hover {
  color: white;
}

/* 无结果提示 */
.no-result {
  margin-bottom: 20px;
}

/* 表单卡片样式 */
.form-card {
  margin-bottom: 20px;
}

/* 表单字段提示 */
.field-hint {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
  display: flex;
  align-items: center;
  gap: 5px;
}

.field-hint i {
  color: #E6A23C;
}

.field-hint.warning {
  color: #E6A23C;
}

/* 表单样式优化 */
.el-form-item {
  margin-bottom: 22px;
}

.el-select {
  width: 100%;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .app-container {
    padding: 10px;
  }
  
  .search-box >>> .el-input-group__append {
    padding: 0 10px;
  }
}
</style>
