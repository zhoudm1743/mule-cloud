/**
 * 扫码和工序上报相关API
 */
import { post, get, put, del } from '@/utils/request'

/**
 * 扫码解析
 * @param {Object} params { qr_code: "二维码内容" }
 */
export function parseScanCode(params) {
	return post('/production/scan/parse', params)
}

/**
 * 提交工序上报
 * @param {Object} params 上报参数
 */
export function submitProcedureReport(params) {
	return post('/production/reports', params)
}

/**
 * 获取上报记录列表
 * @param {Object} params 查询参数
 */
export function getReportList(params) {
	return get('/production/reports', params)
}

/**
 * 获取上报记录详情
 * @param {String} id 上报记录ID
 */
export function getReportDetail(id) {
	return get(`/production/reports/${id}`)
}

/**
 * 删除上报记录
 * @param {String} id 上报记录ID
 */
export function deleteReport(id) {
	return del(`/production/reports/${id}`)
}

/**
 * 获取订单工序进度
 * @param {String} orderId 订单ID
 */
export function getOrderProgress(orderId) {
	return get(`/production/progress/${orderId}`)
}

/**
 * 获取工资统计
 * @param {Object} params { start_date, end_date, worker_id }
 */
export function getSalary(params) {
	return get('/production/salary', params)
}

