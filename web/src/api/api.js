async function makeRequest(context, url, config) {
    const response = await fetch(url, config);
    if (!response.ok) {
        if (response.status === 401) {
            await context.dispatch("auth/updateAuthorizationIfNeeded", {}, {root: true})
        } else {
            throw new Error("Cannot update")
        }
    }
    return response.json();
}

export default makeRequest