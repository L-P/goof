"use strict";

document.goof = document.goof || {};

document.goof.Calendar = Backbone.Collection.extend({
    model: document.goof.Event,
    name: null,
    paletteId: null,

    url: function() {
        return "/calendar/" + this.name;
    },

    parse: function(response, options) {
        if (response.Meta.Errors && response.Meta.Errors.length > 0) {
            console.error(response.Meta.Errors);
        }

        return response.Data.Calendar.Events;
    },

    initialize: function(models, options) {
        this.name = options.name;
        this.paletteId = options.paletteId;
    },
});
