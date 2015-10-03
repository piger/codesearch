$(document).ready(function() {
    var resultList = $("#results-list"),
        searchForm = $("#searchForm");
    
    searchForm.keypress(function(event) {
        // 13 = Enter key
        if (event.which != 13) {
            return true;
        }
        resultList.empty();
        Pace.start();
        $("#searchForm .form-group").removeClass("has-error");
        $("#results-error").text('');
        
        var payload = {
            query: $("#searchText").val()
        };
        var q = {
            url: "/search",
            method: "POST",
            body: JSON.stringify(payload)
        };

        oboe(q)
            .node("results.*", function(result) {
                resultList.append('<li><span class="result-filename">' + result.Filename + '</span> <span class="result-line">' + result.Line + '</span> <span class="result-match">' + result.Match + '</span></li>');
                return oboe.drop;
            })
            .node("errors.*", function(errors) {
                console.log("Uh, errors: " + errors);
                $("#searchForm .form-group").addClass("has-error");
                $("#results-error").append('<div class="alert alert-warning" role="alert">' + errors + '</div>');
            })
            .node("errors", function() {
                Pace.stop();
                this.abort();
            })
            .done(function(things) {
                Pace.stop();
                console.log("results length: " + things.results.length);
            });

        // returning false is like calling event.preventDefault();
        return false;
    });
});
