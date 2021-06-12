import makeRequest, {requestWithAuth} from '@/api/api';
import Cookie from "@/util/cookie";

export const playlistModification = {
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

      const data = await requestWithAuth(context, 'POST', url, json);
      console.log('data');
      console.log(data);
      context.state.message = data.message;
      context.state.success = data.success;
      console.log('inner end');
    } catch (error) {
      console.log('error');
      console.log(error.data);
      const data = error.response.data;
      context.state.message = data.message
      context.state.success = data.success;
    }
  },
  async modifyPlaylist(context, {action, playlistId, videos}) {
    try {
      await context.dispatch('authMod/updateAuthorizationIfNeeded', {}, {root: true});

      const url = process.env.VUE_APP_VIDEO_API + '/api/v1/playlists/' + playlistId + '/modify';
      console.log(url);

      const json = {
        "act": action,
        "videos": videos
      };

      const data = await requestWithAuth(context, 'PUT', url, json);
      context.state.error = data;
      console.log('inner end');
      context.state.message = data.message
      context.state.success = data.success;
    } catch (error) {
      console.log('error');
      console.log(error);
      const data = error.response.data;
      context.state.message = data.message
      context.state.success = data.success;
    }
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

      context.state.playlists = await makeRequest(context, url, config);
    } catch (err) {
      context.state.error = err;
    }
  },
  async getPlaylist(context, {playlistId}) {
    try {
      await context.dispatch('authMod/updateAuthorizationIfNeeded', {}, {root: true});

      const url = process.env.VUE_APP_VIDEO_API + '/api/v1/playlists/' + playlistId;
      console.log(url);

      const config = {
        method: 'GET',
        headers: {
          'Authorization': context.rootGetters['authMod/getTokenHeader'],
        },
      };

      const data = await makeRequest(context, url, config);
      const playlist = data;
      playlist.videos = JSON.parse(data.videos.replace(/\\/g, ''));
      playlist.videos.forEach(item => {
        item.thumbnail = process.env.VUE_APP_VIDEO_API + '/content/' + item.id + '/screen.jpg';
      })
      console.log(playlist);
      context.state.playlist = playlist;

      console.log(context.state.playlist);
    } catch (err) {
      context.state.error = err;
    }
  }
};

export default actions;