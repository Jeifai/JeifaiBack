var urls =  [{"id":3,"linkedinurl":"https://www.linkedin.com/company/deutschebahn"},{"id":4,"linkedinurl":"https://www.linkedin.com/company/microsoft"},{"id":5,"linkedinurl":"https://www.linkedin.com/company/twitter"},{"id":6,"linkedinurl":"https://www.linkedin.com/company/shopify"},{"id":9,"linkedinurl":"https://www.linkedin.com/company/celoorg"},{"id":7,"linkedinurl":"https://www.linkedin.com/company/urbansportsclub"},{"id":8,"linkedinurl":"https://www.linkedin.com/company/blinkist"},{"id":10,"linkedinurl":"https://www.linkedin.com/company/n26"}]

var initialQuery = "INSERT INTO linkedin (targetid, employees, followers, headquarters, industry, companySize, createdat) VALUES "

var headers = {
    "headers": {
        "mode": 'no-cors',
        "User-Agent": 'Mozilla/5.0 Gecko/20100401 Firefox/3.6.3'
    }
}

function sleep() {
   var currentTime = new Date().getTime();
   var sleepRandom = Math.floor(Math.random() * 5000);
   while (currentTime + sleepRandom >= new Date().getTime()) {}
}

function getDataFromResponse(elem) {
    employees = document.getElementById("employee-count").textContent.split(" ")[0].replace(/\D/g,'');
    followers = document.getElementById("followers").textContent.split(" ")[0].replace(/\D/g,'');
    headquarters = document.getElementById("headquarters").textContent;
    industry = document.getElementById("company-industry-data").textContent;
    companySize = document.getElementById("company-size-data").textContent;
    var queryBlock = "(" + elem.id + "," + employees + "," + followers + ",'" + headquarters + "','" + industry + "','" + companySize + "', current_timestamp)";
    console.log(queryBlock);
    return queryBlock
}

var tempQuery = ""
let lastIterations = urls.length - 1;
for (var i = 0; i < urls.length; i++) {
    var response = await fetch(urls[i].linkedinurl, headers);
    var responseText = await response.text();
    document.getElementsByTagName("html")[0].innerHTML = responseText;
    if (i < lastIterations) {
        tempQuery = tempQuery + getDataFromResponse(urls[i]) + ",";
    } else {
        tempQuery = tempQuery + getDataFromResponse(urls[i]) + ";";
    }
    sleep();
}

var fullQuery = initialQuery + tempQuery;
console.log(fullQuery);