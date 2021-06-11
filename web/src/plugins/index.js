import {useVuetify} from "@/plugins/vuetify";
import {useVuex} from "@/store";
import {useVueRouter} from "@/router";

export function registerPlugins(app) {
  useVuex(app);
  useVuetify(app);
  useVueRouter(app);
}