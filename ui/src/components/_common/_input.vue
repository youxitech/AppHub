<template lang="pug">
input._input(
  ref="input"
  class="bg-white focus:outline-none focus:shadow-outline border border-gray-300 rounded-lg py-2 px-4 block w-full appearance-none leading-normal"
  :type="type"
  :placeholder="placeholder"
  :value="value"
  @input="e => $emit('input', e.target.value)"
  @keyup.enter="$emit('enter')"
)
</template>

<script>
export default {
  props: {
    value: {
      type: String,
      default: "",
    },

    type: {
      type: String,
      default: "text",
    },

    placeholder: {
      type: String,
      default: "",
    },
  },

  methods: {
    focus() {
      this.$refs.input.focus()
    },

    blur() {
      this.$refs.input.blur()
    },
  },
}
</script>

<style lang="stylus">

</style>
