<template lang="pug">
.h-screen.bg-gray-100.flex.justify-center.items-center
  div(class="sm:focus-within:shadow-outline rounded-lg")
    div(class="max-w-sm sm:flex sm:shadow sm:rounded-lg sm:overflow-hidden")
      _input(
        class="shadow sm:shadow-none block w-full rounded-lg focus:border-teal-500 border border-transparent sm:rounded-r-none mb-4 sm:mb-0 text-black sm:flex-1 px-6 py-4 lg:py-5 focus:outline-none"
        type="email"
        v-model="token"
        placeholder="Enter your token"
        @enter="login"
      )
      button(
        class="text-sm shadow sm:shadow-none block w-full sm:w-auto rounded-lg sm:rounded-none focus:outline-none bg-teal-500 hover:bg-teal-500-dark focus:bg-teal-500-dark text-white text-shadow uppercase tracking-wide font-semibold px-6 py-4 lg:py-5"
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
</style>
