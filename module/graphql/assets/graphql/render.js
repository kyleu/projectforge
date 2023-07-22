var search = window.location.search;
var parameters = {};
search
  .substr(1)
  .split('&')
  .forEach(function (entry) {
    var eq = entry.indexOf('=');
    if (eq >= 0) {
      const key = decodeURIComponent(entry.slice(0, eq));
      if (key.startsWith("__") || key === "toString") {
        throw "nice try";
      }
      parameters[key] = decodeURIComponent(entry.slice(eq + 1));
    }
  });

function onEditQuery(newQuery) {
  parameters.query = newQuery;
  updateURL();
}

function onEditVariables(newVariables) {
  parameters.variables = newVariables;
  updateURL();
}

function onEditHeaders(newHeaders) {
  parameters.headers = newHeaders;
  updateURL();
}

function onEditOperationName(newOperationName) {
  parameters.operationName = newOperationName;
  updateURL();
}

function onTabChange(tabsState) {
  const activeTab = tabsState.tabs[tabsState.activeTabIndex];
  parameters.query = activeTab.query;
  parameters.variables = activeTab.variables;
  parameters.headers = activeTab.headers;
  parameters.operationName = activeTab.operationName;
  updateURL();
}

function updateURL() {
  var newSearch =
    '?' +
    Object.keys(parameters)
      .filter(function (key) {
        return Boolean(parameters[key]);
      })
      .map(function (key) {
        return (
          encodeURIComponent(key) + '=' + encodeURIComponent(parameters[key])
        );
      })
      .join('&');
  history.replaceState(null, null, newSearch);
}

function getSchemaUrl() {
  let ret = window.location.href;
  if (ret.indexOf("?") > -1) {
    ret = ret.substring(0, window.location.href.lastIndexOf("?"))
  }
  return ret;
}

function getSubscriptionUrl() {
  let ret = window.location.href.substring(0, window.location.href.lastIndexOf("?"));
  return ret.replace("https:", "wss:").replace("http:", "ws:")+"/subscription"
}

(() => {
  function fetchGQL(params) {
    return fetch(getSchemaUrl(), {
      method: "post",
      body: JSON.stringify(params),
      credentials: "include",
      headers: { 'Content-Type': 'application/json' }
    }).then(function (resp) {
      return resp.text();
    }).then(function (body) {
      try { return JSON.parse(body); } catch (error) { return body; }
    });
  }

  let args = {
    fetcher: fetchGQL,
    query: parameters.query,
    variables: parameters.variables,
    headers: parameters.headers,
    operationName: parameters.operationName,
    onEditQuery: onEditQuery,
    onEditVariables: onEditVariables,
    onEditHeaders: onEditHeaders,
    defaultSecondaryEditorOpen: true,
    onEditOperationName: onEditOperationName,
    headerEditorEnabled: true,
    shouldPersistHeaders: true,
    tabs: {onTabChange: onTabChange}
  }
  ReactDOM.render(
    React.createElement(GraphiQL, args),
    document.getElementById('graphiql'),
  );
})();
