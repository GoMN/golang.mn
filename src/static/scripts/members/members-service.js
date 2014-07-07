(function ($ng, $app) {
    'use strict';
    $app.members.MembersService = function () {

    };
    $app.members.MembersService.prototype.getMembers = function () {
        if ($app.bootstrap && $app.bootstrap.members) {
            return $app.bootstrap.members;
        } else {
            $log.error('bootstrap members not present')
        }
        return []
    };

    $ng.module('$app.members').service('members', [$app.members.MembersService]);

}(window.angular, window.$app));