<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

function goBack() {
    router.go(-1)
}
</script>

<template>
    <div class="d-flex">
        <svg class="nav-back" viewBox="0 0 1000 1000" @click="goBack">
            <path fill="#969799" fill-rule="evenodd"
                d="M296.114 508.035c-3.22-13.597.473-28.499 11.079-39.105l333.912-333.912c16.271-16.272 42.653-16.272 58.925 0s16.272 42.654 0 58.926L395.504 498.47l304.574 304.574c16.272 16.272 16.272 42.654 0 58.926s-42.654 16.272-58.926 0L307.241 528.058a41.472 41.472 0 0 1-11.127-20.023z">
            </path>
        </svg>
        <h3>Add Task</h3>
    </div>
</template>

<style lang="scss" scoped>
.d-flex {
    display: flex;
    align-items: center;
}

.nav-back {
    margin-right: 1rem;
    width: 24px;
    height: 24px;
    cursor: pointer;
}
</style>
