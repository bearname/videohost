const actions = {
    async addUser(context, {username}) {
        return await fetch(
            process.env.VUE_APP_USER_SERVER_ADDRESS + "/api/v1/users/" + username,
            {
                headers: {
                    Authorization: context.rootGetters["auth/getTokenHeader"]
                }
            }
        )
            .then(response => {
                if (!response.ok) {
                    if (response.status === 401) {
                        context.dispatch("auth/logout", {}, {root: true});
                    }
                    throw new Error("Cannot get user");
                }
                return response.json();
            })
            .then(data => {
                context.commit("ADD_USER", data);
            })
            .catch(error => {
                console.log(error);
                throw error
            })
    },
    async updateDescription(context, {username, description}) {
        const onfulfilled = async () => {
            return await fetch(
                process.env.VUE_APP_USER_SERVER_ADDRESS + "/api/v1/users/" + username,
                {
                    method: "PATCH",
                    headers: {
                        Authorization: context.rootGetters["auth/getTokenHeader"],
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({
                        username: username,
                        description: description
                    })
                }
            ).then(response => {
                if (!response.ok) {
                    if (response.status === 401) {
                        context.dispatch("auth/logout", {}, {root: true})
                    }
                    throw new Error("Cannot update")
                }
                return response.json()
            }).then(data => {
                this.state.videos = data.videos
                this.state.countAllVideos = data.countAllVideos
                console.log(data)
                context.commit("ADD_USER", data)
            }).catch(error => {
                console.log(error)
                throw error
            })
        };
        context
            .dispatch("auth/updateAuthorizationIfNeeded", {}, {root: true})
            .then(onfulfilled)
    }
}

const getters = {
    getUser(state, userId) {
        if (state.loadedUsers.some(user => user.username === userId)) {
            return state.loadedUsers.find(user => user.username === userId)
        } else {
            //Here I'll have to request from the server!!
            return {}
        }
    },
}
const mutations = {
    ADD_USER(state, user) {
        if (state.loadedUsers.some(u => u.username === user.username)) {
            state.loadedUsers.splice(
                state.loadedUsers.indexOf(u => u.username === user.username),
                1
            )
        }
        state.loadedUsers.push(user)
    }
}
const state1 = {
    loadedUsers: [],
}

export default {
    namespaced: true,
    state: state1,
    mutations: mutations,
    actions: actions,
    getters: getters
}