import {
  React,
  sinon,
  assert,
  expect,
  TestUtils
} from '../testHelper';

import { Movie } from '../../ui/js/components/movie';

describe('Movie components', () => {
  var sandbox, component;

  const movie = {
    imdbID: "ttf10000",
    title: "Ferris Bueller's Day Off",
    year: 1985,
    rating: "7.3",
    plot: "..."
  };

  beforeEach(() => {
    sandbox = sinon.sandbox.create();
    component = TestUtils.renderIntoDocument(<Movie movie={movie} />);
  });

  afterEach(() => {
    sandbox.restore();
  });

  it('should show a movie title', () => {
    const title = TestUtils.findRenderedDOMComponentWithTag(component, "h3");
    expect(title).to.be.ok
  });

});
