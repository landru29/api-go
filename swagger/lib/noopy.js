$(function () {

    var baseUrl = window.location.href.replace(/#.*$/, '').replace(/\?.*$/, '');
    var url = baseUrl + 'swagger.json';

    window.disconnect = function () {
        $.removeCookie('token');
        window.location.href = baseUrl;
    };

    // api url
    var matcher = baseUrl.match(/https?:\/\/[^\/]*/);
    var apiUrl = '';
    if (matcher) {
        apiUrl = matcher[0] + '/';
    }

    // get token
    var token;
    var queries = window.location.search.replace(/^\?/, '').split('&');
    var queryToken = queries.filter(function(elt) { return elt.match(/^api-token=/)});
    if (queryToken.length) {
        token = queryToken[0].replace(/^api-token=/, '');
        $.cookie("token", token, { expires : 10 });
        window.location.href = baseUrl;
    }

  hljs.configure({
    highlightSizeThreshold: 5000
  });

  // Pre load translate...
  if(window.SwaggerTranslator) {
    window.SwaggerTranslator.translate();
  }
  window.swaggerUi = new SwaggerUi({
    url: url,
    dom_id: "swagger-ui-container",
    supportedSubmitMethods: ['get', 'post', 'put', 'delete', 'patch'],
    onComplete: function(swaggerApi, swaggerUi){
        $('a#logo').attr('href', baseUrl);
        $('div#header div.swagger-ui-wrap a#facebook_connect').attr('href', apiUrl + 'auth/facebook/login?redirect=' + encodeURIComponent(baseUrl));
        $('div#header div.swagger-ui-wrap a#google_connect').attr('href', apiUrl + 'auth/google/login?redirect=' + encodeURIComponent(baseUrl));
    },
    onFailure: function(data) {
        console.log(data);
    },
    docExpansion: "none",
    jsonEditor: false,
    defaultModelRendering: 'schema',
    showRequestHeaders: false
  });
  window.swaggerUi.load();
  token = $.cookie("token");
  if (token) {
      var apiKeyAuth = new SwaggerClient.ApiKeyAuthorization( "Authorization", "JWT " + token, "header" );
      window.swaggerUi.api.clientAuthorizations.add("bearer", apiKeyAuth);
  }
});
