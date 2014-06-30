(function ($ng, $app) {
    var ROWS = 6;
    var d = function () {
        return {
            restrict: 'E',
            templateUrl: '/static/scripts/components/calendar/calendar-tmpl.html',
            controller: ['$scope', function ($scope) {
                var i, empty = [], rows = [];

                //at most we'll need 6 empties and 6 rows.  create them here to avoid hashkey issues
                for (i = 0; i < 6; i += 1) {
                    empty.push({
                        id: -1 * i,
                        empty: true
                    });
                    rows.push({
                        id: i,
                        days: []
                    });
                }

                $scope.rows = function (month) {
                    var r, days, t, i;
                    days = empty.slice(0, month.startPos);
                    days = days.concat(month.days);
                    t = days.length;
                    rc = Math.ceil(t / 7);
                    r = rows.slice(0, rc);
                    for (i = 0; i < rc; i += 1) {
                        r[i].days = days.slice(i * 7, (i * 7) + 7)
                    }
                    return r;
                };

                $scope.calendar = $app.bootstrap.calendar;
            }]
        }
    };

    $app.components.calendar = d;
    $app.components.module.directive('calendar', d);
}(window.angular, window.$app));