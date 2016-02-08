import { Record } from 'immutable';

export const Message = Record({
  status: '',
  message: '',
  id: 0,
});

export const Movie = Record({
  Title: '',
  Actors: '',
  Poster: '',
  Year: '',
  Plot: '',
  Director: '',
  imdbID: '',
  imdbRating: '',
  seen: false,
});



