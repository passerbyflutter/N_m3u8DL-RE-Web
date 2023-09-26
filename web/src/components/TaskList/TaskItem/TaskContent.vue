<script setup lang="ts">
import { Field, CellGroup } from 'vant'
import { type DownloadTask } from '@/stores/downloadTaskStore'
import dayjs from 'dayjs'

const props = defineProps<{
    task: DownloadTask,
}>()

function formatDate(date?: Date) {
    if(date != null) {
        return dayjs(date).format('YYYY-MM-DD HH:mm:ss');
    }
    return "";
}

</script>

<template>
    <CellGroup inset>
        <Field label="Url" :model-value="props.task.url" readonly />
        <Field label="Create Time" :model-value="formatDate(props.task.createTime)" readonly />
        <Field label="Start Time" :model-value="formatDate(props.task.startTime)" readonly />
        <Field label="Finish Time" :model-value="formatDate(props.task.finishTime)" readonly />
    </CellGroup>
</template>