import {reactive} from 'vue'
import axios from "axios";
import jwtInterceptor from "../shared/jwt.interceptor";

export const subsystemsStore = reactive({
    subsystems: {
        database: false,
        ntp_timekeeper: false,
        active_trader: false,
        internet_monitor: false,
    },

    getSusbystems() {
        return this.subsystems;
    },

    setSusbystems(data) {
        this.subsystems = data
    },

    async actionSubsystems() {
        const response = await jwtInterceptor.get("/api/sub-status", {
            withCredentials: true,
            credentials: "include",
            headers: {
                'X-Requested-With': 'XMLHttpRequest' // CSRF prevention
            },
        });

        if (response && response.data) {
            this.setSusbystems(response.data)
        }
    },
})