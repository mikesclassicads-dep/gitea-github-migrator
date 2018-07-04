var repo_regex = /^[A-Za-z0-9-.]+\/[A-Za-z0-9-.]+$/;

function openSelectRepoModal() {
   $("#add-repos").modal('setting', {
       onApprove: function () {
           var repos = parseReposInTextArea();
           for (var idx in repos) {
               var repo = repos[idx];
               if (repo_regex.test(repo)){
                   addRepoToList(repo);
               }else {
                   alert(repo + " is not a repository")
               }
           }
           return true;
       }
   }).modal('show');
}

function parseReposInTextArea() {
    var text = $("#repo-textform").val();
    return text.split("\n");
}

function addRepoToList(repo) {
    var item = $("#repo-item").children('.item');
    item.html(item.html().replace(/FULL_REPO_NAME/g, repo));
    console.log(repo, item.html());
    $("#repo-list").append(item);
}