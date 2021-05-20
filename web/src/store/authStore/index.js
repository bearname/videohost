import jwt_decode from "jwt-decode";
import dbUtils from "../dbUtils";
import router from "../../router/index.js";

const actions = {
    async login(context, {username, password}) {
        return await fetch(process.env.VUE_APP_SERVER_ADDRESS + "/api/v1/auth/login", {
            method: "POST",
            body: JSON.stringify({username: username, password: password})
        }).then(response => {
            if (!response.ok) {
                throw new Error("Cannot login!")
            }
            return response.json()
        }).then(data => {
            dbUtils.addUser({username: username, password: password})
            context.commit("LOGIN", {username: username, accessToken: data.accessToken})
            context.commit("SET_REFRESH_TOKEN", {refreshToken: data.refreshToken})
            context.commit("SET_COOKIE", {
                username: username,
                accessToken: data.accessToken,
                refreshToken: data.refreshToken
            });
            context.commit("SET_REFRESH_TOKEN_COOKIE", {refreshToken: data.refreshToken});
            router.push({name: "uploadVideo"})
        }).catch(error => {
            console.log(error)
            context.dispatch("logout")
        })
    },
    async logout(context) {
        return dbUtils.removeUser({
            username: context.getters.getCurrentUser.username
        }).then(() => {
            console.log("logout")
            context.commit("SET_ACCESS_TOKEN", {accessToken: ""});
            context.commit("LOGOUT")
            context.commit("SET_COOKIE", {username: "", accessToken: ""});
            context.commit("SET_REFRESH_TOKEN", {refreshToken: ""})
            router.push({name: "login"})
        }).catch(error => {
            console.log(error)
        })
    },
    async signup(context, {username, password}) {
        return await fetch(
            process.env.VUE_APP_SERVER_ADDRESS + "/api/v1/auth/create-user",
            {
                method: "POST",
                body: JSON.stringify({username: username, password: password})
            }
        )
            .then(response => {
                if (!response.ok) {
                    throw new Error("Cannot signup!");
                }
                return response.json();
            })
            .then(data => {
                dbUtils.addUser({
                    username: username,
                    token: data.token
                });
                context.commit("LOGIN", {
                    username: username,
                    accessToken: data.accessToken
                });
                context.commit("SET_REFRESH_TOKEN", {refreshToken: data.refreshToken});
                context.commit("SET_COOKIE", {username: username, accessToken: data.accessToken});
                context.commit("SET_REFRESH_TOKEN_COOKIE", {refreshToken: data.refreshToken});

                // router.push({name: "home"})
                router.push({name: "uploadVideo"})

                // context.commit("SET_AUTHORIZATION", data.authorization_token);
            })
            .catch(error => {
                console.log(error)

                context.dispatch("logout");
                throw error;
            });
    },
    async loadUser(context) {
        dbUtils.getUser().then(user => {
            if (user && user !== {}) context.commit("LOGIN", user);
        });
    },
    async updateAuthorizationIfNeeded(context) {
        let expiration = 0;
        const getters1 = context.getters;
        if (getters1.getAccessToken !== "") {
            let data = jwt_decode(getters1.getAccessToken);
            //data.expiration expiration in seconds. Date.now is in milliseconds... So just *1000
            expiration = data.exp * 1000;
        }
        console.log(getters1.getCurrentUser)
        if (getters1.getAccessToken === "" || Date.now() > expiration) {
            const refreshToken = context.getters.getRefreshToken;
            return await fetch(process.env.VUE_APP_SERVER_ADDRESS + "/api/v1/auth/token", {
                headers: {
                    'Authorization': "Bearer " + refreshToken
                }
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error("Cannot get token!");
                    }
                    return response.json();
                })
                .then(data => {
                    context.commit("SET_ACCESS_TOKEN", {accessToken: data.accessToken})
                    context.commit("SET_REFRESH_TOKEN", {refreshToken: data.refreshToken})
                })
                .catch(error => {
                    console.log(error)
                    if (error.message === 'Invalid res')
                        context.dispatch("logout");
                    throw error;
                });
        } else {
            return new Promise(success => {
                success([]);
            });
        }
    }
};

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
        return state.user.refreshToken
    },
    getTokenHeader(state) {
        return "Bearer " + state.user.accessToken;
    },
};

const mutations = {
    LOGIN(state, {username, accessToken}) {
        state.user.loggedIn = true
        state.user.username = username
        state.user.accessToken = accessToken
    },
    LOGOUT(state) {
        state.user.loggedIn = false
        state.user.username = ""
        state.user.accessToken = ""
    },
    SET_ACCESS_TOKEN(state, {accessToken}) {
        console.log("SET_ACCESS_TOKEN")
        state.user.accessToken = accessToken
    },
    SET_REFRESH_TOKEN(state, {refreshToken}) {
        console.log("SET_REFRESH_TOKEN")
        state.user.refreshToken = refreshToken
    },
    SET_COOKIE(state, {username, accessToken}) {
        document.cookie = "username=" + username
        document.cookie = "accessToken=" + accessToken
    },
    SET_REFRESH_TOKEN_COOKIE(state, {refreshToken}) {
        document.cookie = "refreshToken=" + refreshToken
    },
};

const state1 = {
    user: {
        username: "",
        loggedIn: false,
        accessToken: "",
        refreshToken: "",
    }
};

export default {
    namespaced: true,
    state: state1,
    mutations: mutations,
    actions: actions,
    getters: getters
}