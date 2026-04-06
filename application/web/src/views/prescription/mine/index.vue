<template>
  <div class="container">
    <el-alert
      class="info-alert"
      type="success"
    >
      <p>账户ID: {{ account_id }}</p>
      <p>用户名: {{ account_name }}</p>
    </el-alert>
    
    <div v-if="prescriptionList.length==0" style="text-align: center;">
      <el-alert
        title="查询不到数据"
        type="warning"
      />
    </div>
    
    <el-row v-loading="loading" :gutter="20">
      <el-col v-for="(val,index) in prescriptionList" :key="index" :xs="24" :sm="12" :md="8" :lg="6">
        <el-card class="prescription-card" shadow="hover" @click.native="showQuickView(val)">
          <div class="card-icon">
            <i class="el-icon-document"></i>
          </div>
          
          <div class="card-title">病历记录</div>
          
          <div class="card-info">
            <div class="info-item">
              <i class="el-icon-user"></i>
              <span>{{ val.diagnosis }}</span>
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

    <!-- 详情对话框 -->
    <el-dialog
      title="病历详情"
      :visible.sync="dialogVisible"
      width="700px"
      :close-on-click-modal="false"
    >
      <div v-if="selectedPrescription" class="detail-content">
        <!-- 病历综合描述 -->
        <div class="summary-section">
          <div class="summary-header">
            <i class="el-icon-document"></i>
            <span>病历摘要</span>
          </div>
          <div class="summary-text">
            您于<strong>{{ selectedPrescription.created }}</strong>在<strong>{{ selectedPrescription.hospital_name || selectedPrescription.organization_name }}</strong>就诊，由<strong>{{ selectedPrescription.doctor_name }}</strong>医生（{{ selectedPrescription.doctor_title || '医师' }}）接诊。主诉：{{ selectedPrescription.chief_complaint || '未填写' }}。现病史：{{ selectedPrescription.present_illness || '未填写' }}。体格检查：{{ selectedPrescription.physical_exam || '未填写' }}。经诊断为<strong class="diagnosis-text-inline">{{ selectedPrescription.diagnosis }}</strong>。<span v-if="selectedPrescription.drug && selectedPrescription.drug.length > 0">处方药品包括：<span v-for="(drug, index) in selectedPrescription.drug" :key="index"><strong>{{ drug.Name }}</strong>（{{ drug.amount }}）<span v-if="index < selectedPrescription.drug.length - 1">、</span></span>。</span>医嘱：{{ selectedPrescription.medical_advice || '未填写' }}。<span v-if="selectedPrescription.comment">备注：{{ selectedPrescription.comment }}</span>
          </div>
        </div>

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
import { queryAccountList } from '@/api/accountV2'
import { queryPrescriptionList } from '@/api/prescription'

export default {
  name: 'MyPrescription',
  data() {
    return {
      loading: true,
      prescriptionList: [],
      dialogVisible: false,
      selectedPrescription: null,
      accountMap: {} // 账户ID到账户名称的映射
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
    // 先加载账户列表，建立映射关系
    queryAccountList().then(response => {
      console.log('账户列表响应:', response)
      // 处理标准响应格式
      const accountList = response && response.code === 200 ? response.data : (Array.isArray(response) ? response : [])
      if (accountList && accountList.length > 0) {
        // 建立账户ID到账户名称的映射
        accountList.forEach(account => {
          this.accountMap[account.account_id] = {
            name: account.account_name,
            username: account.username,
            role: account.role
          }
        })
      }
    })

    // 加载病历列表
    if (this.roles[0] === 'doctor') {
      queryPrescriptionList().then(response => {
        console.log('病历列表响应:', response)
        // 处理标准响应格式
        const prescriptionList = response && response.code === 200 ? response.data : (Array.isArray(response) ? response : [])
        if (prescriptionList && prescriptionList.length > 0) {
          this.prescriptionList = prescriptionList
          console.log('第一条病历数据:', prescriptionList[0])
        }
        this.loading = false
      }).catch(_ => {
        this.loading = false
      })
    } else {
      queryPrescriptionList({ patient: this.account_id }).then(response => {
        console.log('病历列表响应:', response)
        // 处理标准响应格式
        const prescriptionList = response && response.code === 200 ? response.data : (Array.isArray(response) ? response : [])
        if (prescriptionList && prescriptionList.length > 0) {
          this.prescriptionList = prescriptionList
          console.log('第一条病历数据:', prescriptionList[0])
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
      this.selectedPrescription = item
      this.dialogVisible = true
    },
    // 跳转到详情页面
    goToDetail(item) {
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

  .prescription-card {
    margin-bottom: 20px;
    border-radius: 12px;
    transition: all 0.3s;
    cursor: pointer;
    border: 2px solid #EBEEF5;
    overflow: hidden;
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

  /* 药品订单信息样式 */
  .drug-orders-section {
    background: #fff;
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 16px;
    border: 1px solid #EBEEF5;
  }

  .section-title-compact {
    display: flex;
    align-items: center;
    font-size: 15px;
    font-weight: 600;
    color: #303133;
    margin-bottom: 16px;
    padding-bottom: 10px;
    border-bottom: 2px solid #E4E7ED;
  }

  .section-title-compact i {
    margin-right: 8px;
    color: var(--theme-primary, #409EFF);
    font-size: 18px;
  }

  .orders-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .order-item {
    background: #F5F7FA;
    border-radius: 8px;
    padding: 12px;
    border-left: 3px solid var(--theme-primary, #409EFF);
    transition: all 0.3s;
  }

  .order-item:hover {
    background: #ECF5FF;
    box-shadow: 0 2px 8px rgba(64, 158, 255, 0.15);
  }

  .order-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
    padding-bottom: 8px;
    border-bottom: 1px solid #E4E7ED;
  }

  .order-number {
    font-weight: 600;
    color: #303133;
    font-size: 14px;
  }

  .order-content {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .order-row {
    display: flex;
    font-size: 13px;
  }

  .order-label {
    color: #909399;
    width: 90px;
    flex-shrink: 0;
  }

  .order-value {
    color: #606266;
    flex: 1;
    font-weight: 500;
  }

  .no-orders {
    text-align: center;
    padding: 20px;
    color: #909399;
    font-size: 14px;
    background: #F5F7FA;
    border-radius: 8px;
    margin-bottom: 16px;
  }

  .no-orders i {
    font-size: 18px;
    margin-right: 6px;
  }
</style>
