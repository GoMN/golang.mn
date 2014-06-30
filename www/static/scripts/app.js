(function ($ng, $app) {
    var module = $ng.module('$app', [
        'ui.router',
        '$app.components',
        '$app.home',
        '$app.members',
        '$app.metrics'
    ]);

    module.config(['$stateProvider', '$urlRouterProvider', '$locationProvider', function ($stateProvider, $urlRouterProvider, $locationProvider) {
        $locationProvider.html5Mode(true);
        $urlRouterProvider.otherwise("/");
        $stateProvider
          .state('home', {
              url: "/",
              templateUrl: "/static/scripts/home/home-tmpl.html"
          })
          .state('about', {
              url: "/about",
              templateUrl: "/static/scripts/about/about-tmpl.html"
          })
          .state('members', {
              url: "/members",
              templateUrl: "/static/scripts/members/members-tmpl.html"
          })
          .state('metrics', {
              url: "/metrics",
              templateUrl: "/static/scripts/metrics/metrics-tmpl.html"
          });
    }]);

    module.factory('appInterceptor', ['$log', function ($log) {
        return {
            request: function (config) {
                //version cache bust templates
                var url = config.url;
                if (url.indexOf('tmpl.html') > -1) {
                    if (url.indexOf('?') > -1) {
                        url += '&v=' + $app.bootstrap.version;
                    } else {
                        url += '?v=' + $app.bootstrap.version;
                    }
                    config.url = url;
                }
                return config;
            }
        };
    }]);

    module.config(['$httpProvider', function ($httpProvider) {
        $httpProvider.interceptors.push('appInterceptor');
    }]);

    $app.module = module;

    $app.App = ['$rootScope', function ($rootScope) {
        $rootScope.$on('$stateChangeStart',
          function (event, toState, toParams, fromState, fromParams) {
              if (!fromState.name) {
                  //cancel state change?
              }
          });
    }];

}(window.angular, window.$app));