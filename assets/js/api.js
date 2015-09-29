import fetch from 'isomorphic-fetch';

export function getRandomMovie() {
  return fetch("/api/")
    .then(response => response.json());
}

export function getMovies() {
  return fetch("/api/all/")
    .then(response => response.json());
}


export function getMovie(id) {
  return fetch("/api/movie/" + id)
      .then(response => response.json());
}

export function deleteMovie(id) {
  return fetch("/api/movie/" + id, { method: "DELETE" });
}

export function addMovie(title) {
  return fetch("/api/", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({
      title: title
    })
  })
  .then(response => response.json());
}
