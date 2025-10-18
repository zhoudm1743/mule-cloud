/**
 * 质检和返工相关API
 */
import { post, get, put, del } from '@/utils/request'

// ==================== 质检 API ====================

/**
 * 提交质检结果
 * @param {Object} params 质检参数
 */
export function submitInspection(params) {
	return post('/production/inspections', params)
}

/**
 * 获取质检记录列表
 * @param {Object} params 查询参数
 */
export function getInspectionList(params) {
	return get('/production/inspections', params)
}

/**
 * 获取质检记录详情
 * @param {String} id 质检记录ID
 */
export function getInspectionDetail(id) {
	return get(`/production/inspections/${id}`)
}

/**
 * 删除质检记录
 * @param {String} id 质检记录ID
 */
export function deleteInspection(id) {
	return del(`/production/inspections/${id}`)
}

// ==================== 返工 API ====================

/**
 * 创建返工单
 * @param {Object} params 返工参数
 */
export function createRework(params) {
	return post('/production/reworks', params)
}

/**
 * 获取返工列表
 * @param {Object} params 查询参数
 */
export function getReworkList(params) {
	return get('/production/reworks', params)
}

/**
 * 获取返工详情
 * @param {String} id 返工单ID
 */
export function getReworkDetail(id) {
	return get(`/production/reworks/${id}`)
}

/**
 * 完成返工
 * @param {String} id 返工单ID
 * @param {Object} params 完成信息
 */
export function completeRework(id, params) {
	return put(`/production/reworks/${id}/complete`, params)
}

/**
 * 删除返工记录
 * @param {String} id 返工单ID
 */
export function deleteRework(id) {
	return del(`/production/reworks/${id}`)
}

