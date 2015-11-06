"use strict";

document.goof = document.goof || {};

(function() {
    var loadTemplates = function() {
        var templates = {};
        $('.js-template').each(function(_, tpl) {
            var html = $(tpl).html();
            Mustache.parse(html);
            templates[tpl.getAttribute('data-name')] = html;
        });
        return templates;
    };

    var padLeft = function(str, fill, count) {
        return String(String(fill).repeat(count) + str).slice(-count);
    };

    var initCalendars = function() {
        var formatDate = function(date) {
            var day   = padLeft(date.getUTCDate() + 1, "0", 2);
            var month = padLeft(date.getUTCMonth() + 1, "0", 2);

            return [date.getUTCFullYear(), month, day].join("-");
        };

        var delta = 1000 * 3600 * 24 * 20; // 20 days
        var now = new Date().getTime();
        var lower = new Date();
        var upper = new Date();

        lower.setTime(now - delta);
        upper.setTime(now + delta);

        var range = formatDate(lower) + "," + formatDate(upper);

        // TODO: Stub.
        var calendars = {
            "calendar.ics": {}
        };

        var id = 1;
        for (var name in calendars) {
            var calendar = new document.goof.Calendar(null, {
                name: name,
                paletteId: id++
            });
            calendar.fetch({data: {range: range}});
            calendars[name] = calendar;
        }

        return calendars;
    };

    $(function() {
        var calendars = initCalendars();

        var templates = loadTemplates();
        $('body').append(Mustache.render(templates.navbar));
        $('body').append(Mustache.render(templates.calendar));
        $('.js-calendar-body').append(Mustache.render(
            templates["calendar-multiweek"]
        ));

        var nextEventsView = new document.goof.NextEventsView({
            calendars: calendars,
            template: templates["next-events"]
        });
    });
})();
