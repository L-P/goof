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

    var initCalendars = function() {
        var formatDate = function(date) {
            var day = date.getUTCDate() + 1;
            if (("" + day).length < 2)
                day = "0" + day;

            var month = date.getUTCMonth() + 1;
            if (("" + month).length < 2)
                month = "0" + month;

            return [date.getUTCFullYear(), month, day].join("-");
        };

        var now = new Date().getTime();
        var lower = new Date();
        var upper = new Date();
        lower.setTime(now - (1000 * 3600 * 24 * 20));
        upper.setTime(now + (1000 * 3600 * 24 * 20));

        var range = formatDate(lower) + "," + formatDate(upper);
        var calendars = {
            "calendar.ics": {}
        };

        for (var name in calendars) {
            var calendar = new document.goof.Calendar(null, {name: name});
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
        nextEventsView.render();
    });
})();
