(function($ng, $app){
    'use strict';
    $app.members.MembersService = function(){

    };
    $app.members.MembersService.prototype.getMembers = function(){
      return $app.bootstrap.members;
    };

    $ng.module('$app.members').service('members', [$app.members.MembersService]);

}(window.angular, window.$app));