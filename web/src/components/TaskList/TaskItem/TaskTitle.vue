<script setup lang="ts">
import { Tag, Progress, TextEllipsis, Icon, showConfirmDialog } from 'vant'
import { type DownloadTask, downloadTaskStore } from '@/stores/downloadTaskStore'
import { computed } from 'vue'
import axios from 'axios'

const store = downloadTaskStore()

const props = defineProps<{
    task: DownloadTask,
}>()

const statusType = computed(() => {
    switch (props.task.status) {
        case "Pending":
            return "warning"
        case "Downloading":
            return "primary"
        case "Finished":
            return "success"
        default:
            return "default"
    }
})

const isDownloading = computed(() => props.task.status === "Downloading")

function deleteTask() {
    showConfirmDialog({
        title: 'Delete Task?',
        message: props.task.title,
    })
    .then(async () => {
        await axios.delete(`/api/tasks/${props.task.id}`);
        await store.reload()
    })
}

</script>

<template>
    <div class="content">
        <div class="d-flex">
            <div class="flex-grow-1 mr-1">
                <TextEllipsis :content="props.task.title" />
            </div>
            <Tag class="mr-1" :type="statusType">{{ props.task.status }}</Tag>
            <Icon name="delete" color="red" size="1.5rem" @click.stop="deleteTask" />
        </div>
        <div v-if="isDownloading" class="progress-container">
            <Progress :percentage="props.task.progress"></Progress>
        </div>
    </div>
</template>

<style lang="scss" scoped>
.content {
    padding: 0 1rem;
}

.d-flex {
    display: flex;
    align-items: center;
}

.flex-grow-1 {
    flex-grow: 1;
}

.mr-1 {
    margin-right: 1em;
}

.progress-container {
    padding: 1rem 0;
}
</style>
