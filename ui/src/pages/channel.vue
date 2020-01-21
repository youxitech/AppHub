<template lang="pug">
.bg-gray-200.p-6.h-screen.overflow-auto(v-if="app")
  app-info(:app="app.app")

  .text-3xl.text-center.mt-6 {{ channel }}

  pkg-list(
    :pkgs="pkgs.map(item => ({...item.package, version: item.version.version}))"
  )
</template>

<script>
import AppInfo from "app-info"
import PkgList from "pkg-list"

// route params: /:id/channel/:channel
//  id: app alias
//  channel: channel name
export default {
  data() {
    return {
      pkgs: [],
    }
  },

  computed: {
    channel() {
      return this.$route.params.channel
    },

    app() {
      return this.$store.state.app
    },
  },

  mounted() {
    this.fetchChannel()
    this.$store.dispatch("getAppInfo", this.$route.params.id)
  },

  methods: {
    fetchChannel() {
      return axios.get(`/apps/${ this.$route.params.id }/channels/${ this.channel }`)
        .then(res => {
          this.pkgs = res.data.content
        })
        .catch(_displayError)
    },
  },

  components: {
    AppInfo,
    PkgList,
  },
}
</script>

<style lang="stylus">
</style>
