(function() {
    "use strict";
    var calendars = {};
    document.goof = {calendars: calendars};

    $.ajax('/calendar/calendar.ics').done(function(data) {
        calendars["calendar.ics"] = data.Data.Calendars[0];
    });
})();
