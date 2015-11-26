<template>
  <div class="row">

    <div class="col-md-3">
        <img class="img-responsive" 
             v-bind:src="movie.Poster" 
             v-bind:alt="movie.Title" />
    </div>

    <div class="col-md-9">
        <h2>{{movie.Title}}</h2> 
        <h3 v-if="rating">
            <span v-for="star in stars" class="glyphicon glyphicon-star"></span>
            <span v-for="star in emptyStars" class="glyphicon glyphicon-star-empty"></span>
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
            <button v-on:click="getRandom" class="btn btn-primary">
                <span class="glyphicon glyphicon-random"></span> Random
            </button>
            <a v-link="{ path: '/all' }" class="btn btn-default">
                <span class="glyphicon glyphicon-list"></span> See all
            </a>
            <button v-on:click="deleteMovie" class="btn btn-danger">
                <span class="glyphicon glyphicon-trash"></span> Delete
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
            movie: {}
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
                return this.$http.get(`/api/movie/${id}`, movie => {
                    return {
                        movie
                    };
                });
            } else {
                return this.$http.get("/api/", movie => {
                    return {
                        movie
                    };
                });
            }
        }
    },
    methods: {
        getRandom(event) {
            this.$http.get("/api/", movie => { this.movie = movie; });
        },
        deleteMovie(event) {
            this.$http.delete(`/api/movie/${this.movie.imdbID}`, () => {
                store.createAlert('Movie has been deleted', 'info');
                this.$route.router.go("/all");
            });
        }
    }
};
    
</script>

