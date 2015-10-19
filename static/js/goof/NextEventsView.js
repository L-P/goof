"use strict";

document.goof = document.goof || {};

document.goof.NextEventsView = Backbone.View.extend({
    className: "js-next-events",
    calendars: null,
    template: null,

    initialize: function(options) {
        this.calendars = options.calendars;
        this.template = options.template;

        for(var name in this.calendars) {
            this.listenTo(this.calendars[name], "change", this.render);
        }
    },

    render: function() {
        var events = [];
        for(var name in this.calendars) {
            events = events.concat(this.calendars[name].models);
        }

        $(".js-next-events").html(Mustache.render(
            this.template,
            {events: events}
        ));

        return this;
    }
});
