import makeRequest from '@/api/api';

const actions = {

  async createPlaylist(context, {json}) {
    try {
      await context.dispatch('authMod/updateAuthorizationIfNeeded', {}, {root: true});

      const url = process.env.VUE_APP_VIDEO_API + '/api/v1/playlists';
      console.log(url);

      const config = {
        method: 'POST',
        headers: {
          'Authorization': context.rootGetters['authMod/getTokenHeader'],
        },
        body: JSON.stringify(json)
      };

      const data = await makeRequest(context, url, config);
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
};

export default actions;