import Vue from 'vue/dist/vue.js'
import Vuetify from 'vuetify'
import router from './router'
import App from './App'
import 'vuetify/dist/vuetify.min.css'
import VueRouter from "vue-router";

Vue.config.productionTip = false
Vue.use(Vuetify)
Vue.use(VueRouter)

new Vue({
    el: '#app',
    router,
    template: '<App/>',
    components: {App},
})