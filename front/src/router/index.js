import VueRouter from 'vue-router'
import HomePage from '../views/CatalogPage.vue'
import StreamPage from '../views/VideoPage.vue'
import UploadVideo from "../views/UploadVideoPage";
import LoginPage from "../views/LoginPage";
import store from "../store/index.js";

const routes = {
    home: {
        name: 'home',
        path: '/catalog',
        component: HomePage
    },
    uploadVideo: {
        name: 'uploadVideo',
        path: '/uploadVideo',
        component: UploadVideo
    },
    videoStream: {
        name: 'videoStream',
        path: '/videos/:videoId?',
        component: StreamPage
    },
    login: {
        name: "login",
        path: "/login",
        component: LoginPage,
        beforeEnter: (to, from, next) => {
            if (store.getters["auth/isLoggedIn"]) {
                next({ name: "uploadVideo" });
            } else {
                next();
            }
        }
    }
}

export {routes}
const router = new VueRouter({
    mode: 'history',
    routes: [
        routes.videoStream,
        routes.home,
        routes.uploadVideo,
        routes.login,
    ],
});

export default router