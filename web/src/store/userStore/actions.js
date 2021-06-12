import makeRequest, {requestWithAuth} from '@/api/api';
import videosUtil from '@/store/videoStore/videoUtil';
import Cookie from '@/util/cookie';

const actions = {
  async addUser(context, {username}) {
    try {
      const response = await fetch(
        process.env.VUE_APP_VIDEO_API + '/api/v1/users/' + username,
        {
          headers: {
            Authorization: context.rootGetters['authMod/getTokenHeader'],
          },
        },
      );
      if (!response.ok) {
        if (response.status !== 401) {
          console.log(response);
          return;
        }
        await context.dispatch('authMod/logout', {}, {root: true});
      }

      const data = response.json();
      context.commit('ADD_USER', data);
    } catch (error) {
      console.log(error);
      throw error;
    }
  },
  async updateDescription(context, {username, description}) {
    try {
      await context.dispatch('authMod/updateAuthorizationIfNeeded', {}, {root: true});
      const config = {
        method: 'PATCH',
        headers: {
          'Authorization': context.rootGetters['authMod/getTokenHeader'],
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: username,
          description: description,
        }),
      };

      const response = await fetch(process.env.VUE_APP_VIDEO_API + '/api/v1/users/' + username, config);
      if (!response.ok) {
        if (response.status !== 401) {
          console.log(response);
          return;
        }
        await context.dispatch('authMod/logout', {}, {root: true});
      }
      const data = response.json();
      this.state.videos = data.videos;
      this.state.countAllVideos = data.countAllVideos;
      console.log(data);
      context.commit('ADD_USER', data);
    } catch (error) {
      console.log(error);
      throw error;
    }
  },
  async getUserVideos(context, {page, countVideoOnPage}) {
    try {
      const cookie = Cookie.getCookie('userId');
      const url = process.env.VUE_APP_USER_API + '/api/v1/users/' + cookie + '/videos?page=' + page +
        '&countVideoOnPage=' + countVideoOnPage;
      console.log(url);

      const config = {};

      const data = await makeRequest(context, url, config);
      console.log('data');
      console.log(data);
      const {videos, countAllVideos} = data;
      context.state.userVideos = videos;
      context.state.userVideos.forEach(videosUtil.updateThumbnail, context.state.userVideos);
      context.state.countUserVideos = countAllVideos;
      console.log('inner end');
    } catch (error) {
      console.log(error);
      context.state.success = false;
      throw error;
    }
  },
  async getUserLikedVideos(context, {page, countVideoOnPage}) {
    try {
      await context.dispatch('authMod/updateAuthorizationIfNeeded', {}, {root: true});

      const url = process.env.VUE_APP_VIDEO_API + '/api/v1/videos-liked?page=' + page + '&countVideoOnPage=' + countVideoOnPage;
      console.log(url);

      const config = {
        method: 'GET',
        headers: {
          'Authorization': context.rootGetters['authMod/getTokenHeader'],
        },
      };

      const data = await makeRequest(context, url, config);
      console.log('data');
      console.log(data);
      context.state.userVideos = data;
      context.state.userVideos.forEach(videosUtil.updateThumbnail, context.state.userVideos);
      console.log('inner end');
    } catch (error) {
      console.log(error);
      context.state.success = false;
      throw error;
    }
  },
  async subscribe(context, {followingToUserId}) {
    try {
      await context.dispatch('authMod/updateAuthorizationIfNeeded', {}, {root: true});

      const url = process.env.VUE_APP_USER_API + `/api/v1/users/${Cookie.getCookie("userId")}/follow?followingToId=${followingToUserId}`;
      console.log(url);

      const data = await requestWithAuth(context, "POST", url, null);
      console.log('data');
      console.log(data);
      context.state.isSuccess = true;

      console.log('inner end');
    } catch (error) {
      console.log(error);
      context.state.isSuccess = false;
      throw error;
    }
  },
};

export default actions;
