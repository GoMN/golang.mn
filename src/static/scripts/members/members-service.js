(function ($ng, $app) {
    'use strict';
    var _$log;
    $app.members.MembersService = function ($log) {
        _$log = $log
    };
    $app.members.MembersService.prototype.getMembers = function () {
        if ($app.bootstrap && $app.bootstrap.members) {
            return $app.bootstrap.members;
        } else {
            _$log.error('bootstrap members not present')
        }
        return []
    };

    $ng.module('$app.members').service('members', ['$log', $app.members.MembersService]);

}(window.angular, window.$app));