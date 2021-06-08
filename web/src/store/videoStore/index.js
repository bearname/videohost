import actions from "./actions";
import mutations from "./mutations";

const state = {
    videoId: null,
    isProcessing: false,
    success: false,
    code: false,
    error: false,
    video: null,
    videos: null,
    userVideos: null,
    pageCount: 0,
    countUserVideos: 0,
    url: {
        type: String,
        required: false,
        default: "list"
    },
};

const getters = {
    getVideoId(state) {
        return state.videoId;
    },
    getIsProcessing(state) {
        return state.isProcessing;
    },
    getVideos(state) {
        return state.videos;
    },
    getVideo(state) {
        return state.video;
    },
    getUserVideos(state) {
        return {
            videos: state.userVideos,
            countAllVideos: state.countUserVideos
        };
    },
    getPageCount(state) {
        return state.pageCount
    },
    getStatus(state) {
        return state.success
    },
    getCode(state) {
        return state.code
    }
};

export default {
    namespaced: true,
    state: state,
    mutations: mutations,
    actions: actions,
    getters: getters
}