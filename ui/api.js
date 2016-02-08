import axios from 'axios';

axios.interceptors.request.use(config => {
  const headers = Object.assign({}, config.headers, { 'X-CSRF-Token': window.csrfToken });
  return Object.assign({}, config, { headers });
}, (error) => Promise.reject(error));

export function getRandomMovie() {
  return axios.get('/api/');
}

export function getMovies() {
  return axios.get('/api/all/');
}

export function getMovie(id) {
  return axios.get(`/api/movie/${id}`);
}

export function deleteMovie(id) {
  return axios.delete(`/api/movie/${id}`);
}

export function markSeen(id) {
  return axios.patch(`/api/seen/${id}`);
}

export function addMovie(title) {
  return axios.post('/api/', { title });
}
