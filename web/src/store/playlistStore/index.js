import actions from './actions';
import mutations from './mutations';

const state = {
  playlist: null,
  error: null
};

const getters = {
  getPlaylist(state) {
    return state.playlist;
  },
  getError(state) {
    return state.error;
  },
};

export default {
  namespaced: true,
  state: state,
  mutations: mutations,
  actions: actions,
  getters: getters,
};
