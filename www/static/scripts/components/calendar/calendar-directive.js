(function($ng, $app){
    var ROWS = 6;
    var d = function(){
        return {
            link: function($scope, elem, attr){

            }
        }
    };

    function buildCalendar(){

    }

    $app.components.calendar =d;
    $app.components.module.directive('calendar', d);
}(window.angular, window.$app));