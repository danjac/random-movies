<template>
  <div class="row">

    <div class="col-md-3">
        <img class="img-responsive" 
             v-bind:src="movie.Poster" 
             v-bind:alt="movie.Title" />
    </div>

    <div class="col-md-9">
        <h2>{{movie.Title}}</h2> 
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
          </div>
    </div>
  </div>

</template>

<script>
export default {
    
    name: "Movie",
    data() {
        return {
            movie: {}
        };
    },
    route: {
        data({ to }) {
            const id = to.params.id;
            if (id) {
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
            console.log("getrandom!!!")
            this.$http.get("/api/", movie => { this.movie = movie; });
        }
    }
};
    
</script>

