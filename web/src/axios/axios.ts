import axios from 'axios'
import dayjs from 'dayjs'
import qs from 'qs'

export function setAxiosInterceptors() {

    axios.interceptors.request.use((config) => {
        config.paramsSerializer = (params) =>
          qs.stringify(params, {
            serializeDate: (date: Date) => dayjs(date).format('YYYY-MM-DDTHH:mm:ssZ')
          })
        return config
      })
      
      axios.interceptors.response.use((originalResponse) => {
        handleDates(originalResponse.data)
        return originalResponse
      })
}

const isoDateFormat = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{7}Z$/

function isIsoDateString(value: any): boolean {
  return value && typeof value === 'string' && isoDateFormat.test(value)
}

function handleDates(body: any) {
  if (body === null || body === undefined || typeof body !== 'object') return body

  for (const key of Object.keys(body)) {
    const value = body[key]
    if (isIsoDateString(value)) body[key] = dayjs(value).toDate()
    else if (value === '0001-01-01T00:00:00Z') body[key] = null
    else if (typeof value === 'object') handleDates(value)
  }
}
