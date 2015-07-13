(function() {

  function TlsEditor(tlsConfig) {
    this._$mainContent = $('#tls-editor');
    this._$default = generateKeyCert(tlsConfig.default.key, tlsConfig.default.certificate);
    this._$mainContent.append('<h1>Default Key/Cert Pair</h1>', this._$default);
    this._initializeRootCAs(tlsConfig);
    this._initializeNamedCertificates(tlsConfig);
  }

  TlsEditor.prototype._initializeNamedCertificates = function(tlsConfig) {
    var $heading = $('<div class="title"><h1>Certificates</h1><button>Add</button></div>');
    var $certs = $('<div class="named-certificates"></div>');
    var keys = Object.keys(tlsConfig.named).sort();
    for (var i = 0, len = keys.length; i < len; ++i) {
      var name = keys[i];
      var kc = tlsConfig.named[name];
      $certs.append(generateNamedKeyCert(name, kc.key, kc.certificate));
    }
    $heading.find('button').click(function() {
      $certs.append(generateNamedKeyCert('', '', ''));
    }.bind(this));
    this._$certs = $certs;
    this._$mainContent.append($heading, $certs);
  };

  TlsEditor.prototype._initializeRootCAs = function(tlsConfig) {
    var $heading = $('<div class="title"><h1>Root CAs</h1><button>Add</button></div>');
    var $cas = $('<div class="root-cas"></div>');
    for (var i = 0, len = tlsConfig.root_ca.length; i < len; ++i) {
      var ca = tlsConfig.root_ca[i];
      $cas.append($('<textarea class="root-ca"></textarea>').text(ca));
    }
    $heading.find('button').click(function() {
      $cas.append('<textarea class="root-ca"></textarea>');
    }.bind(this));
    this._$cas = $cas;
    this._$mainContent.append($heading, $cas);
  };

  function generateKeyCert(key, cert) {
    var $res = $('<div class="key-cert-pair"><div class="labeled-textarea">' +
      '<label class="labeled-textarea-label">Key</label>' +
      '<textarea class="labeled-textarea-textarea key-value"></textarea>' +
      '</div><div class="labeled-textarea">'+
      '<label class="labeled-textarea-label">Certificate</label>' +
      '<textarea class="labeled-textarea-textarea cert-value"></textarea></div>');
    $res.find('.key-value').text(key);
    $res.find('.cert-value').text(cert);
    return $res;
  }

  function generateNamedKeyCert(name, key, cert) {
    var $res = generateKeyCert(key, cert);
    $res.addClass('named-key-cert-pair');
    $res.prepend('<div class="name-field"><label class="field-label">Name</label>' +
      '<input class="key-cert-name"></div>');
    $res.find('key-cert-name').text(name);
    return $res;
  }

  $(function() {
    new TlsEditor(window.tlsConfiguration);
  });

})();
