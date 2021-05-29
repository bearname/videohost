import jwt_decode from "jwt-decode";
import Cookie from "../../util/cookie";

const mutations = {
    LOGIN(state, {username, accessToken}) {
        state.user.loggedIn = true
        state.user.username = username
        state.user.accessToken = accessToken
        state.user.id = jwt_decode(accessToken).userId
    },
    LOGOUT(state) {
        state.user.loggedIn = false
        state.user.id = ""
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
        if (accessToken !== "") {
            const jwtDecode = jwt_decode(accessToken);
            Cookie.setCookie("userId", jwtDecode.userId, 2)
        }
        Cookie.setCookie("username", username, 2)
        Cookie.setCookie("accessToken", accessToken, 2)
    },
    SET_REFRESH_TOKEN_COOKIE(state, {refreshToken}) {
        Cookie.setCookie("refreshToken", refreshToken, 1100)
    },
    ERASE_COOKIE(state, {cookies}) {
        for (let cookie in cookies) {
            Cookie.eraseCookie(cookie)
        }
    }
};

export default mutations;