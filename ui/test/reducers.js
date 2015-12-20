import { assert } from 'chai';
import { Alert, Actions } from '../constants';
import { messagesReducer } from '../reducers';

describe('Dismiss a message', function() {

  it('Removes a messaage if ID found', function() {
    const state =  [
      {
        id: 1000,
        status: Alert.INFO,
        message: "testing"
      }
    ];
    const action = {
      type: Actions.DISMISS_MESSAGE,
      payload: 1000
    }
    const newState = messagesReducer(state, action);
    assert.equal(newState.length, 0)

  });

  it('Does nothing if no matching ID', function() {
    const state =  [
      {
        id: 1000,
        status: Alert.INFO,
        message: "testing"
      }
    ];
    const action = {
      type: Actions.DISMISS_MESSAGE,
      payload: 1001
    }
    const newState = messagesReducer(state, action);
    assert.equal(newState.length, 1)

  });
});
