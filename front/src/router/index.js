import VueRouter from 'vue-router'
import HomePage from '../views/CatalogPage.vue'
import StreamPage from '../views/VideoPage.vue'
import UploadVideo from "../views/UploadVideoPage";

const routes = {
    home: {name: 'home', path: '/catalog', component: HomePage},
    uploadVideo: {name: 'uploadVideo', path: '/uploadVideo', component: UploadVideo},
    videoStream: {name: 'videoStream', path: '/videos/:videoId?', component: StreamPage},
}

export {routes}
export default new VueRouter({
    mode: 'history',
    routes: [
        routes.videoStream,
        routes.home,
        routes.uploadVideo,
    ],
})