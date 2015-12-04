<template>
  <div class="row">

    <div class="col-md-3" v-if="loaded">
        <span v-if="movie.Poster === 'N/A'">No poster available</span>
        <img v-else 
             class="img-responsive" 
             :src="movie.Poster" 
             :alt="movie.Title" />
    </div>

    <div class="col-md-9" v-if="loaded">
        <h2>{{movie.Title}}</h2> 
        <h3 v-if="rating">
            <glyph v-for="star in stars" icon="star"></glyph>
            <glyph v-for="star in emptyStars" icon="star-empty"></glyph>
            &nbsp; {{rating}} <a target="_blank" href="http://www.imdb.com/title/{{movie.imdbID}}"><small>IMDB</small></a>
        </h3>
        <dl class="dl-unstyled">
            <dt>Year</dt>
            <dd>{{movie.Year}}</dd>
            <dt>Actors</dt>
            <dd>{{movie.Actors}}</dd>
            <dt>Director</dt>
            <dd>{{movie.Director}}</dd>
          </dl>
          <p class="well">{{movie.Plot}}</p>
          <div class="button-group">
            <button @click="getRandom" class="btn btn-primary">
                <glyph icon="random"></glyph>&nbsp; Random
            </button>
            <a v-link="{ path: '/all' }" class="btn btn-default">
                <glyph icon="list"></glyph>&nbsp; See all
            </a>
            <button @click="deleteMovie" class="btn btn-danger">
                <glyph icon="trash"></glyph>&nbsp; Delete
            </button>
          </div>
    </div>
  </div>

</template>

<script>
import _ from 'lodash';
import store from '../store';

export default {
    
    name: "Movie",
    data() {
        return {
            movie: {},
            loaded: false
        };
    },
    computed: {
        rating() {
            if (isNaN(this.movie.imdbRating)) {
                return 0;
            }
            return parseFloat(this.movie.imdbRating);
        },
        stars() {
            return _.range(this.rating);
        },
        emptyStars() {
            return _.range(10 - this.rating);
        }

    },
    route: {
        data({ to }) {
            const id = to.params.id;
            if (id) {
                return this.$api.getMovie(id)
                .then(movie => {
                    return {
                        movie,
                        loaded: true
                    };
                });
            } else {
                return this.$api.getRandomMovie()
                .then(movie => {
                    return {
                        movie,
                        loaded: true
                    };
                });
            }
        },
        deactivate() {
            this.movie = {};
            this.loaded = false;
        }
    },
    methods: {
        getRandom(event) {
            this.$api
            .getRandomMovie()
            .then(movie => this.movie = movie);
        },
        deleteMovie(event) {
            this.$api
            .deleteMovie(this.movie.imdbID)
            .then(() => {
                store.createAlert('Movie has been deleted', 'info');
                this.$route.router.go("/all");
            });
        }
    }
};
    
</script>

