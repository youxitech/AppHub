<template lang="pug">
.flex-1.bg-gray-200.p-6.h-screen(v-if="app")
  app-info(:app="app.app")

  .flex.py-8
    .w-40.bg-white.rounded.shadow-lg.h-auto.overflow-auto.self-start
      .h-10.flex.items-center.justify-center.border-gray-200.border-solid.border-b.cursor-pointer(
        v-for="item in app.versions"
        :key="item.id"
        :class="{ 'bg-blue-400 text-white': item.version === curVersion }"
        @click="curVersion = item.version"
      ) {{ item.version }}

    version(
      :id="curVersion"
      :alias="app.app.alias"
      :isDetail="false"
    )
</template>

<script>
import Version from "./version"
import AppInfo from "app-info"

export default {
  data() {
    return {
      curVersion: null,
    }
  },

  computed: {
    app() {
      return this.$store.state.app
    },
  },

  mounted() {
    this.$store.dispatch("getAppInfo", this.$route.params.id)
      .then(() => {
        this.curVersion = this.app && this.app.versions[0].version
      })
  },

  components: {
    Version,
    AppInfo,
  },
}
</script>

<style lang="stylus">

</style>
