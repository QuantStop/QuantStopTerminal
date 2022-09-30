<template>
  <button @click="disconnect" v-if="status === 'connected'">Disconnect</button>
  <button @click="connect" v-if="status === 'disconnected'">Connect</button> {{ status }}
  <div v-if="status === 'connected'">
    <form @submit.prevent="sendMessage" action="#">
      <input v-model="message"><button type="submit">Send Message</button>
    </form>
    <ul id="logs" ref="logs">
      <li v-for="log in logs" class="log">
        {{ log.event }}: {{ log.data }}
      </li>
    </ul>
  </div>
</template>

<script>


export default {
  name: "TestWebsocket",

  data: function () {
    return {
      message: "",
      logs: [],
      status: "disconnected"
    }
  },
  methods: {
    connect() {
      this.socket = new WebSocket("ws://localhost:8080/api/ws");
      this.socket.onopen = () => {
        this.status = "connected";
        this.logs.push({ event: "Connected to", data: 'ws://localhost:8080/api/ws'})

        this.socket.onmessage = ({data}) => {
          this.logs.push({ event: "Received message", data });

        };

        this.socket.onerror = ({error}) => {
          this.logs.push({ event: "Error", error });

        };

      };
    },
    disconnect() {
      this.socket.close();
      this.status = "disconnected";
      this.logs = [];
    },
    sendMessage(e) {
      this.socket.send(this.message);
      this.logs.push({ event: "Sent message", data: this.message });
      this.message = "";
    },
    scrollLogs() {
      const container = this.$refs.logs;
      // https://developer.mozilla.org/en-US/docs/Web/API/Element/scrollIntoView
      container.lastElementChild.scrollIntoView({behavior: "smooth", block: "end", inline: "nearest"});
    }
  },
  updated() {
    this.scrollLogs()
  },
  unmounted() {
    this.disconnect()
  }
}
</script>

<style scoped>


#logs {
  margin-top: 5px;
  padding: 5px;
  height: 500px;
  overflow: hidden !important;
}
#logs:hover, #logs:active, #logs:focus {
  overflow: auto !important;
}

/* width */
::-webkit-scrollbar {
  width: 10px;
}

/* Track */
::-webkit-scrollbar-track {
  background: var(--background-color-secondary);
  box-shadow: inset 0 0 5px grey;
  border-radius: 10px;
}

/* Handle */
::-webkit-scrollbar-thumb {
  border-color: var(--border-color);
  border-radius: 10px;
  background: var(--background-color-primary);
}

/* Handle on hover */
::-webkit-scrollbar-thumb:hover {
  background: #555;
}

button, form, li, ul {
  color: var(--text-primary-color) !important;
  background-color: var(--background-color-primary) !important;
}
</style>