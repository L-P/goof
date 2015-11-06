"use strict";

document.goof = document.goof || {};

document.goof.NextEventsView = Backbone.View.extend({
    className: "js-next-events",
    calendars: null,
    template: null,

    initialize: function(options) {
        this.calendars = options.calendars;
        this.template = options.template;

        _.each(this.calendars, function(v) {
            this.listenTo(v, "sync", this.render);
        }, this);
    },

    render: function() {
        var events = [];

        _.each(this.calendars, function(v) {
            events = events.concat(_.pluck(v.models, "attributes"));
        });

        $(".js-next-events").html(Mustache.render(
            this.template,
            {events: events}
        ));

        return this;
    }
});
