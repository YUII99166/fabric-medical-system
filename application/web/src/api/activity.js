import request from '@/utils/request'

// 获取联盟活动列表
export function getAllianceActivities(params) {
  return request({
    url: '/getAllianceActivities',
    method: 'post',
    data: params
  })
}

// 获取活动统计数据
export function getActivityStatistics() {
  return request({
    url: '/getActivityStatistics',
    method: 'post'
  })
}
