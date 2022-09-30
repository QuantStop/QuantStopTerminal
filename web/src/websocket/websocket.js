import { reactive } from 'vue'
import { getCurrentInstance } from 'vue'


export const websocket = reactive({
    socket: WebSocket,
    connected: false,
    createWebsocket() {
        const internalInstance = getCurrentInstance();
        const emitter = internalInstance.appContext.config.globalProperties.$bus;
        this.socket = new WebSocket("ws://localhost:8080/api/ws")
        this.socket.onmessage = function(msg){
            emitter.trigger("onWebsocketMessage", msg.data)
        }
        this.socket.onerror = function(err){
            emitter.trigger("onWebsocketError", err)
        }
        this.connected = true
    },
    sendMessage(message) {
        console.log("sending message: " + message)
        this.socket.send(message)
    },
    closeWebsocket() {
        this.connected = false
        this.socket.close()
    },
})