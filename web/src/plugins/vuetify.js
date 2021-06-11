import Vuetify from 'vuetify/lib';
// import {icons} from "./icons";
// import {IS_SERVER} from "@/util/globals";

export function createVuetify() {
  const config = {
    // breakpoint: { mobileBreakpoint: 'md' },
    // icons,
    // theme: {
    //   dark: store.state.user.theme.dark,
    //   options: {
    //     minifyTheme: IS_SERVER ? require('minify-css-string').default : undefined,
    //     variations: false,
    //   },
    //   themes: {
    //     light: {
    //       primary: '#1867C0',
    //       secondary: '#5CBBF6',
    //       tertiary: '#E57373',
    //       accent: '#005CAF',
    //     },
    //   },
    // },
    // rtl: false,
  };

  return new Vuetify(config);
}

export function useVuetify (app) {
  app.use(Vuetify)
}
