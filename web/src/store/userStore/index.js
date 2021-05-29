import actions from "./actions";
import getters from "../authStore/getters";
import mutations from "../authStore/mutaions";

const state = {
    loadedUsers: [],
};

export default {
    namespaced: true,
    state: state,
    mutations: mutations,
    actions: actions,
    getters: getters
};