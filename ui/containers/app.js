import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import 'bootstrap/dist/css/bootstrap.min.css';

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

class App extends React.Component {

  static propTypes = {
    dispatch: PropTypes.func.isRequired,
    children: PropTypes.node
  }

  constructor(props, context) {
    super(props, context);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions, dispatch);
  }

  addMovie(event) {
    event.preventDefault();
    const node = this.refs.title.getInputDOMNode(),
          title = node.value.trim();

    if (title) {
      node.value = "";
      this.actions.addMovie(title);
    }
  }

  renderForm() {
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

  renderHeader() {
    return (
        <div className="page-header">
          <Grid>
            <Row>
              <Col xs={6} md={6}>
                <h1><Glyphicon glyph="film" /> Movie Wishlist</h1>
              </Col>
              <Col xs={6} md={6} className="text-right">
              {this.renderSuggestion()}
              </Col>
            </Row>
        </Grid>
        </div>
    );
  }

  renderSuggestion() {
    const movie = this.props.suggestion;
    if (!movie) return '';
    return (
      <small>
        <b><em>Have you seen?</em></b> <Link to={`/movie/${movie.imdbID}/`}>{movie.Title} ({movie.Year}) </Link>
      </small>
    );
  }

  renderAlerts() {
    return (
      <div className="container">
        {this.props.messages.map(msg => {
          const dismissAlert = () => { this.actions.dismissMessage(msg.id); };
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
  }

  render() {
    return (
      <div className="container">
        {this.renderHeader()}
        {this.renderAlerts()}
        {this.renderForm()}
        {this.props.children}
      </div>
    );
  }
}

export default connect(state => {
  return {
      suggestion: state.suggestion,
      messages: state.messages,
      router: state.router
  };
})(App);

