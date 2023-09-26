import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

import '@vant/touch-emulator'
import { Locale } from 'vant'
import enUS from 'vant/es/locale/lang/en-US'

import { setAxiosInterceptors } from './axios/axios'

Locale.use('en-US', enUS)
setAxiosInterceptors()

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
