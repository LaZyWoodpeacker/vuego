<template>
  <v-app id="auth">
    <v-content>
      <v-container>
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field v-model="authForm.login" :counter="20" label="Login" required></v-text-field>
          </v-col>

          <v-col cols="12" md="6">
            <v-text-field v-model="authForm.pass" :counter="20" label="Password" required></v-text-field>
          </v-col>

          <v-col cols="12" md="12" class="text-center">
            <v-btn @click="onAuthClick">Inter</v-btn>
          </v-col>
        </v-row>
      </v-container>
    </v-content>

    <v-footer app>
      <span>&copy; 2019</span>
    </v-footer>
  </v-app>
</template>

<script>
import axios from "axios";
export default {
  props: {
    source: String
  },
  data: () => ({
    name: "Auth",
    authForm: {
      login: "admin",
      pass: "123"
    }
  }),
  created() {},
  methods: {
    onAuthClick() {
      axios
        .post("/api/get", {
          login: this.authForm.login,
          pass: this.authForm.pass
        })
        .then(res => {})
        .catch(req => {
          alert(req);
        });
    }
  }
};
</script>
