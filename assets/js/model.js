var m = require('mithril');

module.exports = {
  Movie: {
    getRandom: function() {
      return m.request({ method: "GET", url: "/api/" });
    },
    getMovie: function(id) {
      return m.request({ method: "GET", url: "/api/movie/" + id});
    },
    addNew: function(title) {
      return m.request({ method: "POST", url: "/api/", data: { title: title } });
    },
    getList: function() {
      return m.request({ method: "GET", url: "/api/all/" });
    }

  }
};



