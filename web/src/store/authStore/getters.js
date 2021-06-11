const getters = {
  username: state => state.username,

  getCurrentUser(state) {
    return state.user;
  },
  isLoggedIn(state) {
    console.log(state.user);
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
