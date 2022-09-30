<template>
  <div>
    <input
      @change="toggleTheme"
      @click="isActive = !isActive"
      id="checkbox"
      type="checkbox"
      class="switch-checkbox"
    />
    <label for="checkbox" class="switch-label">
      <span class="switch-toggle" :class="{ 'switch-toggle-checked': userTheme === 'dark' }">
        <svg class="sun-and-moon" aria-hidden="true" width="22.4" height="22.4" viewBox="0 0 24 24">
          <circle class="sun" cx="12" cy="12" r="6" mask="url(#moon-mask)" stroke="currentColor" :style="{'fill': userTheme === 'dark' ? 'currentColor': 'none'}" />
          <g class="sun-beams" stroke="currentColor">
            <line x1="12" y1="2" x2="12" y2="4" />
            <line x1="12" y1="20" x2="12" y2="22" />
            <line x1="5.22" y1="5.22" x2="6.64" y2="6.64" />
            <line x1="17.36" y1="17.36" x2="18.78" y2="18.78" />
            <line x1="2" y1="12" x2="4" y2="12" />
            <line x1="20" y1="12" x2="22" y2="12" />
            <line x1="5.22" y1="18.78" x2="6.64" y2="17.36" />
            <line x1="17.36" y1="6.64" x2="18.78" y2="5.22" />
          </g>
          <mask class="moon" id="moon-mask">
            <rect x="0" y="0" width="100%" height="100%" fill="white" />
            <circle cx="22.4" cy="7" r="5" fill="black" />
          </mask>
        </svg>
      </span>
    </label>
  </div>
</template>

<script>
import {themeStore} from "../store/themeStore";
export default {
  name: 'ThemeButton',
  mounted() {
    const initUserTheme = this.getTheme() || this.getMediaPreference();
    this.setTheme(initUserTheme);
    themeStore.setTheme(initUserTheme)
  },

  data() {
    return {
      themeStore,
      userTheme: "light",
      isActive: false,
    };
  },
  emits: ['customChange'],
  methods: {
    toggleTheme() {
      const activeTheme = localStorage.getItem("user-theme");
      if (activeTheme === "light") {
        this.setTheme("dark");
        themeStore.theme = "dark"
      } else {
        this.setTheme("light");
        themeStore.theme = "light"
      }
      this.$emit("customChange", event.target.value)

    },

    getTheme() {
      return localStorage.getItem("user-theme");
    },

    setTheme(theme) {
      localStorage.setItem("user-theme", theme);
      this.userTheme = theme;
      document.documentElement.className = theme + ' h-100';
    },

    getMediaPreference() {
      const hasDarkPreference = window.matchMedia(
          "(prefers-color-scheme: dark)"
      ).matches;
      if (hasDarkPreference) {
        return "dark";
      } else {
        return "light";
      }
    },
  },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

/* helpful resources:
  https://web.dev/building-a-theme-switch-component/
*/
@import "https://unpkg.com/open-props/easings.min.css";

.switch-checkbox {
  display: none;
}
.switch-label {
  align-items: center;
  background: var(--theme-switch-background-color);
  border: calc(var(--element-size) * 0.025) solid var(--border-color);
  border-radius: var(--element-size);
  cursor: pointer;
  display: flex;
  font-size: calc(var(--element-size) * 0.3);
  height: calc(var(--element-size) * 0.4);
  position: relative;
  transition: background 0.5s ease;
  justify-content: space-between;
  width: calc(var(--element-size) * 0.8);
  z-index: 1;
}

.switch-label:hover {
  border: calc(var(--element-size) * 0.025) solid var(--border-color-hover);
}
.switch-toggle {
  position: absolute;
  background-color: var(--background-color-primary);
  border-radius: 50%;
  height: calc(var(--element-size) * 0.35);
  width: calc(var(--element-size) * 0.35);
  transform: translateX(0);
  transition: transform 0.25s ease, background-color 0.5s ease;
  z-index: 2;
}
.switch-toggle-checked {
  transform: translateX(calc(var(--element-size) * 0.4)) !important;

}

/* styles for sun/moon icon */
.switch-toggle-checked .sun {
  transition-timing-function: var(--ease-3);
  transition-duration: .25s;
}
.switch-toggle-checked .sun-beams {
  opacity: 0;
  transform: rotateZ(-25deg);
  transition-duration: .15s;
}
.switch-toggle-checked .moon circle {
  transform: translateX(-7px);
  transition: transform .25s var(--ease-out-5);
  transition-delay: .25s;
  transition-duration: .5s;
}
.sun-and-moon {
  display:flex;
  align-items:center;
  justify-content:center;
}
.sun {
  transition: transform .5s var(--ease-elastic-3);
}
.sun-beams {
  transition:
    transform .5s var(--ease-elastic-4),
    opacity .5s var(--ease-3)
  ;
}

</style>
