import actions from './actions';
import mutations from './mutations';

const state = {
  playlist: null,
  playlists: null,
  success: null,
  message: null,
  error: null,
};

const getters = {
  getPlaylists(state) {
    return state.playlists;
  },
  getPlaylist(state) {
    return state.playlist;
  },
  getError(state) {
    return state.error;
  },
  getSuccess(state) {
    return state.success;
  },
  getResponseMessage(state) {
    return state.message;
  },
};

export default {
  namespaced: true,
  state: state,
  mutations: mutations,
  actions: actions,
  getters: getters,
};
