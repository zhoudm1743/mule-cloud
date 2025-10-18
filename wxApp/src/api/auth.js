/**
 * 认证相关API
 */
import { post, get, put, del } from '@/utils/request'

/**
 * 微信登录
 * @param {Object} params 登录参数（code + 用户信息）
 */
export function wechatLogin(params) {
	return post('/miniapp/wechat/login', params)
}

/**
 * 绑定租户
 * @param {String} userId 用户ID
 * @param {String} inviteCode 邀请码
 */
export function bindTenant(userId, inviteCode) {
	return post('/miniapp/wechat/bind-tenant', {
		user_id: userId,
		invite_code: inviteCode
	})
}

/**
 * 选择租户
 * @param {String} userId 用户ID
 * @param {String} tenantId 租户ID
 */
export function selectTenant(userId, tenantId) {
	return post('/miniapp/wechat/select-tenant', {
		user_id: userId,
		tenant_id: tenantId
	})
}

/**
 * 切换租户
 * @param {String} tenantId 租户ID
 */
export function switchTenant(tenantId) {
	return post('/miniapp/wechat/switch-tenant', {
		tenant_id: tenantId
	})
}

/**
 * 获取用户信息
 */
export function getUserInfo() {
	return get('/miniapp/user/info')
}

/**
 * 更新用户信息
 * @param {Object} data 用户信息（包括基本信息和企业信息）
 */
export function updateUserInfo(data) {
	return put('/miniapp/user/info', data)
}

/**
 * 更新租户成员信息
 * @param {Object} data 成员信息（工号、部门、岗位等）
 */
export function updateMemberInfo(data) {
	return put('/miniapp/member/info', data)
}

/**
 * 绑定手机号
 * @param {String} code 微信返回的code
 */
export function bindPhone(code) {
	return post('/miniapp/wechat/phone', { code })
}

/**
 * 解绑手机号
 */
export function unbindPhone() {
	return del('/miniapp/wechat/phone')
}

/**
 * 上传文件（头像等）
 * @param {String} filePath 本地文件路径
 * @param {String} businessType 业务类型（如：avatar, document等）
 */
export function uploadFile(filePath, businessType = 'avatar') {
	return new Promise((resolve, reject) => {
		const token = uni.getStorageSync('token')
		const currentTenant = uni.getStorageSync('currentTenant')
		
		uni.showLoading({
			title: '上传中...',
			mask: true
		})
		
		uni.uploadFile({
			url: 'https://dev.inzj.cn/api/admin/common/files/upload',
			filePath: filePath,
			name: 'file',
			formData: {
				business_type: businessType
			},
			header: {
				'Authorization': 'Bearer ' + token,
				'X-Tenant-Code': currentTenant?.tenant_code || ''
			},
			success: (res) => {
				uni.hideLoading()
				if (res.statusCode === 200) {
					const result = JSON.parse(res.data)
					// 后端统一响应格式：{ code, data, message }
					if (result.code === 200 || result.code === 0) {
						resolve(result.data)
					} else {
						reject(new Error(result.msg || result.message || '上传失败'))
					}
				} else {
					reject(new Error('上传失败'))
				}
			},
			fail: (error) => {
				uni.hideLoading()
				reject(error)
			}
		})
	})
}

