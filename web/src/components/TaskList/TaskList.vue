<script setup lang="ts">
import { Empty, PullRefresh, Collapse, CollapseItem } from 'vant'
import { downloadTaskStore } from '@/stores/downloadTaskStore'
import { computed, ref, type Ref } from 'vue'
import TaskTitle from './TaskItem/TaskTitle.vue'
import TaskContent from './TaskItem/TaskContent.vue'
import { onMounted, onUnmounted } from 'vue'

const store = downloadTaskStore()
const loading = ref(false)
const activeNames: Ref<string[]> = ref([])

// Computed Properties
const hasTasks = computed(() => store.downloadTasks.length > 0)

// Lifecycle Hooks
onMounted(setupTaskList)
onUnmounted(cleanup)

// Component Functions
async function setupTaskList() {
    initializeSSEConnection()
    await loadInitialTasks()
}

function initializeSSEConnection() {
    store.initEventSource()
}

async function loadInitialTasks() {
    await store.reload()
}

function cleanup() {
    store.closeEventSource()
}

async function onRefresh() {
    try {
        await store.reload()
    } finally {
        resetUI()
    }
}

function resetUI() {
    activeNames.value = []
    loading.value = false
}
</script>

<template>
    <PullRefresh v-model="loading" @refresh="onRefresh">
        <Empty v-if="!hasTasks" description="No Download Tasks" />
        <div v-else class="collapse-container">
            <Collapse v-model="activeNames">
                <CollapseItem 
                    v-for="task in store.downloadTasks" 
                    :key="task.id" 
                    :name="task.id"
                >
                    <template #title>
                        <TaskTitle :task="task" />
                    </template>
                    <TaskContent :task="task" />
                </CollapseItem>
            </Collapse>
        </div>
    </PullRefresh>
</template>

<style lang="scss" scoped>
.van-pull-refresh {
    height: calc(100vh - 64px);
}

.collapse-container {
    height: calc(100vh - 64px);
    overflow: auto;
}
</style>
