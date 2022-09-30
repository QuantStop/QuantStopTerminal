<template>
  <Form @submit="handleUpdateSettings" :validation-schema="schema">
  <div class="list-group mb-5 shadow pt-3">
    <div class="list-group-item m-3">

      <div class="row align-items-center">

        <div class="col-auto">
          <strong class="mb-0">Go Max Processors</strong>
          <p class="text-muted mb-0">Set the number of logical processors that the system will use.</p>
        </div>
        <div class="col">
          <Field name="maxProcs" type="number" class="form-control" />
          <ErrorMessage name="maxProcs" class="error-feedback" />
        </div>

      </div>

      <hr class="dropdown-divider mt-3 mb-3">

      <div class="row align-items-center">

        <div class="col-auto">
          <strong class="mb-0">Log Level</strong>
          <p class="text-muted mb-0">Logging levels.</p>
        </div>
        <div class="col">
          <Field name="logLevel" type="text" class="form-control" />
          <ErrorMessage name="logLevel" class="error-feedback" />
        </div>

      </div>

      <hr class="dropdown-divider mt-3 mb-3">

      <div class="row align-items-center">

        <div class="col-auto">
          <strong class="mb-0">Log Output</strong>
          <p class="text-muted mb-0">Logging output types.</p>
        </div>
        <div class="col">
          <Field name="logOutput" type="text" class="form-control" />
          <ErrorMessage name="logOutput" class="error-feedback" />
        </div>

      </div>

      <hr class="dropdown-divider mt-3 mb-3">

      <div class="row align-items-center">

        <div class="col-auto">
          <strong class="mb-0">Log Filename</strong>
          <p class="text-muted mb-0">Log filename.</p>
        </div>
        <div class="col">
          <Field name="logFilename" type="text" class="form-control" />
          <ErrorMessage name="logFilename" class="error-feedback" />
        </div>

      </div>

      <hr class="dropdown-divider mt-3 mb-3">

      <!--      <div class="row align-items-center">

      <div class="col-auto">
        <strong class="mb-0">Log Rotate</strong>
        <p class="text-muted mb-0">Rotate the log file?</p>
      </div>
      <div class="col">
        <Field name="logRotate" type="" class="form-control" />
        <ErrorMessage name="apiUrl" class="error-feedback" />
      </div>

      </div>

      <hr class="dropdown-divider mt-3 mb-3">-->

      <div class="row align-items-center">
        <div class="form-group">
          <button class="btn btn-primary btn-block" :disabled="loading">
            <span
                v-show="loading"
                class="spinner-border spinner-border-sm"
            ></span>
            <span>Update</span>
          </button>
        </div>

        <div class="form-group">
          <div v-if="message" class="alert alert-danger" role="alert">
            {{ message }}
          </div>
        </div>
      </div>

    </div>

  </div>
  </Form>
</template>

<script>
import {ErrorMessage, Field, Form} from "vee-validate";
import * as yup from "yup";
import jwtInterceptor from "../../shared/jwt.interceptor";

export default {
  name: "SystemSettings",
  components: {
    Form,
    Field,
    ErrorMessage,
  },
  data() {
    const schema = yup.object().shape({
      apiUrl: yup
          .string()
          .required("url is required!")
          .min(3, "Must be at least 3 characters!")
          .max(50, "Must be maximum 50 characters!"),
      maxProcs: yup
          .number()
          .required("max processors is required!")
          .min(-1, "Must be at least -1!")
          .max(20, "Must be maximum 20!"),
    });

    return {
      successful: false,
      loading: false,
      message: "",
      schema,
    };
  },
  /*computed: {
    ...mapGetters("auth", {
      getUserProfile: "getUserProfile",
    }),
  },*/
  methods: {
    handleUpdateSettings(settings) {
      const payload = {
        apiUrl: settings.apiUrl,
        maxProcs: settings.maxProcs
      };
      const response = jwtInterceptor.post("/api/set-sysconfig", payload, {
        withCredentials: true,
        credentials: "include",
        headers: {
          'X-Requested-With': 'XMLHttpRequest' // CSRF prevention
        },
      });
      if (response && response.data) {

      }

    },
    getSettings() {

    }
  },
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