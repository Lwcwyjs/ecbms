import request from '../utils/request'

export function getSystemInfo() {
  return request({
    url: '/system/info',
    method: 'get'
  })
}

export function getSystemStats() {
  return request({
    url: '/system/stats',
    method: 'get'
  })
}

export function getSystemLogs() {
  return request({
    url: '/system/logs',
    method: 'get'
  })
}

export function rebootSystem() {
  return request({
    url: '/system/reboot',
    method: 'post'
  })
}

export function shutdownSystem() {
  return request({
    url: '/system/shutdown',
    method: 'post'
  })
}
