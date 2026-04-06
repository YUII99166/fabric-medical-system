<template>
  <div class="app-container">
    <el-card class="box-card">
      <div slot="header" class="clearfix">
        <span style="font-size: 18px; font-weight: bold;">
          <i class="el-icon-document"></i> 病历详情
        </span>
        <el-button style="float: right;" size="small" @click="goBack">返回</el-button>
      </div>

      <el-row :gutter="20">
        <!-- 左侧：时间线导航 -->
        <el-col :span="6">
          <div class="timeline-container">
            <h3>诊疗记录</h3>
            <el-timeline>
              <!-- 原始病历 -->
              <el-timeline-item
                :timestamp="medicalHistory.original_prescription.created"
                placement="top"
                :color="selectedRecordId === 'original' ? '#409EFF' : '#909399'"
              >
                <el-card
                  class="timeline-card"
                  :class="{ 'active-card': selectedRecordId === 'original' }"
                  @click.native="selectRecord('original')"
                >
                  <div class="timeline-card-header">
                    <i class="el-icon-document"></i>
                    <span>原始病历</span>
                  </div>
                  <div class="timeline-card-body">
                    <p>{{ medicalHistory.original_prescription.hospital_name }}</p>
                    <p>{{ medicalHistory.original_prescription.doctor_name }}</p>
                  </div>
                  <el-tag size="mini" type="success">本院</el-tag>
                </el-card>
              </el-timeline-item>

              <!-- 补充记录 -->
              <el-timeline-item
                v-for="(record, index) in medicalHistory.supplement_records"
                :key="record.id"
                :timestamp="record.created"
                placement="top"
                :color="selectedRecordId === record.id ? '#409EFF' : '#909399'"
              >
                <el-card
                  class="timeline-card"
                  :class="{ 'active-card': selectedRecordId === record.id }"
                  @click.native="selectRecord(record.id)"
                >
                  <div class="timeline-card-header">
                    <i class="el-icon-edit-outline"></i>
                    <span>补充记录 {{ index + 1 }}</span>
                  </div>
                  <div class="timeline-card-body">
                    <p>{{ record.hospital_name }}</p>
                    <p>{{ record.doctor_name }}</p>
                  </div>
                  <el-tag size="mini" :type="getRecordTypeTag(record.record_type)">
                    {{ getRecordTypeName(record.record_type) }}
                  </el-tag>
                </el-card>
              </el-timeline-item>
            </el-timeline>

            <!-- 添加补充记录按钮 -->
            <el-button
              v-if="canAddSupplement"
              type="primary"
              icon="el-icon-plus"
              size="small"
              style="width: 100%; margin-top: 20px;"
              @click="showAddDialog"
            >
              添加补充记录
            </el-button>
          </div>
        </el-col>

        <!-- 右侧：详细内容 -->
        <el-col :span="18">
          <div class="detail-container">
            <!-- 原始病历详情 -->
            <div v-if="selectedRecordId === 'original'" class="record-detail">
              <h2>原始病历</h2>
              <el-form label-width="120px" class="detail-form">
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="病历编号">
                      <span>{{ medicalHistory.original_prescription.prescription_no }}</span>
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="创建时间">
                      <span>{{ medicalHistory.original_prescription.created }}</span>
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="患者姓名">
                      <span>{{ medicalHistory.original_prescription.patient_name }}</span>
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="医生">
                      <span>{{ medicalHistory.original_prescription.doctor_name }} ({{ medicalHistory.original_prescription.doctor_title }})</span>
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="医院">
                      <span>{{ medicalHistory.original_prescription.hospital_name }}</span>
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="科室">
                      <span>{{ medicalHistory.original_prescription.department }}</span>
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-form-item label="主诉">
                  <span>{{ medicalHistory.original_prescription.chief_complaint || '无' }}</span>
                </el-form-item>
                <el-form-item label="现病史">
                  <span>{{ medicalHistory.original_prescription.present_illness || '无' }}</span>
                </el-form-item>
                <el-form-item label="体格检查">
                  <span>{{ medicalHistory.original_prescription.physical_exam || '无' }}</span>
                </el-form-item>
                <el-form-item label="诊断">
                  <el-tag type="danger">{{ medicalHistory.original_prescription.diagnosis }}</el-tag>
                </el-form-item>
                <el-form-item label="用药">
                  <div v-if="medicalHistory.original_prescription.drug && medicalHistory.original_prescription.drug.length > 0">
                    <el-tag
                      v-for="(drug, index) in medicalHistory.original_prescription.drug"
                      :key="index"
                      style="margin-right: 10px; margin-bottom: 5px;"
                    >
                      {{ drug.Name }} × {{ drug.amount }}
                    </el-tag>
                  </div>
                  <span v-else>无</span>
                </el-form-item>
                <el-form-item label="医嘱">
                  <span>{{ medicalHistory.original_prescription.medical_advice || '无' }}</span>
                </el-form-item>
                <el-form-item label="备注">
                  <span>{{ medicalHistory.original_prescription.comment || '无' }}</span>
                </el-form-item>
                
                <!-- 药品订单信息 - 可折叠 -->
                <el-form-item label="药品订单">
                  <div v-loading="loadingOrders">
                    <el-collapse v-if="drugOrders.length > 0" v-model="activeOrderCollapse">
                      <el-collapse-item name="orders">
                        <template slot="title">
                          <span style="font-weight: 600; color: #303133;">
                            <i class="el-icon-shopping-cart-2" style="margin-right: 8px; color: #409EFF;"></i>
                            共 {{ drugOrders.length }} 个药品订单
                          </span>
                        </template>
                        <div class="drug-orders-section">
                          <div v-for="(order, index) in drugOrders" :key="order.id" class="order-card">
                            <div class="order-header">
                              <span class="order-title">
                                <i class="el-icon-shopping-cart-2"></i>
                                订单 #{{ index + 1 }}
                              </span>
                              <el-tag :type="getOrderStatusType(order.status)" size="small">
                                {{ getOrderStatusText(order.status) }}
                              </el-tag>
                            </div>
                            <div class="order-body">
                              <div class="order-item">
                                <span class="order-label">订单ID：</span>
                                <span class="order-value">{{ order.id }}</span>
                              </div>
                              <div class="order-item">
                                <span class="order-label">药品名称：</span>
                                <span class="order-value highlight">{{ order.Name }}</span>
                              </div>
                              <div class="order-item">
                                <span class="order-label">数量：</span>
                                <span class="order-value">{{ order.amount }}</span>
                              </div>
                              <div class="order-item">
                                <span class="order-label">药店：</span>
                                <span class="order-value">{{ order.drug_store_name || order.DrugStoreName || getPharmacyName(order.drug_store) }}</span>
                              </div>
                              <div class="order-item">
                                <span class="order-label">创建时间：</span>
                                <span class="order-value">{{ order.created }}</span>
                              </div>
                            </div>
                          </div>
                        </div>
                      </el-collapse-item>
                    </el-collapse>
                    <span v-else class="no-orders-text">暂无药品订单</span>
                  </div>
                </el-form-item>
                
                <el-form-item label="区块链交易ID">
                  <el-tag type="info" size="mini">{{ medicalHistory.original_prescription.tx_id }}</el-tag>
                </el-form-item>
              </el-form>
            </div>

            <!-- 补充记录详情 -->
            <div v-else class="record-detail">
              <h2>补充诊疗记录</h2>
              <el-form label-width="120px" class="detail-form">
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="记录类型">
                      <el-tag :type="getRecordTypeTag(selectedRecord.record_type)">
                        {{ getRecordTypeName(selectedRecord.record_type) }}
                      </el-tag>
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="创建时间">
                      <span>{{ selectedRecord.created }}</span>
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="医生">
                      <span>{{ selectedRecord.doctor_name }} ({{ selectedRecord.doctor_title }})</span>
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="医院">
                      <span>{{ selectedRecord.hospital_name }}</span>
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="科室">
                      <span>{{ selectedRecord.department }}</span>
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="组织">
                      <span>{{ selectedRecord.organization_name }}</span>
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-form-item label="主诉">
                  <span>{{ selectedRecord.chief_complaint || '无' }}</span>
                </el-form-item>
                <el-form-item label="现病史">
                  <span>{{ selectedRecord.present_illness || '无' }}</span>
                </el-form-item>
                <el-form-item label="体格检查">
                  <span>{{ selectedRecord.physical_exam || '无' }}</span>
                </el-form-item>
                <el-form-item label="诊断">
                  <el-tag type="danger">{{ selectedRecord.diagnosis }}</el-tag>
                </el-form-item>
                <el-form-item label="治疗方案">
                  <span>{{ selectedRecord.treatment || '无' }}</span>
                </el-form-item>
                <el-form-item label="用药">
                  <div v-if="selectedRecord.drug && selectedRecord.drug.length > 0">
                    <el-tag
                      v-for="(drug, index) in selectedRecord.drug"
                      :key="index"
                      style="margin-right: 10px; margin-bottom: 5px;"
                    >
                      {{ drug.Name }} × {{ drug.amount }}
                    </el-tag>
                  </div>
                  <span v-else>无</span>
                </el-form-item>
                <el-form-item label="医嘱">
                  <span>{{ selectedRecord.medical_advice || '无' }}</span>
                </el-form-item>
                <el-form-item label="备注">
                  <span>{{ selectedRecord.comment || '无' }}</span>
                </el-form-item>
                <el-form-item label="区块链交易ID">
                  <el-tag type="info" size="mini">{{ selectedRecord.tx_id }}</el-tag>
                </el-form-item>
              </el-form>
            </div>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 添加补充记录对话框 -->
    <el-dialog
      title="添加补充诊疗记录"
      :visible.sync="addDialogVisible"
      width="60%"
      :close-on-click-modal="false"
    >
      <el-form ref="supplementForm" :model="supplementForm" :rules="supplementRules" label-width="100px">
        <el-form-item label="记录类型" prop="record_type">
          <el-select v-model="supplementForm.record_type" placeholder="请选择记录类型">
            <el-option label="复诊" value="followup"></el-option>
            <el-option label="会诊" value="consultation"></el-option>
            <el-option label="急诊" value="emergency"></el-option>
            <el-option label="转院" value="transfer"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="主诉">
          <el-input v-model="supplementForm.chief_complaint" type="textarea" :rows="2" placeholder="请输入主诉"></el-input>
        </el-form-item>
        <el-form-item label="现病史">
          <el-input v-model="supplementForm.present_illness" type="textarea" :rows="2" placeholder="请输入现病史"></el-input>
        </el-form-item>
        <el-form-item label="体格检查">
          <el-input v-model="supplementForm.physical_exam" type="textarea" :rows="2" placeholder="请输入体格检查结果"></el-input>
        </el-form-item>
        <el-form-item label="诊断" prop="diagnosis">
          <el-input v-model="supplementForm.diagnosis" placeholder="请输入诊断结果"></el-input>
        </el-form-item>
        <el-form-item label="治疗方案">
          <el-input v-model="supplementForm.treatment" type="textarea" :rows="2" placeholder="请输入治疗方案"></el-input>
        </el-form-item>
        <el-form-item label="药品名称">
          <el-input v-model="supplementForm.drug_name" placeholder="多个药品用逗号分隔，如：阿莫西林,布洛芬"></el-input>
        </el-form-item>
        <el-form-item label="药品用量">
          <el-input v-model="supplementForm.drug_amount" placeholder="多个用量用逗号分隔，如：3盒,2盒"></el-input>
        </el-form-item>
        <el-form-item label="医嘱">
          <el-input v-model="supplementForm.medical_advice" type="textarea" :rows="2" placeholder="请输入医嘱"></el-input>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="supplementForm.comment" type="textarea" :rows="2" placeholder="请输入备注"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="addDialogVisible = false">取 消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitSupplement">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { queryFullMedicalHistory, addSupplementRecord } from '@/api/supplement'
import { queryDrugOrderList } from '@/api/drugOrder'
import { getAccountMap } from '@/utils/accountCache'
import { recordPrescriptionAccess } from '@/api/accessTrace'

export default {
  name: 'PrescriptionDetail',
  data() {
    return {
      prescriptionId: '',
      currentUser: null,
      medicalHistory: {
        original_prescription: {
          prescription_no: '',
          created: '',
          patient_name: '',
          doctor_name: '',
          doctor_title: '',
          hospital_name: '',
          department: '',
          chief_complaint: '',
          present_illness: '',
          physical_exam: '',
          diagnosis: '',
          drug: [],
          medical_advice: '',
          comment: '',
          tx_id: ''
        },
        supplement_records: [],
        total_records: 0
      },
      selectedRecordId: 'original',
      selectedRecord: {},
      addDialogVisible: false,
      submitting: false,
      supplementForm: {
        record_type: 'followup',
        chief_complaint: '',
        present_illness: '',
        physical_exam: '',
        diagnosis: '',
        treatment: '',
        drug_name: '',
        drug_amount: '',
        medical_advice: '',
        comment: ''
      },
      supplementRules: {
        record_type: [{ required: true, message: '请选择记录类型', trigger: 'change' }],
        diagnosis: [{ required: true, message: '请输入诊断结果', trigger: 'blur' }]
      },
      drugOrders: [], // 药品订单列表
      loadingOrders: false, // 加载订单状态
      accountMap: {}, // 账户ID到账户信息的映射
      activeOrderCollapse: ['orders'] // 默认展开订单折叠面板
    }
  },
  computed: {
    canAddSupplement() {
      // 医生角色且有权限查看该病历才能添加补充记录
      const result = this.currentUser && this.currentUser.role === '医生'
      console.log('canAddSupplement 检查:', {
        currentUser: this.currentUser,
        role: this.currentUser?.role,
        result: result
      })
      return result
    }
  },
  created() {
    // 从 sessionStorage 读取用户信息
    const userInfoStr = sessionStorage.getItem('userInfo')
    if (userInfoStr) {
      try {
        this.currentUser = JSON.parse(userInfoStr)
        console.log('当前用户信息:', this.currentUser)
      } catch (e) {
        console.error('解析用户信息失败:', e)
      }
    }

    // 使用缓存的账户列表，避免重复查询
    this.loadAccountListFromCache()

    this.prescriptionId = this.$route.query.id
    console.log('路由参数:', this.$route.query)
    console.log('病历ID:', this.prescriptionId)
    
    if (this.prescriptionId) {
      // 并行加载数据，提升性能
      Promise.all([
        this.loadMedicalHistory(),
        this.loadDrugOrders()
      ]).catch(err => {
        console.error('加载数据失败:', err)
      })
    } else {
      console.error('缺少病历ID参数，完整路由:', this.$route)
      this.$message.error('缺少病历ID参数')
      setTimeout(() => {
        this.goBack()
      }, 1500)
    }
  },
  methods: {
    // 从缓存加载账户列表（快速）
    loadAccountListFromCache() {
      getAccountMap().then(accountMap => {
        this.accountMap = accountMap
        console.log('✅ 从缓存加载账户映射:', Object.keys(accountMap).length, '个账户')
      }).catch(err => {
        console.error('加载账户映射失败:', err)
        // 缓存失败不影响主要功能
      })
    },
    // 获取药店名称
    getPharmacyName(pharmacyId) {
      if (!pharmacyId) return '未知药店'
      const pharmacy = this.accountMap[pharmacyId]
      if (pharmacy) {
        return pharmacy.organization_name || pharmacy.name || pharmacyId
      }
      return pharmacyId
    },
    // 加载药品订单
    loadDrugOrders() {
      this.loadingOrders = true
      queryDrugOrderList({}).then(response => {
        const orderList = response && response.code === 200 ? response.data : (Array.isArray(response) ? response : [])
        // 过滤出当前病历的订单
        this.drugOrders = orderList.filter(order => order.prescription === this.prescriptionId)
        this.loadingOrders = false
        console.log('加载到的药品订单:', this.drugOrders)
      }).catch(err => {
        console.error('加载药品订单失败:', err)
        this.loadingOrders = false
      })
    },
    // 获取订单状态类型
    getOrderStatusType(status) {
      const statusMap = {
        'pending': 'warning',
        'processing': 'info',
        'completed': 'success',
        'cancelled': 'danger'
      }
      return statusMap[status] || 'info'
    },
    // 获取订单状态文本
    getOrderStatusText(status) {
      const statusMap = {
        'pending': '待处理',
        'processing': '处理中',
        'completed': '已完成',
        'cancelled': '已取消'
      }
      return statusMap[status] || status
    },
    loadMedicalHistory() {
      const loading = this.$loading({
        lock: true,
        text: '加载中...',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      })

      console.log('开始加载病历历史，ID:', this.prescriptionId)

      // 带静默重试的加载函数
      const loadWithRetry = async (retryCount = 0) => {
        const maxRetries = 3
        
        try {
          const response = await queryFullMedicalHistory({ prescription_id: this.prescriptionId })
          
          console.log('API响应:', response)
          
          if (response.code === 200) {
            loading.close()
            const data = response.data || {}
            console.log('返回的数据:', data)
            
            // 确保数据结构完整，处理 null 值
            this.medicalHistory = {
              original_prescription: data.original_prescription || {},
              supplement_records: data.supplement_records || [],
              total_records: data.total_records || 0
            }
            
            console.log('处理后的病历数据:', this.medicalHistory)
            
            // 默认选中原始病历
            this.selectedRecordId = 'original'
            
            // 记录访问日志（如果不是患者本人）
            if (this.currentUser && this.medicalHistory.original_prescription) {
              const prescription = this.medicalHistory.original_prescription
              const patientId = prescription.patient_id || prescription.patient || ''
              
              console.log('检查是否需要记录访问:', {
                currentUserId: this.currentUser.account_id,
                patientId: patientId,
                shouldRecord: this.currentUser.account_id !== patientId
              })
              
              if (this.currentUser.account_id !== patientId) {
                const accessData = {
                  prescription_id: this.prescriptionId,
                  prescription_no: prescription.prescription_no || '',
                  patient_id: patientId,
                  patient_name: prescription.patient_name || '',
                  accessor_id: this.currentUser.account_id,
                  accessor_name: this.currentUser.account_name || this.currentUser.name || '',
                  accessor_role: this.currentUser.role || '',
                  accessor_organization: this.currentUser.organization || '',
                  accessor_organization_name: this.currentUser.organization_name || '',
                  access_type: 'view',
                  access_reason: '查看病历详情'
                }
                console.log('准备记录访问日志:', accessData)
                
                recordPrescriptionAccess(accessData).then(response => {
                  console.log('✅ 访问日志记录成功:', response)
                }).catch(err => {
                  console.error('❌ 记录访问日志失败:', err)
                  // 不影响主流程，静默失败
                })
              }
            }
          } else {
            // 如果失败且还有重试次数，静默重试
            if (retryCount < maxRetries) {
              console.log(`加载病历失败，正在静默重试 (${retryCount + 1}/${maxRetries})...`)
              await new Promise(resolve => setTimeout(resolve, 500 * (retryCount + 1)))
              return loadWithRetry(retryCount + 1)
            } else {
              // 所有重试都失败后才显示错误
              loading.close()
              console.error('加载失败，错误信息:', response.msg)
              this.$message.error(response.msg || '加载失败')
            }
          }
        } catch (error) {
          // 只在控制台记录错误
          if (retryCount === 0) {
            console.error('加载病历历史失败:', error)
          }
          
          // 如果失败且还有重试次数，静默重试
          if (retryCount < maxRetries) {
            console.log(`加载病历失败，正在静默重试 (${retryCount + 1}/${maxRetries})...`)
            await new Promise(resolve => setTimeout(resolve, 500 * (retryCount + 1)))
            return loadWithRetry(retryCount + 1)
          } else {
            // 所有重试都失败后才显示错误
            loading.close()
            this.$message.error('加载失败，请重试')
          }
        }
      }
      
      // 开始加载
      loadWithRetry()
    },
    selectRecord(recordId) {
      this.selectedRecordId = recordId
      if (recordId !== 'original') {
        const records = this.medicalHistory.supplement_records || []
        this.selectedRecord = records.find(r => r.id === recordId) || {}
      }
    },
    showAddDialog() {
      this.addDialogVisible = true
      this.$nextTick(() => {
        this.$refs.supplementForm.resetFields()
      })
    },
    submitSupplement() {
      this.$refs.supplementForm.validate(valid => {
        if (valid) {
          if (!this.currentUser || !this.currentUser.account_id) {
            this.$message.error('无法获取用户信息，请重新登录')
            return
          }

          this.submitting = true
          const data = {
            original_prescription_id: this.prescriptionId,
            doctor_id: this.currentUser.account_id,
            ...this.supplementForm
          }

          addSupplementRecord(data)
            .then(response => {
              this.submitting = false
              if (response.code === 200) {
                this.$message.success('添加补充记录成功')
                this.addDialogVisible = false
                // 重新加载病历历史
                this.loadMedicalHistory()
              } else {
                this.$message.error(response.msg || '添加失败')
              }
            })
            .catch(error => {
              this.submitting = false
              console.error('添加补充记录失败:', error)
              this.$message.error('添加失败，请重试')
            })
        }
      })
    },
    getRecordTypeName(type) {
      const typeMap = {
        'consultation': '会诊',
        'followup': '复诊',
        'emergency': '急诊',
        'transfer': '转院'
      }
      return typeMap[type] || '其他'
    },
    getRecordTypeTag(type) {
      const tagMap = {
        'consultation': 'warning',
        'followup': 'success',
        'emergency': 'danger',
        'transfer': 'info'
      }
      return tagMap[type] || ''
    },
    goBack() {
      this.$router.go(-1)
    }
  }
}
</script>

<style scoped>
.app-container {
  padding: 20px;
}

.box-card {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.08);
}

.timeline-container {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 25px;
  border-radius: 8px;
  max-height: 800px;
  overflow-y: auto;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.2);
}

.timeline-container::-webkit-scrollbar {
  width: 6px;
}

.timeline-container::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
}

.timeline-container::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.3);
  border-radius: 3px;
}

.timeline-container::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.5);
}

.timeline-container h3 {
  margin-top: 0;
  margin-bottom: 25px;
  color: #ffffff;
  font-size: 18px;
  font-weight: 600;
  text-align: center;
  letter-spacing: 1px;
}

.timeline-card {
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 10px;
  border-radius: 8px;
  overflow: hidden;
}

.timeline-card:hover {
  transform: translateX(8px) scale(1.02);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
}

.active-card {
  border: 2px solid #ffffff;
  box-shadow: 0 6px 25px rgba(255, 255, 255, 0.4);
  transform: translateX(8px);
}

.timeline-card >>> .el-card__body {
  padding: 15px;
}

.timeline-card-header {
  font-weight: 600;
  margin-bottom: 10px;
  font-size: 14px;
  color: #303133;
  display: flex;
  align-items: center;
}

.timeline-card-header i {
  margin-right: 8px;
  font-size: 16px;
  color: #409EFF;
}

.timeline-card-body p {
  margin: 5px 0;
  font-size: 13px;
  color: #606266;
  line-height: 1.6;
}

.detail-container {
  background: #ffffff;
  padding: 30px;
  border-radius: 8px;
  min-height: 600px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.record-detail {
  background: #ffffff;
  border: 2px solid #409EFF;
  border-radius: 12px;
  padding: 30px;
  box-shadow: 0 4px 16px rgba(64, 158, 255, 0.15);
  transition: all 0.3s ease;
}

.record-detail:hover {
  box-shadow: 0 8px 24px rgba(64, 158, 255, 0.25);
  transform: translateY(-2px);
}

.record-detail h2 {
  margin-top: 0;
  margin-bottom: 25px;
  color: #303133;
  border-bottom: 3px solid #409EFF;
  padding-bottom: 12px;
  font-size: 20px;
  font-weight: 600;
  display: flex;
  align-items: center;
}

.record-detail h2::before {
  content: '';
  width: 4px;
  height: 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  margin-right: 12px;
  border-radius: 2px;
}

.detail-form {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 25px;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
}

.detail-form .el-form-item {
  margin-bottom: 18px;
  background: #ffffff;
  padding: 12px 15px;
  border-radius: 6px;
  transition: all 0.3s;
}

.detail-form .el-form-item:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.detail-form .el-form-item__label {
  font-weight: 600;
  color: #606266;
  font-size: 14px;
}

.detail-form .el-form-item__content span {
  color: #303133;
  font-size: 14px;
  line-height: 1.6;
}

.detail-form .el-tag {
  margin: 2px 5px 2px 0;
}

/* 时间线样式优化 */
.timeline-container >>> .el-timeline-item__timestamp {
  color: rgba(255, 255, 255, 0.9);
  font-size: 12px;
  font-weight: 500;
}

.timeline-container >>> .el-timeline-item__node {
  background-color: #ffffff;
}

.timeline-container >>> .el-timeline-item__tail {
  border-left: 2px solid rgba(255, 255, 255, 0.3);
}

/* 添加按钮样式 */
.timeline-container .el-button {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: 2px solid #ffffff;
  color: #ffffff;
  font-weight: 600;
  transition: all 0.3s;
}

.timeline-container .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(255, 255, 255, 0.3);
  background: linear-gradient(135deg, #764ba2 0%, #667eea 100%);
}

/* 对话框样式优化 */
>>> .el-dialog__header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

>>> .el-dialog__title {
  color: #ffffff;
  font-weight: 600;
  font-size: 18px;
}

>>> .el-dialog__headerbtn .el-dialog__close {
  color: #ffffff;
  font-size: 20px;
}

>>> .el-dialog__body {
  padding: 25px;
}

/* 响应式优化 */
@media (max-width: 768px) {
  .timeline-container {
    max-height: 400px;
  }
  
  .detail-container {
    padding: 15px;
  }
}

/* 药品订单样式 */
.drug-orders-section {
  display: flex;
  flex-direction: column;
  gap: 15px;
  margin-top: 10px;
}

.order-card {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  border-radius: 10px;
  padding: 16px;
  border-left: 4px solid #409EFF;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.order-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 16px rgba(64, 158, 255, 0.25);
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 10px;
  border-bottom: 2px solid rgba(64, 158, 255, 0.2);
}

.order-title {
  font-weight: 600;
  color: #303133;
  font-size: 15px;
  display: flex;
  align-items: center;
}

.order-title i {
  margin-right: 8px;
  color: #409EFF;
  font-size: 18px;
}

.order-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.order-item {
  display: flex;
  font-size: 14px;
  line-height: 1.6;
  background: #ffffff;
  padding: 8px 12px;
  border-radius: 6px;
}

.order-label {
  color: #909399;
  width: 100px;
  flex-shrink: 0;
  font-weight: 500;
}

.order-value {
  color: #606266;
  flex: 1;
}

.order-value.highlight {
  color: #409EFF;
  font-weight: 600;
  font-size: 15px;
}

.no-orders-text {
  color: #909399;
  font-size: 14px;
  font-style: italic;
}

/* 折叠面板样式优化 */
.detail-form >>> .el-collapse {
  border: none;
  background: transparent;
}

.detail-form >>> .el-collapse-item__header {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  border: none;
  padding: 12px 16px;
  border-radius: 8px;
  font-size: 15px;
  transition: all 0.3s;
}

.detail-form >>> .el-collapse-item__header:hover {
  background: linear-gradient(135deg, #e8eaf0 0%, #b8c5d8 100%);
}

.detail-form >>> .el-collapse-item__wrap {
  border: none;
  background: transparent;
}

.detail-form >>> .el-collapse-item__content {
  padding: 16px 0 0 0;
}
</style>
