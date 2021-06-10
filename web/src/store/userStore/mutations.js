const mutations = {
  ADD_USER(state, user) {
    if (state.loadedUsers.some((u) => u.username === user.username)) {
      state.loadedUsers.splice(
          state.loadedUsers.indexOf((u) => u.username === user.username),
          1,
      );
    }
    state.loadedUsers.push(user);
  },
};

export default mutations;
