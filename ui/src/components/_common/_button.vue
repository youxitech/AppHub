<template lang="pug">
button._button(
  :class="klass"
  @click="$emit('click')"
)
  slot
</template>

<script>
export default {
  computed: {
    klass() {
      return "bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
    },
  },
}
</script>

<style lang="stylus">

</style>
