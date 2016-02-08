import React, { PropTypes } from 'react';
import { Link } from 'react-router';

//import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootswatch/paper/bootstrap.min.css';

import {
  Input,
  Button,
  ButtonInput,
  Glyphicon,
  Alert,
  Grid,
  Row,
  Col
} from 'react-bootstrap';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import * as actions from '../actions';

const Suggestion = props => {
    const { movie } = props;
    if (!movie.imdbID) return <span></span>;
    return (
      <small>
        <b><em>Have you seen?</em></b> <Link to={`/movie/${movie.imdbID}/`}>{movie.Title} ({movie.Year}) </Link>
      </small>
    );

};

const Header = props => {
    return (
        <div className="page-header">
          <Grid>
            <Row>
              <Col xs={6} md={6}>
                <h1><Glyphicon glyph="film" /> Movie Wishlist</h1>
              </Col>
              <Col xs={6} md={6} className="text-right">
                <Suggestion movie={props.suggestion} />
              </Col>
            </Row>
        </Grid>
        </div>
    );
};

const Alerts = props => {
  return (
    <div className="container">
      {props.messages.map(msg => {
        const dismissAlert = () => { props.dismissMessage(msg.id); };
        return (
          <Alert key={msg.id}
                 bsStyle={msg.status}
                 onDismiss={dismissAlert}
                 dismissAfter={3000}>
            <p>{msg.message}</p>
          </Alert>
        );
      })}
    </div>
  );
};

class AddMovieForm extends React.Component {

  addMovie(event) {
    event.preventDefault();
    const node = this.refs.title.getInputDOMNode(),
          title = node.value.trim();

    if (title) {
      node.value = "";
      this.props.addMovie(title);
    }
  }

  render() {
    return (
      <div className="container">
        <form className="form form-horizontal" onSubmit={this.addMovie.bind(this)}>
        <Input
          type="text"
          ref="title"
          placeholder="Add another title" />
        <Button
          bsStyle="primary"
          className="form-control"
          type="submit"><Glyphicon glyph="plus" /> Add</Button>
        </form>
      </div>
    );
  }

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
        <Header suggestion={this.props.suggestion} />
        <Alerts dismissMessage={this.actions.dismissMessage}
                messages={this.props.messages} />
        <AddMovieForm addMovie={this.actions.addMovie}  />
        {this.props.children}
      </div>
    );
  }
}

App.propTypes = {
  dispatch: PropTypes.func.isRequired,
  messages: PropTypes.object,
  suggestion: PropTypes.object,
  children: PropTypes.node
}

export default connect(state => {
  return {
      suggestion: state.suggestion,
      messages: state.messages
  };
})(App);

