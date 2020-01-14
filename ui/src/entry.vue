<template lang="pug">
#entry
  router-view
  notifications
  modal(name="confirm")
    .flex.items-center.justify-center.w-full.h-full.flex-col.p-6
      div {{ confirmHint }}
      .flex.mt-10
        button.mr-10(class="mr-8 hover:text-teal-500" @click="onConfirm") Confirm
        button(class="mr-8 hover:text-teal-500" @click="onCancel") Cancel
</template>

<script>
export default {
  data() {
    return {
      confirmHint: "",
    }
  },

  mounted() {
    window._showConfirm = (hint) => new Promise((resolve, reject) => {
      this.confirmHint = hint
      this.$modal.show("confirm")
      this.reject = reject
      this.resolve = resolve
    })
  },

  methods: {
    onConfirm() {
      this.$modal.hide("confirm")
      this.resolve()
    },

    onCancel() {
      this.$modal.hide("confirm")
      this.reject()
    },
  },
}
</script>

<style type="text/css">
</style>

<style lang="stylus">
</style>
