const getters = {
  getCurrentUser(state) {
    return state.user;
  },
  isLoggedIn(state) {
    if (!state.user) {
      return false;
    }
    return state.user.loggedIn;
  },
  getAccessToken(state) {
    return state.user.accessToken;
  },
  getRefreshToken(state) {
    return state.user.refreshToken;
  },
  getTokenHeader(state) {
    return 'Bearer ' + state.user.accessToken;
  },
};

export default getters;
