(function ($app) {
    $app.metrics.Metrics = function () {
        function initialize() {
            var mapOptions = {
                center: new google.maps.LatLng(44.9871011, -93.2717069),
                zoom: 9
            };
            var map = new google.maps.Map(document.getElementById("map-canvas"),
              mapOptions);
            addMarkers(map);
        }

        function addMarkers(map) {
            var i, memberCoords = $app.bootstrap.member_coords, l = memberCoords.length;
            for (i = 0; i < l; i += 1) {
                var memberCoord = memberCoords[i];
                var coords = new google.maps.LatLng(memberCoord.lat, memberCoord.lon);
                addMarker(coords, map, memberCoord.title);
            }

        }

        function addMarker(coords, map, title) {
            var marker = new google.maps.Marker({
                position: coords,
                map: map,
                title: title
            });
        }

        initialize()
    };
}(window.$app));