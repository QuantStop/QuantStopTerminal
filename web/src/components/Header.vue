<template>

  <header class="qst-header-nav navbar navbar-expand-md sticky-top p-0">

    <div class="container-fluid p-0">

      <router-link to="/" class="navbar-brand col-md-3 col-lg-2 me-0 px-3">
        <span>Q</span>
        <span>u</span>
        <span>a</span>
        <span>n</span>
        <span>t</span>
        <span style="color:#d21515">s</span>
        <span style="color:#d21515">t</span>
        <span style="color:#d21515">o</span>
        <span style="color:#d21515">p</span>
        <span>T</span>
        <span>e</span>
        <span>r</span>
        <span>m</span>
        <span>i</span>
        <span>n</span>
        <span>a</span>
        <span>l</span>
      </router-link>

      <button ref="qstNavToggle" class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target=".qst-collapse" aria-controls="sidebarMenu headerMenu" aria-expanded="false" aria-label="Toggle navigation">
        <svg class="m-auto" id="menu-bars" viewBox="0 0 20 20" width="20" height="20" fill="currentColor">
          <rect width="20" height="4"></rect>
          <rect y="6.5" width="20" height="4"></rect>
          <rect y="13" width="20" height="4"></rect>
        </svg>
      </button>

      <div class="qst-collapse collapse navbar-collapse w-100" id="headerMenu">

        <input v-if="this.userStore.isAuthed" class="form-control w-100 me-3" type="text" placeholder="Search" aria-label="Search">


        <div class="nav-item m-auto me-3 theme-button d-inline-flex justify-content-end">
          <theme-button @custom-change="toggleTheme" />
        </div>

        <div class="nav-item dropdown m-auto me-3 user-dropdown d-inline-flex justify-content-end" v-if="this.userStore.isAuthed">
          <a href="#" class="d-flex align-items-center text-decoration-none dropdown-toggle" id="userDropdown" data-bs-toggle="dropdown" aria-expanded="false">
            <img src="//ssl.gstatic.com/accounts/ui/avatar_2x.png" alt="" width="32" height="32" class="rounded-circle me-2">
            <strong>{{ this.userStore.userProfile.username }}</strong>
          </a>
          <ul class="qst-user-dropdown dropdown-menu dropdown-menu-end text-small shadow" aria-labelledby="userDropdown">
            <li><router-link to="/settings" class="dropdown-item">Settings</router-link></li>
            <li><router-link to="/profile" class="dropdown-item">Profile</router-link></li>
            <li><hr class="dropdown-divider"></li>
            <li><a class="dropdown-item" @click.prevent="logOut" href="/logout"><font-awesome-icon icon="sign-out-alt" /> Sign out</a></li>
          </ul>
        </div>

      </div>



    </div>
  </header>

</template>

<script>
import {userStore} from "../store/userStore";
import ThemeButton from "./ThemeButton.vue"
export default {
  name: "Header",
  components: {
    ThemeButton,
  },
  data() {
    return {
      userStore
    }
  },
  emits: ['customChange2'],
  computed: {

    showAdminBoard() {
      if (this.userStore.userProfile && this.userStore.userProfile.roles) {
        return this.userStore.userProfile.roles.includes('admin');
      }
      return false;
    },
    showModeratorBoard() {
      if (this.userStore.userProfile && this.userStore.userProfile.roles) {
        return this.userStore.userProfile.roles.includes('moderator');
      }
      return false;
    }
  },
  methods: {
    async logOut() {
      await this.userStore.actionLogout();
      await this.$router.push("/login")
    },
    toggleTheme() {
      //console.log("toggle theme from header")
      this.$emit("customChange2", event.target.value)
    }
  },
}
</script>

<style scoped>
.qst-header-nav {
  color: var(--text-primary-color) !important;
  background-color: var(--background-color-primary) !important;
  border-bottom: 1px solid var(--border-color);
}
.navbar-brand {
  padding-top: 0.75rem;
  padding-bottom: 0.75rem;
  text-align: center;
  font-family: Impact, sans-serif;
  font-style: italic;
  font-size: large;
  color: var(--text-primary-color) !important;
}
.qst-user-dropdown {
  color: var(--text-primary-color) !important;
  background-color: var(--background-color-primary) !important;
}
.navbar-toggler {
  height: 30px;
  width: 30px;
  text-align: center;
  margin-right: 10px;
  opacity: 1;
  color: var(--text-primary-color) !important;
  background-color: var(--theme-switch-background-color) !important;
}
#menu-bars {
  text-align: center;
  position: relative;
  right: 8px;
  bottom: 2px;
}

@media screen and (max-width: 767px) {
  .qst-collapse {
  }
  .nav-item {
    width: 100%;
  }
  .user-dropdown {
    height: 48px;
  }
  .theme-button {
    height: 48px;
  }
}
</style>