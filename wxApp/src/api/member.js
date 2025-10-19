/**
 * 员工档案相关API
 */
import request from '@/utils/request'

/**
 * 获取个人档案
 * @param {Boolean} silent - 是否静默加载（不显示loading）
 */
export function getProfile(silent = false) {
	return request({
		url: '/miniapp/member/profile',
		method: 'GET',
		loading: !silent  // 如果 silent 为 true，则不显示 loading
	})
}

/**
 * 更新基本信息
 * @param {Object} data - 基本信息数据
 */
export function updateBasicInfo(data) {
	return request({
		url: '/miniapp/member/profile/basic',
		method: 'PUT',
		data
	})
}

/**
 * 更新联系信息
 * @param {Object} data - 联系信息数据
 */
export function updateContactInfo(data) {
	return request({
		url: '/miniapp/member/profile/contact',
		method: 'PUT',
		data
	})
}

/**
 * 上传照片
 * @param {Object} data - 照片数据 {type: 'avatar'|'photo', url: '图片URL'}
 */
export function uploadPhoto(data) {
	return request({
		url: '/miniapp/member/profile/photo',
		method: 'POST',
		data
	})
}

