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

    /**
     * Return all events from all calendars set in the future.
     *
     * @return Event[]
     */
    getNextEvents: function() {
        var now = new Date().getTime();

        return _.chain(this.calendars)
            .pluck("models")
            .reduce(function(memo, v) {
                return memo.concat(v);
            })
            .filter(function(v) {
                var start = Date.parse(v.attributes.Start);
                return start > now;
            })
            .value()
        ;
    },

    render: function() {
        $(".js-next-events").html(Mustache.render(
            this.template,
            {events: this.getNextEvents()}
        ));

        return this;
    }
});
