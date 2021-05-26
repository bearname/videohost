import VueRouter from 'vue-router'
import HomePage from '../views/CatalogPage.vue'
import StreamPage from '../views/VideoPage.vue'
import UploadVideo from "../views/UploadVideoPage";
import LoginPage from "../views/LoginPage";
import UserPage from "../views/UserPage";
import store from "../store/index.js";
import SearchPage from "../views/SearchPage";

const routes = {
    home: {
        name: 'home',
        path: '/catalog',
        component: HomePage
    },
    user: {
        name: 'user',
        path: '/user/',
        component: UserPage,
        beforeEnter: (to, from, next) => {
            if (store.getters["auth/isLoggedIn"]) {
                next();
            } else {
                next({name: "login"});
            }
        }
    },
    uploadVideo: {
        name: 'uploadVideo',
        path: '/uploadVideo',
        component: UploadVideo,
        beforeEnter: (to, from, next) => {
            if (store.getters["auth/isLoggedIn"]) {
                next();
            } else {
                next({name: "login"});
            }
        }
    },
    videoStream: {
        name: 'videoStream',
        path: '/videos/:videoId?',
        component: StreamPage
    },
    search: {
        name: 'search',
        path: '/search/:searchQuery',
        component: SearchPage,
    },
    // editVideo: {
    //     name: 'videoStream',
    //     path: '/videos/:videoId/edit',
    //     component: VideoEditPage,
    //     beforeEnter: (to, from, next) => {
    //         if (store.getters["auth/isLoggedIn"]) {
    //             next();
    //         } else {
    //             next({name: "login"});
    //         }
    //     }
    // },
    login: {
        name: "login",
        path: "/login",
        component: LoginPage,
        beforeEnter: (to, from, next) => {
            if (store.getters["auth/isLoggedIn"]) {
                next({name: "uploadVideo"});
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
        routes.home,
        routes.user,
        routes.videoStream,
        routes.uploadVideo,
        routes.login,
        routes.search,
        // routes.editVideo,
    ],
});

export default router