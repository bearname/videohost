import videosUtil from "./video";

const mutations = {
    SET_VIDEOS(state, {videos}) {
        state.videos = videos
        state.videos.forEach(videosUtil.updateThumbnail, state.videos)
    },
}
export default mutations