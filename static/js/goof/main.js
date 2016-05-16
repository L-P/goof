"use strict";

document.goof = document.goof || {};

$(function() {
    $('.js-calendar-body').fullCalendar({
        // TODO: fetch dynamically
        eventSources: [
            '/calendar/calendar.ics',
        ],
    });
});
