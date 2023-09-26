<script setup lang="ts">
import { Button, Form, Field, CellGroup } from 'vant'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()

const url = ref('')
const title = ref('')

async function onSubmit() {
    await axios.post(`/api/tasks/`, {
        url: url.value,
        title: title.value,
    })
    router.push("/")
}
</script>

<template>
    <Form @submit="onSubmit">
        <CellGroup inset>
            <Field v-model="url" name="Url" label="Url" placeholder="Url"
                :rules="[{ required: true, message: 'Url is required' }]" />
            <Field v-model="title" name="Title" label="Title" placeholder="Title" />
        </CellGroup>
        <div style="margin: 16px;">
            <Button round block type="primary" native-type="submit">
                Add
            </Button>
        </div>
    </Form>
</template>

<style lang="scss" scoped></style>
