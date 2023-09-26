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

export const downloadTaskStore = defineStore('downloadTasks', () => {
  const downloadTasks: Ref<DownloadTask[]> = ref([])
  async function reload() {
    downloadTasks.value = (await axios.get<DownloadTask[]>('/api/tasks')).data
  }

  return { downloadTasks, reload }
})
