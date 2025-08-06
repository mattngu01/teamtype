<script lang="ts">
import { defineComponent } from 'vue';

interface LobbyInfoMsg {
	type: "LobbyInfo";
	data: {
		lobbyId: string;
		username: string;
		players: string[];
	};
}

export default defineComponent({
	name: "LobbyComponent",
	data() {
		return {
			socket: null as WebSocket | null,
			lobbyId: "" as string,
			username: "" as string,
			players: [] as string[],
			connectionError: null as string | null,
		}
	},
	mounted() {
	},
	methods: {
		connectWebsocket() {
			this.socket = new WebSocket("ws://localhost:8080/ws");
			this.socket.onopen = function (event: any) {
				console.log("Connected to server", event);
				this.socket.send(JSON.stringify({ "type": "LobbyInfo" }));
			}
			this.socket.onmessage = this.parseMessage;

			this.socket.onclose = () => {
				console.log("Connection closed.")
			}

			this.socket.onerror = (error) => {
				console.error("WebSocket error:", error);
				this.connectionError = "Failed to connect to the server. Please try again.";
			};
		},
		parseMessage(event: any) {
			console.log("Received event", event);
			try {
				let eventPayload = JSON.parse(event.data) as LobbyInfoMsg;
				if (eventPayload["type"] == "LobbyInfo") {
					this.lobbyId = eventPayload.data.lobbyId;
					this.username = eventPayload.data.username;
					this.players = eventPayload.data.players;
				}
			} catch (error) {
				let connectionError = "Failed to parse msg or invalid payload" + JSON.stringify(err, Object.getOwnPropertyNames(error));
				this.connectionError = connectionError;
				console.error(connectionError);
			}
		}
	},
})
</script>

<template>
	<p v-if="connectionError" class="error">{{ connectionError }}</p>
	<p v-if="lobbyId">Lobby ID: {{ lobbyId }}</p>
	<p v-if="username">Username: {{ username }}</p>
	<p v-if="players">Players: {{ players }}</p>

</template>
