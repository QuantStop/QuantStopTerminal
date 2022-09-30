
import { reactive } from 'vue'

export const themeStore = reactive({
    theme: "light",
    value: false,
    toggle() {
        this.value = !this.value
    },
    getTheme() {
        return this.theme
    },
    setTheme(theme) {
        this.theme = theme
    }
})