import axios from "axios";
import Cookie from "../../util/cookie";
import videosUtil from "./video";
import logError from "../../util/logger";
import makeRequest from "../../api/api";
import VideoStatus from "./videoStatus";

const actions = {
    async getVideoById(context, {videoId}) {
        try {
            const url = process.env.VUE_APP_VIDEO_API + '/api/v1/videos/' + videoId;
            const response = await axios.get(url);
            if (response.status !== 200) {
                throw new Error("failed get video by id")
            }
            const data = response.data;
            console.log(data);
            context.state.video = data;
            context.state.videoStatus = VideoStatus.intToStatus(data.status);
            context.state.video.uploaded = videosUtil.getElapsedString(data.uploaded);
        } catch (error) {
            logError(error);
        }
    },
    async uploadVideo(context, {file, title, description}) {
        try {
            await context.dispatch("auth/updateAuthorizationIfNeeded", {}, {root: true});
            console.log(title, description)
            const formData = new FormData();
            formData.append("file", file)
            formData.append("title", title)
            formData.append("description", description)

            const config = {
                headers: {
                    'Content-Type': 'video/mp4',
                    'Authorization': context.rootGetters["auth/getTokenHeader"]
                }
            }

            const url = process.env.VUE_APP_VIDEO_API + "/api/v1/videos/";
            console.log('upload video' + url)

            const response = await axios.post(url, formData, config);
            console.log(response);
            const status = response.status
            if (status !== 200) {
                context.state.isProcessing = false
                throw new Error('Failed upload video')
            } else {

                console.log('SUCCESS!!')
                context.state.isProcessing = true
                context.state.videoId = response.data
            }
        } catch (error) {
            context.state.isProcessing = false;
            console.error(error);
            console.log('FAILURE!!');
        }
    },
    // async sendFIle(context, file, title, description) {
    // const formData = new FormData();
    // formData.append("file", file)
    // formData.append("title", title)
    // formData.append("description", description)
    //
    // const config = {
    //     headers: {
    //         'Content-Type': 'video/mp4',
    //         'Authorization': context.rootGetters["auth/getTokenHeader"]
    //     }
    // }
    //
    // const url = process.env.VUE_APP_VIDEO_API + "/api/v1/videos/";
    // console.log('upload video' + url)
    //
    // const response = await axios.post(url, formData, config);
    // console.log(response);
    // const status = response.status
    // if (status !== 200) {
    //     context.state.isProcessing = false
    //     throw new Error('Failed upload video')
    // } else {
    //     console.log('SUCCESS!!')
    //     context.state.isProcessing = true
    //     context.state.videoId = response.data
    // }
    // },
    async getVideoOnPage(context, page = '1', countVideoOnPage = '10') {
        try {
            const url = process.env.VUE_APP_VIDEO_API + '/api/v1/videos/?page=' + page + '&countVideoOnPage=' + countVideoOnPage;
            console.log(url)
            const response = await axios.get(url);
            console.log(response.data)
            const {pageCount, videos} = response.data;

            if (Object.keys(response.data).includes("pageCount")) {
                context.state.pageCount = pageCount
            }
            if (Object.keys(response.data).includes("videos")) {
                context.commit("SET_VIDEOS", {videos: videos});
            }
        } catch (error) {
            logError(error)
        }
    },
    async getUserVideos(context, {page, countVideoOnPage}) {
        try {
            await context.dispatch("auth/updateAuthorizationIfNeeded", {}, {root: true});
            const url = process.env.VUE_APP_USER_API + "/api/v1/users/" + Cookie.getCookie("username") + "/videos?page=" + page + "&countVideoOnPage=" + countVideoOnPage;
            console.log(url)

            const config = {
                headers: {
                    'Authorization': context.rootGetters["auth/getTokenHeader"]
                }
            };

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
    async updateTitleAndDescription(context, {videoId, name, description}) {
        try {
            await context.dispatch("auth/updateAuthorizationIfNeeded", {}, {root: true});
            const url = process.env.VUE_APP_VIDEO_API + "/api/v1/videos/" + videoId;
            const config = {
                method: "PUT",
                headers: {
                    'Authorization': context.rootGetters["auth/getTokenHeader"]
                },
                body: JSON.stringify({"title": name, "description": description})
            };

            const data = await makeRequest(context, url, config);

            console.log(data);
            const {success, message} = data;
            context.state.success = success;
            context.state.message = message;
        } catch (error) {
            console.log(error);
            context.state.success = false;
            throw error;
        }
    },
    async deleteVideoPermanent(context, {videoId}) {
        try {
            await context.dispatch("auth/updateAuthorizationIfNeeded", {}, {root: true});
            const url = process.env.VUE_APP_VIDEO_API + "/api/v1/videos/" + videoId;
            const config = {
                method: "DELETE",
                headers: {
                    'Authorization': context.rootGetters["auth/getTokenHeader"]
                }
            };
            const data = await makeRequest(context, url, config);

            console.log(data);
            const {success, message} = data;

            context.state.success = !success;
            context.state.message = message;
        } catch (error) {
            console.log(error);
            context.state.success = false;
            throw error;
        }
    },
    async likeVideo(context, {videoId, isLike}) {
        try {
            await context.dispatch("auth/updateAuthorizationIfNeeded", {}, {root: true});
            const videoServerAddress = process.env.VUE_APP_VIDEO_API;
            const url = videoServerAddress + `/api/v1/videos/${videoId}/like/${isLike ? 1 : 0}`;
            const config = {
                headers: {
                    'Content-Type': 'video/mp4',
                    'Authorization': context.rootGetters["auth/getTokenHeader"]
                }
            }

            const response = await axios.post(url, null, config);

            console.log('like response');
            console.log(response);
            context.state.success = (response.status === 200);
            context.state.code = response.data.code;
            // console.log('response.data');
            // console.log(data);

            //
            // if (Object.keys(data).includes("pageCount")) {
            //     this.countPage = data.pageCount;
            // }
            // if (Object.keys(data).includes("videos")) {
            //     context.state.videos = data.videos;
            //     context.state.videos.forEach(videosUtil.updateThumbnail, context.state.videos);
            //     context.state.countUserVideos = data.countAllVideos;
            // }
        } catch (error) {
            logError(error);
        }
    },
    async searchVideos(context, {searchString, page = '1', countVideoOnPage = '10'}) {
        try {
            const videoServerAddress = process.env.VUE_APP_VIDEO_API;
            const url = videoServerAddress + '/api/v1/videos/search?page=' + page + '&limit=' + countVideoOnPage + '&search=' + searchString;
            const response = await axios.get(url);

            let data = response.data;
            console.log('response.data');
            console.log(data);

            if (Object.keys(data).includes("pageCount")) {
                this.countPage = data.pageCount;
            }
            if (Object.keys(data).includes("videos")) {
                context.state.videos = data.videos;
                context.state.videos.forEach(videosUtil.updateThumbnail, context.state.videos);
                context.state.countUserVideos = data.countAllVideos;
            }
        } catch (error) {
            logError(error);
        }
    }
};

export default actions;