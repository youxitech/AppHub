<template lang="pug">
.admin-version.text-center(v-if="version != null")
  .text-3xl
    span {{ app && app.app.name }}
    router-link.text-blue-500(
      class="hover:underline"
      :to="`/${ $route.params.id }/${ $route.params.version }`"
      target="_blank"
    ) {{ version.version.version}}
  .text-sm.mt-2.text-gray-600 bundleID: {{ app && app.app.bundleID }}
  pkg-list(
    :pkgs="version.packages.map(item => ({...item, version: version.version.version }))"
    )
</template>

<script>
import PkgList from "pkg-list"

export default {
  data() {
    return {
      version: null,
    }
  },

  computed: {
    app() {
      return this.$store.state.app
    },
  },

  mounted() {
    this.$store.dispatch("getAppInfo", this.$route.params.id)
    this.fetchVersion()
  },

  methods: {
    fetchVersion() {
      axios.get(`/apps/${ this.$route.params.id }/${ this.$route.params.version }`)
        .then(res => {
          this.version = res.data
        })
        .catch(_displayError)
    },
  },

  components: {
    PkgList,
  },
}
</script>

<style lang="stylus">

</style>
