import { ref, type Ref } from 'vue'
import { defineStore } from 'pinia'
import axios from 'axios'

export interface DownloadTask {
  id: string
  url: string
  title: string
  progress: number
  status: string
  createTime: Date
  startTime?: Date
  finishTime?: Date
}

const RETRY_INTERVAL = 5000 // 5-second retry interval

export const downloadTaskStore = defineStore('downloadTasks', () => {
  const downloadTasks: Ref<DownloadTask[]> = ref([])
  let eventSource: EventSource | null = null

  const eventHandlers = {
    onMessage: (event: MessageEvent) => {
      const data = JSON.parse(event.data)
      downloadTasks.value = data
    },
    
    onError: () => {
      closeEventSource()
      scheduleReconnect()
    }
  }

  function scheduleReconnect() {
    setTimeout(initEventSource, RETRY_INTERVAL)
  }

  function initEventSource() {
    closeEventSource()
    setupEventSource()
  }

  function setupEventSource() {
    eventSource = new EventSource('/api/tasks/events')
    eventSource.onmessage = eventHandlers.onMessage
    eventSource.onerror = eventHandlers.onError
  }

  function closeEventSource() {
    if (eventSource) {
      eventSource.close()
      eventSource = null
    }
  }

  async function reload() {
    try {
      const response = await axios.get<DownloadTask[]>('/api/tasks')
      downloadTasks.value = response.data
    } catch (error) {
      console.error('Failed to reload tasks:', error)
    }
  }

  return { 
    downloadTasks, 
    reload, 
    initEventSource, 
    closeEventSource 
  }
})
