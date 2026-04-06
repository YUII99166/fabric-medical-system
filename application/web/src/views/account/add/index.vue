<template>
  <div class="app-container">
    <el-card class="form-card">
      <div slot="header" class="card-header">
        <i class="el-icon-user-solid"></i>
        <span>创建新用户</span>
      </div>

      <el-form 
        ref="ruleForm" 
        v-loading="loading" 
        :model="ruleForm" 
        :rules="rules"
        label-width="120px"
        class="create-form"
      >
        <!-- 基本信息 -->
        <div class="form-section">
          <h3 class="section-title">基本信息</h3>
          
          <el-form-item label="用户名" prop="username">
            <el-input 
              v-model="ruleForm.username" 
              placeholder="请输入用户名（2-20个字符）"
              maxlength="20"
            />
            <span class="form-tip">用于登录系统，不可重复</span>
          </el-form-item>

          <el-form-item label="密码" prop="password">
            <el-input 
              v-model="ruleForm.password" 
              type="password"
              placeholder="请输入密码（6-20位）"
              maxlength="20"
              show-password
            />
          </el-form-item>

          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input 
              v-model="ruleForm.confirmPassword" 
              type="password"
              placeholder="请再次输入密码"
              maxlength="20"
              show-password
            />
          </el-form-item>

          <el-form-item label="姓名" prop="name">
            <el-input 
              v-model="ruleForm.name" 
              placeholder="请输入真实姓名"
              maxlength="20"
            />
          </el-form-item>

          <el-form-item label="角色" prop="role">
            <el-select v-model="ruleForm.role" placeholder="请选择角色" @change="handleRoleChange">
              <el-option label="管理员" value="管理员"></el-option>
              <el-option label="医生" value="医生"></el-option>
              <el-option label="病人" value="病人"></el-option>
              <el-option label="药店" value="药店"></el-option>
            </el-select>
          </el-form-item>
        </div>

        <!-- 医生和管理员专属信息 -->
        <div v-if="ruleForm.role === '医生' || ruleForm.role === '管理员'" class="form-section">
          <h3 class="section-title">{{ ruleForm.role === '管理员' ? '管理员信息' : '医生信息' }}</h3>
          
          <el-form-item label="所属组织" prop="organization">
            <el-select v-model="ruleForm.organization" placeholder="请选择所属组织" @change="handleOrgChange">
              <el-option v-if="ruleForm.role === '医生'" label="协和医院" value="TaobaoMSP"></el-option>
              <el-option v-if="ruleForm.role === '医生'" label="301医院" value="JDMSP"></el-option>
              <el-option v-if="ruleForm.role === '医生'" label="温江医疗中心" value="WenjinMSP"></el-option>
              <el-option v-if="ruleForm.role === '管理员'" label="监管中心" value="RegCenterMSP"></el-option>
            </el-select>
            <span v-if="ruleForm.role === '管理员'" class="form-tip">监管中心账户具有管理员权限</span>
          </el-form-item>

          <el-form-item v-if="ruleForm.role === '医生'" label="科室" prop="department">
            <el-input 
              v-model="ruleForm.department" 
              placeholder="请输入科室名称"
              maxlength="50"
            />
          </el-form-item>

          <el-form-item v-if="ruleForm.role === '医生'" label="职称" prop="doctor_title">
            <el-select v-model="ruleForm.doctor_title" placeholder="请选择职称">
              <el-option label="住院医师" value="住院医师"></el-option>
              <el-option label="主治医师" value="主治医师"></el-option>
              <el-option label="副主任医师" value="副主任医师"></el-option>
              <el-option label="主任医师" value="主任医师"></el-option>
            </el-select>
          </el-form-item>
        </div>

        <!-- 按钮 -->
        <el-form-item class="form-buttons">
          <el-button type="primary" @click="submitForm" :loading="loading">
            <i class="el-icon-check"></i>
            立即创建
          </el-button>
          <el-button @click="resetForm">
            <i class="el-icon-refresh-left"></i>
            重置
          </el-button>
          <el-button @click="goBack">
            <i class="el-icon-back"></i>
            返回
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { register } from '@/api/accountV2'

export default {
  name: 'AddAccount',
  data() {
    // 自定义验证规则
    const validateUsername = (rule, value, callback) => {
      if (!value || value.trim() === '') {
        callback(new Error('请输入用户名'))
      } else if (value.trim().length < 2 || value.trim().length > 20) {
        callback(new Error('用户名长度应为2-20个字符'))
      } else if (!/^[a-zA-Z0-9_\u4e00-\u9fa5]+$/.test(value)) {
        callback(new Error('用户名只能包含字母、数字、下划线和中文'))
      } else {
        callback()
      }
    }

    const validatePassword = (rule, value, callback) => {
      if (!value) {
        callback(new Error('请输入密码'))
      } else if (value.length < 6 || value.length > 20) {
        callback(new Error('密码长度应为6-20位'))
      } else {
        if (this.ruleForm.confirmPassword !== '') {
          this.$refs.ruleForm.validateField('confirmPassword')
        }
        callback()
      }
    }

    const validateConfirmPassword = (rule, value, callback) => {
      if (!value) {
        callback(new Error('请再次输入密码'))
      } else if (value !== this.ruleForm.password) {
        callback(new Error('两次输入密码不一致'))
      } else {
        callback()
      }
    }

    return {
      ruleForm: {
        username: '',
        password: '',
        confirmPassword: '',
        name: '',
        role: '',
        organization: '',
        organization_name: '',
        department: '',
        doctor_title: ''
      },
      rules: {
        username: [
          { required: true, validator: validateUsername, trigger: 'blur' }
        ],
        password: [
          { required: true, validator: validatePassword, trigger: 'blur' }
        ],
        confirmPassword: [
          { required: true, validator: validateConfirmPassword, trigger: 'blur' }
        ],
        name: [
          { required: true, message: '请输入姓名', trigger: 'blur' },
          { min: 2, max: 20, message: '姓名长度应为2-20个字符', trigger: 'blur' }
        ],
        role: [
          { required: true, message: '请选择角色', trigger: 'change' }
        ],
        organization: [
          { required: true, message: '请选择所属医院', trigger: 'change' }
        ],
        department: [
          { required: true, message: '请输入科室', trigger: 'blur' }
        ],
        doctor_title: [
          { required: true, message: '请选择职称', trigger: 'change' }
        ]
      },
      loading: false,
      organizationMap: {
        'TaobaoMSP': '协和医院',
        'JDMSP': '301医院',
        'WenjinMSP': '温江医疗中心',
        'RegCenterMSP': '监管中心'
      }
    }
  },
  computed: {
    ...mapGetters([
      'account_id'
    ])
  },
  methods: {
    handleRoleChange(role) {
      // 切换角色时清空组织相关字段
      if (role !== '医生' && role !== '管理员') {
        this.ruleForm.organization = ''
        this.ruleForm.organization_name = ''
        this.ruleForm.department = ''
        this.ruleForm.doctor_title = ''
      } else if (role === '管理员') {
        // 管理员自动选择监管中心
        this.ruleForm.organization = 'RegCenterMSP'
        this.ruleForm.organization_name = '监管中心'
        this.ruleForm.department = ''
        this.ruleForm.doctor_title = ''
      } else {
        // 医生角色清空组织，让用户选择
        this.ruleForm.organization = ''
        this.ruleForm.organization_name = ''
      }
    },
    
    handleOrgChange(org) {
      this.ruleForm.organization_name = this.organizationMap[org] || ''
    },

    submitForm() {
      this.$refs.ruleForm.validate((valid) => {
        if (valid) {
          this.$confirm('确认创建该用户吗？', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            this.createUser()
          }).catch(() => {
            this.$message.info('已取消创建')
          })
        } else {
          this.$message.warning('请填写完整的表单信息')
          return false
        }
      })
    },

    async createUser() {
      this.loading = true
      
      try {
        // 构建账户名称
        const accountName = `${this.ruleForm.role}-${this.ruleForm.name}`
        
        // 调用注册接口
        const response = await register({
          account_name: accountName,
          username: this.ruleForm.username.trim(),
          password: this.ruleForm.password,
          role: this.ruleForm.role,
          organization: this.ruleForm.organization || '',
          organization_name: this.ruleForm.organization_name || '',
          department: this.ruleForm.department || '',
          doctor_title: this.ruleForm.doctor_title || ''
        })

        this.loading = false

        if (response && response.code === 200) {
          this.$message.success('用户创建成功！')
          // 延迟跳转，让用户看到成功提示
          setTimeout(() => {
            this.$router.push('/account/all')
          }, 1500)
        } else {
          this.$message.error(response.msg || '创建失败，请重试')
        }
      } catch (error) {
        this.loading = false
        console.error('创建用户失败:', error)
        
        // 根据错误类型显示不同的提示
        if (error.response && error.response.data) {
          this.$message.error(error.response.data.msg || '创建失败')
        } else if (error.message) {
          this.$message.error(error.message)
        } else {
          this.$message.error('创建失败，请检查网络连接或联系管理员')
        }
      }
    },

    resetForm() {
      this.$refs.ruleForm.resetFields()
      this.ruleForm = {
        username: '',
        password: '',
        confirmPassword: '',
        name: '',
        role: '',
        organization: '',
        organization_name: '',
        department: '',
        doctor_title: ''
      }
    },

    goBack() {
      this.$router.back()
    }
  }
}
</script>

<style lang="scss" scoped>
.app-container {
  padding: 20px;
}

.form-card {
  max-width: 800px;
  margin: 0 auto;
  border-radius: 12px;

  .card-header {
    font-size: 18px;
    font-weight: 600;
    color: #303133;
    display: flex;
    align-items: center;
    gap: 10px;

    i {
      color: #409EFF;
      font-size: 20px;
    }
  }
}

.create-form {
  padding: 20px 0;
}

.form-section {
  margin-bottom: 30px;
  padding: 20px;
  background: #f5f7fa;
  border-radius: 8px;

  .section-title {
    margin: 0 0 20px 0;
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    padding-bottom: 10px;
    border-bottom: 2px solid #409EFF;
  }
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-left: 10px;
}

.form-buttons {
  margin-top: 30px;
  text-align: center;
  padding-top: 20px;
  border-top: 1px solid #EBEEF5;

  .el-button {
    min-width: 120px;
  }
}

::v-deep .el-form-item__label {
  font-weight: 500;
}

::v-deep .el-input,
::v-deep .el-select {
  width: 100%;
}
</style>
