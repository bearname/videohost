const actions = {
    async addUser(context, {username}) {
        try {
            const response = await fetch(
                process.env.VUE_APP_VIDEO_API + "/api/v1/users/" + username,
                {
                    headers: {
                        Authorization: context.rootGetters["auth/getTokenHeader"]
                    }
                }
            );
            if (!response.ok) {
                if (response.status !== 401) {
                    throw new Error("Cannot get user");
                }
                await context.dispatch("auth/logout", {}, {root: true});
            }

            const data = response.json();
            context.commit("ADD_USER", data);
        } catch (error) {
            console.log(error);
            throw error;
        }
    },
    async updateDescription(context, {username, description}) {
        try {
            await context.dispatch("auth/updateAuthorizationIfNeeded", {}, {root: true})
            const config = {
                method: "PATCH",
                headers: {
                    Authorization: context.rootGetters["auth/getTokenHeader"],
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    username: username,
                    description: description
                })
            };

            const response = await fetch(process.env.VUE_APP_VIDEO_API + "/api/v1/users/" + username, config);
            if (!response.ok) {
                if (response.status !== 401) {
                    throw new Error("Cannot update");
                }
                await context.dispatch("auth/logout", {}, {root: true});
            }
            const data = response.json();
            this.state.videos = data.videos;
            this.state.countAllVideos = data.countAllVideos;
            console.log(data);
            context.commit("ADD_USER", data);
        } catch (error) {
            console.log(error);
            throw error;
        }
    }
}

export default actions;