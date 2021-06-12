import makeRequest, {sendWithAuth} from '@/api/api';
import Cookie from "@/util/cookie";

const action = {
  AddVideo: 0,
  RemoveVideo: 1,
  ReorderVideo: 2
}

const actions = {

  async createPlaylist(context, {name, privacy, videos}) {
    try {
      await context.dispatch('authMod/updateAuthorizationIfNeeded', {}, {root: true});

      const url = process.env.VUE_APP_VIDEO_API + '/api/v1/playlists';
      console.log(url);

      const json = {
        "name": name,
        "privacy": privacy,
        "videos": videos
      };

      const data = await sendWithAuth(context, 'POST', url, json);
      console.log('data');
      console.log(data);
      context.state.error = data;
      console.log('inner end');
    } catch (error) {
      console.log(error);
      context.state.success = false;
      throw error;
    }
  },

  async saveToPlaylist(context, {playlistId, videos}) {
    await context.dispatch('authMod/updateAuthorizationIfNeeded', {}, {root: true});

    const url = process.env.VUE_APP_VIDEO_API + '/api/v1/playlists/' + playlistId + '/modify';
    console.log(url);

    const json = {
      "act": action.AddVideo,
      "videos": videos
    };

    const data = await sendWithAuth(context, 'PUT', url, json);
    console.log('data');
    console.log(data);
    context.state.error = data;
    console.log('inner end');
  },

  async getUserPlaylists(context) {
    try {
      await context.dispatch('authMod/updateAuthorizationIfNeeded', {}, {root: true});
      const userId = Cookie.getCookie("userId");

      const url = process.env.VUE_APP_VIDEO_API + '/api/v1/playlists?ownerId=' + userId;
      console.log(url);

      const config = {
        method: 'GET',
        headers: {
          'Authorization': context.rootGetters['authMod/getTokenHeader'],
        },
      };

      const data = await makeRequest(context, url, config);

      console.log('playlists');
      console.log(data);
      context.state.playlist = data;
      console.log('inner end');
    } catch (err) {
      context.state.error = err;
    }
  },
};

export default actions;