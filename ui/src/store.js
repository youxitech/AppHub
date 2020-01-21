import Vue from "vue"
import Vuex from "vuex"

Vue.use(Vuex)

const store = new Vuex.Store({
  state: {
    app: null,
  },

  mutations: {
    updateApp(state, data) {
      state.app = data
    },
  },

  actions: {
    getAppInfo(context, id) {
      const requestURL = _db.token === "" ?
        `/apps/${ id }` : `/admin/apps/${ id }`

      return axios.get(requestURL)
        .then(res => {
          context.commit("updateApp", res.data)
        })
        .catch(_displayError)
    },
  },
})

export default store
