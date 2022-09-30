<template>

  <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
    <h2>Profile</h2>
  </div>

  <div class="container rounded mt-5 mb-5 py-4 card">
    <div class="row p-5 mb-4 rounded-3">
      <div class="col-md-5 border-right">
        <div class="d-flex flex-column align-items-center text-center p-3 py-5">
          <img
              id="profile-img"
              src="//ssl.gstatic.com/accounts/ui/avatar_2x.png"
              class="profile-img-card"
              alt=""
          />
          <span class="font-weight-bold">{{userStore.userProfile.username}}</span>

        </div>
      </div>
      <Form @submit="handleUpdateProfile" :validation-schema="schema">
        <div class="form-group">
          <label for="id"></label>
          <Field name="id" type="hidden" class="form-control" v-model="userStore.userProfile.id" />
        </div>
        <div class="form-group">
          <label for="username">Username</label>
          <Field name="username" type="text" class="form-control" v-model="userStore.userProfile.username" />
          <ErrorMessage name="username" class="error-feedback" />
        </div>


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
      </Form>
    </div>
  </div>
</template>

<script>
import { Form, Field, ErrorMessage } from "vee-validate";
import * as yup from "yup";
import {userStore} from "../store/userStore";
export default {
  components: {
    Form,
    Field,
    ErrorMessage,
  },
  data() {
    const schema = yup.object().shape({
      username: yup
          .string()
          .required("Username is required!")
          .min(3, "Must be at least 3 characters!")
          .max(20, "Must be maximum 20 characters!"),
    });

    return {
      userStore,
      successful: false,
      loading: false,
      message: "",
      schema,
    };
  },

  methods: {
    handleUpdateProfile(user) {

      /*this.message = "";
      this.successful = false;
      this.loading = true;
      //let userid = this.$store.state.auth.user.id
      //console.log("profile ID: " + userid)
      /!*this.$store.dispatch("auth/update", user).then(
          (data) => {
            this.message = data.response;
            this.successful = true;
            this.loading = false;
          },
          (error) => {
            this.message =
                (error.response &&
                    error.response.data &&
                    error.response.data.message) ||
                error.message ||
                error.toString();
            this.successful = false;
            this.loading = false;
          }
      );*!/
      UserService.updateProfile(user).then(
          (response) => {
            this.content = response.data;
          },
          (error) => {
            this.content =
                (error.response &&
                    error.response.data &&
                    error.response.data.message) ||
                error.message ||
                error.toString();
          }
      );*/
    },
  },
};
</script>

<style scoped>

.form-control:focus {
  box-shadow: none;
  border-color: #BA68C8
}

.profile-button {
  background: rgb(99, 39, 120);
  box-shadow: none;
  border: none
}

.profile-button:hover {
  background: #682773
}

.profile-button:focus {
  background: #682773;
  box-shadow: none
}

.profile-button:active {
  background: #682773;
  box-shadow: none
}

.back:hover {
  color: #682773;
  cursor: pointer
}

.labels {
  font-size: 11px
}

.add-experience:hover {
  background: #BA68C8;
  color: #fff;
  cursor: pointer;
  border: solid 1px #BA68C8
}
</style>