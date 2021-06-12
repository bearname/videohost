export async function requestWithAuth(context, method, url, json) {
  const config = {
    method: method,
    headers: {
      'Authorization': context.rootGetters['authMod/getTokenHeader'],
    },
    body: JSON.stringify(json)
  };

  return await makeRequest(context, url, config);
}

export async function makeRequest(context, url, config) {
  const response = await fetch(url, config);
  if (!response.ok) {
    if (response.status === 401) {
      await context.dispatch('authMod/updateAuthorizationIfNeeded', {}, {root: true});
    }
  }
  return response.json();
}

export default makeRequest;
