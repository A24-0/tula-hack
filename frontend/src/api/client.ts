import axios from 'axios'

export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '',
  timeout: 10 * 60 * 1000,
})

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    const message = error.response?.data?.error || error.response?.data?.message || error.message
    return Promise.reject(new Error(message || 'Сервис временно недоступен'))
  },
)
