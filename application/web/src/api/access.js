import request from '@/utils/request'

// 申请授权
export function requestAccess(data) {
  return request({
    url: '/requestAccess',
    method: 'post',
    data
  })
}

// 审批授权
export function approveAccess(data) {
  return request({
    url: '/approveAccess',
    method: 'post',
    data
  })
}

// 查询授权请求列表
export function queryAccessRequests(data) {
  return request({
    url: '/queryAccessRequests',
    method: 'post',
    data
  })
}

// 根据患者姓名或ID查询病历
export function queryPrescriptionsByPatient(data) {
  return request({
    url: '/queryPrescriptionsByPatient',
    method: 'post',
    data
  })
}
