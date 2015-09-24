import React, { PropTypes } from 'react';
import _ from 'lodash';
import { Link } from 'react-router';
import { bindActionCreators } from 'redux';

import { connect } from 'react-redux';

import * as actions from '../actions';


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

@connect(state => {
  let movies = state.movies;
  movies.sort((left, right) => left.Title > right.Title ? 1 : (left.Title < right.Title ? -1 : 0));
  return {
      movies: movies
  };
})
export default class MovieList extends React.Component {

  static propTypes = {
    dispatch: PropTypes.func.isRequired
  }

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions, dispatch);
  }

  componentDidMount() {
    this.actions.getMovies();
  }

  render() {
    const movies = this.props.movies;
    const groups = _.groupBy(movies, movie => getInitial(movie.Title));
    const cols = _.chunk(_.sortBy(Object.keys(groups)), 4);
    const rows = _.chunk(cols, 4);
    return (
      <div>
        <h3>Total: {movies.length}</h3>
        {rows.map((row, index) => {
          return (
            <div key={index} className="row">
              {row.map((col, index) => {
                return (
                  <div key={index} className="col-md-3">
                    {col.map((initial) => {
                      return (
                      <div key={initial}>
                        <h3>{initial}</h3>
                        <ul className="list-unstyled">
                          {groups[initial].map(movie => {
                          return (
                            <li key={movie.imdbID}>
                              <Link to={`/movie/${movie.imdbID}/`}>{movie.Title}</Link>
                            </li>
                            );
                          })}
                        </ul>
                      </div>
                      );
                    })}
                  </div>
                );
              })}
            </div>
          );
        })}
      </div>
    );
  }


}

