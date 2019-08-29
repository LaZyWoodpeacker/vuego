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
          <v-card v-for="i in list" :key="i.Name" width="200" outlined :raised="true">
            <v-img class="white--text" height="200px" :src="'/api/images/'+i.Name"></v-img>
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
              <v-btn text @click="onRemove(i.Name)">Delete</v-btn>
              <v-btn
                text
                @click="chdialogValues.oldFileName=i.Name;chdialogValues.newFileName=i.Name;chdialogValues.act=!chdialogValues.act"
              >Edit</v-btn>
            </v-card-actions>
          </v-card>
        </v-layout>
        <v-btn fixed dark fab bottom right color="pink" @click="onAddFile">
          <v-icon>mdi-file-image</v-icon>
        </v-btn>
      </v-container>
      <v-dialog v-model="chdialogValues.act" persistent max-width="600px" value="chdialogValue">
        <v-card>
          <v-card-title>
            <span class="headline">Change filename</span>
          </v-card-title>
          <v-card-text>
            <v-container>
              <v-text-field label="Filename*" required v-model="chdialogValues.newFileName"></v-text-field>
            </v-container>
            <small>*indicates required field</small>
          </v-card-text>
          <v-card-actions>
            <div class="flex-grow-1"></div>
            <v-btn color="blue darken-1" text @click="chdialogValues.act = false">Close</v-btn>
            <v-btn color="blue darken-1" text @click="onRename">Rename</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-snackbar v-model="snackbar.act" :timeout="2000">
        {{ snackbar.message }}
        <v-btn color="pink" text @click="snackbar.act = false">Close</v-btn>
      </v-snackbar>
    </v-content>

    <v-footer app>
      <span>&copy; 2019</span>
    </v-footer>
    <input
      type="file"
      class="hide"
      name="uploadfile"
      id="uploadfile"
      accept="image/*, .pdf, video/*"
    />
  </v-app>
</template>

<script>
import axios from "axios";
import io from "socket.io-client";
export default {
  props: {
    source: String
  },
  data: () => ({
    name: "Media files",
    menuVisible: false,
    list: [],
    SOC: false,
    snackbar: {
      act: false,
      message: ""
    },
    chdialogValues: {
      act: false,
      oldFileName: "",
      newFileName: ""
    },
    auth: false
  }),
  created() {
    this.auth = this.authCheck();
    if (this.auth) {
      this.createSocket();
      this.onInter();
    }
    this.$vuetify.theme.dark = true;
  },
  methods: {
    createSocket() {
      this.SOC = io(location.host, {
        path: "/api/socket"
      });
      this.SOC.on("connect", () => (this.api_connected = true));
      this.SOC.on("disconnect", () => (this.api_connected = false));
      this.SOC.on("CHANGEFS", data => {
        this.onInter();
      });
    },
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
        this.list = res.data
          .filter(e => {
            return e.Name.match(/(png|bmp|jpeg|jpg|bmp|mp4|mpeg|mov)$/i);
          })
          .sort((a, b) => b.ModTime - a.ModTime);
      });
    },
    onExit() {
      document.cookie = "auth=;expires=Thu, 01 Jan 1970 00:00:00 GMT";
      this.$router.push("/auth");
    },
    onAddFile() {
      this.menuVisible = false;
      var imagefile = document.getElementById("uploadfile");
      imagefile.click();
      setInterval(() => {
        if (imagefile.files.length > 0) {
          var formData = new FormData();
          formData.append("image", imagefile.files[0]);
          document.getElementById("uploadfile").value = "";
          axios
            .post("/api/upload", formData, {
              headers: {
                "Content-Type": "multipart/form-data"
              }
            })
            .then(res => {
              // this.$refs.listRef.onInter();
            });
        }
      }, 1000);
    },
    onRemove(val) {
      axios.post("/api/remove", { image: val }).then(res => {});
    },
    onRename(val, val2) {
      if (this.chdialogValues.oldFileName !== this.chdialogValues.newFileName) {
        axios
          .post("/api/edit", {
            oldname: this.chdialogValues.oldFileName,
            newname: this.chdialogValues.newFileName
          })
          .then(res => {
            this.snackbar.message = res.data.message;
            this.snackbar.act = true;
          });
      }
      this.chdialogValues.act = false;
    }
  }
};
</script>

<style lang="scss" scoped>
.hide {
  visibility: hidden;
  width: 0;
  height: 0;
}
</style>