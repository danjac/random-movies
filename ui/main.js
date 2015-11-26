import Vue from 'vue';
import Router from 'vue-router';

import App from './components/App.vue';
import Movie from './components/Movie.vue';

Vue.use(Router);

const router = new Router();

router.map({
  '/': {
    component: Movie
  }
});

router.start(App, '#app');
