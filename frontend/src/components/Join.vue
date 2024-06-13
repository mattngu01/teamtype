<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
    data() {
        return {
            socket: null as WebSocket | null,
            lobbyId: "",
        }
    },
    mounted() {
        this.socket = new WebSocket("ws://localhost:8080/ws");
        this.socket.onopen = function(event: any) {
            console.log("Connected to server", event);
            this.send(JSON.stringify({"type": "LobbyInfo"}));
        }
        this.socket.onmessage = this.parseMessage;
    },
    methods: {
        parseMessage(event: any) {
            console.log("Received event", event);
            let eventPayload = JSON.parse(event.data);
            if (eventPayload["type"] == "LobbyInfo") {
                this.lobbyId = eventPayload.data["lobbyId"];
            }
        }
    },
})
</script>

<template>
<p v-if="lobbyId">Lobby ID: {{ lobbyId }}</p>
</template>