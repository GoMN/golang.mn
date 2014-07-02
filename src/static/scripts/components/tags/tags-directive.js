(function ($ng, $app) {
    'use strict';
    var d = function () {
        return {
            restrict: 'E',
            templateUrl: '/static/scripts/components/tags/tags-tmpl.html',
            controller: ['$scope', function ($scope) {

                $scope.tags = build($app.bootstrap.topics);

                function build(tags) {
                    tags = _.groupBy(tags, function (t) {
                        return t.name
                    });
                    var t, cloud = [], cent = new Array(tags.length);
                    for (t in tags) {
                        if (tags.hasOwnProperty(t)) {
                            cloud.push({
                                count: tags[t].length,
                                name: t
                            });
                        }
                    }
                    cloud.sort(function (a, b) {
                        return a.count - b.count;
                    });

                    var i = 0, l = cloud.length, h = Math.ceil(l / 2), y = (l - 1);
                    var cidx = 0;
                    for (i; cidx < h; i += 2) {
                        cent[cidx] = cloud[i];
                        if (cidx !== (y - cidx)) {
                            cent[y - cidx] = cloud[i + 1];
                        }
                        cidx += 1
                    }
                    return cent;
                }
            }]
        };
    };

    $app.components.calendar = d;
    $app.components.module.directive('tags', d);
}(window.angular, window.$app));