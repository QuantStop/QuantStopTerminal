import {reactive} from "vue";
import jwtInterceptor from "../shared/jwt.interceptor";
const defaultVersion = {
    copyright: "Copyright (c) 2021-2022 QuantStop.com",
    github: "GitHub: https://github.com/QuantStop/QuantStopTerminal",
    isdaemon: false,
    isdevelopment: true,
    isrelease: false,
    issues: "Issues: https://github.com/QuantStop/QuantStopTerminal/issues",
    prereleaseblurb: "This version is pre-release and is not intended to be used as a production ready trading framework or bot - use at your own risk.",
    version: "0.0.1"
}
export const versionStore = reactive({
    version: {
        copyright: "",
        github: "",
        isdaemon: false,
        isdevelopment: true,
        isrelease: false,
        issues: "",
        prereleaseblurb: "",
        version: ""
    },

    getVersion() {
        return this.version
    },
    setVersion(data) {
        this.version = data
    },

    async actionGetVersion() {
        const response = await jwtInterceptor.get("/api/version", {
            withCredentials: true,
            credentials: "include",
            headers: {
                'X-Requested-With': 'XMLHttpRequest' // CSRF prevention
            },
        });

        if (response && response.data) {
            this.setVersion(response.data)
        } else {
            this.setVersion(defaultVersion)
        }
    },


})