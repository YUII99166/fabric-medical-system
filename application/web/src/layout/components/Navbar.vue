<template>
  <div class="navbar">
    <hamburger :is-active="sidebar.opened" class="hamburger-container" @toggleClick="toggleSideBar" />

    <breadcrumb class="breadcrumb-container" />

    <div class="right-menu">
      <el-dropdown class="avatar-container" trigger="click">
        <div class="avatar-wrapper">
          <img :src="avatarUrl" class="user-avatar">
          <div class="user-info">
            <span class="user-name">{{ displayAccountName }}</span>
            <span class="user-role">{{ roleText }}</span>
          </div>
          <i class="el-icon-caret-bottom" />
        </div>
        <el-dropdown-menu slot="dropdown" class="user-dropdown">
          <el-dropdown-item>
            <i class="el-icon-user"></i> {{ displayAccountName }}
          </el-dropdown-item>
          <el-dropdown-item>
            <i class="el-icon-postcard"></i> ID: {{ account_id }}
          </el-dropdown-item>
          <el-dropdown-item>
            <i class="el-icon-s-custom"></i> {{ roleText }}
          </el-dropdown-item>
          <el-dropdown-item divided @click.native="logout">
            <i class="el-icon-switch-button"></i> 切换账户
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import Breadcrumb from '@/components/Breadcrumb'
import Hamburger from '@/components/Hamburger'
import { applyTheme } from '@/utils/theme'

export default {
  components: {
    Breadcrumb,
    Hamburger
  },
  data() {
    return {
      syncCheckInterval: null,
      syncAttempts: 0,
      maxSyncAttempts: 10 // 最多尝试10次自动同步
    }
  },
  computed: {
    ...mapGetters([
      'sidebar',
      'account_id',
      'account_name',
      'roles'
    ]),
    // 显示用的账户名称（根据角色调整）
    displayAccountName() {
      // 如果是管理员角色，将"医生-xxx"替换为"管理员-xxx"
      if (this.roles && this.roles[0] === 'admin' && this.account_name) {
        return this.account_name.replace(/^医生-/, '管理员-')
      }
      return this.account_name
    },
    // 根据角色返回对应的头像
    avatarUrl() {
      if (!this.roles || this.roles.length === 0) {
        return '/image/病人-男.png'
      }
      
      const role = this.roles[0]
      const avatarMap = {
        'admin': '/image/系统管理员.png',
        'doctor': '/image/男医生.png',
        'patient': '/image/病人-男.png',
        'drugstore': '/image/药店-copy.png'
      }
      
      return avatarMap[role] || '/image/病人-男.png'
    },
    // 角色中文名称
    roleText() {
      if (!this.roles || this.roles.length === 0) {
        return '未知'
      }
      
      const role = this.roles[0]
      const roleMap = {
        'admin': '管理员',
        'doctor': '医生',
        'patient': '病人',
        'drugstore': '药店'
      }
      
      return roleMap[role] || '未知'
    }
  },
  watch: {
    // 监听角色变化，应用对应主题
    roles: {
      handler(newRoles) {
        if (newRoles && newRoles.length > 0) {
          applyTheme(newRoles[0])
        }
      },
      immediate: true
    },
    // 监听account_id变化，如果未同步则自动尝试同步
    account_id: {
      handler(newAccountId) {
        if (newAccountId === '未同步到区块链' && this.syncAttempts < this.maxSyncAttempts) {
          console.log('检测到账户未同步，将在3秒后尝试自动同步...')
          // 延迟3秒后尝试同步，给区块链一些时间完成注册
          setTimeout(() => {
            this.attemptAutoSync()
          }, 3000)
        }
      },
      immediate: true
    }
  },
  mounted() {
    // 组件挂载后，如果账户未同步，启动定期检查
    if (this.account_id === '未同步到区块链') {
      this.startSyncCheck()
    }
  },
  beforeDestroy() {
    // 组件销毁前清除定时器
    this.stopSyncCheck()
  },
  methods: {
    toggleSideBar() {
      this.$store.dispatch('app/toggleSideBar')
    },
    async logout() {
      await this.$store.dispatch('account/logout')
      this.$router.push(`/login?redirect=${this.$route.fullPath}`)
    },
    // 启动定期同步检查
    startSyncCheck() {
      // 如果已经有定时器在运行，先清除
      this.stopSyncCheck()
      
      // 每10秒检查一次，最多尝试10次（总共100秒）
      this.syncCheckInterval = setInterval(() => {
        if (this.account_id === '未同步到区块链' && this.syncAttempts < this.maxSyncAttempts) {
          this.attemptAutoSync()
        } else {
          // 如果已同步或达到最大尝试次数，停止检查
          this.stopSyncCheck()
        }
      }, 10000) // 10秒
    },
    // 停止定期同步检查
    stopSyncCheck() {
      if (this.syncCheckInterval) {
        clearInterval(this.syncCheckInterval)
        this.syncCheckInterval = null
      }
    },
    // 尝试自动同步
    async attemptAutoSync() {
      this.syncAttempts++
      console.log(`尝试自动同步账户ID (第${this.syncAttempts}次)...`)
      
      try {
        const success = await this.$store.dispatch('account/syncAccountId')
        if (success) {
          console.log('✅ 自动同步成功')
          this.$message.success('账户已成功同步到区块链')
          this.stopSyncCheck()
        } else {
          console.log(`⚠️ 自动同步失败 (第${this.syncAttempts}/${this.maxSyncAttempts}次)`)
          
          // 如果达到最大尝试次数，提示用户
          if (this.syncAttempts >= this.maxSyncAttempts) {
            console.warn('已达到最大同步尝试次数')
            this.$message.warning('账户同步失败，部分功能可能受限。请稍后刷新页面重试。')
          }
        }
      } catch (error) {
        console.error('自动同步出错:', error)
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.navbar {
  height: 60px;
  overflow: hidden;
  position: relative;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0,0,0,.08);
  display: flex;
  align-items: center;

  .hamburger-container {
    line-height: 60px;
    height: 100%;
    float: left;
    cursor: pointer;
    transition: background .3s;
    -webkit-tap-highlight-color:transparent;
    padding: 0 15px;

    &:hover {
      background: rgba(0, 0, 0, .025)
    }
  }

  .breadcrumb-container {
    float: left;
    line-height: 60px;
  }

  .right-menu {
    float: right;
    height: 100%;
    display: flex;
    align-items: center;
    padding-right: 20px;

    &:focus {
      outline: none;
    }

    .right-menu-item {
      display: inline-block;
      padding: 0 12px;
      height: 100%;
      font-size: 18px;
      color: #5a5e66;
      vertical-align: text-bottom;

      &.hover-effect {
        cursor: pointer;
        transition: background .3s;

        &:hover {
          background: rgba(0, 0, 0, .025)
        }
      }
    }

    .avatar-container {
      margin-right: 10px;
      cursor: pointer;

      .avatar-wrapper {
        position: relative;
        display: flex;
        align-items: center;
        padding: 5px 10px;
        border-radius: 8px;
        transition: all .3s;

        &:hover {
          background: var(--theme-light-bg, #ecf5ff);
        }

        .user-avatar {
          width: 40px;
          height: 40px;
          border-radius: 8px;
          border: 2px solid var(--theme-primary, #409EFF);
          object-fit: contain;
          object-position: center;
          background: #fff;
          transition: all .3s;

          &:hover {
            transform: scale(1.1);
          }
        }

        .user-info {
          margin-left: 10px;
          display: flex;
          flex-direction: column;
          
          .user-name {
            font-size: 14px;
            font-weight: 500;
            color: #303133;
            line-height: 1.2;
          }

          .user-role {
            font-size: 12px;
            color: #909399;
            margin-top: 2px;
          }
        }

        .el-icon-caret-bottom {
          margin-left: 8px;
          font-size: 12px;
          color: #909399;
          transition: transform .3s;
        }

        &:hover .el-icon-caret-bottom {
          transform: rotate(180deg);
        }
      }
    }
  }
}

// 下拉菜单样式优化
::v-deep .user-dropdown {
  margin-top: 10px !important;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0,0,0,.1);

  .el-dropdown-menu__item {
    padding: 12px 20px;
    font-size: 14px;
    
    &:hover {
      background: var(--theme-light-bg, #ecf5ff);
      color: var(--theme-primary, #409EFF);
    }

    &.is-divided {
      border-top: 1px solid #ebeef5;
      margin-top: 6px;
    }
  }
}
</style>
