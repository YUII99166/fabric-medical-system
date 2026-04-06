<template>
  <div class="app-container">
    <el-form ref="ruleForm" v-loading="loading" :model="ruleForm" :rules="rules" label-width="100px">

      <el-form-item label="查询患者" prop="searchName">
        <el-input 
          v-model="searchName" 
          placeholder="请输入患者姓名或用户名" 
          style="width: 300px"
        >
          <el-button slot="append" icon="el-icon-search" @click="searchPatient">查询</el-button>
        </el-input>
        <div style="margin-top: 10px; color: #909399; font-size: 12px;">
          已加载账户数: {{ accountList.length }} | 病人账户数: {{ accountList.filter(item => item.role === '病人').length }}
        </div>
      </el-form-item>

      <el-form-item label="选择患者" prop="patient" v-if="searchResults.length > 0">
        <el-select 
          v-model="ruleForm.patient" 
          placeholder="请从查询结果中选择患者" 
          style="width: 400px"
          @change="selectPatient"
        >
          <el-option
            v-for="item in searchResults"
            :key="item.account_id"
            :label="`${item.account_name} (${item.username}) - ID: ${item.account_id}`"
            :value="item.account_id"
          />
        </el-select>
      </el-form-item>

      <el-form-item v-if="patientInfo">
        <el-card shadow="hover" style="width: 400px;">
          <div slot="header">
            <i class="el-icon-user"></i>
            <span style="margin-left: 10px;">已选择患者</span>
          </div>
          <div>
            <p><strong>姓名：</strong>{{ patientInfo.account_name }}</p>
            <p><strong>用户名：</strong>{{ patientInfo.username }}</p>
            <p><strong>账户ID：</strong>{{ patientInfo.account_id }}</p>
          </div>
        </el-card>
      </el-form-item>

      <el-form-item v-if="patientError">
        <el-alert
          :title="patientError"
          type="error"
          :closable="false"
          show-icon
        />
      </el-form-item>

      <el-divider content-position="left">病历信息</el-divider>

      <el-form-item label="主诉" prop="chief_complaint">
        <el-input 
          v-model="ruleForm.chief_complaint" 
          type="textarea"
          :rows="2"
          placeholder="患者主要症状描述，如：发热3天，咳嗽伴咳痰2天"
          style="width: 500px"
        />
      </el-form-item>

      <el-form-item label="现病史" prop="present_illness">
        <el-input 
          v-model="ruleForm.present_illness" 
          type="textarea"
          :rows="3"
          placeholder="本次疾病的发展过程"
          style="width: 500px"
        />
      </el-form-item>

      <el-form-item label="体格检查" prop="physical_exam">
        <el-input 
          v-model="ruleForm.physical_exam" 
          type="textarea"
          :rows="2"
          placeholder="体温、血压、心率等，如：T:38.5℃ BP:120/80mmHg HR:88次/分"
          style="width: 500px"
        />
      </el-form-item>

      <el-form-item label="诊断" prop="diagnosis">
        <el-input 
          v-model="ruleForm.diagnosis" 
          placeholder="诊断结果"
          style="width: 500px"
        />
      </el-form-item>

      <el-divider content-position="left">处方信息</el-divider>

      <el-form-item label="药品名" prop="drug_name">
        <el-input v-model="ruleForm.drug_name" placeholder="多个药品用逗号分隔" style="width: 500px" />
      </el-form-item>
      <el-form-item label="药品数量" prop="drug_amount">
        <el-input v-model="ruleForm.drug_amount" placeholder="对应药品数量，用逗号分隔" style="width: 500px" />
      </el-form-item>

      <el-form-item label="医嘱" prop="medical_advice">
        <el-input 
          v-model="ruleForm.medical_advice" 
          type="textarea"
          :rows="2"
          placeholder="用药指导和注意事项"
          style="width: 500px"
        />
      </el-form-item>

      <el-form-item label="备注" prop="comment">
        <el-input 
          v-model="ruleForm.comment" 
          type="textarea"
          :rows="2"
          placeholder="其他备注信息"
          style="width: 500px"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submitForm('ruleForm')">立即创建</el-button>
        <el-button @click="resetForm('ruleForm')">重置</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { queryAccountList } from '@/api/accountV2'
import { createPrescription } from '@/api/prescription'

export default {
  name: 'AddPrescription',
  data() {
    var checkArea = (rule, value, callback) => {
      if (value <= 0) {
        callback(new Error('必须大于0'))
      } else {
        callback()
      }
    }
    return {
      ruleForm: {
        patient: '',
        chief_complaint: '',
        present_illness: '',
        physical_exam: '',
        diagnosis: '',
        drug_name: '',
        drug_amount: '',
        medical_advice: '',
        comment:'',
      },
      searchName: '',
      accountList: [],
      searchResults: [],
      patientInfo: null,
      patientError: '',
      rules: {
        patient: [
          { required: true, message: '请选择患者', trigger: 'change' }
        ],
        diagnosis: [
          { required: true, message: '请输入诊断', trigger: 'blur' }
        ],
        drug_name: [
          { required: true, message: '请输入药品名', trigger: 'blur' }
        ],
        drug_amount: [
          { required: true, message: '请输入药品数量', trigger: 'blur' }
        ]
      },
      loading: false,
    }
  },
  computed: {
    ...mapGetters([
      'account_id'
    ])
  },
  created() {
    // 加载所有账户列表用于验证
    queryAccountList().then(response => {
      console.log('账户列表响应:', response)
      if (response && response.data) {
        this.accountList = response.data
        console.log('加载的账户数量:', this.accountList.length)
        console.log('病人账户:', this.accountList.filter(item => item.role === '病人'))
      } else if (Array.isArray(response)) {
        // 兼容直接返回数组的情况
        this.accountList = response
        console.log('加载的账户数量:', this.accountList.length)
      } else {
        console.error('响应格式错误:', response)
        this.$message({
          type: 'error',
          message: '加载账户列表失败：响应格式错误'
        })
      }
    }).catch(error => {
      console.error('加载账户列表失败:', error)
      this.$message({
        type: 'error',
        message: '加载账户列表失败'
      })
    })
  },
  methods: {
    searchPatient() {
      const keyword = this.searchName.trim()
      
      console.log('开始搜索，关键词:', keyword)
      console.log('当前账户列表数量:', this.accountList.length)
      
      // 清空之前的状态
      this.searchResults = []
      this.patientInfo = null
      this.patientError = ''
      this.ruleForm.patient = ''
      
      if (!keyword) {
        this.patientError = '请输入患者姓名或用户名'
        return
      }
      
      // 从账户列表中搜索病人（支持姓名和用户名模糊匹配）
      this.searchResults = this.accountList.filter(item => {
        console.log('检查账户:', item)
        if (item.role !== '病人') return false
        
        const nameMatch = item.account_name && item.account_name.toLowerCase().includes(keyword.toLowerCase())
        const usernameMatch = item.username && item.username.toLowerCase().includes(keyword.toLowerCase())
        
        console.log(`账户 ${item.username}: nameMatch=${nameMatch}, usernameMatch=${usernameMatch}`)
        
        return nameMatch || usernameMatch
      })
      
      console.log('搜索结果:', this.searchResults)
      
      if (this.searchResults.length === 0) {
        this.patientError = '未找到匹配的患者，请尝试其他关键词'
        this.$message({
          type: 'warning',
          message: '未找到匹配的患者'
        })
      } else {
        this.$message({
          type: 'success',
          message: `找到 ${this.searchResults.length} 位患者`
        })
      }
    },
    selectPatient(accountId) {
      // 从搜索结果中找到选中的患者
      const patient = this.searchResults.find(item => item.account_id === accountId)
      
      if (patient) {
        this.patientInfo = patient
        this.patientError = ''
        this.$message({
          type: 'success',
          message: '患者选择成功'
        })
      }
    },
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          // 检查患者是否已验证
          if (!this.patientInfo) {
            this.$message({
              type: 'warning',
              message: '请先验证患者ID'
            })
            return
          }
          
          this.$confirm('是否立即创建?', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'success'
          }).then(() => {
            this.loading = true
            createPrescription({
              doctor: String(this.account_id),
              patient: String(this.ruleForm.patient),
              chief_complaint: this.ruleForm.chief_complaint,
              present_illness: this.ruleForm.present_illness,
              physical_exam: this.ruleForm.physical_exam,
              diagnosis: this.ruleForm.diagnosis,
              drug_name: this.ruleForm.drug_name,
              drug_amount: this.ruleForm.drug_amount,
              medical_advice: this.ruleForm.medical_advice,
              hospital:'0feceb66ffc1',
              comment: this.ruleForm.comment,
            }).then(response => {
              this.loading = false
              if (response !== null) {
                this.$message({
                  type: 'success',
                  message: '创建成功!'
                })
                // 重置表单
                this.resetForm(formName)
                this.patientInfo = null
                this.patientError = ''
              } else {
                this.$message({
                  type: 'error',
                  message: '创建失败!'
                })
              }
            }).catch(_ => {
              this.loading = false
            })
          }).catch(() => {
            this.loading = false
            this.$message({
              type: 'info',
              message: '已取消创建'
            })
          })
        } else {
          return false
        }
      })
    },
    resetForm(formName) {
      this.$refs[formName].resetFields()
      this.searchName = ''
      this.searchResults = []
      this.patientInfo = null
      this.patientError = ''
    },
    selectGet(account_id) {
      this.ruleForm.patient = account_id
    },
  }
}
</script>

<style scoped>
</style>
