import Vue from 'vue/dist/vue.js'
import Vuex from 'vuex'
import VueRouter from "vue-router"

import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import router from './router'
import store from "./store"
// import vuetify from "./plugins/vuetify";
import App from './App'

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