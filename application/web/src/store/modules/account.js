import {
  login,
  loginWithPassword,
  getUserInfo
} from '@/api/accountV2'
import {
  getToken,
  setToken,
  removeToken
} from '@/utils/auth'
import {
  resetRouter
} from '@/router'

const getDefaultState = () => {
  return {
    token: getToken(),
    account_id: '',
    account_name: '',
    roles: []
  }
}

const state = getDefaultState()

const mutations = {
  RESET_STATE: (state) => {
    Object.assign(state, getDefaultState())
  },
  SET_TOKEN: (state, token) => {
    state.token = token
  },
  SET_ACCOUNTID: (state, account_id) => {
    state.account_id = account_id
  },
  SET_USERNAME: (state, account_name) => {
    state.account_name = account_name
  },
  SET_ROLES: (state, roles) => {
    state.roles = roles
  }
}

const actions = {
  login({
    commit
  }, account_id) {
    return new Promise((resolve, reject) => {
      login({
        args: [{
          account_id: account_id
        }]
      }).then(response => {
        commit('SET_TOKEN', response[0].account_id)
        setToken(response[0].account_id)
        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },
  // 密码登录
  loginWithPassword({
    commit
  }, userInfo) {
    const { username, password } = userInfo
    return new Promise((resolve, reject) => {
      loginWithPassword({
        username: username.trim(),
        password: password
      }).then(response => {
        console.log('登录响应:', response)
        // response 现在是完整的响应对象 {code, msg, data}
        if (response && response.code === 200 && response.data && response.data.id) {
          commit('SET_TOKEN', response.data.id)
          setToken(response.data.id)
          resolve()
        } else {
          reject(new Error('登录失败：无效的响应'))
        }
      }).catch(error => {
        reject(error)
      })
    })
  },
  // get user info
  getInfo({
    commit,
    state
  }) {
    return new Promise((resolve, reject) => {
      // 使用新的 getUserInfo API，通过数据库 ID 获取用户信息
      getUserInfo({
        id: parseInt(state.token)
      }).then(async response => {
        console.log('getUserInfo 响应:', response)
        
        // response 现在是完整的响应对象 {code, msg, data}
        if (!response || response.code !== 200 || !response.data) {
          reject('获取用户信息失败')
          return
        }
        
        const userInfo = response.data
        
        var roles
        // 根据 role 字段判断角色
        const role = userInfo.role
        
        // 特殊处理：监管中心的医生角色映射为管理员
        // 需要先查询区块链获取组织信息
        let isRegCenter = false
        try {
          const { queryAccountList } = await import('@/api/accountV2')
          const accountListResponse = await queryAccountList()
          const accountList = accountListResponse.code === 200 ? accountListResponse.data : accountListResponse
          
          if (accountList && Array.isArray(accountList)) {
            const userAccount = accountList.find(acc => 
              acc.username === userInfo.username || acc.account_name === userInfo.account_name
            )
            if (userAccount && userAccount.organization === 'RegCenterMSP') {
              isRegCenter = true
            }
          }
        } catch (error) {
          console.error('查询组织信息失败:', error)
        }
        
        // 角色映射逻辑
        if (isRegCenter) {
          // 监管中心用户统一映射为管理员
          roles = ['admin']
          console.log('✅ 监管中心用户，映射为管理员角色')
        } else if (role === '管理员') {
          roles = ['admin']
        } else if (role === '医生'){
          roles = ['doctor']
        } else if (role === '病人'){
          roles = ['patient']
        } else if (role === '药店'){
          roles = ['drugstore']
        } else {
          // 默认角色
          roles = ['patient']
        }
        
        commit('SET_ROLES', roles)
        commit('SET_USERNAME', userInfo.account_name)
        
        // 从区块链查询该用户的account_id（异步，不阻塞登录）
        // 使用重试机制，最多尝试5次
        try {
          const { queryAccountList } = await import('@/api/accountV2')
          console.log('开始查询区块链账户列表...')
          
          let accountListResponse = null
          let lastError = null
          const maxRetries = 5
          
          for (let attempt = 1; attempt <= maxRetries; attempt++) {
            try {
              accountListResponse = await queryAccountList()
              console.log(`区块链账户列表响应(第${attempt}次尝试):`, accountListResponse)
              // 如果成功获取数据，跳出重试循环
              if (accountListResponse && (accountListResponse.code === 200 || Array.isArray(accountListResponse))) {
                break
              }
            } catch (error) {
              lastError = error
              console.warn(`查询区块链失败(第${attempt}次尝试):`, error.message)
              
              // 如果不是最后一次尝试，等待后重试
              if (attempt < maxRetries) {
                const delay = 500 * attempt // 递增延迟: 500ms, 1000ms, 1500ms, 2000ms, 2500ms
                console.log(`等待 ${delay}ms 后重试...`)
                await new Promise(resolve => setTimeout(resolve, delay))
              }
            }
          }
          
          // 如果所有重试都失败，抛出最后一个错误
          if (!accountListResponse || (!accountListResponse.code && !Array.isArray(accountListResponse))) {
            throw lastError || new Error('查询区块链失败')
          }
          
          console.log('当前用户 username:', userInfo.username, 'account_name:', userInfo.account_name)
          
          // accountListResponse 现在也是完整的响应对象
          const accountList = accountListResponse.code === 200 ? accountListResponse.data : accountListResponse
          
          if (accountList && Array.isArray(accountList)) {
            console.log('账户列表数量:', accountList.length)
            console.log('账户列表示例:', accountList.slice(0, 3))
            
            // 优先用 username 匹配，如果找不到则用 account_name 匹配
            let userAccount = accountList.find(acc => {
              console.log(`比对username: acc.username="${acc.username}" vs userInfo.username="${userInfo.username}"`)
              return acc.username === userInfo.username
            })
            
            if (!userAccount) {
              console.log('username匹配失败，尝试用account_name匹配')
              userAccount = accountList.find(acc => {
                console.log(`比对account_name: acc.account_name="${acc.account_name}" vs userInfo.account_name="${userInfo.account_name}"`)
                return acc.account_name === userInfo.account_name
              })
            }
            
            // 如果还是找不到，尝试更宽松的匹配（忽略大小写和空格）
            if (!userAccount) {
              console.log('精确匹配失败，尝试宽松匹配')
              const normalizeString = (str) => str ? str.toLowerCase().trim() : ''
              const normalizedUsername = normalizeString(userInfo.username)
              const normalizedAccountName = normalizeString(userInfo.account_name)
              
              userAccount = accountList.find(acc => {
                const accUsername = normalizeString(acc.username)
                const accAccountName = normalizeString(acc.account_name)
                console.log(`宽松匹配: "${accUsername}" vs "${normalizedUsername}" 或 "${accAccountName}" vs "${normalizedAccountName}"`)
                return accUsername === normalizedUsername || accAccountName === normalizedAccountName
              })
            }
            
            console.log('找到的用户账户:', userAccount)
            
            if (userAccount && userAccount.account_id) {
              commit('SET_ACCOUNTID', userAccount.account_id)
              console.log('✅ 从区块链获取account_id成功:', userAccount.account_id)
              
              // 保存完整用户信息到sessionStorage，包含区块链account_id和组织信息
              const fullUserInfo = {
                id: userInfo.id,
                username: userInfo.username,
                account_name: userInfo.account_name,
                role: userInfo.role,
                account_id: userAccount.account_id,
                organization: userAccount.organization || '',
                organization_name: userAccount.organization_name || '',
                department: userAccount.department || '',
                doctor_title: userAccount.doctor_title || '',
                age: userAccount.age || 0,
                gender: userAccount.gender || ''
              }
              sessionStorage.setItem('userInfo', JSON.stringify(fullUserInfo))
              console.log('用户信息已保存到sessionStorage:', fullUserInfo)
            } else {
              // 如果区块链中找不到，显示提示信息
              console.warn('⚠️ 区块链中未找到用户')
              console.warn('数据库用户信息:', userInfo)
              console.warn('区块链账户列表:', accountList.map(acc => ({
                username: acc.username,
                account_name: acc.account_name,
                account_id: acc.account_id
              })))
              
              const placeholderAccountId = '未同步到区块链'
              commit('SET_ACCOUNTID', placeholderAccountId)
              
              // 保存用户信息（标记为未同步）
              const fullUserInfo = {
                id: userInfo.id,
                username: userInfo.username,
                account_name: userInfo.account_name,
                role: userInfo.role,
                account_id: placeholderAccountId
              }
              sessionStorage.setItem('userInfo', JSON.stringify(fullUserInfo))
            }
          } else {
            console.warn('⚠️ 查询区块链账户列表失败，返回数据格式错误')
            console.warn('accountListResponse:', accountListResponse)
            
            const placeholderAccountId = '未同步到区块链'
            commit('SET_ACCOUNTID', placeholderAccountId)
            
            // 保存用户信息（标记为未同步）
            const fullUserInfo = {
              id: userInfo.id,
              username: userInfo.username,
              account_name: userInfo.account_name,
              role: userInfo.role,
              account_id: placeholderAccountId
            }
            sessionStorage.setItem('userInfo', JSON.stringify(fullUserInfo))
          }
        } catch (error) {
          console.error('❌ 查询区块链账户ID失败:', error)
          
          const placeholderAccountId = '未同步到区块链'
          commit('SET_ACCOUNTID', placeholderAccountId)
          
          // 保存用户信息（标记为未同步）
          const fullUserInfo = {
            id: userInfo.id,
            username: userInfo.username,
            account_name: userInfo.account_name,
            role: userInfo.role,
            account_id: placeholderAccountId
          }
          sessionStorage.setItem('userInfo', JSON.stringify(fullUserInfo))
        }
        
        resolve(roles)
      }).catch(error => {
        reject(error)
      })
    })
  },
  logout({
    commit
  }) {
    return new Promise(resolve => {
      removeToken()
      resetRouter()
      commit('RESET_STATE')
      sessionStorage.removeItem('userInfo') // 清除用户信息
      resolve()
    })
  },

  resetToken({
    commit
  }) {
    return new Promise(resolve => {
      removeToken()
      commit('RESET_STATE')
      resolve()
    })
  },

  // 从区块链同步账户ID（用于修复"未同步到区块链"状态）
  async syncAccountId({ commit, state }) {
    try {
      const userInfo = JSON.parse(sessionStorage.getItem('userInfo') || '{}')
      if (!userInfo.username) {
        console.warn('无法同步：用户信息不完整')
        return false
      }

      // 如果已经有有效的account_id，不需要同步
      if (state.account_id && state.account_id !== '未同步到区块链') {
        console.log('账户已同步，无需重新同步')
        return true
      }

      console.log('开始同步账户ID...')
      const { queryAccountList } = await import('@/api/accountV2')
      
      // 使用重试机制
      let accountListResponse = null
      const maxRetries = 3
      
      for (let attempt = 1; attempt <= maxRetries; attempt++) {
        try {
          accountListResponse = await queryAccountList()
          if (accountListResponse && (accountListResponse.code === 200 || Array.isArray(accountListResponse))) {
            break
          }
        } catch (error) {
          console.warn(`同步失败(第${attempt}次尝试):`, error.message)
          if (attempt < maxRetries) {
            await new Promise(resolve => setTimeout(resolve, 500 * attempt))
          }
        }
      }

      if (!accountListResponse) {
        console.error('同步失败：无法查询区块链')
        return false
      }

      const accountList = accountListResponse.code === 200 ? accountListResponse.data : accountListResponse
      
      if (accountList && Array.isArray(accountList)) {
        const userAccount = accountList.find(acc => 
          acc.username === userInfo.username || acc.account_name === userInfo.account_name
        )
        
        if (userAccount && userAccount.account_id) {
          commit('SET_ACCOUNTID', userAccount.account_id)
          
          // 更新sessionStorage
          const updatedUserInfo = {
            ...userInfo,
            account_id: userAccount.account_id,
            organization: userAccount.organization || userInfo.organization,
            organization_name: userAccount.organization_name || userInfo.organization_name
          }
          sessionStorage.setItem('userInfo', JSON.stringify(updatedUserInfo))
          
          console.log('✅ 账户ID同步成功:', userAccount.account_id)
          return true
        }
      }
      
      console.warn('同步失败：区块链中未找到用户')
      return false
    } catch (error) {
      console.error('同步账户ID失败:', error)
      return false
    }
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
