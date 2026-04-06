<template>
  <div class="container">
    <el-alert class="info-alert" type="success">
      <p>账户ID: {{ account_id }}</p>
      <p>用户名: {{ account_name }}</p>
    </el-alert>

    <!-- 搜索框 -->
    <el-card class="search-card" shadow="hover">
      <div class="search-box">
        <el-input
          v-model="searchKey"
          placeholder="请输入患者姓名，例如：周杰伦 或 病人-周杰伦"
          clearable
          @keyup.enter.native="searchPatient"
        >
          <el-button 
            slot="append" 
            icon="el-icon-search" 
            :loading="searching"
            @click="searchPatient"
          >
            {{ searching ? '搜索中...' : '搜索' }}
          </el-button>
        </el-input>
        <div class="search-hint">
          <i class="el-icon-info"></i>
          <span>提示：请输入患者姓名，例如"周杰伦"或"病人-周杰伦"，系统会自动搜索相关病历</span>
        </div>
      </div>
    </el-card>

    <!-- 搜索结果 -->
    <div v-if="searched && prescriptionList.length === 0" style="text-align: center; margin-top: 20px;">
      <el-alert title="未找到该患者的病历记录" type="warning" />
    </div>

    <el-row v-loading="loading" :gutter="20" style="margin-top: 20px;">
      <el-col v-for="(val, index) in prescriptionList" :key="index" :xs="24" :sm="12" :md="8" :lg="6">
        <el-card class="prescription-card" shadow="hover">
          <!-- 标签：本院/他院 -->
          <div class="card-tag">
            <el-tag v-if="val.is_same_org" type="success" size="small">本院</el-tag>
            <el-tag v-else type="warning" size="small">他院</el-tag>
          </div>

          <div class="card-icon">
            <i class="el-icon-document"></i>
          </div>

          <div class="card-title">病历记录</div>

          <div class="card-info">
            <div class="info-item">
              <i class="el-icon-user"></i>
              <span>{{ val.patient_name }}</span>
            </div>
            <div class="info-item">
              <i class="el-icon-hospital"></i>
              <span>{{ val.hospital_name }}</span>
            </div>
            <!-- 诊断结果 - 重点显示 -->
            <div class="info-item diagnosis-item">
              <i class="el-icon-document-checked"></i>
              <div class="diagnosis-content">
                <div class="diagnosis-label">诊断：</div>
                <div class="diagnosis-value">{{ val.diagnosis || '未填写' }}</div>
              </div>
            </div>
            <div class="info-item time">
              <i class="el-icon-time"></i>
              <span>{{ val.created }}</span>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="card-action">
            <el-button
              v-if="val.can_view"
              type="primary"
              size="small"
              @click="openDialog(val)"
            >
              查看详情
            </el-button>
            <el-button
              v-else
              type="warning"
              size="small"
              @click="openRequestDialog(val)"
            >
              申请授权
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 病历详情对话框 -->
    <el-dialog
      title="病历详情"
      :visible.sync="dialogVisible"
      width="700px"
      :close-on-click-modal="false"
    >
      <div v-if="selectedPrescription" class="detail-content">
        <!-- 患者信息 -->
        <div class="info-section">
          <div class="section-header">
            <i class="el-icon-user"></i>
            <span>患者信息</span>
          </div>
          <div class="section-text">
            <p><strong>姓名：</strong>{{ selectedPrescription.patient_name }}</p>
            <p><strong>就诊医院：</strong>{{ selectedPrescription.hospital_name }}</p>
            <p><strong>就诊科室：</strong>{{ selectedPrescription.department }}</p>
            <p><strong>就诊时间：</strong>{{ selectedPrescription.created }}</p>
          </div>
        </div>

        <!-- 医生信息 -->
        <div class="info-section">
          <div class="section-header">
            <i class="el-icon-user-solid"></i>
            <span>医生信息</span>
          </div>
          <div class="section-text">
            <p><strong>医生：</strong>{{ selectedPrescription.doctor_name }} {{ selectedPrescription.doctor_title }}</p>
          </div>
        </div>

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

        <!-- 诊断结果 -->
        <div class="diagnosis-highlight">
          <div class="diagnosis-header">
            <i class="el-icon-document-checked"></i>
            <span>诊断结果</span>
          </div>
          <div class="diagnosis-text">{{ selectedPrescription.diagnosis }}</div>
        </div>

        <!-- 处方药品 -->
        <div class="drug-highlight">
          <div class="drug-header">
            <i class="el-icon-medicine-box"></i>
            <span>处方药品</span>
          </div>
          <div class="drug-list-compact">
            <div v-for="(drug, idx) in selectedPrescription.drug" :key="idx" class="drug-item-compact">
              <span class="drug-name">{{ drug.Name }}</span>
              <span class="drug-amount">× {{ drug.amount }}</span>
            </div>
          </div>
        </div>

        <!-- 医嘱 -->
        <div class="info-section">
          <div class="section-header">
            <i class="el-icon-document"></i>
            <span>医嘱</span>
          </div>
          <div class="section-text">{{ selectedPrescription.medical_advice || '未填写' }}</div>
        </div>

        <!-- 备注 -->
        <div class="info-section">
          <div class="section-header">
            <i class="el-icon-edit-outline"></i>
            <span>备注</span>
          </div>
          <div class="section-text">{{ selectedPrescription.comment || '无' }}</div>
        </div>
      </div>
    </el-dialog>

    <!-- 申请授权对话框 -->
    <el-dialog
      title="申请查看授权"
      :visible.sync="requestDialogVisible"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="requestForm" label-width="100px">
        <el-form-item label="患者姓名">
          <el-input v-model="requestForm.patientName" disabled />
        </el-form-item>
        <el-form-item label="就诊医院">
          <el-input v-model="requestForm.hospitalName" disabled />
        </el-form-item>
        <el-form-item label="申请理由" required>
          <el-input
            v-model="requestForm.reason"
            type="textarea"
            :rows="4"
            placeholder="请说明申请查看该病历的理由"
          />
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="requestDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitRequest">提交申请</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { queryPrescriptionsByPatient, requestAccess } from '@/api/access'
import { mapGetters } from 'vuex'

export default {
  name: 'PatientSearch',
  data() {
    return {
      searchKey: '',
      searched: false,
      loading: false,
      searching: false,
      prescriptionList: [],
      dialogVisible: false,
      selectedPrescription: null,
      requestDialogVisible: false,
      requestForm: {
        prescriptionId: '',
        patientName: '',
        hospitalName: '',
        reason: ''
      }
    }
  },
  computed: {
    ...mapGetters([
      'account_id',
      'account_name',
      'roles'
    ])
  },
  methods: {
    searchPatient() {
      if (!this.searchKey.trim()) {
        this.$message.warning('请输入患者姓名')
        return
      }

      // 检查account_id是否有效
      let doctorId = this.account_id
      if (!doctorId || doctorId === '未同步到区块链') {
        // 尝试从sessionStorage获取
        const userInfo = JSON.parse(sessionStorage.getItem('userInfo') || '{}')
        if (userInfo.account_id && userInfo.account_id !== '未同步到区块链') {
          doctorId = userInfo.account_id
          // 同时更新store
          this.$store.commit('account/SET_ACCOUNTID', userInfo.account_id)
        } else {
          this.$message.error('无法获取医生ID，请重新登录')
          console.error('医生ID为空或未同步，无法搜索')
          return
        }
      }

      this.searching = true
      this.searched = true
      this.prescriptionList = []

      // 准备搜索关键词列表
      const searchKeys = []
      const originalKey = this.searchKey.trim()
      
      // 添加原始搜索词
      searchKeys.push(originalKey)
      
      // 如果不包含"病人-"前缀，添加带前缀的版本
      if (!originalKey.startsWith('病人-')) {
        searchKeys.push('病人-' + originalKey)
      }
      
      // 如果包含"病人-"前缀，添加不带前缀的版本
      if (originalKey.startsWith('病人-')) {
        searchKeys.push(originalKey.substring(3))
      }

      console.log('开始搜索患者:', originalKey)
      console.log('医生ID:', doctorId)
      console.log('搜索关键词列表:', searchKeys)

      // 依次尝试不同的搜索关键词
      this.trySearchWithKeys(searchKeys, 0, doctorId)
    },

    // 尝试使用不同关键词搜索（带静默重试）
    async trySearchWithKeys(searchKeys, index, doctorId, retryCount = 0) {
      const maxRetries = 3
      
      if (index >= searchKeys.length) {
        // 所有搜索都失败了
        this.searching = false
        this.$message.info('未找到该患者的病历记录')
        return
      }

      const currentKey = searchKeys[index]
      if (retryCount === 0) {
        console.log(`尝试搜索关键词 ${index + 1}/${searchKeys.length}: "${currentKey}"`)
      }

      try {
        const res = await queryPrescriptionsByPatient({
          search_key: currentKey,
          doctor_id: doctorId
        })
        
        console.log(`搜索关键词 "${currentKey}" 的响应:`, res)
        
        if (res.code === 200 && res.data && res.data.length > 0) {
          // 找到了结果
          this.searching = false
          this.prescriptionList = res.data
          console.log('找到病历数量:', this.prescriptionList.length)
          this.$message.success(`找到 ${this.prescriptionList.length} 条病历记录`)
        } else {
          // 当前关键词没找到，尝试下一个
          console.log(`搜索关键词 "${currentKey}" 未找到结果，尝试下一个`)
          this.trySearchWithKeys(searchKeys, index + 1, doctorId, 0)
        }
      } catch (err) {
        // 只在控制台记录错误
        if (retryCount === 0) {
          console.error(`搜索关键词 "${currentKey}" 出错:`, err)
        }
        
        // 如果失败且还有重试次数，静默重试当前关键词
        if (retryCount < maxRetries) {
          console.log(`搜索关键词 "${currentKey}" 出错，正在静默重试 (${retryCount + 1}/${maxRetries})...`)
          await new Promise(resolve => setTimeout(resolve, 500 * (retryCount + 1)))
          return this.trySearchWithKeys(searchKeys, index, doctorId, retryCount + 1)
        } else {
          // 当前关键词所有重试都失败，尝试下一个关键词
          console.log(`搜索关键词 "${currentKey}" 所有重试都失败，尝试下一个`)
          this.trySearchWithKeys(searchKeys, index + 1, doctorId, 0)
        }
      }
    },
    openDialog(prescription) {
      // 跳转到病历详情页面
      this.$router.push({
        name: 'PrescriptionDetail',
        query: { id: prescription.id }
      })
    },
    openRequestDialog(prescription) {
      this.requestForm = {
        prescriptionId: prescription.id,
        patientName: prescription.patient_name,
        hospitalName: prescription.hospital_name,
        reason: ''
      }
      this.requestDialogVisible = true
    },
    // 提交授权申请（带静默重试）
    async submitRequest(retryCount = 0) {
      const maxRetries = 3
      
      if (!this.requestForm.reason.trim()) {
        this.$message.warning('请填写申请理由')
        return
      }

      try {
        const res = await requestAccess({
          prescription_id: this.requestForm.prescriptionId,
          doctor_id: this.account_id,
          reason: this.requestForm.reason
        })
        
        if (res.code === 200) {
          this.$message.success('申请已提交，请等待患者审批')
          this.requestDialogVisible = false
          this.searchPatient() // 刷新列表
        } else {
          // 如果失败且还有重试次数，静默重试
          if (retryCount < maxRetries) {
            console.log(`提交申请失败，正在静默重试 (${retryCount + 1}/${maxRetries})...`)
            await new Promise(resolve => setTimeout(resolve, 500 * (retryCount + 1)))
            return this.submitRequest(retryCount + 1)
          } else {
            // 所有重试都失败后才显示错误
            this.$message.error(res.msg || '申请失败')
          }
        }
      } catch (err) {
        // 只在控制台记录错误
        if (retryCount === 0) {
          console.error('提交申请失败:', err)
        }
        
        // 如果失败且还有重试次数，静默重试
        if (retryCount < maxRetries) {
          console.log(`提交申请失败，正在静默重试 (${retryCount + 1}/${maxRetries})...`)
          await new Promise(resolve => setTimeout(resolve, 500 * (retryCount + 1)))
          return this.submitRequest(retryCount + 1)
        } else {
          // 所有重试都失败后才显示错误
          this.$message.error('申请失败：' + err.message)
        }
      }
    }
  }
}
</script>

<style scoped>
.container {
  padding: 20px;
}

.info-alert {
  margin-bottom: 20px;
}

.search-card {
  margin-bottom: 20px;
}

.search-box {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.search-hint {
  font-size: 13px;
  color: #909399;
  display: flex;
  align-items: center;
  gap: 5px;
  padding-left: 5px;
}

.search-hint i {
  color: #409EFF;
}

.search-box >>> .el-input-group__append {
  background-color: var(--theme-primary, #409EFF);
  border-color: var(--theme-primary, #409EFF);
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

.prescription-card {
  margin-bottom: 20px;
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
}

.prescription-card:hover {
  transform: translateY(-5px);
}

.card-tag {
  position: absolute;
  top: 10px;
  right: 10px;
}

.card-icon {
  text-align: center;
  font-size: 48px;
  color: #409EFF;
  margin: 20px 0;
}

.card-title {
  text-align: center;
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 15px;
  color: #303133;
}

.card-info {
  padding: 0 10px;
}

.info-item {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
  color: #606266;
}

.info-item i {
  margin-right: 8px;
  color: #409EFF;
}

.info-item.time {
  font-size: 12px;
  color: #909399;
}

.info-item.diagnosis-item {
  background: linear-gradient(135deg, #667eea15 0%, #764ba215 100%);
  padding: 12px;
  border-radius: 8px;
  margin: 10px 0;
  display: block;
}

.diagnosis-content {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.diagnosis-label {
  font-size: 12px;
  color: #909399;
  font-weight: 500;
}

.diagnosis-value {
  font-size: 14px;
  color: #303133;
  font-weight: 600;
  line-height: 1.6;
}

.card-action {
  text-align: center;
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #EBEEF5;
}

/* 详情对话框样式 */
.detail-content {
  max-height: 600px;
  overflow-y: auto;
}

.info-section {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.section-header {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 10px;
}

.section-header i {
  margin-right: 8px;
  color: #409EFF;
}

.section-text {
  color: #606266;
  line-height: 1.8;
}

.section-text p {
  margin: 5px 0;
}

.diagnosis-highlight {
  margin-bottom: 20px;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  color: white;
}

.diagnosis-header {
  display: flex;
  align-items: center;
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 12px;
}

.diagnosis-header i {
  margin-right: 8px;
}

.diagnosis-text {
  font-size: 16px;
  line-height: 1.8;
}

.drug-highlight {
  margin-bottom: 20px;
  padding: 20px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  border-radius: 8px;
  color: white;
}

.drug-header {
  display: flex;
  align-items: center;
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 12px;
}

.drug-header i {
  margin-right: 8px;
}

.drug-list-compact {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.drug-item-compact {
  background-color: rgba(255, 255, 255, 0.2);
  padding: 8px 15px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.drug-name {
  font-weight: 500;
}

.drug-amount {
  font-size: 14px;
  opacity: 0.9;
}
</style>
