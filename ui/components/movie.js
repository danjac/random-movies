import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import _ from 'lodash';

import {
  Button,
  ButtonGroup,
  Glyphicon,
  Badge,
  Row,
  Col,
} from 'react-bootstrap';

import { bindActionCreators } from 'redux';

import { connect } from 'react-redux';

import { Movie } from '../records';
import * as actions from '../actions';


const Stars = props => {
  const { movie } = props;
  const isRated = !(isNaN(movie.imdbRating));

  if (!isRated) {
    return <h3><em>This movie has not been rated yet.</em></h3>;
  }

  const rating = parseFloat(movie.imdbRating);
  const stars = Math.round(rating);

  return (
      <h3>
        {_.range(stars).map(index => <Glyphicon key={index} glyph="star" />)}
        {_.range(10 - stars).map(index => <Glyphicon key={index} glyph="star-empty" />)}
        &nbsp; {rating} <a target="_blank" href={`http://www.imdb.com/title/${movie.imdbID}/`}><small>IMDB</small></a>
      </h3>
  );
};

Stars.propTypes = {
  movie: PropTypes.instanceOf(Movie).isRequired,
};

const Controls = props => {
  const { movie } = props;
  return (
    <ButtonGroup>
      <Link className="btn btn-default" to="/">
        <Glyphicon glyph="list" /> See all
      </Link>
      <Button bsStyle="danger" onClick={props.deleteMovie}>
        <Glyphicon glyph="trash" /> Delete
      </Button>
      {movie.seen ? '' :
      <Button bsStyle="primary" onClick={props.markSeen}>
        <Glyphicon glyph="ok" /> Seen it!
      </Button>}
    </ButtonGroup>
  );
};

Controls.propTypes = {
  deleteMovie: PropTypes.func.isRequired,
  markSeen: PropTypes.func.isRequired,
  movie: PropTypes.instanceOf(Movie).isRequired,
};

export class MovieDetail extends React.Component {

  constructor(props) {
    super(props);
    this.actions = bindActionCreators(actions, this.props.dispatch);
    this.deleteMovie = this.deleteMovie.bind(this);
    this.markSeen = this.markSeen.bind(this);
  }

  componentWillMount() {
    this.getMovie(this.props.params.id);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.params.id !== this.props.params.id) {
      this.getMovie(nextProps.params.id);
    }
  }

  getMovie(id) {
    this.actions.getMovie(id);
  }

  deleteMovie(event) {
    event.preventDefault();
    this.actions.deleteMovie(this.props.movie);
  }

  markSeen(event) {
    event.preventDefault();
    this.actions.markSeen(this.props.movie);
  }

  render() {
    const movie = this.props.movie;
    if (!movie || !movie.imdbID) {
      return <div></div>;
    }

    return (
      <div>
        <h2>{movie.Title}&nbsp;
          {movie.seen ?
            <Badge><Glyphicon glyph="ok" /> Seen it!</Badge> : ''}
          </h2>
          <Row>
          <Col md={3}>
            {movie.Poster === 'N/A' ?
            'No poster available' :
            <img
              className="img-responsive"
              src={`/static/images/${movie.Poster}`}
              alt={movie.Title}
            />}
          </Col>
          <Col md={9}>
            <Stars movie={movie} />
            <dl className="dl-unstyled">
              <dt>Year</dt>
              <dd>{movie.Year}</dd>
              <dt>Actors</dt>
              <dd>{movie.Actors}</dd>
              <dt>Director</dt>
              <dd>{movie.Director}</dd>
            </dl>
            <p className="well">{movie.Plot}</p>
            <Controls
              movie={movie}
              deleteMovie={this.deleteMovie}
              markSeen={this.markSeen}
            />
          </Col>
        </Row>
      </div>
    );
  }

}

MovieDetail.propTypes = {
  params: PropTypes.object.isRequired,
  dispatch: PropTypes.func.isRequired,
  movie: PropTypes.instanceOf(Movie).isRequired,
};

export default connect(state => {
  return {
    movie: state.movie,
  };
})(MovieDetail);
