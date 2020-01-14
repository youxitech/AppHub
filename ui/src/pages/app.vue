<template lang="pug">
.flex-1.bg-gray-200.p-6(v-if="app")
  .rounded.overflow-hidden.shadow-lg.p-6.flex.flex-col.bg-white.h-full
    .flex
      img.w-24.h-24.mr-2.rounded(:src="_getAsset('icon', app.app.platform, app.app.bundleID)")
      .ml-5
        .font-semibold {{ app.app.name }}
        .text-sm.text-gray-600.mt-3 Platform: {{ app.app.platform }}
        .text-sm.text-gray-600 BundleID: {{ app.app.bundleID }}
        .text-sm.text-gray-600 Download: {{ app.app.downloadCount }}
      .flex.ml-auto.items-center APP alias: {{ app.app.alias }}

  version(:id="String(app.versions[0].id)")
</template>

<script>
import Version from "./version"

export default {
  data() {
    return {
      app: null,
    }
  },

  mounted() {
    this.fetchApp()
  },

  methods: {
    fetchApp() {
      return axios.get(`/${ this.$route.params.id }`)
        .then(res => {
          this.app = res.data
        })
        .catch(_displayError)
    },
  },

  components: {
    Version,
  },
}
</script>

<style lang="stylus">

</style>
