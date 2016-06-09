"use strict";

document.goof = document.goof || {};

$(function() {
    function transformEvent(event) {
        return {
            title: event.Summary,
            start: event.Start,
            end: event.End,
            raw: event
        }
    }

    function getCalendarEventSource(calendarUrl) {
        return function(start, end, timezone, callback) {
            $.ajax({
                url: calendarUrl,
                data: {
                    start: start.format("YYYY-MM-DD"),
                    end: end.format("YYYY-MM-DD"),
                },
                success: function(data) {
                    callback(_.map(data.Data.Calendar.Events, transformEvent));
                }
            });
        };
    }

    $('.js-calendar-body').fullCalendar({
        eventSources: getCalendarEventSource('/calendar/calendar.ics')
    });
});
