var urls = [{"id":232,"linkedinurl":"https://www.linkedin.com/company/gritspot"},{"id":234,"linkedinurl":"https://www.linkedin.com/company/caspar-health"},{"id":236,"linkedinurl":"https://www.linkedin.com/company/fortocompany"},{"id":238,"linkedinurl":"https://www.linkedin.com/company/joblift"},{"id":240,"linkedinurl":"https://www.linkedin.com/company/medloopltd"},{"id":242,"linkedinurl":"https://www.linkedin.com/company/merantix"},{"id":244,"linkedinurl":"https://www.linkedin.com/company/zenjob"},{"id":246,"linkedinurl":"https://www.linkedin.com/company/coachhub-io"},{"id":72,"linkedinurl":"https://www.linkedin.com/company/blacklane-gmbh"},{"id":245,"linkedinurl":"https://www.linkedin.com/company/peat-ug-haftungsbeschränkt-"},{"id":239,"linkedinurl":"https://www.linkedin.com/company/kontist"},{"id":233,"linkedinurl":"https://www.linkedin.com/school/careerfoundry"},{"id":225,"linkedinurl":"https://www.linkedin.com/company/paintgun"},{"id":208,"linkedinurl":"https://www.linkedin.com/company/deloitte"},{"id":202,"linkedinurl":"https://www.linkedin.com/company/vodafone"},{"id":166,"linkedinurl":"https://www.linkedin.com/company/github"},{"id":4,"linkedinurl":"https://www.linkedin.com/company/microsoft"},{"id":5,"linkedinurl":"https://www.linkedin.com/company/twitter"},{"id":6,"linkedinurl":"https://www.linkedin.com/company/shopify"},{"id":9,"linkedinurl":"https://www.linkedin.com/company/celoorg"},{"id":7,"linkedinurl":"https://www.linkedin.com/company/urbansportsclub"},{"id":8,"linkedinurl":"https://www.linkedin.com/company/blinkist"},{"id":10,"linkedinurl":"https://www.linkedin.com/company/n26"},{"id":21,"linkedinurl":"https://www.linkedin.com/company/contentful"},{"id":23,"linkedinurl":"https://www.linkedin.com/company/hometogo"},{"id":25,"linkedinurl":"https://www.linkedin.com/company/lana-labs"},{"id":292,"linkedinurl":"https://www.linkedin.com/company/adjustcom"},{"id":170,"linkedinurl":"https://www.linkedin.com/company/getyourguide-ag"},{"id":172,"linkedinurl":"https://www.linkedin.com/company/celonis"},{"id":174,"linkedinurl":"https://www.linkedin.com/company/about-you-gmbh"},{"id":176,"linkedinurl":"https://www.linkedin.com/company/taxfix"},{"id":178,"linkedinurl":"https://www.linkedin.com/company/fincompare"},{"id":180,"linkedinurl":"https://www.linkedin.com/company/pair-finance"},{"id":11,"linkedinurl":"https://www.linkedin.com/company/kununu"},{"id":53,"linkedinurl":"https://www.linkedin.com/company/revolut"},{"id":182,"linkedinurl":"https://www.linkedin.com/company/liqid"},{"id":294,"linkedinurl":"https://www.linkedin.com/company/bonify-germany"},{"id":12,"linkedinurl":"https://www.linkedin.com/company/mitte®"},{"id":61,"linkedinurl":"https://www.linkedin.com/company/circleci"},{"id":13,"linkedinurl":"https://www.linkedin.com/company/babelforce"},{"id":14,"linkedinurl":"https://www.linkedin.com/company/zalando"},{"id":126,"linkedinurl":"https://www.linkedin.com/company/datadog"},{"id":128,"linkedinurl":"https://www.linkedin.com/company/stripe"},{"id":310,"linkedinurl":"https://www.linkedin.com/company/wirsindconstruyo"},{"id":306,"linkedinurl":"https://www.linkedin.com/company/candis-gmbh"},{"id":247,"linkedinurl":"https://www.linkedin.com/company/join-raisin"},{"id":243,"linkedinurl":"https://www.linkedin.com/company/ninox-software-gmbh"},{"id":241,"linkedinurl":"https://www.linkedin.com/company/medwing"},{"id":237,"linkedinurl":"https://www.linkedin.com/company/idagio"},{"id":235,"linkedinurl":"https://www.linkedin.com/company/ecosia"},{"id":231,"linkedinurl":"https://www.linkedin.com/company/pitchhq"},{"id":229,"linkedinurl":"https://www.linkedin.com/company/chatterbug"},{"id":227,"linkedinurl":"https://www.linkedin.com/company/nen-energia"},{"id":216,"linkedinurl":"https://www.linkedin.com/company/facebook"},{"id":214,"linkedinurl":"https://www.linkedin.com/company/subito-it"},{"id":212,"linkedinurl":"https://www.linkedin.com/company/roche"},{"id":206,"linkedinurl":"https://www.linkedin.com/company/bendingspoons"},{"id":204,"linkedinurl":"https://www.linkedin.com/company/glickon-srl"},{"id":198,"linkedinurl":"https://www.linkedin.com/company/freeda-media"},{"id":80,"linkedinurl":"https://www.linkedin.com/company/docker"},{"id":78,"linkedinurl":"https://www.linkedin.com/company/quora"},{"id":16,"linkedinurl":"https://www.linkedin.com/company/pentabanking"},{"id":15,"linkedinurl":"https://www.linkedin.com/company/imusician"},{"id":74,"linkedinurl":"https://www.linkedin.com/company/flixbus"}]

var initialQuery = "INSERT INTO linkedin (targetid, employees, followers, headquarters, industry, companySize, createdat) VALUES "

var headers = {
    "headers": {
        "mode": 'no-cors',
        "User-Agent": 'Mozilla/5.0 Gecko/20100401 Firefox/3.6.3'
    }
}

function sleep() {
   var currentTime = new Date().getTime();
   var sleepRandom = Math.floor(Math.random() * 5000) + 4000;
   while (currentTime + sleepRandom >= new Date().getTime()) {}
}

function getDataFromResponse(elem) {
    try {
        employees = document.getElementById("employee-count").textContent.split(" ")[0].replace(/\D/g,'');
    } catch {
        employees = ""
    }
    try {
        followers = document.getElementById("followers").textContent.split(" ")[0].replace(/\D/g,'');
    } catch {
        followers = ""
    }
    try {
        headquarters = document.getElementById("headquarters").textContent;
    } catch {
        headquarters = ""
    }
    try {
        industry = document.getElementById("company-industry-data").textContent;
    } catch {
        industry = ""
    }
    try {
        companySize = document.getElementById("company-size-data").textContent;
    } catch {
        companySize = ""
    }
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