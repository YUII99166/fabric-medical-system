import request from '@/utils/request'

// 记录病历访问
export function recordPrescriptionAccess(data) {
  return request({
    url: '/recordPrescriptionAccess',
    method: 'post',
    data
  })
}

// 获取我的访问日志
export function getMyAccessLogs(params) {
  return request({
    url: '/getMyAccessLogs',
    method: 'post',
    data: params
  })
}

// 获取访问统计
export function getAccessStatistics(patientId) {
  return request({
    url: '/getAccessStatistics',
    method: 'get',
    params: { patient_id: patientId }
  })
}
