import Vue from 'vue/dist/vue.js';
import {createVuetify} from './plugins/vuetify';
import App from './App';
import 'material-design-icons-iconfont/dist/material-design-icons.css';
import {createStore} from './store';
import {registerPlugins} from "@/plugins";
import routerModule from "@/router";
import { sync } from 'vuex-router-sync'

Vue.config.productionTip = false;
registerPlugins(Vue);

const store = createStore();
const vuetify = createVuetify();
const router = routerModule.getInstance(store)
// createRouter(vuetify, store)
sync(store, router)

new Vue({
  el: '#app',
  router,
  store,
  vuetify,
  template: '<App/>',
  components: {App},
});