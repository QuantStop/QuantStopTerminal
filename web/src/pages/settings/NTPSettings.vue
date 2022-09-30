<template>

  <div class="list-group mb-5 shadow pt-3">
    <div class="list-group-item m-3">
      <div class="row align-items-center">
        <div class="col">
          <strong class="mb-0">Enable subsystem</strong>
          <p class="text-muted mb-0">Turn the NTP subsystem on or off.</p>
        </div>
        <div class="col-auto">
          <ToggleSwitch
              id="ntp_enable"
              :value="this.value"
              @input="setNtp"
              :width="100"
              :height="50"
          >
          </ToggleSwitch>
        </div>
      </div>
    </div>

  </div>

</template>

<script>
import ToggleSwitch from "../../components/ToggleSwitch";
import jwtInterceptor from "../../shared/jwt.interceptor";
export default {
  name: "NTPSettings",
  components: {ToggleSwitch},
  props: {
    value:{
      type: Boolean,
      required: true,
      default: false
    },
  },

  methods: {
    async setNtp(value) {
      //console.log(value)
      //const response = api.SetSubsystem("ntp_timekeeper", value)
      //this.subStatus = value
      let payload = {
        subsystem: "ntp_timekeeper",
        enable: value
      }
      await jwtInterceptor.post("/api/set-subsystem", payload, {
        withCredentials: true,
        credentials: "include",
        headers: {
          'X-Requested-With': 'XMLHttpRequest' // CSRF prevention
        },
      }).then(function (response) {
        // handle success
        //console.log(response);
      }).catch(function (error) {
        // handle error
        //console.log(error);
      })
    }
  }
}
</script>

<style scoped>
.list-group {
  border-right: 1px solid var(--border-color) !important;
  border-left: 1px solid var(--border-color) !important;
  border-bottom: 1px solid var(--border-color) !important;
  border-radius: 0 0 .25rem .25rem;
}
</style>