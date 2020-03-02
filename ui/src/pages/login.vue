<template lang="pug">
.h-screen.bg-gray-100.flex.justify-center.items-center
  div(class="focus-within:shadow-outline rounded-lg")
    div.max-w-sm.flex.shadow.rounded-lg.overflow-hidden
      _input.login__input(
        type="email"
        v-model="token"
        placeholder="Enter your token"
        @enter="login"
      )
      button.login__btn(
        @click="login"
      ) Login
</template>

<script>
export default {
  data() {
    return {
      token: "",
    }
  },

  methods: {
    login() {
      axios.post("/login", {
        token: this.token,
      })
        .then(() => {
          _db.token = "admin"
          this.$router.push("/admin")
        })
        .catch(_displayError)
    },
  },
}
</script>

<style lang="stylus">
.login__btn
  @apply text-sm block bg-teal-500 text-white \
    uppercase tracking-wide font-semibold \
    px-6 py-4 w-auto rounded-none shadow-none

.login__btn:focus
  @apply outline-none bg-teal-400

.login__btn:hover
  @apply bg-teal-400

.login__input
  @apply shadow-none block w-full rounded-lg border \
    border-transparent rounded-r-none mb-0 \
    text-black flex-1 px-6 py-4

.login__input:focus
  @apply border-teal-500 outline-none
</style>
