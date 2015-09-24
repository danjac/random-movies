const initialState = {
  movie: null,
  newMovie: null,
  movies: []
};

export default function(state=initialState, action) {
  switch(action.type) {
    case 'RESET_MOVIE':
      state.movie = null;
      return state;
    case 'GET_MOVIES':
      state.movies = action.payload;
      return state;
    case 'ADD_MOVIE':
      state.newMovie = action.payload;
      return state;
    case 'GET_MOVIE':
    case 'GET_RANDOM_MOVIE':
      state.movie = action.payload;
      state.newMovie = null;
      return state;
  }
  return state;
}
