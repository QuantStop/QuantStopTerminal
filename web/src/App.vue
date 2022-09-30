<template>

  <!--  Header / Navbar if logged in  -->
  <div v-if="userStore.isAuthed">
    <Header @custom-change2="toggleTheme" />
  </div>

  <!--  Header / Navbar if not logged in (theme button only) -->
  <div v-if="!userStore.isAuthed" class="d-flex">
    <div class="ms-auto p-3" style="height: 52px;">
      <ThemeButton @custom-change="toggleTheme" />
    </div>
  </div>

  <!--  Main container and row, full height minus header height  -->
  <div class="container-fluid qst-body m-0 p-0">
    <div class="row h-100 m-0 p-0">

      <!--  Sidebar Navigation / Logged in only  -->
      <div v-if="userStore.isAuthed" class="sidebar-container">
        <SidebarNav></SidebarNav>
      </div>

      <!--  App View  -->
      <div id="routerView" :class="userStore.isAuthed ? 'col-md-9 ms-sm-auto col-lg-10 p-0 m-0 h-100' : ''">
        <router-view />
      </div>

    </div>
  </div>

  <!--  Footer  -->
<!--  <div v-if="getUserProfile.id !== 0">
    <Footer />
  </div>-->


</template>

<script>
import Header from "./components/Header.vue"
import SidebarNav from "./components/SidebarNav";
import Footer from "./components/Footer";
import {userStore} from "./store/userStore";
import ThemeButton from "./components/ThemeButton";
export default {
  name: 'App',
  components: {
    ThemeButton,
    Header,
    SidebarNav,
    Footer
  },
  data() {
    return {
      userStore
    }
  },
  mounted() {
    /*document.body.classList.add('d-flex', 'flex-column', 'h-100')*/
    document.body.classList.add('h-100')
    document.documentElement.classList.add('h-100')
  },
  methods: {
    toggleTheme() {
      //console.log("toggle theme from app")
    }
  }
}
</script>

<style>

/* Define styles for the default root window element */
:root {
  --background-color-primary: #ffffff;
  --background-color-secondary: #fafafa;
  --accent-color: #cacaca;
  --text-primary-color: #222;
  --element-size: 4rem;
  --theme-switch-background-color: #f1f1f1; /* Theme switch background color */
  --border-color: rgb(60 60 60 / 29%);
  --border-color-hover: rgba(59, 59, 59, 0.29);
  --table-border-hover: rgba(65, 64, 64, 0.29);

  --slider-color-left: #b7b7b7;
  --slider-color-right: #fafafa;

  --orderbook-background-color: #ffffff;
}

/* Define styles for the root window with dark - mode preference */
:root.dark {
  --background-color-primary: #1a1a1a;
  --background-color-secondary: #2d2d30;
  --accent-color: #3f3f3f;
  --text-primary-color: rgb(255 255 255 / 87%);
  --theme-switch-background-color: #2f2f2f; /* Theme switch background color */
  --border-color: rgb(84 84 84 / 65%);
  --border-color-hover: rgba(196, 193, 193, 0.29);
  --table-border-hover: rgba(65, 64, 64, 0.29);

  --slider-color-left: #8c8c8c;
  --slider-color-right: #2d2d30;

  --orderbook-background-color: #141414;
}

body {
  font-family: Segoe UI,serif !important;
  font-size: .875rem;
  background-color: var(--background-color-primary) !important;
}
.qst-body {
  height: calc(100% - 52px);
}

footer {
  background-color: var(--background-color-primary);
}

a, p, i, .nav-link, label, h1, h2, h3, h4, button {
  color: var(--text-primary-color) !important;
}

ul {
  list-style-type: none !important;
  margin-block-start: 1em;
  margin-block-end: 1em;
  margin-inline-start: 0px;
  margin-inline-end: 0px;
  padding-inline-start: 40px;

}
li {
  color: var(--text-primary-color) !important;
  background-color: var(--background-color-primary) !important;
}
input, select {
  color: var(--text-primary-color) !important;
  background-color: var(--background-color-primary) !important;
}

input {
  border-color: var(--border-color) !important;
}
input:hover {
  border-color: var(--border-color-hover) !important;
}
.form-control {
  border: 1px solid var(--border-color) !important;
}

.feather {
  width: 16px;
  height: 16px;
  vertical-align: text-bottom;
}
.error-feedback {
  color: red;
}


/************** Sidebar **************/
.sidebar-container {
  /* remove bootstrap styles inherited from row */
  width: 0 !important;
  padding: 0 !important;
  margin: 0 !important;
}
.sidebar {
  position: fixed;
  top: 0;
  /* rtl:raw:
  right: 0;
  */
  bottom: 0;
  /* rtl:remove */
  left: 0;
  z-index: 100; /* Behind the navbar */
  padding: 52px 0 0; /* Height of navbar */

}

@media (max-width: 767.98px) {
  .sidebar {
    top: 5rem;
  }
}

.sidebar-sticky {
  position: relative;
  top: 0;
  height: calc(100vh - 52px);
  padding-top: .5rem;
  overflow-x: hidden;
  overflow-y: auto; /* Scrollable contents if viewport is shorter than content. */
}

.sidebar .nav-link {
  font-weight: 500;
  color: #333;
}

.sidebar .nav-link .feather {
  margin-right: 4px;
  color: #727272;
}

.sidebar .nav-link.active {
  color: #2470dc !important;
}

.sidebar .nav-link:hover .feather,
.sidebar .nav-link.active .feather {
  color: inherit;
}

.sidebar-heading {
  font-size: .75rem;
  text-transform: uppercase;
}

/************** Navbar **************/

.navbar-brand {
  padding-top: .75rem;
  padding-bottom: .75rem;
  font-size: 1rem;
  color: var(--text-primary-color) !important;
  background-color: var(--background-color-primary) !important;
}

.navbar .navbar-toggler {
  top: .25rem;
  right: 1rem;
}

.navbar .form-control {
  padding: .75rem 1rem;
  border-width: 0;
  border-radius: 0;
}


.dropdown-item:hover {
  background-color: var(--background-color-secondary) !important;
}


.card {
  color: var(--text-primary-color) !important;
  background-color: var(--background-color-primary) !important;
  border-color: var(--border-color) !important;
}
.card-header {
  color: var(--text-primary-color) !important;
  background-color: var(--background-color-secondary) !important;
  border-color: var(--border-color) !important;
}
.list-group-item {
  color: var(--text-primary-color) !important;
  background-color: var(--background-color-primary) !important;
  border-color: var(--border-color) !important;
}

.modal-header, .modal-body, .modal-footer, .dropdown-menu {
  color: var(--text-primary-color) !important;
  background-color: var(--background-color-primary) !important;
  border-color: var(--border-color) !important;
}


</style>
