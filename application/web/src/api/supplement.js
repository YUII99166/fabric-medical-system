import request from '@/utils/request'

// 添加补充诊疗记录
export function addSupplementRecord(data) {
  return request({
    url: '/addSupplementRecord',
    method: 'post',
    data
  })
}

// 查询补充记录列表
export function querySupplementRecords(data) {
  return request({
    url: '/querySupplementRecords',
    method: 'post',
    data
  })
}

// 查询完整病历历史（原始病历+补充记录）
export function queryFullMedicalHistory(data) {
  return request({
    url: '/queryFullMedicalHistory',
    method: 'post',
    data
  })
}
