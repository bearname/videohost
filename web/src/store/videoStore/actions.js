import axios from "axios";
import VideoStatus from "./videoStatus";
import videos from "./video";
import Cookie from "../../util/cookie";
import logError from "../../util/logger";

const actions = {
    async getVideoById(context, {videoId}) {
        try {
            const url = process.env.VUE_APP_VIDEO_API + '/api/v1/videos/' + videoId;
            const response = await axios.get(url);

            const data = response.data;
            console.log(data);
            context.state.video = data;
            context.state.videoStatus = VideoStatus.intToStatus(data.status);
            context.state.video.uploaded = videos.getElapsedString(data.uploaded);
        } catch (error) {
            logError(error);
        }
    },
    async uploadVideo(context, {file, title, description}) {
        try {
            await context.dispatch("auth/updateAuthorizationIfNeeded", {}, {root: true});
            await this.submitFile(context, file, title, description);
        } catch (error) {
            context.state.isProcessing = false;
            console.error(error);
            console.log('FAILURE!!');
        }
    },
    async submitFile(context, file, title, description) {
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
    },
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
            const url = process.env.VUE_APP_VIDEO_API + "/api/v1/users/" + Cookie.getCookie("username") + "/videos?page=" + page + "&countVideoOnPage=" + countVideoOnPage;
            console.log(url)

            const config = {
                headers: {
                    'Authorization': context.rootGetters["auth/getTokenHeader"]
                }
            };

            const data = await this.makeRequest(context, url, config)
            console.log(data)
            const {videos, countAllVideos} = data;
            context.state.userVideos = videos
            context.state.userVideos.forEach(videos.updateThumbnail, context.state.userVideos)
            context.state.countUserVideos = countAllVideos
        } catch (error) {
            console.log(error)
            throw error
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

            const data = await this.makeRequest(context, url, config)

            console.log(data)
            const {success, message} = data;
            context.state.success = success
            context.state.message = message
        } catch (error) {
            console.log(error)
            context.state.success = false
            throw error
        }
    },
    async makeRequest(context, url, config) {
        const response = await fetch(url, config);
        if (!response.ok) {
            if (response.status === 401) {
                await context.dispatch("auth/updateAuthorizationIfNeeded", {}, {root: true})
            } else {
                throw new Error("Cannot update")
            }
        }
        return response.json();
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
            const data = await this.makeRequest(context, url, config)

            console.log(data)
            const {success, message} = data;

            context.state.success = !success
            context.state.message = message
        } catch (error) {
            console.log(error)
            context.state.success = true
            throw error
        }
    },
    async searchVideos(context, {searchString, page = '1', countVideoOnPage = '10'}) {
        try {
            const videoServerAddress = process.env.VUE_APP_VIDEO_API;
            const url = videoServerAddress + '/api/v1/videos/search?page=' + page + '&limit=' + countVideoOnPage + '&search=' + searchString;
            const response = await axios.get(url);

            let data = response.data;
            console.log(data)

            if (Object.keys(data).includes("pageCount")) {
                this.countPage = data.pageCount
            }
            if (Object.keys(data).includes("videos")) {
                context.state.userVideos = data.videos
                context.state.userVideos.forEach(videos.updateThumbnail, context.state.userVideos)
                context.state.countUserVideos = data.countAllVideos
            }
        } catch (error) {
            logError(error);
        }
    }
};

export default actions;