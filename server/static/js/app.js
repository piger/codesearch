// Code Search js app

$(document).ready(function() {
    var resultList = $("#results-list");

    // Oboe callbacks
    var onResult = function(result) {
        var out = [];
        out.push('<li><span class="result-filename">');
        out.push(result.Filename);
        out.push('</span>');
        out.push('<span class="result-line"><small>');
        out.push(':' + result.Line);
        out.push('</small></span> ');
        out.push('<span class="result-match"><code>');
        out.push(result.Match);
        out.push('</code></span></li>');
        resultList.append(out.join(''));

        return oboe.drop;
    };

    var onError = function(error) {
        $("#searchForm .form-group").addClass("has-error");
        $("#results-error").append('<div class="alert alert-warning" role="alert">' + error + '</div>');
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
