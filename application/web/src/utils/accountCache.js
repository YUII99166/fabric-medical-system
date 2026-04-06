/**
 * 账户信息缓存工具
 * 用于减少重复的账户列表查询请求
 */

import { queryAccountList } from '@/api/accountV2'

// 缓存对象
const cache = {
  data: null, // 账户映射数据
  timestamp: null, // 缓存时间戳
  loading: false, // 是否正在加载
  callbacks: [] // 等待回调队列
}

// 缓存有效期（5分钟）
const CACHE_DURATION = 5 * 60 * 1000

/**
 * 获取账户映射
 * @returns {Promise<Object>} 账户ID到账户信息的映射
 */
export function getAccountMap() {
  return new Promise((resolve, reject) => {
    // 检查缓存是否有效
    if (cache.data && cache.timestamp && (Date.now() - cache.timestamp < CACHE_DURATION)) {
      console.log('使用缓存的账户数据')
      resolve(cache.data)
      return
    }

    // 如果正在加载，加入等待队列
    if (cache.loading) {
      console.log('账户数据加载中，加入等待队列')
      cache.callbacks.push({ resolve, reject })
      return
    }

    // 开始加载
    cache.loading = true
    console.log('开始加载账户数据...')

    queryAccountList()
      .then(response => {
        let accountList = []
        if (response && response.code === 200 && response.data) {
          accountList = response.data
        } else if (Array.isArray(response)) {
          accountList = response
        }

        if (Array.isArray(accountList)) {
          // 建立账户ID到账户信息的映射
          const accountMap = {}
          accountList.forEach(account => {
            accountMap[account.account_id] = {
              name: account.account_name,
              username: account.username,
              role: account.role,
              organization_name: account.organization_name,
              organization: account.organization
            }
          })

          // 更新缓存
          cache.data = accountMap
          cache.timestamp = Date.now()
          console.log('账户数据加载完成，缓存已更新:', Object.keys(accountMap).length, '个账户')

          // 通知所有等待的回调
          cache.callbacks.forEach(cb => cb.resolve(accountMap))
          cache.callbacks = []

          resolve(accountMap)
        } else {
          const error = new Error('账户数据格式错误')
          cache.callbacks.forEach(cb => cb.reject(error))
          cache.callbacks = []
          reject(error)
        }
      })
      .catch(error => {
        console.error('加载账户数据失败:', error)
        cache.callbacks.forEach(cb => cb.reject(error))
        cache.callbacks = []
        reject(error)
      })
      .finally(() => {
        cache.loading = false
      })
  })
}

/**
 * 清除缓存
 */
export function clearAccountCache() {
  cache.data = null
  cache.timestamp = null
  console.log('账户缓存已清除')
}

/**
 * 获取账户名称
 * @param {string} accountId 账户ID
 * @param {Object} accountMap 账户映射（可选，如果已有）
 * @returns {string} 账户名称
 */
export function getAccountName(accountId, accountMap = null) {
  if (!accountId) return '未知'
  if (accountMap && accountMap[accountId]) {
    return accountMap[accountId].name || accountId
  }
  return accountId
}
