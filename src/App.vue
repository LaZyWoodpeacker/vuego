<template>
  <v-app id="inspire">
    <v-navigation-drawer v-model="menuVisible" app clipped>
      <v-list dense>
        <v-list-item>
          <v-list-item-action>
            <v-icon>mdi-view-dashboard</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Dashboard</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item>
          <v-list-item-action>
            <v-icon>mdi-settings</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Settings</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item>
          <v-list-item-action>
            <v-icon>mdi-settings</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Settings</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar app clipped-left collapse-on-scroll="true">
      <v-app-bar-nav-icon @click.stop="menuVisible = !menuVisible"></v-app-bar-nav-icon>
      <v-toolbar-title>{{name}}</v-toolbar-title>
    </v-app-bar>

    <v-content>
      <v-container v-if="this.auth">
        <v-layout row wrap justify-center>
          <v-card v-for="i in list" :key="i.Name" width="200" outlined>
            <v-card-title>{{ i.Name }}</v-card-title>
            <v-card-text>
              {{ new Date(i.ModTime*1000).toLocaleString("ru-RU",{ year: '2-digit',
              month: 'short',
              day: 'numeric',
              weekday:'long',
              hour:"numeric",
              minute:"numeric",
              timeZoneName:"short"
              }) }}
            </v-card-text>
            <v-card-actions>
              <v-btn text>Click</v-btn>
            </v-card-actions>
          </v-card>
        </v-layout>
        <v-btn fixed dark fab bottom right color="pink">
          <v-icon>mdi-plus</v-icon>
        </v-btn>
      </v-container>
      <v-container v-if="!this.auth">
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="authForm.name"
              :rules="nameRules"
              :counter="20"
              label="Login"
              required
            ></v-text-field>
          </v-col>

          <v-col cols="12" md="6">
            <v-text-field
              v-model="authForm.pass"
              :rules="nameRules"
              :counter="20"
              label="Password"
              required
            ></v-text-field>
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
    name: "Media files",
    menuVisible: false,
    list: [],
    authForm: {
      name: "admin",
      pass: "123"
    },
    auth: false
  }),
  created() {
    this.auth = this.authCheck();
    if (this.auth) {
      this.onInter();
    }
    this.$vuetify.theme.dark = true;
  },
  methods: {
    authCheck() {
      return document.cookie
        .split("; ")
        .filter(e => e.length > 0)
        .reduce((acc, em) => {
          let [key, val] = em.split("=");
          acc[key] = val;
          return acc;
        }, {}).auth;
    },
    onAuthClick() {
      axios
        .post("/api/get", {
          login: this.authForm.name,
          pass: this.authForm.pass
        })
        .then(res => {
          this.auth = this.authCheck();
        })
        .catch(req => {
          alert("Error");
        });
    },
    onInter() {
      axios.get("/api/files").then(res => {
        this.list = res.data.sort((a, b) => b.ModTime - a.ModTime);
      });
    }
  }
};
</script>
