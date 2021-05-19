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
            dbUtils.addUser({
                username: username,
                password: password
            })
            context.commit("LOGIN", {
                username: username,
                token: data.token
            })
        }).catch(error => {
            console.log(error)
            context.dispatch("logout")
        })
    },
    async logout(context) {
        return dbUtils.removeUser({
            username: context.getters.currentUser.username
        }).then(() => {
            context.commit("SET_AUTHORIZATION", "");
        })
        // return dbUtils.removeUser({
        //     username: context.getters.currentUser.username
        // }).then(() => {
        //     return fetch(process.env.VUE_APP_SERVER_ADDRESS + "/api/v1/auth/logout", {
        //         method: "POST",
        //         mode: 'no-cors',
        //         headers: {
        //             Authorization:
        //                 "Bearer " + context.getters.getTokenAuthorization
        //         }
        //     }).then(() => {
        //         context.commit("LOGOUT")
        //     }).then(() => {
        //         context.commit("LOGOUT")
        //         router.push({name: "login"})
        //     })
        // })
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
                    token: data.token
                });
                router.push({name: "catalog"})

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
        if (context.getters.getTokenAuthorization !== "") {
            let data = jwt_decode(context.getters.getTokenAuthorization);
            //data.expiration expiration in seconds. Date.now is in milliseconds... So just *1000
            expiration = data.exp * 1000;
        }
        if (context.getters.getTokenAuthorization === "" || Date.now() > expiration) {
            return await fetch(process.env.VUE_APP_SERVER_ADDRESS + "/api/v1/auth/token", {
                headers: {
                    Authorization: "Bearer " + context.getters.getTokenAuthorization
                }
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error("Cannot get token!");
                    }
                    return response.json();
                })
                .then(data => {
                    console.log(data)

                    context.commit("SET_AUTHORIZATION", data.getTokenAuthorization);
                })
                .catch(error => {
                    console.log(error)
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
    currentUser(state) {
        return state.user;
    },
    isLoggedIn(state) {
        if (!state.user) {
            return false;
        }

        return state.user.loggedIn;
    },
    getTokenAuthorization(state) {
        return state.user.token;
    },
    getTokenHeader(state) {
        return "Bearer " + state.user.token;
    },
};

const mutations = {
    LOGIN(state, {username, token}) {
        state.user.loggedIn = true;
        state.user.username = username;
        state.user.token = token;
    },
    LOGOUT(state) {
        state.user.loggedIn = false;
        state.user.username = "";
        state.user.token = "";
    },
    SET_AUTHORIZATION(state, token) {
        state.user.token = token
    }
};

const state1 = {
    user: {
        username: "",
        loggedIn: false,
        token: "",
    }
};

export default {
    namespaced: true,
    state: state1,
    mutations: mutations,
    actions: actions,
    getters: getters
}