<template>
    <div>
        <h3>Total {{total}} movies</h3>
        <div v-for="row in rows" class="row">
           <div v-for="col in row" class="col-md-3">
                <div v-for="group in col">
                    <h3>{{group.initial}}</h3>
                    <ul class="list-unstyled">
                        <li v-for="movie in group.movies">
                            <a v-link="{ path: '/movie/' + movie.imdbID}">{{movie.Title}}</a>
                        </li>
                    </ul>
                </div>
           </div>
        </div>
    </div>
</template>

<script>
import _ from 'lodash';

function getInitial(title) {
  if (title.match(/^the\s/i)) {
    title = title.substring(4);
  }
  var upCase = title.charAt(0).toUpperCase();
  if (upCase.toLowerCase() !== upCase) { // ASCII letter
    return upCase;
  }
  return '-';
}

function regroup(movies) {
    const groups = _.groupBy(movies, movie => getInitial(movie.Title));
    const cols = _.chunk(_.sortBy(Object.keys(groups)), 4);
    const rows = _.chunk(cols, 4);
    return rows.map(col => {
        return col.map(initials => {
            return initials.map(initial => {
                return {
                    initial: initial,
                    movies: groups[initial]
                };
            });
        });
    });
}

export default {
    name: "MovieList",
    data() {
        return {
            rows: [],
            total: 0
        };
    },
    route: {
        data() {
            return this.$api
            .getMovies()
            .then(data => {
                return {
                    total: data.length,
                    rows: regroup(data)
                };
            });
        }
    }
};
</script>
