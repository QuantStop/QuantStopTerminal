<template>
  <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
    <h2>Settings</h2>
<!--    <div class="btn-toolbar mb-2 mb-md-0">
      <div class="btn-group me-2">
        <button type="button" class="btn btn-sm btn-outline-secondary">Share</button>
        <button type="button" class="btn btn-sm btn-outline-secondary">Export</button>
      </div>
      <button type="button" class="btn btn-sm btn-outline-secondary dropdown-toggle">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-calendar" aria-hidden="true"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
        This week
      </button>
    </div>-->
  </div>

  <ul class="nav nav-tabs" id="myTab" role="tablist">
    <li class="nav-item" role="presentation">
      <button class="nav-link active" id="home-tab" data-bs-toggle="tab" data-bs-target="#home" type="button" role="tab" aria-controls="home" aria-selected="true">NTP</button>
    </li>
    <li class="nav-item" role="presentation">
      <button class="nav-link" id="profile-tab" data-bs-toggle="tab" data-bs-target="#profile" type="button" role="tab" aria-controls="profile" aria-selected="false">Database</button>
    </li>
    <li class="nav-item" role="presentation">
      <button class="nav-link" id="contact-tab" data-bs-toggle="tab" data-bs-target="#contact" type="button" role="tab" aria-controls="contact" aria-selected="false">System</button>
    </li>
  </ul>
  <div class="tab-content" id="myTabContent">
    <div class="tab-pane fade show active" id="home" role="tabpanel" aria-labelledby="home-tab"><NTPSettings :value="this.subsystemsStore.subsystems.ntp_timekeeper"></NTPSettings></div>
    <div class="tab-pane fade" id="profile" role="tabpanel" aria-labelledby="profile-tab">Database</div>
    <div class="tab-pane fade" id="contact" role="tabpanel" aria-labelledby="contact-tab"><SystemSettings></SystemSettings></div>
  </div>

</template>

<script>
import NTPSettings from "./NTPSettings";
import SystemSettings from "./SystemSettings";
import jwtInterceptor from "../../shared/jwt.interceptor";
import {subsystemsStore} from "../../store/subsystemsStore";
export default {
  name: "Settings",
  data: () => ({
    subsystemsStore,
  }),
  components: {SystemSettings, NTPSettings},

}
</script>

<style scoped>
.nav-link.active {
  border-right: 1px solid var(--border-color) !important;
  border-left: 1px solid var(--border-color) !important;
  border-top: 1px solid var(--border-color) !important;
  background-color: var(--background-color-primary) !important;
}
</style>