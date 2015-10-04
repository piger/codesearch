// Code Search js app

$(document).ready(function() {
    var resultList = $("#results-list");

    // Oboe callbacks
    var onResult = function(result) {
        resultList.append('<li><span class="result-filename">' + result.Filename + '</span> <span class="result-line">' + result.Line + '</span> <span class="result-match">' + result.Match + '</span></li>');
        return oboe.drop;
    };

    var onError = function(error) {
        console.log("Uh, errors: " + errors);
        $("#searchForm .form-group").addClass("has-error");
        $("#results-error").append('<div class="alert alert-warning" role="alert">' + errors + '</div>');
    };

    var onErrors = function() {
        Pace.stop();
        this.abort();
    };

    var onDone = function(things) {
        Pace.stop();
        console.log("results length: " + things.results.length);
    };

    // form submit event
    $("#searchForm").keypress(function(event) {
        // 13 = Enter key
        if (event.which != 13) {
            return true;
        }

        resultList.empty();
        Pace.start();
        $("#searchForm .form-group").removeClass("has-error");
        $("#results-error").text('');

        // POST query payload
        var payload = {
            query: $("#searchText").val()
        };

        // Oboe parameters
        var q = {
            url: "/search",
            method: "POST",
            body: JSON.stringify(payload)
        };

        oboe(q)
            .node("results.*", onResult)
            .node("errors.*", onError)
            .node("errors", onErrors)
            .done(onDone);

        // returning false is like calling event.preventDefault();
        return false;
    });
});
