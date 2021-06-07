import Vue from 'vue/dist/vue.js'
import Vuex from 'vuex'
import VueRouter from "vue-router"

import router from './router'
import store from "./store"
import App from './App'
import vuetify from './plugins/vuetify'
import 'material-design-icons-iconfont/dist/material-design-icons.css'
import Vuetify from 'vuetify/lib';
Vue.use(Vuex)
Vue.config.productionTip = false
Vue.use(VueRouter)
Vue.use(Vuetify);

new Vue({
    el: '#app',
    router,
    store,
    template: '<App/>',
    vuetify,
    components: {App}
})