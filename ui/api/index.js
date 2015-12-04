import { Promise } from 'es6-promise';

class API {
  constructor(client) {
    this.client = client;
  }
  addMovie(title) {
    return new Promise((resolve, reject) => {
      this.client.post("/api/", { title: title }, resolve, {
        error: reject
      });
    });
  }

  getMovie(imdbID) {
    return new Promise((resolve, reject) => {
      this.client.get(`/api/movie/${imdbID}`, resolve, reject);
    });
  }

  deleteMovie(imdbID) {
    return new Promise((resolve, reject) => {
      this.client.delete(`/api/movie/${imdbID}`, resolve, reject);
    });
  }

  getRandomMovie() {
    return new Promise((resolve, reject) => {
      this.client.get("/api/", resolve, reject);
    });
  }

  getMovies() {
    return new Promise((resolve, reject) => {
      this.client.get("/api/all/", resolve, reject);
    });
  }

}

export default {
  install(Vue) {
    Vue.prototype.$api = new API(Vue.http);
  }
};
