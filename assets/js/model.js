module.exports = {
  Movie: {
    getRandom: function() {
      return m.request({ method: "GET", url: "/api/" });
    },
    getMovie: function(title) {
      return m.request({ method: "GET", url: "/api/?title=" + title });
    },
    addNew: function(title) {
      return m.request({ method: "POST", url: "/api/", data: { title: title } });
    },
    getList: function() {
      return m.request({ method: "GET", url: "/api/titles/" });
    }

  }
};



