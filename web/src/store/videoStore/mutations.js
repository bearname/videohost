import videosUtil from './videoUtil';

const mutations = {
  SET_VIDEOS(state, {videos}) {
    state.videos = videos;
    state.videos.forEach(videosUtil.updateThumbnail, state.videos);
  },
};
export default mutations;
