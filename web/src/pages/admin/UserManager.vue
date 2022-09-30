<template>

  <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
    <h2>User Manager</h2>
    <div class="btn-toolbar mb-2 mb-md-0">
      <button type="button" class="btn btn-sm btn-outline-secondary" data-bs-toggle="modal" data-bs-target="#exampleModal">New User</button>
    </div>
  </div>

  <table>
    <thead>
    <tr>
      <th scope="col">ID</th>
      <th scope="col">Username</th>
      <th scope="col">Password</th>
      <th scope="col">Salt</th>
      <th scope="col">Roles</th>
      <th scope="col">Actions</th>
    </tr>
    </thead>
    <tbody>
    <tr v-for="user in users">
      <td data-label="ID">{{user.ID}}</td>
      <td data-label="Username">{{user.Username}}</td>
      <td data-label="Password">{{user.Password}}</td>
      <td data-label="Salt">{{user.Salt}}</td>
      <td data-label="Roles">
        <span v-for="role in user.Roles" class="badge bg-primary m-1">{{ role }}</span>
      </td>
      <td data-label="Actions">
        <button type="button" class="btn btn-sm btn-warning m-1">Edit</button>
        <button type="button" class="btn btn-sm btn-danger m-1">Delete</button>
      </td>
    </tr>
    </tbody>
  </table>

  <!-- Create User Modal -->
  <div class="modal fade" id="exampleModal" tabindex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
    <div class="modal-dialog">
      <div class="modal-content">
        <Form @submit="register" :validation-schema="schema">

          <div class="modal-header">
            <h5 class="modal-title" id="exampleModalLabel">Create a new User</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">

            <div v-if="!successful">
              <div class="form-group">
                <label for="username">Username</label>
                <Field name="username" type="text" class="form-control" />
                <ErrorMessage name="username" class="error-feedback" />
              </div>
              <div class="form-group">
                <label for="password">Password</label>
                <Field name="password" type="password" class="form-control" />
                <ErrorMessage name="password" class="error-feedback" />
              </div>
            </div>

            <div v-if="message" class="alert" :class="successful ? 'alert-success' : 'alert-danger'">
              {{ message }}
            </div>
          </div>
          <div class="modal-footer">
            <button v-if="successful" type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
            <div v-if="!successful" class="form-group">
              <button class="btn btn-primary btn-block" :disabled="loading">
                <span v-show="loading" class="spinner-border spinner-border-sm"></span>
                Create User
              </button>
            </div>
          </div>
        </Form>
      </div>
    </div>
  </div>

</template>

<script>
import {userStore} from "../../store/userStore";
import { Form, Field, ErrorMessage } from "vee-validate";
import * as yup from "yup";
import jwtInterceptor from "../../shared/jwt.interceptor";
import axios from "axios";
export default {
  name: "UserManager",
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
      password: yup
          .string()
          .required("Password is required!")
          .min(6, "Must be at least 6 characters!")
          .max(40, "Must be maximum 40 characters!"),
    });
    return {
      userStore,
      users: [

      ],
      successful: false,
      loading: false,
      message: "",
      schema,
    };
  },

  methods: {

    async getAll() {
      let response = await jwtInterceptor.get("/api/get-users", {
        withCredentials: true,
        credentials: "include",
        headers: {
          'X-Requested-With': 'XMLHttpRequest' // CSRF prevention
        },
      });

      if (response && response.data) {
        this.users = response.data
      } else {

      }
    },

    async register(user) {
      this.loading = true;
      const payload = {
        username: user.username,
        email: user.email,
        password: user.password
      };
      /*await this.actionRegisterApi(payload).then(
          () => {
            this.message = "Success!"
            this.successful = true;
            this.loading = false;
          },
          (error) => {
            this.loading = false;
            this.successful = false;
            this.message = error.toString() + " | " + error.response.data.error;
          }
      );*/

    },
  },
  beforeMount() {
    this.getAll();
  },
}
</script>

<style scoped>
table {
  color: var(--text-primary-color) !important;
  background-color: var(--background-color-secondary) !important;
  border-color: var(--border-color) !important;
  border-collapse: collapse;
  margin: 0;
  padding: 0;
  width: 100%;
  table-layout: fixed;

}
table tr {
  color: var(--text-primary-color) !important;
  background-color: var(--background-color-secondary);
  border: 1px solid var(--border-color);
  padding: .35em;
}
table tr:hover {
  border: 1px solid var(--table-border-hover) !important;
}
table th,
table td {
  overflow: hidden;
  text-overflow: ellipsis;
  padding: .625em;
  /*text-align: center;*/
  color: var(--text-primary-color) !important;
  background-color: var(--background-color-secondary) !important;
  border: 1px solid var(--border-color) !important;
}
table th {
  font-size: .85em;
  letter-spacing: .1em;
  text-transform: uppercase;
}
.modal-content, .modal-body, .modal-header {
  color: var(--text-primary-color) !important;
  background-color: var(--background-color-secondary) !important;
}


@media screen and (max-width: 600px) {
  table {
    border: 0;
    background-color: var(--background-color-primary) !important;
  }
  table caption {
    font-size: 1.3em;
  }
  table thead {
    border: none;
    clip: rect(0 0 0 0);
    height: 1px;
    margin: -1px;
    overflow: hidden;
    text-overflow: ellipsis;
    padding: 0;
    position: absolute;
    width: 1px;
  }
  table tr {
    /*border: 3px solid #ddd;*/
    display: block;
    margin-bottom: .625em;
  }
  table td {
    /*border-bottom: 1px solid #ddd;*/
    display: block;
    font-size: .8em;
    text-align: right;
  }
  table td::before {
    /*
    * aria-label has no advantage, it won't be read inside a table
    content: attr(aria-label);
    */
    content: attr(data-label);
    float: left;
    font-weight: bold;
    text-transform: uppercase;
  }
  table td:last-child {
    border-bottom: 0;
  }
}



</style>