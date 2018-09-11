var done = false;

function update() {
    $.getJSON("/status", function(data) {
        handleData(data);
    }).always(function() {
        if(!done){
            window.setTimeout(update, 1000);
        }
    });
}
$(function() {
   $(".repo-progress").progress({
       text: {
           active  : 'Migrated or failed {value} of {total} repositories',
           success : '{total} repositories migrated or failed!'
       },
       total: $(".repo-progress").data("total"),
       value: 0
   });
});

function handleData(data) {
    if(Object.keys(data.finished).length + Object.keys(data.failed).length === $(".repo-progress").progress('get total')) {
        $(".repo-progress").progress('complete');
        done = true;
    } else {
        $(".repo-progress").progress('set progress', Object.keys(data.finished).length + Object.keys(data.failed).length);
    }
    data.pending.forEach(function(repo) {
        var content = contentFromRepo(repo);
        if (!content.hasClass("pending")) {
            content.html(renderPending().html());
            content.addClass("pending");
        }
    });
    forEach(data.failed, function (repo, errormsg) {
        var content = contentFromRepo(repo);
        if (!content.hasClass("failed")) {
            content.html(renderFailed(errormsg).html());
            content.addClass("failed");
        }
    });
    forEach(data.running, function (repo, report) {
        var content = handleNonPending(repo, report);
        if (content.find(".comment-progress").progress('get total') !== report.total_comments) {
            content.find(".comment-progress").progress('set total', report.total_comments)
        }
        if (content.find(".issue-progress").progress('get total') !== report.total_issues) {
            content.find(".issue-progress").progress('set total', report.total_issues)
        }
        content.find(".comment-progress").progress('set progress', report.migrated_comments + report.failed_comments);
        content.find(".issue-progress").progress('set progress', report.migrated_issues + report.failed_issues);
        content.find('.log-content').val(report.log);
        content.find('.log-content').scrollTop(content.find('.log-content')[0].scrollHeight)
    });
    forEach(data.finished, function (repo, report) {
        var content = handleNonPending(repo, report);
        if (content.find(".comment-progress").progress('get total') !== report.total_comments) {
            content.find(".comment-progress").progress('set total', report.total_comments)
        }
        if (content.find(".issue-progress").progress('get total') !== report.total_issues) {
            content.find(".issue-progress").progress('set total', report.total_issues)
        }
        content.find(".comment-progress").progress('set progress', report.migrated_comments + report.failed_comments);
        content.find(".issue-progress").progress('set progress', report.migrated_issues + report.failed_issues);
        content.find(".issue-progress").progress('complete');
        content.find(".comment-progress").progress('complete');
        content.find('.log-content').val(report.log);
    });
}
function forEach(object, callback) {
    for(var prop in object) {
        if(object.hasOwnProperty(prop)) {
            callback(prop, object[prop]);
        }
    }
}

function handleNonPending(repo, report) {
    var content = contentFromRepo(repo);
    if(!content.hasClass("non-pending")) {
        content.html(renderNonPending().html());
        content.find(".accordion").accordion();
        content.find(".issue-progress").progress({
            text: {
                active  : 'Migrated {value} of {total} issues',
                success : '{total} issues migrated!'
            },
            total: report.total_issues,
            value: report.migrated_issues + report.failed_issues,
        });
        content.find(".comment-progress").progress({
            text: {
                active  : 'Migrated {value} of {total} comments',
                success : '{total} comments migrated!'
            },
            total: report.total_comments+1,
            value: report.migrated_comments + report.failed_comments,
        });
        content.addClass("non-pending");
    }
    content.find(".failed-issues").text(report.failed_issues);
    content.find(".failed-comments").text(report.failed_comments);
    return content
}

function contentFromRepo(repo) {
    return $(".repo-content[data-repo='" + repo + "']")
}

function renderPending() {
    return $("#content-pending").clone();
}

function renderFailed(errormsg) {
    var failed = $("#content-failed").clone();
    failed.find(".errormsg").text(errormsg);
    return failed
}
function renderNonPending() {
    return $("#content-nonpending").clone();
}

$(update());