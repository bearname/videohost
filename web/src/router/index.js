import VueRouter from 'vue-router';

import HomePage from '../views/CatalogPage.vue';
import VideoPage from '../views/VideoPage.vue';
import UploadVideo from '../views/UploadVideoPage';
import LoginPage from '../views/LoginPage';
import UserPage from '../views/UserPage';
import SearchPage from '../views/SearchPage';
import Playlist from '../views/Playlist';

function routes(store) {
  return {
    home: {
      name: 'home',
      path: '/catalog',
      component: HomePage,
    },
    user: {
      name: 'user',
      path: '/user/',
      component: UserPage,
      beforeEnter: (to, from, next) => {
        if (store.getters['authMod/isLoggedIn']) {
          next();
        } else {
          next({name: 'login'});
        }
      },
    },
    uploadVideo: {
      name: 'uploadVideo',
      path: '/uploadVideo',
      component: UploadVideo,
      beforeEnter: (to, from, next) => {
        const getter = store.getters['authMod/isLoggedIn'];
        console.log("authmod isloggedin");
        console.log(getter);

        if (getter) {
          next();
        } else {
          next({name: 'login'});
        }
      },
    },
    videoStream: {
      name: 'videoStream',
      path: '/videos/:videoId?',
      component: VideoPage,
    },
    search: {
      name: 'search',
      path: '/search/:searchQuery',
      component: SearchPage,
    },
    playlist: {
      name: 'playlist',
      path: '/playlists',
      component: Playlist,
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
      name: 'login',
      path: '/login',
      component: LoginPage,
      beforeEnter: (to, from, next) => {
        if (store.getters['authMod/isLoggedIn']) {
          next({name: 'uploadVideo'});
        } else {
          next();
        }
      },
    },
  };
}

export function useVueRouter(app) {
  app.use(VueRouter);
}

export {routes};

let routerModule = (function () {
  let instance;
  let createInstance = function (store) {
    return new VueRouter({
      mode: 'history',
      routes: [
        routes(store).home,
        routes(store).user,
        routes(store).videoStream,
        routes(store).uploadVideo,
        routes(store).login,
        routes(store).search,
        routes(store).playlist,
        // routes.editVideo,
      ],
    });
  }

  return {
    getInstance: function (store) {
      return instance || (instance = createInstance(store));
    }
  }
})();

export default routerModule;
// let router
//
// export function createRouter(vuetify, store) {
//   if (router == null && router === undefined) {
//     router = new VueRouter({
//       mode: 'history',
//       routes: [
//         routes(store).home,
//         routes(store).user,
//         routes(store).videoStream,
//         routes(store).uploadVideo,
//         routes(store).login,
//         routes(store).search,
//         routes(store).playlist,
//         // routes.editVideo,
//       ],
//     });
//   }
//   console.log(router);
//   return router;
// }
//
// console.log(router);
//
// export default router;

