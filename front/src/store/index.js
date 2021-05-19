import Vuex from 'vuex'
import authModule from './authStore/index.js'
import videoModule from './videoStore/index.js'
import Vue from "vue";
Vue.use(Vuex)

export default new Vuex.Store({
    state: {},
    mutations: {},
    actions: {},
    modules: {
        auth: authModule,
        video: videoModule
    },
})