<template lang="pug">
.text-center.flex-1(
  :class="{ 'mt-8': isDetail }"
)
  template(v-if="version != null")
    router-link.text-3xl.text-blue-500(
      v-if="!isDetail"
      :to="'/' + alias + '/' + id"
      class="hover:underline"
    ) {{ version.version.version }}
    .text-3xl(v-else) {{ app && app.app.name }} {{ version.version.version}}
    .text-sm.mt-2.text-gray-600(v-if="isDetail") bundleID: {{ app && app.app.bundleID }}

    pkg-list(:pkgs="version.packages.map(item => ({...item, version: version.version.version }))")
</template>

<script>
import PkgList from "pkg-list"

export default {
  props: {
    id: {
      type: String,
      default: "",
    },

    alias: {
      type: String,
      default: "",
    },

    isDetail: {
      type: Boolean,
      default: true,
    },
  },

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

    if(this.id == null) return

    this.fetchVersion()
  },

  watch: {
    id() {
      if(this.id == null) return

      this.version = null
      this.fetchVersion()
    },
  },

  methods: {
    fetchVersion() {
      axios.get(`/apps/${ this.alias || this.$route.params.id }/${ this.id || this.$route.params.version }`)
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
.version__hover:hover
  transform: scale(2)
  transition: 0.5s

.version__modal
  padding: 12px
  background: white
  height: 400px !important
</style>
