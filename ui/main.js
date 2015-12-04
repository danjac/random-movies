import Vue from 'vue';
import Resource from 'vue-resource';
import Router from 'vue-router';

import API from './api';
import App from './components/App.vue';
import Movie from './components/Movie.vue';
import MovieList from './components/MovieList.vue';
import Glyphicon from './components/widgets/Glyphicon.vue';

Vue.use(Router);
Vue.use(Resource);
Vue.use(API);

// widgets
//
Vue.component('glyph', Glyphicon);

Vue.http.headers.common['X-CSRF-Token'] = window.csrfToken;

const router = new Router();

router.map({
  '/': {
    component: Movie
  },
  '/movie/:id/': {
    component: Movie
  },
  '/all': {
    component: MovieList
  }
});

router.start(App, '#app');
