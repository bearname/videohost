const getters = {
  getUser(state, userId) {
    if (state.loadedUsers.some((user) => user.username === userId)) {
      return state.loadedUsers.find((user) => user.username === userId);
    } else {
      // Here I'll have to request from the server!!
      return {};
    }
  },
  getUserVideos(state) {
    return {
      videos: state.userVideos,
      countAllVideos: state.countUserVideos,
    };
  },
  isSuccess(state) {
    return state.isSuccess;
  },
};

export default getters;
