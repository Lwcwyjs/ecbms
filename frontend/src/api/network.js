import request from '../utils/request'

export function getNetworkInterfaces() {
  return request({
    url: '/network/interfaces',
    method: 'get'
  })
}

export function getNetworkConfigs() {
  return request({
    url: '/network/configs',
    method: 'get'
  })
}

export function configureNetwork(data) {
  return request({
    url: '/network/configure',
    method: 'post',
    data
  })
}

export function applyNetworkConfig(data) {
  return request({
    url: '/network/apply',
    method: 'post',
    data
  })
}

export function pingTest(data) {
  return request({
    url: '/network/ping',
    method: 'post',
    data
  })
}

export function getDNSConfig() {
  return request({
    url: '/network/dns',
    method: 'get'
  })
}
