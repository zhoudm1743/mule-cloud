/**
 * 网络请求封装
 */

// API基础地址（根据实际情况修改）
const BASE_URL = 'http://localhost:8080/api'; // 网关地址

/**
 * 发起请求
 * @param {Object} options 请求配置
 */
function request(options) {
	return new Promise((resolve, reject) => {
		// 获取token
		const token = uni.getStorageSync('token');
		
		// 完整URL
		const url = BASE_URL + options.url;
		
		// 请求头
		const header = {
			'Content-Type': 'application/json',
			...options.header
		};
		
		// 添加token
		if (token) {
			header['Authorization'] = `Bearer ${token}`;
		}
		
		// 显示加载提示
		if (options.loading !== false) {
			uni.showLoading({
				title: options.loadingText || '加载中...',
				mask: true
			});
		}
		
		uni.request({
			url,
			method: options.method || 'GET',
			data: options.data || {},
			header,
			timeout: options.timeout || 30000,
			success: (res) => {
				// 隐藏加载提示
				if (options.loading !== false) {
					uni.hideLoading();
				}
				
				// 统一处理响应
				if (res.statusCode === 200) {
					// 后端统一响应格式：{ code, data, message }
					// 兼容 code: 0 和 code: 200 两种成功状态
					if (res.data.code === 200 || res.data.code === 0) {
						resolve(res.data.data);
					} else {
						// 业务错误
						const errMsg = res.data.msg || res.data.message || '请求失败';
						uni.showToast({
							title: errMsg,
							icon: 'none',
							duration: 2000
						});
						reject(new Error(errMsg));
					}
				} else if (res.statusCode === 401) {
					// 未认证，跳转到登录页
					uni.showToast({
						title: '请先登录',
						icon: 'none'
					});
					uni.removeStorageSync('token');
					uni.removeStorageSync('userInfo');
					uni.reLaunch({
						url: '/pages/login/login'
					});
					reject(new Error('未认证'));
				} else {
					// HTTP错误
					uni.showToast({
						title: `请求失败(${res.statusCode})`,
						icon: 'none'
					});
					reject(new Error(`HTTP ${res.statusCode}`));
				}
			},
			fail: (err) => {
				// 隐藏加载提示
				if (options.loading !== false) {
					uni.hideLoading();
				}
				
				// 网络错误
				uni.showToast({
					title: '网络连接失败',
					icon: 'none'
				});
				reject(err);
			}
		});
	});
}

/**
 * GET请求
 */
export function get(url, data, options = {}) {
	return request({
		url,
		method: 'GET',
		data,
		...options
	});
}

/**
 * POST请求
 */
export function post(url, data, options = {}) {
	return request({
		url,
		method: 'POST',
		data,
		...options
	});
}

/**
 * PUT请求
 */
export function put(url, data, options = {}) {
	return request({
		url,
		method: 'PUT',
		data,
		...options
	});
}

/**
 * DELETE请求
 */
export function del(url, data, options = {}) {
	return request({
		url,
		method: 'DELETE',
		data,
		...options
	});
}

export default request;

