import React, { PropTypes } from 'react';
import { pushState } from 'redux-router';

import 'bootstrap/dist/css/bootstrap.min.css';

import {
  Input,
  Button,
  ButtonInput,
  Glyphicon,
  Alert
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
      this.actions = bindActionCreators({
        pushState,
        ...actions
      }, dispatch);
  }

  addMovie(event) {
    event.preventDefault();
    const node = this.refs.title.getInputDOMNode(),
          title = node.value.trim();

    if (title) {
      node.value = "";
      const onSuccess = (result) => {
          this.actions.pushState(null, `/movie/${result.data.imdbID}/`);
      }
      this.actions.addMovie(title, onSuccess);
    }
  }

  renderForm() {
    return (
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
    );
  }

  render() {
    return (
      <div className="container">
        {this.props.messages.map((msg, index) => {
        const dismissAlert = (index) => { this.actions.dismissMessage(index); };
        return (
        <Alert key={index} bsStyle={msg.status} onDismiss={dismissAlert} dismissAfter={3000}>
          <p>{msg.msg}</p>
        </Alert>
          );
        })}
        <h1>Random movies</h1>
        {this.renderForm()}
        {this.props.children}
      </div>
    );
  }
}

export default connect(state => {
  return {
      messages: state.main.messages,
      router: state.router
  };
})(App);

