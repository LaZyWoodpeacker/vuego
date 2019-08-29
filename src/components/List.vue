<template>
  <v-app id="inspire">
    <v-navigation-drawer v-model="menuVisible" app clipped>
      <v-list dense>
        <v-list-item @click="onExit">
          <v-list-item-action>
            <v-icon>mdi-logout</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Exit</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar app clipped-left :collapse-on-scroll="true">
      <v-app-bar-nav-icon @click.stop="menuVisible = !menuVisible"></v-app-bar-nav-icon>
      <v-toolbar-title>{{name}}</v-toolbar-title>
    </v-app-bar>

    <v-content>
      <v-container>
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
          <v-icon>mdi-file-image</v-icon>
        </v-btn>
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
    onInter() {
      axios.get("/api/files").then(res => {
        this.list = res.data.sort((a, b) => b.ModTime - a.ModTime);
      });
    },
    onExit() {
      document.cookie = "auth=;expires=Thu, 01 Jan 1970 00:00:00 GMT";
      this.auth = false;
    }
  }
};
</script>
