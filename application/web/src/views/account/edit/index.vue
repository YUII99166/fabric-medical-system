<template>
  <div class="app-container">
    <el-card class="form-card">
      <div slot="header" class="card-header">
        <i class="el-icon-edit"></i>
        <span>编辑用户</span>
      </div>

      <el-form 
        ref="ruleForm" 
        v-loading="loading" 
        :model="ruleForm" 
        :rules="rules"
        label-width="120px"
        class="edit-form"
      >
        <!-- 基本信息 -->
        <div class="form-section">
          <h3 class="section-title">基本信息</h3>
          
          <el-form-item label="用户名" prop="username">
            <el-input 
              v-model="ruleForm.username" 
              disabled
              placeholder="用户名不可修改"
            />
            <span class="form-tip">用户名创建后不可修改</span>
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

        <!-- 修改密码 -->
        <div class="form-section">
          <h3 class="section-title">修改密码（可选）</h3>
          <el-alert
            title="如果不需要修改密码，请保持密码字段为空"
            type="info"
            :closable="false"
            style="margin-bottom: 20px;"
          />
          
          <el-form-item label="新密码" prop="password">
            <el-input 
              v-model="ruleForm.password" 
              type="password"
              placeholder="留空表示不修改密码"
              maxlength="20"
              show-password
              clearable
            />
          </el-form-item>

          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input 
              v-model="ruleForm.confirmPassword" 
              type="password"
              placeholder="请再次输入新密码"
              maxlength="20"
              show-password
              clearable
            />
          </el-form-item>
        </div>

        <!-- 按钮 -->
        <el-form-item class="form-buttons">
          <el-button type="primary" @click="submitForm" :loading="loading">
            <i class="el-icon-check"></i>
            保存修改
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
import { getUserDetail, updateUser } from '@/api/accountV2'

export default {
  name: 'EditAccount',
  data() {
    // 自定义验证规则
    const validatePassword = (rule, value, callback) => {
      // 如果密码为空，表示不修改密码，跳过验证
      if (!value || value === '') {
        callback()
        return
      }
      
      if (value.length < 6 || value.length > 20) {
        callback(new Error('密码长度应为6-20位'))
      } else {
        if (this.ruleForm.confirmPassword !== '') {
          this.$refs.ruleForm.validateField('confirmPassword')
        }
        callback()
      }
    }

    const validateConfirmPassword = (rule, value, callback) => {
      // 如果新密码为空，确认密码也应该为空
      if (!this.ruleForm.password || this.ruleForm.password === '') {
        if (value && value !== '') {
          callback(new Error('请先输入新密码'))
        } else {
          callback()
        }
        return
      }
      
      // 如果新密码不为空，确认密码必须一致
      if (!value) {
        callback(new Error('请再次输入密码'))
      } else if (value !== this.ruleForm.password) {
        callback(new Error('两次输入密码不一致'))
      } else {
        callback()
      }
    }

    return {
      userId: null,
      loading: false,
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
      originalData: null,
      rules: {
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
        ],
        password: [
          { validator: validatePassword, trigger: 'blur' }
        ],
        confirmPassword: [
          { validator: validateConfirmPassword, trigger: 'blur' }
        ]
      },
      organizationMap: {
        'TaobaoMSP': '协和医院',
        'JDMSP': '301医院',
        'WenjinMSP': '温江医疗中心',
        'RegCenterMSP': '监管中心'
      }
    }
  },
  created() {
    this.loadUserDetail()
  },
  methods: {
    async loadUserDetail() {
      const id = this.$route.query.id
      const username = this.$route.query.username
      
      if (!id && !username) {
        this.$message.error('缺少用户信息')
        this.goBack()
        return
      }
      
      this.loading = true
      
      try {
        // 优先使用ID，否则使用用户名
        const response = await getUserDetail(id || username)
        if (response && response.code === 200) {
          const user = response.data
          this.userId = user.id
          this.originalData = { ...user }
          
          // 填充表单
          this.ruleForm.username = user.username
          this.ruleForm.name = user.account_name.replace(/^(医生|病人|药店|管理员)-/, '')
          this.ruleForm.role = user.role
          this.ruleForm.organization = user.organization || ''
          this.ruleForm.organization_name = user.organization_name || ''
          this.ruleForm.department = user.department || ''
          this.ruleForm.doctor_title = user.doctor_title || ''
        } else {
          this.$message.error(response.msg || '加载用户信息失败')
          this.goBack()
        }
      } catch (error) {
        console.error('加载用户信息失败:', error)
        this.$message.error('加载用户信息失败')
        this.goBack()
      } finally {
        this.loading = false
      }
    },
    
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
          this.$confirm('确认保存修改吗？', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            this.updateUserInfo()
          }).catch(() => {
            this.$message.info('已取消修改')
          })
        } else {
          this.$message.warning('请填写完整的表单信息')
          return false
        }
      })
    },

    async updateUserInfo() {
      this.loading = true
      
      try {
        // 构建账户名称
        const accountName = `${this.ruleForm.role}-${this.ruleForm.name}`
        
        // 构建更新数据
        const updateData = {
          id: this.userId,
          account_name: accountName,
          role: this.ruleForm.role,
          organization: this.ruleForm.organization || '',
          organization_name: this.ruleForm.organization_name || '',
          department: this.ruleForm.department || '',
          doctor_title: this.ruleForm.doctor_title || ''
        }
        
        // 如果修改了密码，添加密码字段
        if (this.ruleForm.password && this.ruleForm.password !== '') {
          updateData.password = this.ruleForm.password
        }
        
        const response = await updateUser(updateData)

        this.loading = false

        if (response && response.code === 200) {
          this.$message.success('用户信息更新成功！')
          setTimeout(() => {
            this.$router.push('/account/all')
          }, 1500)
        } else {
          this.$message.error(response.msg || '更新失败，请重试')
        }
      } catch (error) {
        this.loading = false
        console.error('更新用户失败:', error)
        
        if (error.response && error.response.data) {
          this.$message.error(error.response.data.msg || '更新失败')
        } else if (error.message) {
          this.$message.error(error.message)
        } else {
          this.$message.error('更新失败，请检查网络连接或联系管理员')
        }
      }
    },

    resetForm() {
      if (this.originalData) {
        this.ruleForm.name = this.originalData.account_name.replace(/^(医生|病人|药店|管理员)-/, '')
        this.ruleForm.role = this.originalData.role
        this.ruleForm.organization = this.originalData.organization || ''
        this.ruleForm.organization_name = this.originalData.organization_name || ''
        this.ruleForm.department = this.originalData.department || ''
        this.ruleForm.doctor_title = this.originalData.doctor_title || ''
        this.ruleForm.password = ''
        this.ruleForm.confirmPassword = ''
        this.$refs.ruleForm.clearValidate()
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

.edit-form {
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
