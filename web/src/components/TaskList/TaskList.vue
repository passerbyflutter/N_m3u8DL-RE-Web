<script setup lang="ts">
import { Empty, PullRefresh, Collapse, CollapseItem } from 'vant'
import { downloadTaskStore } from '@/stores/downloadTaskStore'
import { computed, ref, type Ref } from 'vue'
import TaskTitle from './TaskItem/TaskTitle.vue'
import TaskContent from './TaskItem/TaskContent.vue'
import { onMounted } from 'vue'

const store = downloadTaskStore()

let loading = ref(false);
let activeNames: Ref<string[]> = ref([]);

onMounted(async () => {
    await store.reload()
})

const hasTasks = computed(() => store.downloadTasks.length != 0)

async function onRefresh() {
    try {
        await store.reload()
    } finally {
        activeNames.value = [];
        loading.value = false;
    }
}

</script>

<template>
    <PullRefresh v-model="loading" @refresh="onRefresh">
        <Empty v-if="!hasTasks" description="No Download Tasks" />
        <div class="collapse-container" v-if="hasTasks">
            <Collapse v-model="activeNames">
                <CollapseItem v-for="task in store.downloadTasks" :key="task.id" :name="task.id">
                    <template #title>
                        <TaskTitle :task="task"></TaskTitle>
                    </template>
                    <TaskContent :task="task"></TaskContent>
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
