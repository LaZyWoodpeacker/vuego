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
    if (this.authCheck()) {
      this.$router.push("/list");
    } else {
      this.$router.push("/auth");
    }
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
    }
  }
};
</script>
