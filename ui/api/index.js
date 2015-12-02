import { Promise } from 'es6-promise';

export function addMovie(client, title) {
  return new Promise((resolve, reject) => {
    client.post("/api/", { title: title }, resolve, {
      error: reject
    });
  });
}

export function getMovie(client, imdbID) {
  return new Promise((resolve, reject) => {
    client.get(`/api/movie/${imdbID}`, resolve, reject);
  });
}

export function deleteMovie(client, imdbID) {
  return new Promise((resolve, reject) => {
    client.delete(`/api/movie/${imdbID}`, resolve, reject);
  });
}

export function getRandomMovie(client) {
  return new Promise((resolve, reject) => {
    client.get("/api/", resolve, reject);
  });
}


