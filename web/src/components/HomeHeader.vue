<script setup lang="ts">
import { Button } from 'vant'
import { useRouter } from 'vue-router'

const router = useRouter()

function goToAddTask() {
    router.push('/add-task')
}
</script>

<template>
    <div class="d-flex">
        <div class="flex-grow-1">
            <h3>N_m3u8DL-RE-Web</h3>
        </div>
        <div>
            <Button size="small" type="primary" icon="plus" @click="goToAddTask">Add Task</Button>
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
