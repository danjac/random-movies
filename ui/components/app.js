import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import { List } from 'immutable';

// import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootswatch/cyborg/bootstrap.min.css';

import { Glyphicon, Alert, Button } from 'react-bootstrap';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import { Movie } from '../records';
import * as actions from '../actions';
import AddMovieForm from './add';

const Suggestion = props => {
  const { movie } = props;
  if (!movie.imdbID) return <span></span>;
  return (
    <p className="text-center">
      <small>
        <b><em>Have you seen?</em></b>&nbsp;
        <Link to={`/movie/${movie.imdbID}/`}>{movie.Title} ({movie.Year}) </Link>
      </small>
    </p>
  );
};

Suggestion.propTypes = {
  movie: PropTypes.instanceOf(Movie).isRequired,
};


const Header = props => {
  return (
    <div className="page-header text-center">
      <Suggestion movie={props.suggestion} />
      <h1><Link to="/" style={{ color: '#fff' }}>Movie Wishlist</Link></h1>
      <Button className="form-control" bsStyle="primary" onClick={props.getRandomMovie}>
        <Glyphicon glyph="random" /> Random
      </Button>
    </div>
  );
};

Header.propTypes = {
  getRandomMovie: PropTypes.func.isRequired,
  suggestion: PropTypes.instanceOf(Movie).isRequired,
};


const Alerts = props => {
  return (
    <div className="container">
      {props.messages.map(msg => {
        const dismissAlert = () => { props.dismissMessage(msg.id); };
        return (
          <Alert
            key={msg.id}
            bsStyle={msg.status}
            onDismiss={dismissAlert}
            dismissAfter={3000}
          >
            <p>{msg.message}</p>
          </Alert>
        );
      })}
    </div>
  );
};

Alerts.propTypes = {
  messages: PropTypes.instanceOf(List).isRequired,
};


class App extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions, dispatch);
  }

  render() {
    return (
      <div className="container">
        <Header suggestion={this.props.suggestion} getRandomMovie={this.actions.getRandomMovie} />
        <Alerts
          dismissMessage={this.actions.dismissMessage}
          messages={this.props.messages}
        />
        <AddMovieForm addMovie={this.actions.addMovie} />
        {this.props.children}
      </div>
    );
  }
}

App.propTypes = {
  dispatch: PropTypes.func.isRequired,
  messages: PropTypes.instanceOf(List).isRequired,
  suggestion: PropTypes.instanceOf(Movie).isRequired,
  children: PropTypes.node.isRequired,
};

export default connect(state => {
  return {
    suggestion: state.suggestion,
    messages: state.messages,
  };
})(App);
