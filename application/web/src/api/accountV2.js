import request from '@/utils/request'

// 获取登录界面角色选择列表（从区块链）
export function queryAccountList() {
    return request({
        url: '/queryAccountV2List',
        method: 'post'
    })
}

// 从数据库获取用户列表（用于用户管理页面）
export function queryAccountListFromDB(status = 1) {
    return request({
        url: '/queryAccountListFromDB',
        method: 'post',
        data: { status }
    })
}

// 登录
export function login(data) {
    return request({
        url: '/queryAccountV2List',
        method: 'post',
        data
    })
}

// 密码登录
export function loginWithPassword(data) {
    return request({
        url: '/loginWithPassword',
        method: 'post',
        data
    })
}

// 获取用户信息
export function getUserInfo(data) {
    return request({
        url: '/getUserInfo',
        method: 'post',
        data
    })
}

// 获取用户详情
export function getUserDetail(idOrUsername) {
    // 如果是数字，作为ID查询；否则作为用户名查询
    const isId = typeof idOrUsername === 'number' || !isNaN(idOrUsername)
    return request({
        url: '/getUserDetail',
        method: 'post',
        data: isId ? { id: parseInt(idOrUsername) } : { username: idOrUsername }
    })
}

// 更新用户
export function updateUser(data) {
    return request({
        url: '/updateUser',
        method: 'post',
        data
    })
}

// 删除用户（软删除）
export function deleteUser(id) {
    return request({
        url: '/deleteUser',
        method: 'post',
        data: { id }
    })
}

// 恢复用户
export function restoreUser(id) {
    return request({
        url: '/restoreUser',
        method: 'post',
        data: { id }
    })
}

// 批量删除用户
export function batchDeleteUsers(ids) {
    return request({
        url: '/batchDeleteUsers',
        method: 'post',
        data: { ids }
    })
}

// 注册
export function register(data) {
    return request({
        url: '/register',
        method: 'post',
        data
    })
}

// 创建角色
export function createAccount(data) {
    return request({
        url: '/createAccountV2',
        method: 'post',
        data
    })
}

// 获取统计数据
export function getStatistics() {
    return request({
        url: '/statistics',
        method: 'get'
    })
}

// 获取最近活动
export function getRecentActivities() {
    return request({
        url: '/recentActivities',
        method: 'get'
    })
}

// 从区块链同步账户ID
export function syncAccountFromBlockchain(username) {
    return request({
        url: '/syncAccountFromBlockchain',
        method: 'post',
        data: { username }
    })
}
