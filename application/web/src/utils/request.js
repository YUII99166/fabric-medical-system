import axios from 'axios'
import {
  MessageBox,
  Message
} from 'element-ui'

const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API,
  timeout: 5000
})

service.interceptors.response.use(
  response => {
    const res = response.data
    console.log('✅ 拦截器收到响应:', res)
    
    // 如果响应有 code 字段，说明是标准格式
    if (res.code !== undefined) {
      if (res.code !== 200) {
        // 错误响应：提取 msg 字段作为错误消息
        const errorMsg = res.msg || res.data || 'Error'
        console.error('❌ API错误 (code=' + res.code + '):', errorMsg)
        return Promise.reject(new Error(errorMsg))
      } else {
        // 成功响应：返回完整的响应对象（包含 code, msg, data）
        return res
      }
    } else {
      // 没有 code 字段，直接返回数据（兼容旧的 API）
      return res
    }
  },
  error => {
    console.error('❌ 请求错误:', error)
    
    if (!error.response) {
      // 网络错误（没有响应）
      console.error('❌ 网络错误，无响应')
      return Promise.reject(new Error('网络连接失败，请检查网络后重试'))
    }
    
    // 服务器返回错误响应（HTTP 状态码 4xx 或 5xx）
    const res = error.response.data
    console.error('❌ 服务器错误响应:', res)
    
    let errorMsg = '请求失败'
    
    // 优先使用 msg 字段
    if (res && res.msg) {
      errorMsg = res.msg
      console.error('❌ 提取错误消息 (msg):', errorMsg)
    } else if (res && res.data && typeof res.data === 'string') {
      errorMsg = res.data
      console.error('❌ 提取错误消息 (data):', errorMsg)
    } else if (error.message) {
      errorMsg = error.message
      console.error('❌ 提取错误消息 (message):', errorMsg)
    }
    
    return Promise.reject(new Error(errorMsg))
  }
)

export default service
