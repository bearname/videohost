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

                const onSuccess = response => {
                    console.log(response);
                    const status = response.status
                    if (status !== 200) {
                        context.state.isProcessing = false
                    } else {
                        console.log('SUCCESS!!')
                        context.state.isProcessing = true
                        context.state.videoId = response.data
                    }
                };

                const onFail = error => {
                    context.state.isProcessing = false

                    console.error(error)
                    console.log('FAILURE!!')
                };

                const config = {
                    headers: {
                        'Content-Type': 'video/mp4',
                        Authorization: context.rootGetters["auth/getTokenHeader"]
                    }
                }

                return await axios.post("http://localhost:8000/api/v1/video", formData, config).then(onSuccess).catch(onFail)
            })
    },
};

export default {
    namespaced: true,
    state: {
        videoId: null,
        isProcessing: false
    },
    mutations: {
    },
    actions: actions,
    getters: {
        getVideoId(state) {
            return state.videoId
        },
        getIsProcessing(state) {
            return state.isProcessing
        }
    }
}