import Vuex from 'vuex'
import authModule from './authStore/index.js'
import videoModule from './videoStore/index.js'
import userModule from './userStore/index.js'
import Vue from "vue";

import createPersistedState from 'vuex-persistedstate'
import * as Cookies from 'js-cookie'

Vue.use(Vuex)

const modules = {
    video: videoModule,
    auth: authModule,
    user: userModule
};

const plugins = [
    createPersistedState({
        getState: (key) => Cookies.getJSON(key),
        setState: (key, state) => Cookies.set(key, state, {expires: 3, secure: true})
    })
];

export default new Vuex.Store({
    state: {},
    mutations: {},
    actions: {},
    modules: modules,
    plugins: plugins
})