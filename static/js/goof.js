(function() {
    "use strict";
    var loadTemplates = function() {
        var templates = {};
        $('.js-template').each(function(_, tpl) {
            var html = $(tpl).html();
            Mustache.parse(html);
            templates[tpl.getAttribute('data-name')] = html;
        });
        return templates;
    };

    $(function() {
        var templates = loadTemplates();
        $('body').append(Mustache.render(templates.navbar));
        $('body').append(Mustache.render(templates.calendar));
        $('.js-calendar-body').append(Mustache.render(
            templates["calendar-multiweek"]
        ));
    });
})();
