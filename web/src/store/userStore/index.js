import actions from './actions';
import getters from './getters';
import mutations from './mutations';

const state = {
  loadedUsers: [],
  userVideos: null,
  countUserVideos: 0,
};

export default {
  namespaced: true,
  state: state,
  mutations: mutations,
  actions: actions,
  getters: getters,
};
