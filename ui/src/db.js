import Vue from "vue"

const DB_KEY = "app-hub:v1"

const db = new Vue({
  data: {
    token: "",
  },

  created() {
    const data = window.localStorage.getItem(DB_KEY)

    if(data != null) {
      Object.assign(this.$data, JSON.parse(data))
    }
  },
})

db.$watch(function() { return this.$data }, function() {
  window.localStorage.setItem(DB_KEY, JSON.stringify(this.$data))
}, {
  deep: true,
  immediate: true,
})

export default db
