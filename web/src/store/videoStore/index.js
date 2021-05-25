import axios from "axios";
import Cookie from "../../util/cookie.js";
import VideoStatus from "./videoStatus";
import getElapsedString from "./video";

const actions = {
    getVideoById(context, {videoId}) {
        let url = process.env.VUE_APP_VIDEO_SERVER_ADDRESS + '/api/v1/videos/' + videoId;
        axios.get(url)
            .then(response => {
                console.log(response.data)
                this.video = response.data
                this.videoStatus = VideoStatus.intToStatus(this.video.status)
                this.video.uploaded = getElapsedString(this.video.uploaded)
            })
            .catch(function (error) {
                if (error.response) {
                    console.log(error.response.data);
                    console.log(error.response.status);
                    console.log(error.response.headers);
                } else if (error.request) {
                    console.log(error.request);
                } else {
                    console.log('Error', error.message);
                }
                this.error = true
            });
    },
    async uploadVideo(context, {file, title, description}) {
        const promise = context
            .dispatch("auth/updateAuthorizationIfNeeded", {}, {root: true});
        promise
            .then(async () => {
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
                        'Authorization': context.rootGetters["auth/getTokenHeader"]
                    }
                }

                const url = process.env.VUE_APP_VIDEO_SERVER_ADDRESS + "/api/v1/videos/";
                console.log('upload video' + url)
                return await axios.post(url, formData, config)
                    .then(onSuccess)
                    .catch(onFail)
            })

        return promise
    },
    async getVideoOnPage(context, page = '1', countVideoOnPage = '10') {
        let url = process.env.VUE_APP_VIDEO_SERVER_ADDRESS + '/api/v1/videos/?page=' + page + '&countVideoOnPage=' + countVideoOnPage;
        console.log(url)
        return await axios.get(url)
            .then(response => {
                console.log(response.data)
                if (Object.keys(response.data).includes("pageCount")) {
                    context.state.pageCount = response.data.pageCount
                }
                if (Object.keys(response.data).includes("videos")) {
                    context.commit("SET_VIDEOS", {videos: response.data.videos});
                }
            })
            .catch(function (error) {
                if (error.response) {
                    console.log(error.response.data);
                    console.log(error.response.status);
                    console.log(error.response.headers);
                } else if (error.request) {
                    console.log(error.request);
                } else {
                    console.log('Error', error.message);
                }
                context.state.error = true
            });
    },
    getUserVideos(context, {page, countVideoOnPage}) {
        const promise = context.dispatch("auth/updateAuthorizationIfNeeded", {}, {root: true});

        promise.then(() => {
            const url = process.env.VUE_APP_USER_SERVER_ADDRESS + "/api/v1/users/" + Cookie.getCookie("username") + "/videos?page=" + page + "&countVideoOnPage=" + countVideoOnPage;
            console.log(url)
            const config = {
                headers: {
                    'Authorization': context.rootGetters["auth/getTokenHeader"]
                }
            };
            return fetch(url, config).then(response => {
                if (!response.ok) {
                    if (response.status === 401) {
                        context.dispatch("auth/updateAuthorizationIfNeeded", {}, {root: true})
                    } else {
                        throw new Error("Cannot update")
                    }
                }
                return response.json()
            }).then(data => {
                console.log(data)
                context.state.userVideos = data.videos
                context.state.userVideos.forEach(updateThumbnail, context.state.userVideos)
                context.state.countUserVideos = data.countAllVideos
            }).catch(error => {
                console.log(error)
                throw error
            })
        })
        return promise
    },
    updateTitleAndDescription(context, {videoId, name, description}) {
        const promise = context.dispatch("auth/updateAuthorizationIfNeeded", {}, {root: true});

        promise.then(() => {
            const url = process.env.VUE_APP_VIDEO_SERVER_ADDRESS + "/api/v1/videos/" + videoId;
            return fetch(url, {
                method: "PUT",
                headers: {
                    'Authorization': context.rootGetters["auth/getTokenHeader"]
                },
                body: JSON.stringify({"title": name, "description": description})
            }).then(response => {
                if (!response.ok) {
                    if (response.status === 401) {
                        context.dispatch("auth/updateAuthorizationIfNeeded", {}, {root: true})
                    } else {
                        throw new Error("Cannot update")
                    }
                }
                return response.json()
            }).then(data => {
                context.state.success = data.success
                context.state.message = data.message
                console.log(data)
            }).catch(error => {
                console.log(error)
                context.state.success = false
                throw error
            })
        })
    },
    deleteVideoPermanent(context, {videoId}) {
        const promise = context.dispatch("auth/updateAuthorizationIfNeeded", {}, {root: true});

        promise.then(() => {
            const url = process.env.VUE_APP_VIDEO_SERVER_ADDRESS + "/api/v1/videos/" + videoId;
            return fetch(url, {
                method: "DELETE",
                headers: {
                    'Authorization': context.rootGetters["auth/getTokenHeader"]
                }
            }).then(response => {
                if (!response.ok) {
                    if (response.status === 401) {
                        context.dispatch("auth/updateAuthorizationIfNeeded", {}, {root: true})
                    } else {
                        throw new Error("Cannot delete")
                    }
                }
                return response.json()
            }).then(data => {
                context.state.success = !data.success
                context.state.message = data.message
                console.log(data)
            }).catch(error => {
                console.log(error)
                context.state.success = true
                throw error
            })
        })
        return promise
    }
};

const updateThumbnail = function (part, index) {
    this[index].thumbnail = process.env.VUE_APP_VIDEO_SERVER_ADDRESS + "/" + this[index].thumbnail
};

export default {
    namespaced: true,
    state: {
        videoId: null,
        isProcessing: false,
        success: false,
        error: false,
        videos: null,
        userVideos: null,
        pageCount: 0,
        countUserVideos: 0,
        url: {
            type: String,
            required: false,
            default: "list"
        },
    },
    mutations: {
        SET_VIDEOS(state, {videos}) {
            state.videos = videos
            state.videos.forEach(updateThumbnail, state.videos)
        },
    },
    actions: actions,
    getters: {
        getVideoId(state) {
            return state.videoId
        },
        getIsProcessing(state) {
            return state.isProcessing
        },
        getVideos(state) {
            return state.videos
        },
        getUserVideos(state) {
            return {
                videos: state.userVideos,
                countAllVideos: state.countUserVideos
            }
        },
        getPageCount(state) {
            return state.pageCount
        },
        getStatus(state) {
            return state.success
        }
    }
}