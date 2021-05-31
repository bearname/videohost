import dbUtils from "../dbUtils";
import router from "../../router";
import jwt_decode from "jwt-decode";

const actions = {
    async login(context, {username, password}) {
        try {
            const loginResponse = await fetch(process.env.VUE_APP_USER_API + "/api/v1/auth/login", {
                method: "POST",
                body: JSON.stringify({username: username, password: password})
            });
            if (!loginResponse.ok) {
                throw new Error("Cannot signup!");
            }
            const data = await loginResponse.json();
            const {accessToken, refreshToken} = data
            await dbUtils.addUser({username, password})

            context.commit("LOGIN", {username, accessToken})
            context.commit("SET_REFRESH_TOKEN", {refreshToken})
            context.commit("SET_COOKIE", {
                username,
                accessToken,
                refreshToken
            });
            context.commit("SET_REFRESH_TOKEN_COOKIE", {refreshToken});
            await router.push({name: "uploadVideo"})
        } catch (err) {
            console.log(err)
            await context.dispatch("logout")
        }
    },
    async logout(context) {
        try {
            const user = {
                username: context.getters.getCurrentUser.username
            };

            await dbUtils.removeUser(user)
            console.log("logout")
            context.commit("SET_ACCESS_TOKEN", {accessToken: ""});
            context.commit("LOGOUT")
            context.commit("ERASE_COOKIE", {cookies: ["username", "accessToken"]});
            context.commit("SET_REFRESH_TOKEN", {refreshToken: ""})
            await router.push({name: "login"})
        } catch (error) {
            console.log(error)
        }
    },
    async signup(context, {username, password, email, isSubscribed}) {
        try {
            const config = {
                method: "POST",
                body: JSON.stringify({
                    username: username,
                    password: password,
                    email: email,
                    isSubscribed: isSubscribed
                }),
            };

            const response = await fetch(process.env.VUE_APP_USER_API + "/api/v1/auth/create-user", config);
            if (!response.ok) {
                throw new Error("Cannot signup!");
            }

            const {token, accessToken, refreshToken} = await response.json();
            await dbUtils.addUser({username, token});

            context.commit("LOGIN", {username, accessToken});
            context.commit("SET_REFRESH_TOKEN", {refreshToken: refreshToken});
            context.commit("SET_COOKIE", {username: username, accessToken: accessToken});
            context.commit("SET_REFRESH_TOKEN_COOKIE", {refreshToken: refreshToken});

            await router.push({name: "uploadVideo"})
        } catch (err) {
            console.log(err)

            await context.dispatch("logout");
            throw err;
        }
    },
    async loadUser(context) {
        const user = await dbUtils.getUser();
        if (user && user !== {}) {
            context.commit("LOGIN", user);
        }
    },
    async updateAuthorizationIfNeeded(context) {
        let expiration = 0;
        const getters = context.getters;
        if (getters.getAccessToken !== "") {
            let data = jwt_decode(getters.getAccessToken);
            expiration = data.exp * 1000;
        }

        console.log(getters.getCurrentUser)
        if (getters.getAccessToken === "" || Date.now() > expiration) {
            try {
                const config = {
                    headers: {
                        'Authorization': "Bearer " + context.getters.getRefreshToken
                    }
                };

                const response = await fetch(process.env.VUE_APP_USER_API + "/api/v1/auth/token", config);
                if (!response.ok) {
                    throw new Error("Cannot get token!");
                }

                const {accessToken, refreshToken} = await response.json();
                context.commit("SET_ACCESS_TOKEN", {accessToken})
                context.commit("SET_REFRESH_TOKEN", {refreshToken})

                return response;
            } catch (error) {
                console.log(error)
                if (error.message === 'Invalid res') {
                    await context.dispatch("logout");
                    alert("logout")
                }
                throw error;
            }
        } else {
            return new Promise(success => {
                success([]);
            });
        }
    }
};

export default actions;