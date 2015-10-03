import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import _ from 'lodash';

import {
  Input,
  Button,
  ButtonInput,
  ButtonGroup,
  Glyphicon
} from 'react-bootstrap';

import { bindActionCreators } from 'redux';

import { connect } from 'react-redux';
import { pushState } from 'redux-router';

import * as actions from '../actions';

export class Movie extends React.Component {

  static propTypes = {
    dispatch: PropTypes.func.isRequired
  }

  constructor(props) {
    super(props);
    this.actions = bindActionCreators({ pushState, ...actions}, this.props.dispatch);
  }

  deleteMovie() {
    this.actions.deleteMovie(this.props.movie.imdbID);
    this.actions.pushState(null, "/all/");
  }

  render() {
    const movie = this.props.movie;
    if (!movie || !movie.imdbID) {
      return <div></div>;
    }
    const rating = this.props.movie.imdbRating ? parseFloat(this.props.movie.imdbRating) : 0;
    const stars = Math.round(rating);

    return (
      <div className="row">
        <div className="col-md-3">
          {movie.Poster === 'N/A'? 'No poster available' : <img className="img-responsive" src={movie.Poster} alt={movie.Title} />}
        </div>
        <div className="col-md-9">
          <h2>{movie.Title}</h2>
          <h3>
            {_.range(stars).map(index => <Glyphicon key={index} glyph="star" />)}
            {_.range(10 - stars).map(index => <Glyphicon key={index} glyph="star-empty" />)}
            &nbsp; {rating} <a target="_blank" href={`http://www.imdb.com/title/${movie.imdbID}/?ref_=fn_al_tt_1`}><small>IMDB</small></a>
          </h3>
          <dl className="dl-unstyled">
            <dt>Year</dt>
            <dd>{movie.Year}</dd>
            <dt>Actors</dt>
            <dd>{movie.Actors}</dd>
            <dt>Director</dt>
            <dd>{movie.Director}</dd>
          </dl>
          <p className="well">{movie.Plot}</p>
          <ButtonGroup>
            <Button bsStyle="primary" onClick={this.actions.getRandomMovie.bind(this)}><Glyphicon glyph="random" /> Random</Button>
          <Link className="btn btn-default" to="/all/"><Glyphicon glyph="list" /> See all</Link>
          <Button bsStyle="danger" onClick={this.deleteMovie.bind(this)}><Glyphicon glyph="trash" /> Delete</Button>
          </ButtonGroup>
        </div>
      </div>
    );
  }

}

export default connect(state => {
  return {
    movie: state.main.movie,
  };
})(Movie);