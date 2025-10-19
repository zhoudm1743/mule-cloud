/** 不同请求服务的环境配置 */
export const serviceConfig: Record<ServiceEnvType, Record<string, string>> = {
  dev: {
    url: 'http://localhost:8080',
  },
  production: {
    url: 'https://mock.apifox.cn/m1/4071143-0-default',
  },
}
