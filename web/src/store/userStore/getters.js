const getters = {
    getUser(state, userId) {
        if (state.loadedUsers.some(user => user.username === userId)) {
            return state.loadedUsers.find(user => user.username === userId)
        } else {
            //Here I'll have to request from the server!!
            return {}
        }
    },
};

export default getters