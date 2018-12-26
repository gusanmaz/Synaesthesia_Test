
$(document).ready(function() {
    console.log(window.location.pathname );
    if (window.location.pathname == "/st/initTest") {
        $(window).bind("load", ajax_code);
    }
})

var startTime;

function ajax_code() {
    console.log("In ajax code")
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            startTime = new Date().valueOf();
            var testData = JSON.parse(this.responseText);
            console.log(testData.Answer);
            drawDigits(testData.Rectangles, testData.Class1Color, testData.Class2Color);
            drawChoicesWidget(testData.Choices, testData.Answer);
        }
    };
    xhttp.open("POST", "js/data.json", true);
    xhttp.send();

    var textReq = new XMLHttpRequest();
    textReq.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            var textData = JSON.parse(this.responseText);
            console.log(textData);

            ///var testHeader = document.getElementById("test-header");
            var questionText = document.getElementById("question-text");
            var progressText = document.getElementById("progress-bar-text");
            var graphicsElem = document.getElementById("progress-bar-graphics");

            app.progressbar.set(graphicsElem, 10 * (textData.TestNo - 1), 0);
            //app.progressbar.
            progressText.innerText = textData.ProgressBar;
            questionText.textContent = textData.Question;
        }
    };
    textReq.open("POST", "js/data_text.json", true);
    textReq.send();
}

function drawChoicesWidget(choices, answer) {
    var outerElem = document.getElementById("choice_button_container");
    while (outerElem.firstChild) {
        outerElem.removeChild(outerElem.firstChild);
    }

    var choicesNum = choices.length;
    choices.forEach(function (value, index) {
        var buttonElem = document.createElement("button")
        buttonElem.setAttribute("class", "button col button-fill")
        buttonElem.setAttribute("onclick", "processAnswer(this," + answer + ")");
        var textNode = document.createTextNode(value + "");
        buttonElem.appendChild(textNode);
        outerElem.appendChild(buttonElem);
    });
}

function renderProgressBar(testNo) {
    var graphicsElem = document.getElementById("progress-bar-graphics");
    var textElem = document.getElementById("progress-bar-text");

    app.progressbar.set(graphicsElem, 10 * (testNo - 1), 0)
    textElem.innerText = 10 * (testNo - 1) + "% complete!"
}


// 16 30
function drawDigits(rectArr, class1, class2) {
    console.log("I am called!")
    var cookieVal = getCookie("test_no")
    console.log("CookieVal: " + cookieVal);
    $("svg").empty();
    d3.select("svg").selectAll("image").data(rectArr).enter().append("image").attr("x", function (d, i) {
        return d.X1 * 500
    }).attr("y", function (d, i) {
        return d.Y1 * 500
    }).attr("width", function (d, i) {
        return 24
    }).attr("height", function (d, i) {
        return 45
    }).attr("xlink:href", function (d, i) {
        return d.Class == 1 ? 'assets/digits/' + class1 : 'assets/digits/' + class2
    });
}

function processAnswer(obj, answer) {
    obj.disabled = true;
    console.log("In function");
    var btnAns = parseInt(obj.textContent);
    var correctAns = false;

    if (btnAns == answer) {
        correctAns = true;
    }

    endTime = new Date().valueOf();
    console.log("Elapsed time ", endTime - startTime)

    var response = {}
    response.correct = correctAns;
    response.answer = btnAns;
    response.elapsedTime = (endTime - startTime);

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "processUserResponse", true);
    xhr.setRequestHeader("Content-type", "application/json");
    xhr.send(JSON.stringify(response));

    xhr.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            var testNo = getCookie("test_no");
            if (testNo <= 10) {
                ajax_code();
            }
            else {
                window.location.replace("processLastResponse");
            }
        }
    };
}

function getCookie(name) {
    function escape(s) {
        return s.replace(/([.*+?\^${}()|\[\]\/\\])/g, '\\$1');
    };
    var match = document.cookie.match(RegExp('(?:^|;\\s*)' + escape(name) + '=([^;]*)'));
    return match ? match[1] : null;
}
