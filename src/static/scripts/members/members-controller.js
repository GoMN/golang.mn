(function ($app) {
    'use strict';
    $app.members.Members = ['members', function (members) {
        this.members = members.getMembers();
    }];
}(window.$app));