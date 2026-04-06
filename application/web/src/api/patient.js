import request from '@/utils/request'
import store from '@/store'

/**
 * 获取健康档案统计
 * @param {string} timeRange - 时间范围：1m, 3m, 1y, all
 */
export function getHealthProfile(timeRange = 'all') {
  const accountId = store.state.account.account_id
  return request({
    url: '/patient/health-profile',
    method: 'get',
    params: { 
      accountId: accountId,
      timeRange 
    }
  })
}
