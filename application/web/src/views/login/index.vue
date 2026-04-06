<template>
  <div class="login-container">
    <!-- 背景装饰 -->
    <div class="background-decoration">
      <div class="circle circle-1"></div>
      <div class="circle circle-2"></div>
      <div class="circle circle-3"></div>
    </div>

    <!-- 左侧信息区域 -->
    <div class="info-section">
      <div class="info-content">
        <div class="logo-section">
          <div class="logo-icon">
            <!-- 详细版Logo - 用于大尺寸展示 -->
            <img src="/image/共建联盟链 (1).png" alt="MedChain 详细版Logo" class="logo-image detailed-logo" />
            <!-- 精简版Logo - 用于小尺寸场景 -->
            <img src="/image/联盟链_数据可视化备份.png" alt="MedChain 精简版Logo" class="logo-image simple-logo" />
          </div>
          <h1 class="system-title">MedChain</h1>
          <p class="system-subtitle">基于区块链的社区医疗管理系统</p>
        </div>
        
        <div class="features">
          <div class="feature-item">
            <i class="el-icon-lock"></i>
            <div class="feature-text">
              <h4>数据安全</h4>
              <p>区块链加密存储，不可篡改</p>
            </div>
          </div>
          <div class="feature-item">
            <i class="el-icon-document"></i>
            <div class="feature-text">
              <h4>病历管理</h4>
              <p>电子病历，随时查阅</p>
            </div>
          </div>
          <div class="feature-item">
            <i class="el-icon-key"></i>
            <div class="feature-text">
              <h4>授权管理</h4>
              <p>跨院诊疗，安全授权</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧登录表单 -->
    <div class="form-section">
      <el-form ref="loginForm" :model="formData" :rules="formRules" class="login-form" auto-complete="on" label-position="left">
        
        <div class="form-header">
          <div class="form-logo">
            <img src="/image/联盟链_数据可视化备份.png" alt="MedChain Logo" class="form-logo-image" />
          </div>
          <h2 class="form-title">欢迎使用 MedChain</h2>
          <p class="form-subtitle">请登录或注册您的账户</p>
        </div>

        <el-tabs v-model="activeTab" class="login-tabs">
        <el-tab-pane label="登录" name="login">
          <el-form-item prop="username">
            <span class="svg-container">
              <svg-icon icon-class="user" />
            </span>
            <el-input
              ref="username"
              v-model="formData.username"
              placeholder="用户名"
              name="username"
              type="text"
              tabindex="1"
              auto-complete="on"
            />
          </el-form-item>

          <el-form-item prop="password">
            <span class="svg-container">
              <svg-icon icon-class="password" />
            </span>
            <el-input
              :key="passwordType"
              ref="password"
              v-model="formData.password"
              :type="passwordType"
              placeholder="密码"
              name="password"
              tabindex="2"
              auto-complete="on"
              @keyup.enter.native="handleLogin"
            />
            <span class="show-pwd" @click="showPwd">
              <svg-icon :icon-class="passwordType === 'password' ? 'eye' : 'eye-open'" />
            </span>
          </el-form-item>

          <el-button :loading="loading" type="primary" style="width:100%;margin-bottom:30px;" @click.native.prevent="handleLogin">登录</el-button>
        </el-tab-pane>

        <el-tab-pane label="注册" name="register">
          <el-form-item prop="regUsername">
            <span class="svg-container">
              <svg-icon icon-class="user" />
            </span>
            <el-input
              v-model="formData.regUsername"
              placeholder="用户名"
              name="regUsername"
              type="text"
              tabindex="3"
              auto-complete="off"
            />
          </el-form-item>

          <el-form-item prop="regPassword">
            <span class="svg-container">
              <svg-icon icon-class="password" />
            </span>
            <el-input
              v-model="formData.regPassword"
              type="password"
              placeholder="密码"
              name="regPassword"
              tabindex="4"
              auto-complete="off"
            />
          </el-form-item>

          <el-form-item prop="confirmPassword">
            <span class="svg-container">
              <svg-icon icon-class="password" />
            </span>
            <el-input
              v-model="formData.confirmPassword"
              type="password"
              placeholder="确认密码"
              name="confirmPassword"
              tabindex="5"
              auto-complete="off"
            />
          </el-form-item>

          <el-form-item prop="role">
            <span class="svg-container">
              <svg-icon icon-class="peoples" />
            </span>
            <el-select v-model="formData.role" placeholder="请选择角色" style="width: 100%;" @change="handleRoleChange">
              <el-option label="医生" value="doctor" />
              <el-option label="病人" value="patient" />
              <el-option label="药店" value="drugstore" />
            </el-select>
          </el-form-item>

          <!-- 医生需要选择组织（不包含监管中心） -->
          <el-form-item v-if="formData.role === 'doctor'" prop="organization">
            <span class="svg-container">
              <svg-icon icon-class="tree" />
            </span>
            <el-select v-model="formData.organization" placeholder="请选择所属医院" style="width: 100%;">
              <el-option label="协和医院（Taobao组织）" value="TaobaoMSP" />
              <el-option label="301医院（JD组织）" value="JDMSP" />
              <el-option label="成都温江社区医疗中心" value="WenjinMSP" />
              <el-option label="华西医院" value="HuaxiMSP" />
            </el-select>
          </el-form-item>

          <!-- 医生需要填写科室和职称 -->
          <el-form-item v-if="formData.role === 'doctor'" prop="department">
            <span class="svg-container">
              <svg-icon icon-class="education" />
            </span>
            <el-input
              v-model="formData.department"
              placeholder="科室（如：内科）"
              name="department"
              type="text"
              auto-complete="off"
            />
          </el-form-item>

          <el-form-item v-if="formData.role === 'doctor'" prop="doctorTitle">
            <span class="svg-container">
              <svg-icon icon-class="star" />
            </span>
            <el-select v-model="formData.doctorTitle" placeholder="请选择职称" style="width: 100%;">
              <el-option label="主任医师" value="主任医师" />
              <el-option label="副主任医师" value="副主任医师" />
              <el-option label="主治医师" value="主治医师" />
              <el-option label="住院医师" value="住院医师" />
            </el-select>
          </el-form-item>

          <!-- 病人需要填写年龄和性别 -->
          <el-form-item v-if="formData.role === 'patient'" prop="age">
            <span class="svg-container">
              <i class="el-icon-user" style="font-size: 16px;"></i>
            </span>
            <el-input
              v-model.number="formData.age"
              placeholder="年龄"
              name="age"
              type="number"
              auto-complete="off"
            />
          </el-form-item>

          <el-form-item v-if="formData.role === 'patient'" prop="gender">
            <span class="svg-container">
              <i class="el-icon-male" style="font-size: 16px;"></i>
            </span>
            <el-select v-model="formData.gender" placeholder="请选择性别" style="width: 100%;">
              <el-option label="男" value="男" />
              <el-option label="女" value="女" />
            </el-select>
          </el-form-item>

          <el-button :loading="loading" type="primary" style="width:100%;margin-bottom:30px;" @click.native.prevent="handleRegister">注册</el-button>
        </el-tab-pane>
      </el-tabs>

      <div class="tips">
        <i class="el-icon-info"></i>
        <span>{{ activeTab === 'login' ? '请输入用户名和密码登录' : '请填写信息注册账号' }}</span>
      </div>

      <div class="footer-text">
        <p>© 2026 MedChain | 基于区块链的社区医疗管理系统</p>
      </div>

    </el-form>
    </div>
  </div>
</template>

<script>
import { register, loginWithPassword } from '@/api/accountV2'

export default {
  name: 'Login',
  data() {
    const validatePassword = (rule, value, callback) => {
      if (value && value.length < 6) {
        callback(new Error('密码不能少于6位'))
      } else {
        callback()
      }
    }
    const validateConfirmPassword = (rule, value, callback) => {
      if (value && value !== this.formData.regPassword) {
        callback(new Error('两次输入密码不一致'))
      } else {
        callback()
      }
    }
    return {
      activeTab: 'login',
      formData: {
        // 登录字段
        username: '',
        password: '',
        // 注册字段
        regUsername: '',
        regPassword: '',
        confirmPassword: '',
        role: '',
        organization: '',
        department: '',
        doctorTitle: '',
        age: null,
        gender: ''
      },
      formRules: {
        // 登录验证规则
        username: [{ required: true, trigger: 'blur', message: '请输入用户名' }],
        password: [{ required: true, trigger: 'blur', message: '请输入密码' }],
        // 注册验证规则
        regUsername: [{ required: true, trigger: 'blur', message: '请输入用户名' }],
        regPassword: [
          { required: true, trigger: 'blur', message: '请输入密码' },
          { validator: validatePassword, trigger: 'blur' }
        ],
        confirmPassword: [
          { required: true, trigger: 'blur', message: '请确认密码' },
          { validator: validateConfirmPassword, trigger: 'blur' }
        ],
        role: [{ required: true, trigger: 'change', message: '请选择角色' }]
      },
      loading: false,
      passwordType: 'password',
      redirect: undefined
    }
  },
  watch: {
    $route: {
      handler: function(route) {
        this.redirect = route.query && route.query.redirect
      },
      immediate: true
    }
  },
  methods: {
    showPwd() {
      if (this.passwordType === 'password') {
        this.passwordType = ''
      } else {
        this.passwordType = 'password'
      }
      this.$nextTick(() => {
        this.$refs.password.focus()
      })
    },
    handleLogin() {
      // 手动验证登录字段
      if (!this.formData.username) {
        this.$message.error('请输入用户名')
        return
      }
      if (!this.formData.username.trim()) {
        this.$message.error('用户名不能为空格')
        return
      }
      if (this.formData.username.length < 2) {
        this.$message.error('用户名至少2个字符')
        return
      }
      if (!this.formData.password) {
        this.$message.error('请输入密码')
        return
      }
      if (this.formData.password.length < 6) {
        this.$message.error('密码至少6位')
        return
      }
      
      this.loading = true
      this.$store.dispatch('account/loginWithPassword', {
        username: this.formData.username.trim(),
        password: this.formData.password
      }).then(() => {
        this.$message.success('登录成功，欢迎回来！')
        this.$router.push({ path: this.redirect || '/' })
        this.loading = false
      }).catch((error) => {
        console.error('登录错误:', error)
        this.loading = false
        
        // 根据错误信息提供详细反馈
        const errorMsg = error.message || error.toString()
        if (errorMsg.includes('用户不存在')) {
          this.$message.error('用户不存在，请先注册')
        } else if (errorMsg.includes('密码错误')) {
          this.$message.error('密码错误，请重新输入')
        } else if (errorMsg.includes('账户已被停用')) {
          this.$message.error('该账户已被停用，请联系管理员')
        } else if (errorMsg.includes('用户名或密码错误')) {
          this.$message.error('用户名或密码错误，请检查后重试')
        } else if (errorMsg.includes('网络') || errorMsg.includes('timeout') || errorMsg.includes('Network')) {
          this.$message.error('网络连接失败，请检查网络后重试')
        } else if (errorMsg.includes('服务器') || errorMsg.includes('500')) {
          this.$message.error('服务器错误，请稍后重试')
        } else {
          this.$message.error(errorMsg || '登录失败，请重试')
        }
      })
    },
    handleRegister() {
      // 手动验证注册字段
      if (!this.formData.regUsername) {
        this.$message.error('请输入用户名')
        return
      }
      if (!this.formData.regUsername.trim()) {
        this.$message.error('用户名不能为空格')
        return
      }
      if (this.formData.regUsername.length < 2) {
        this.$message.error('用户名至少2个字符')
        return
      }
      if (this.formData.regUsername.length > 20) {
        this.$message.error('用户名不能超过20个字符')
        return
      }
      // 用户名只能包含字母、数字、下划线
      if (!/^[a-zA-Z0-9_\u4e00-\u9fa5]+$/.test(this.formData.regUsername)) {
        this.$message.error('用户名只能包含字母、数字、下划线和中文')
        return
      }
      
      if (!this.formData.regPassword) {
        this.$message.error('请输入密码')
        return
      }
      if (this.formData.regPassword.length < 6) {
        this.$message.error('密码不能少于6位')
        return
      }
      if (this.formData.regPassword.length > 20) {
        this.$message.error('密码不能超过20位')
        return
      }
      // 密码强度检查：至少包含字母和数字
      if (!/[a-zA-Z]/.test(this.formData.regPassword) || !/[0-9]/.test(this.formData.regPassword)) {
        this.$message.warning('建议密码同时包含字母和数字，更安全')
      }
      
      if (!this.formData.confirmPassword) {
        this.$message.error('请确认密码')
        return
      }
      if (this.formData.regPassword !== this.formData.confirmPassword) {
        this.$message.error('两次输入密码不一致，请重新输入')
        return
      }
      
      if (!this.formData.role) {
        this.$message.error('请选择角色')
        return
      }
      
      // 医生需要选择组织和填写科室、职称
      if (this.formData.role === 'doctor') {
        if (!this.formData.organization) {
          this.$message.error('医生必须选择所属医院')
          return
        }
        if (!this.formData.department || !this.formData.department.trim()) {
          this.$message.error('医生必须填写科室信息')
          return
        }
        if (!this.formData.doctorTitle) {
          this.$message.error('医生必须选择职称')
          return
        }
      }
      
      // 病人需要填写年龄和性别
      if (this.formData.role === 'patient') {
        if (!this.formData.age || this.formData.age <= 0) {
          this.$message.error('请填写正确的年龄')
          return
        }
        if (this.formData.age > 150) {
          this.$message.error('年龄不能超过150岁')
          return
        }
        if (!this.formData.gender) {
          this.$message.error('请选择性别')
          return
        }
      }
      
      this.loading = true
      const roleMap = {
        'doctor': '医生',
        'patient': '病人',
        'drugstore': '药店'
      }
      
      const orgMap = {
        'TaobaoMSP': '协和医院',
        'JDMSP': '301医院',
        'WenjinMSP': '成都温江社区医疗中心',
        'HuaxiMSP': '华西医院',
        'RegCenterMSP': '监管中心'
      }
      
      const accountName = `${roleMap[this.formData.role]}-${this.formData.regUsername}`
      const organizationName = this.formData.organization ? orgMap[this.formData.organization] : ''
      
      register({
        account_name: accountName,
        username: this.formData.regUsername.trim(),
        password: this.formData.regPassword,
        role: roleMap[this.formData.role],
        organization: this.formData.organization || '',
        organization_name: organizationName,
        department: this.formData.department ? this.formData.department.trim() : '',
        doctor_title: this.formData.doctorTitle || '',
        age: this.formData.age || 0,
        gender: this.formData.gender || '',
        operator: 'system'
      }).then(() => {
        this.$message.success('注册成功！请使用新账号登录')
        this.activeTab = 'login'
        this.formData.username = this.formData.regUsername
        this.formData.regUsername = ''
        this.formData.regPassword = ''
        this.formData.confirmPassword = ''
        this.formData.role = ''
        this.formData.organization = ''
        this.formData.department = ''
        this.formData.doctorTitle = ''
        this.loading = false
      }).catch((error) => {
        console.error('注册错误:', error)
        this.loading = false
        
        // 根据错误信息提供详细反馈
        const errorMsg = error.message || error.toString()
        if (errorMsg.includes('用户名已存在') || errorMsg.includes('already exists')) {
          this.$message.error('该用户名已被注册，请更换用户名')
        } else if (errorMsg.includes('参数') || errorMsg.includes('空值')) {
          this.$message.error('请填写完整的注册信息')
        } else if (errorMsg.includes('组织')) {
          this.$message.error('请选择所属组织')
        } else if (errorMsg.includes('网络') || errorMsg.includes('timeout') || errorMsg.includes('Network')) {
          this.$message.error('网络连接失败，请检查网络后重试')
        } else if (errorMsg.includes('区块链') || errorMsg.includes('blockchain')) {
          this.$message.error('区块链服务暂时不可用，请稍后重试或联系管理员')
        } else if (errorMsg.includes('服务器') || errorMsg.includes('500')) {
          this.$message.error('服务器错误，请稍后重试')
        } else {
          this.$message.error(errorMsg || '注册失败，请重试')
        }
      })
    },
    handleRoleChange() {
      // 切换角色时清空组织相关字段
      if (this.formData.role !== 'doctor') {
        this.formData.organization = ''
        this.formData.department = ''
        this.formData.doctorTitle = ''
      }
    }
  }
}
</script>

<style lang="scss" scoped>
$primary: #409EFF;
$bg: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
$dark_gray: #889aa4;
$light_gray: #eee;

.login-container {
  min-height: 100vh;
  width: 100%;
  display: flex;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);

  // 背景装饰
  .background-decoration {
    position: absolute;
    width: 100%;
    height: 100%;
    overflow: hidden;
    z-index: 0;

    .circle {
      position: absolute;
      border-radius: 50%;
      background: rgba(255, 255, 255, 0.1);
      animation: float 20s infinite ease-in-out;

      &.circle-1 {
        width: 300px;
        height: 300px;
        top: -100px;
        left: -100px;
        animation-delay: 0s;
      }

      &.circle-2 {
        width: 200px;
        height: 200px;
        bottom: -50px;
        right: 10%;
        animation-delay: 5s;
      }

      &.circle-3 {
        width: 150px;
        height: 150px;
        top: 50%;
        left: 5%;
        animation-delay: 10s;
      }
    }
  }

  @keyframes float {
    0%, 100% {
      transform: translateY(0) rotate(0deg);
    }
    50% {
      transform: translateY(-20px) rotate(180deg);
    }
  }

  // 左侧信息区域
  .info-section {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 60px;
    position: relative;
    z-index: 1;

    .info-content {
      max-width: 500px;
      color: white;

      .logo-section {
        margin-bottom: 60px;
        text-align: center;

        .logo-icon {
          width: 120px;
          height: 120px;
          display: flex;
          align-items: center;
          justify-content: center;
          margin: 0 auto 20px;
          position: relative;

          .logo-image {
            width: 100px;
            height: 100px;
            object-fit: contain;
            filter: drop-shadow(0 4px 12px rgba(0, 0, 0, 0.2));
            transition: all 0.3s ease;
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);

            &:hover {
              transform: translate(-50%, -50%) scale(1.05);
              filter: drop-shadow(0 6px 16px rgba(0, 0, 0, 0.3));
            }
          }

          // 默认显示详细版Logo
          .detailed-logo {
            display: block;
          }

          .simple-logo {
            display: none;
          }

          // 小屏幕时切换到精简版Logo
          @media (max-width: 1200px) {
            .detailed-logo {
              display: none;
            }

            .simple-logo {
              display: block;
            }
          }

          svg {
            width: 90px;
            height: 90px;
            filter: drop-shadow(0 4px 12px rgba(0, 0, 0, 0.2));
            transition: all 0.3s ease;

            &:hover {
              transform: scale(1.05);
              filter: drop-shadow(0 6px 16px rgba(0, 0, 0, 0.3));
            }
          }

          i {
            font-size: 40px;
            color: white;
          }
        }

        .system-title {
          font-size: 36px;
          font-weight: 700;
          margin: 0 0 10px 0;
          text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);

          // 响应式字体大小
          @media (max-width: 1400px) {
            font-size: 32px;
          }

          @media (max-width: 1200px) {
            font-size: 28px;
          }
        }

        .system-subtitle {
          font-size: 16px;
          opacity: 0.9;
          margin: 0;

          // 响应式字体大小
          @media (max-width: 1400px) {
            font-size: 15px;
          }

          @media (max-width: 1200px) {
            font-size: 14px;
          }
        }
      }

      .features {
        .feature-item {
          display: flex;
          align-items: center;
          margin-bottom: 30px;
          padding: 20px;
          background: rgba(255, 255, 255, 0.1);
          border-radius: 12px;
          backdrop-filter: blur(10px);
          transition: all 0.3s;

          &:hover {
            background: rgba(255, 255, 255, 0.15);
            transform: translateX(10px);
          }

          i {
            font-size: 32px;
            margin-right: 20px;
            opacity: 0.9;
          }

          .feature-text {
            h4 {
              margin: 0 0 5px 0;
              font-size: 18px;
              font-weight: 600;
            }

            p {
              margin: 0;
              font-size: 14px;
              opacity: 0.8;
            }
          }
        }
      }
    }
  }

  // 右侧表单区域
  .form-section {
    width: 480px;
    background: white;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: -10px 0 30px rgba(0, 0, 0, 0.1);
    position: relative;
    z-index: 1;
  }

  .login-form {
    width: 100%;
    max-width: 400px;
    padding: 40px;

    .form-header {
      margin-bottom: 40px;
      text-align: center;

      .form-logo {
        margin-bottom: 20px;
        padding: 8px;
        background: linear-gradient(135deg, rgba(102, 126, 234, 0.1), rgba(118, 75, 162, 0.1));
        border-radius: 12px;
        display: inline-block;
        
        svg {
          width: 45px;
          height: 45px;
          filter: drop-shadow(0 2px 6px rgba(102, 126, 234, 0.3));
          transition: all 0.3s ease;

          &:hover {
            transform: scale(1.1);
            filter: drop-shadow(0 4px 8px rgba(102, 126, 234, 0.4));
          }
        }
        
        .form-logo-image {
          width: 45px;
          height: 45px;
          object-fit: contain;
          filter: brightness(1.1) contrast(1.1) drop-shadow(0 2px 6px rgba(102, 126, 234, 0.3));
          transition: all 0.3s ease;

          &:hover {
            transform: scale(1.1);
            filter: brightness(1.2) contrast(1.2) drop-shadow(0 4px 8px rgba(102, 126, 234, 0.4));
          }
        }

        // 在表单区域始终使用精简版Logo
        // 这样保持表单区域的简洁性
      }

      .form-title {
        font-size: 28px;
        font-weight: 700;
        color: #303133;
        margin: 0 0 10px 0;
      }

      .form-subtitle {
        font-size: 14px;
        color: #909399;
        margin: 0;
      }
    }
  }

  .login-tabs {
    margin-bottom: 30px;
    
    ::v-deep .el-tabs__header {
      margin: 0 0 30px;
    }
    
    ::v-deep .el-tabs__nav-wrap::after {
      background-color: #EBEEF5;
    }
    
    ::v-deep .el-tabs__item {
      color: #909399;
      font-size: 16px;
      font-weight: 500;
      padding: 0 30px;
      
      &.is-active {
        color: $primary;
      }

      &:hover {
        color: $primary;
      }
    }
    
    ::v-deep .el-tabs__active-bar {
      background-color: $primary;
      height: 3px;
    }
  }

  ::v-deep .el-form-item {
    margin-bottom: 24px;
    border: 1px solid #DCDFE6;
    background: #F5F7FA;
    border-radius: 8px;
    transition: all 0.3s;

    &:hover {
      border-color: $primary;
      background: white;
    }

    &:focus-within {
      border-color: $primary;
      background: white;
      box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
    }
  }

  ::v-deep .el-input {
    display: inline-block;
    height: 50px;
    width: 85%;

    input {
      background: transparent;
      border: 0px;
      -webkit-appearance: none;
      border-radius: 0px;
      padding: 12px 5px 12px 15px;
      color: #303133;
      height: 50px;
      caret-color: $primary;
      font-size: 14px;

      &::placeholder {
        color: #C0C4CC;
      }

      &:-webkit-autofill {
        box-shadow: 0 0 0px 1000px #F5F7FA inset !important;
        -webkit-text-fill-color: #303133 !important;
      }
    }
  }

  ::v-deep .el-select {
    width: 100%;
    
    .el-input__inner {
      background: transparent;
      border: 0px;
      color: #303133;
      height: 50px;
      line-height: 50px;
    }
  }

  ::v-deep .el-button--primary {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border: none;
    height: 50px;
    font-size: 16px;
    font-weight: 500;
    border-radius: 8px;
    transition: all 0.3s;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 16px rgba(102, 126, 234, 0.4);
    }

    &:active {
      transform: translateY(0);
    }
  }

  .tips {
    font-size: 13px;
    color: #909399;
    margin-bottom: 20px;
    text-align: center;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;

    i {
      font-size: 16px;
    }
  }

  .footer-text {
    text-align: center;
    margin-top: 30px;
    padding-top: 20px;
    border-top: 1px solid #EBEEF5;

    p {
      font-size: 12px;
      color: #C0C4CC;
      margin: 0;
    }
  }

  .svg-container {
    padding: 6px 5px 6px 15px;
    color: #909399;
    vertical-align: middle;
    width: 30px;
    display: inline-block;
  }

  .show-pwd {
    position: absolute;
    right: 15px;
    top: 15px;
    font-size: 18px;
    color: #909399;
    cursor: pointer;
    user-select: none;
    transition: color 0.3s;

    &:hover {
      color: $primary;
    }
  }

  // 响应式设计
  @media (max-width: 1024px) {
    .info-section {
      display: none;
    }

    .form-section {
      width: 100%;
    }
  }

  @media (max-width: 768px) {
    .login-form {
      padding: 20px;
    }

    .form-header {
      .form-logo {
        .form-logo-image {
          width: 40px;
          height: 40px;
        }
      }

      .form-title {
        font-size: 24px;
      }
    }
  }

  @media (max-width: 480px) {
    .login-form {
      padding: 15px;
    }

    .form-header {
      margin-bottom: 30px;

      .form-logo {
        padding: 6px;

        .form-logo-image {
          width: 35px;
          height: 35px;
        }
      }

      .form-title {
        font-size: 22px;
      }

      .form-subtitle {
        font-size: 13px;
      }
    }
  }
}
</style>
