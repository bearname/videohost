const DB_NAME = 'authDB';
const DB_VERSION = 1;
let DB;

export default {
  async getDb() {
    return new Promise((resolve) => {
      if (DB) {
        return resolve(DB);
      }
      const request = window.indexedDB.open(DB_NAME, DB_VERSION);

      request.onerror = (e) => {
        console.log(e);
        Promise.reject(new Error('Error'));
      };

      request.onsuccess = (e) => {
        DB = e.target.result;
        resolve(DB);
      };

      request.onupgradeneeded = (e) => {
        const db = e.target.result;
        db.createObjectStore('userAuth', {
          autoIncrement: true,
          keyPath: 'username',
        });
      };
    });
  },
  async addUser(user) {
    const db = await this.getDb();

    return new Promise((resolve) => {
      const trans = db.transaction(['userAuth'], 'readwrite');
      trans.oncomplete = () => {
        resolve();
      };

      const store = trans.objectStore('userAuth');
      store.openCursor().onsuccess = (e) => {
        const cursor = e.target.result;
        if (cursor) {
          if (cursor.value.username !== user.username) {
            store.delete(cursor.value.username);
          }
          cursor.continue();
        }
      };
      store.put(user);
    });
  },
  async removeUser(user) {
    if (!user || !user.username) {
      return;
    }
    const db = await this.getDb();

    return new Promise((resolve) => {
      const trans = db.transaction(['userAuth'], 'readwrite');
      trans.oncomplete = () => {
        resolve();
      };

      const store = trans.objectStore('userAuth');
      store.delete(user.username);
    });
  },
  async getUser() {
    const db = await this.getDb();

    return new Promise((resolve) => {
      const trans = db.transaction(['userAuth'], 'readwrite');
      trans.oncomplete = () => {
        resolve(user);
      };
      let user = null;
      const store = trans.objectStore('userAuth');
      store.openCursor().onsuccess = (e) => {
        const cursor = e.target.result;

        if (!cursor) return;
        user = cursor.value;
      };
    });
  },
};
