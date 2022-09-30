import {reactive} from 'vue'
import axios from "axios";
import jwtInterceptor from "../shared/jwt.interceptor";
const defaultUserProfile = {
    id: 0,
    username: "",
    roles: "",
}
export const userStore = reactive({
    loginStatus: "",
    registerExchangeStatus: "",
    isAuthed: false,
    userProfile: {
        id: 0,
        username: "",
        roles: "",
    },
    logOut: false,

    /* Getters */
    getLoginStatus() {
        return this.loginStatus;
    },
    getRegisterExchangeStatus() {
        return this.registerExchangeStatus;
    },
    getIsAuthed() {
        if (localStorage.getItem('isAuthenticated')) {
            try {
                this.isAuthed = JSON.parse(localStorage.getItem('isAuthenticated'));
            } catch(e) {
                localStorage.removeItem('isAuthenticated');
            }
        } else {
            // first time no local storage saved yet, save the defaults
            this.setIsAuthed(this.isAuthed)
        }
        return this.isAuthed;
    },
    getUserProfile() {
        if (localStorage.getItem('userProfile')) {
            try {
                this.userProfile = JSON.parse(localStorage.getItem('userProfile'));
            } catch(e) {
                localStorage.removeItem('userProfile');
            }
        } else {
            // first time no local storage saved yet, save the defaults
            this.setUserProfile(this.userProfile)
        }
        return this.userProfile;
    },
    getLogout() {
        return this.logOut;
    },

    /* Setters */
    setLoginStatus(data) {
        this.loginStatus = data;
    },
    setRegisterExchangeStatus(data) {
        this.registerExchangeStatus = data;
    },
    setIsAuthed(data) {
        this.isAuthed = data
        localStorage.setItem("isAuthenticated", data);
    },
    setUserProfile(data) {
        this.userProfile = {
            id: data.id,
            username: data.username,
            roles: data.roles
        };
        localStorage.setItem("userProfile", JSON.stringify(this.userProfile));
    },
    setLogout(data) {
        this.logOut = data;
    },

    /* Actions */
    async actionLogin(payload) {
        const response = await axios.post("/api/session", payload, {
            withCredentials: true,
            credentials: "include",
            headers: {
                'X-Requested-With': 'XMLHttpRequest' // CSRF prevention
            },
        });
        if (response.status === 200 && response.data) {
            localStorage.setItem("isAuthenticated", "true");
            this.setUserProfile(response.data)
            this.setLoginStatus("success")
        } else {
            this.setLoginStatus("failed")
        }
    },
    async actionRegisterExchange({ commit }, payload) {
        const response = await axios.post("/api/exchange", payload, {
            withCredentials: true,
            credentials: "include",
            headers: {
                'X-Requested-With': 'XMLHttpRequest' // CSRF prevention
            },
        });
        if (response && response.data) {
            this.setRegisterExchangeStatus("success")
        } else {
            this.setRegisterExchangeStatus("failed")
        }
    },
    async actionGetUserProfile() {
        const response = await jwtInterceptor.get("/api/user", {
            withCredentials: true,
            credentials: "include",
            headers: {
                'X-Requested-With': 'XMLHttpRequest' // CSRF prevention
            },
        });

        if (response && response.data) {
            this.setUserProfile(response.data)
        } else {
            this.setUserProfile(defaultUserProfile)
        }
    },
    async actionLogout() {
        await axios.delete("/api/session", {
            withCredentials: true,
            credentials: "include",
            headers: {
                'X-Requested-With': 'XMLHttpRequest' // CSRF prevention
            },
        });

        this.setLogout(true)
        this.setIsAuthed(false)
        this.setUserProfile(defaultUserProfile)
    },



})