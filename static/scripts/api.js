// Generated by CoffeeScript 1.7.1
(function() {
  if (window.goule == null) {
    window.goule = {};
  }

  window.goule.api = function(name, object, callback) {
    var match, path;
    path = window.location.pathname;
    match = /^(.*)index.html$/.exec(path);
    if (match != null) {
      path = match[1];
    }
    $.ajax("" + path + "api/" + name, {
      type: 'POST',
      data: JSON.stringify(object),
      contentType: 'application/json',
      cache: false,
      dataType: 'json',
      error: function() {
        return callback('Error making API call.', null);
      },
      success: function(data) {
        return callback(null, data);
      }
    });
  };

  window.goule.boolApi = function(name, object, callback) {
    return window.goule.api(name, object, function(err, obj) {
      return callback(err == null);
    });
  };

  window.goule.auth = function(password, callback) {
    return window.goule.boolApi('auth', password, callback);
  };

  window.goule.listServices = function(callback) {
    return window.goule.api('services', null, callback);
  };

  window.goule.changePassword = function(newPassword, callback) {
    return window.goule.boolApi('change_password', newPassword, callback);
  };

}).call(this);
