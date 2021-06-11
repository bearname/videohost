import actions from './actions';
import getters from './getters';
import mutations from './mutations';

const state = {
  loadedUsers: [],
  userVideos: null,
  countUserVideos: 0,
  rtl: false,
  theme: {
    dark: false,
    system: false,
    // Provides a 3rd state for the
    // light theme w/ dark fences
    mixed: false,
  },
};

export default {
  namespaced: true,
  state: state,
  mutations: mutations,
  actions: actions,
  getters: getters,
};
