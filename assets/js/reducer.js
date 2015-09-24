const initialState = {
  movie: null,
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
    case 'GET_MOVIE':
    case 'GET_RANDOM_MOVIE':
      state.movie = action.payload;
      return state;
  }
  return state;
}
