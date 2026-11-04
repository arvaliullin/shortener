import axios from 'axios'
import { keycloak } from '../auth'

export const api = axios.create({
  baseURL: '/api',
  timeout: 5000,
})

api.interceptors.request.use(async (config) => {
  await keycloak.updateToken(30)
  config.headers.Authorization = `Bearer ${keycloak.token}`
  return config
})
