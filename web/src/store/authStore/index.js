import actions from "./actions";
import mutations from "./mutaions";
import getters from "./getters";

const state = {
    user: {
        id: "",
        username: "",
        loggedIn: false,
        accessToken: "",
        refreshToken: "",
    }
};

export default {
    namespaced: true,
    state: state,
    mutations: mutations,
    actions: actions,
    getters: getters
}