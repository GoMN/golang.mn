(function ($app) {
    $app.members.Members = ['members', function (members) {
        var self = this;
        this.members = members.getMembers();
    }];
}(window.$app));