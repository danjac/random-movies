import { assert } from 'chai';
import sinon from 'sinon';
import * as api from '../api';
import * as actions from '../actions';


describe('Get random movie', function() {

  beforeEach(function() {
    sinon.spy(api, 'getRandomMovie');
  });

  afterEach(function() {
    sinon.restore(api.getRandomMovie);
  });

  it('Gets a random movie', function() {
    actions.getRandomMovie()();
    assert.ok(api.getRandomMovie.calledOnce);
  });
});
