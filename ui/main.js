import Vue from 'vue';
import Resource from 'vue-resource';
import Router from 'vue-router';

import App from './components/App.vue';
import Movie from './components/Movie.vue';

Vue.use(Router);
Vue.use(Resource);

const router = new Router();

router.map({
  '/': {
    component: Movie
  }
});

router.start(App, '#app');
