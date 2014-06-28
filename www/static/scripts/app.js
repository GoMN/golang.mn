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
              //controller: $app.members.Members
          })
          .state('metrics', {
              url: "/metrics",
              templateUrl: "/static/scripts/metrics/metrics-tmpl.html"
              //controller: $app.metrics.Metrics
          });
    }]);

    $app.module = module;

    $app.App = ['$rootScope', function ($rootScope) {
        $rootScope.$on('$stateChangeStart',
          function (event, toState, toParams, fromState, fromParams) {
              if (!fromState.name) {
                  //cancel state change?
              }
              console.log('state changing');
          });
    }];

}(window.angular, window.$app));