<script setup lang="ts">
import { Button, showToast, showConfirmDialog } from 'vant'
import { useRouter } from 'vue-router'
import { deleteCompletedTasks } from '@/axios/axios'
import { downloadTaskStore } from '@/stores/downloadTaskStore'

const router = useRouter()
const store = downloadTaskStore()

function goToAddTask() {
    router.push('/add-task')
}

async function handleDeleteCompletedTasks() {
    const result = await showConfirmDialog({
        title: 'Confirm Delete',
        message: 'Are you sure you want to delete all completed download tasks?',
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel'
    })

    if (result) {
        try {
            await deleteCompletedTasks()
            await store.reload()
            showToast({
                type: 'success',
                message: 'Completed tasks deleted successfully!',
                position: 'top'
            })
        } catch (error) {
            showToast({
                type: 'fail',
                message: 'Failed to delete completed tasks',
                position: 'top'
            })
        }
    }
}
</script>

<template>
    <div class="d-flex">
        <div class="flex-grow-1">
            <h3>N_m3u8DL-RE-Web</h3>
        </div>
        <div>
            <Button size="small" type="primary" icon="plus" @click="goToAddTask">Add Task</Button>
            <Button size="small" type="danger" icon="delete" @click="handleDeleteCompletedTasks">Delete Completed
                Tasks</Button>
        </div>
    </div>
</template>

<style lang="scss" scoped>
.d-flex {
    display: flex;
    align-items: center;
}

.flex-grow-1 {
    flex-grow: 1;
}
</style>
