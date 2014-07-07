(function ($app, google) {
    'use strict';
    $app.metrics.Metrics = ['$log', function ($log) {
        var disabled = true;

        function initialize() {
            var mapOptions = {
                center: new google.maps.LatLng(44.9871011, -93.2717069),
                zoom: 9,
                scrollwheel: false
            };

            var map = new google.maps.Map(document.getElementById('map-canvas'),
              mapOptions);

            if($app.bootstrap && $app.bootstrap.memberCoords) {
                addMarkers(map);
            }else{
                $log.error('bootstrap memberCoords not present')
            }
        }

        function addMarkers(map) {
            var i, memberCoords = $app.bootstrap.memberCoords, l = memberCoords.length;
            for (i = 0; i < l; i += 1) {
                var memberCoord = memberCoords[i];
                var coords = new google.maps.LatLng(memberCoord.lat, memberCoord.lon);
                addMarker(coords, map, memberCoord.title);
            }

        }

        function addMarker(coords, map, title) {
            new google.maps.Marker({
                position: coords,
                map: map,
                title: title
            });
        }

        initialize();
    }];
}(window.$app, window.google));