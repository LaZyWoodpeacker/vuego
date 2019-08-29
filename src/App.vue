<template>
  <router-view></router-view>
</template>

<script>
export default {
  props: {
    source: String
  },
  data: () => ({}),
  created() {
    this.$vuetify.theme.dark = true;
    this.getPage();
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
    getPage() {
      if (this.authCheck()) {
        this.$router.replace("/list");
      } else {
        this.$router.replace("/auth");
      }
    }
  }
};
</script>
