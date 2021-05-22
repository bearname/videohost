import Vue from 'vue/dist/vue.js'
import Vuex from 'vuex'
import VueRouter from "vue-router"

import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import router from './router'
import store from "./store"
import App from './App'

import Cookie from "./util/cookie";
Cookie.eraseCookie("userId")
document.cookie = "userId=; Max-Age=-99999999;";
Vue.use(Vuex)
Vue.config.productionTip = false
Vue.use(Vuetify)
Vue.use(VueRouter)

const options = {};
new Vue({
    el: '#app',
    router,
    store,
    vuetify: new Vuetify(options),
    template: '<App/>',
    components: {App},
})