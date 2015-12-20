import _ from 'lodash';
import React from 'react';
import TestUtils from 'react-addons-test-utils';
import jsdom from 'mocha-jsdom';
import { assert } from 'chai';


//import mockDom from '../mock-dom';
import { Movie } from '../../components/movie';

describe('Movie components', () => {

  jsdom({ skipWindowCheck: true });

  const movie = {
    imdbID: "ttf10000",
    Title: "Ferris Bueller's Day Off",
    Year: 1985,
    Rating: "7.3",
    Plot: "..."
  };

  beforeEach(() => {
  });

  afterEach(() => {
  });

  it('should show a movie title', () => {
    const component = <Movie movie={movie} dispatch={_.noop} />;
    const rendered = TestUtils.renderIntoDocument(component, 'div');
    const header = TestUtils.findRenderedDOMComponentWithTag(rendered, "h2");
    assert.include(header.textContent, "Ferris Bueller's Day Off");
  });

});
