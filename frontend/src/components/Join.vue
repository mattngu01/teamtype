<script lang="ts">
import { defineComponent } from 'vue';
export default defineComponent({
    data() {
        return {
            socket: null as WebSocket | null,
        }
    },
    mounted() {
        this.socket = new WebSocket("ws://localhost:8080/ws");
        console.log("Connected to server")
    },
    methods: {
    },
})
</script>

<template>
</template>