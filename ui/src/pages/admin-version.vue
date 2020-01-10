<template lang="pug">
.admin-version
  version(:key="key" :id="$route.params.version")
</template>

<script>
import Version from "./version"

export default {
  data() {
    return {
      key: 0,
    }
  },

  watch: {
    $route() {
      this.key = this.$route.query.t
    },
  },

  components: {
    Version,
  },
}
</script>

<style lang="stylus">

</style>
