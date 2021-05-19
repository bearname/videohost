// import dbUtils from "../dbUtils";
import axios from "axios";

const actions = {
    async uploadVideo(context, {file, title, description}) {
        console.log("Upload file")
        context
            .dispatch("auth/updateAuthorizationIfNeeded", {}, {root: true})
            .then(async () => {
                console.log("Upload file")
                const formData = new FormData();
                formData.append("file", file)
                formData.append("title", title)
                formData.append("description", description)
                // const config = {
                //     method: 'POST',
                //     headers: {
                //         // 'Content-Type': 'application/x-www-form-urlencoded',
                //         'Content-Type': 'video/mp4',
                //         Authorization: context.rootGetters["auth/getTokenHeader"]
                //     },
                //     body: formData
                // }
                const onSuccess = response => {
                    console.log(response);
                    const status = response.status
                    if (status !== 200) {
                        this.processing = false
                    } else {
                        console.log('SUCCESS!!')
                        this.processing = true
                        this.videoId = response.data
                    }
                };
                const onFail = error => {
                    this.processing = false
                    console.error(error)
                    console.log('FAILURE!!')
                };

                // return await fetch("http://localhost:8000/api/v1/video", config).then(onSuccess).catch(onFail)
                const config = {
                    headers: {
                        // 'Content-Type': 'application/x-www-form-urlencoded',
                        'Content-Type': 'video/mp4',
                        Authorization: context.rootGetters["auth/getTokenHeader"]
                    }
                }

                //
                return await axios.post("http://localhost:8000/api/v1/video", formData, config).then(onSuccess).catch(onFail)

                // return fetch(process.env.VUE_APP_SERVER_ADDRESS + "/api/v1/video", {
                //     method: "POST",
                //     headers: {
                //         "Content-Type": "application/json",
                //         Authorization: context.rootGetters["auth/getTokenHeader"]
                //     },
                //     body: JSON.stringify({username: username, password: password})
                // }).then(response => {
                //     if (!response.ok) {
                //         throw new Error("Cannot login!")
                //     }
                //     return response.json()
                // }).then(data => {
                //     // dbUtils.addUser({
                //     //     username: username,
                //     //     password: password
                //     // })
                //     // context.commit("LOGIN", {
                //     //     username: username,
                //     //     token: data.token
                //     // })
                // }).catch(error => {
                //     // console.log(error)
                //     // context.dispatch("logout")
                // })
            });
    },
};

export default {
    namespaced: true,
    state: {},
    mutations: {},
    actions: actions,
    getters: {}
}