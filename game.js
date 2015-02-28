var WS_URI = "ws://" + location.hostname + ":" + location.port + "/socket/";
var MAP_WIDTH = 32;
var MAP_HEIGHT = 32;
var GRID_WIDTH = 1;
var ws;
var canvas;
var playerId;
var canMove;

$(function() {
    ws = new WebSocket(WS_URI);
    ws.onopen = function(evt) { socketOpen(evt); };
    ws.onclose = function(evt) { socketClose(evt); };
    ws.onmessage = function(evt) { socketMessage(evt); };
    ws.onerror = function(evt) { socketError(evt); };
});

$(document).keydown(function(evt) {
    if (playerId !== null && canMove) {
        switch(evt.keyCode) {
            case 38:
                sendKey("k");
                break;
            case 40:
                sendKey("j");
                break;
            case 37:
                sendKey("h");
                break;
            case 39:
                sendKey("l");
                break;
        }
        canMove = false;
    }
});

function socketOpen(evt) {
    console.log("socket opened");
    console.log(evt);
}

function socketClose(evt) {
    console.log("socket closed");
    console.log(evt);
}

function socketMessage(evt) {
    var a = evt.data.split(" ");
    switch(a[0]) {
        case "redraw":
            var coordStrs = a[1].split("|");
            for (var i = 0; i < coordStrs.length; i++) {
                var coordStr = coordStrs[i].split(",");
                var coord = {
                    x: coordStr[0],
                    y: coordStr[1],
                    id: coordStr[2]
                };

                clearCoord(coord);
                if (coord.id !== "-") {
                    drawCoord(coord);
                }
            }

            canMove = true;
            break;
        case "start":
            playerId = a[1];
            startGame();
            break;
        case "death":
            endPlayer();
            break;
        case "end":
            endGame(a[1]);
            break;
    }
}

function socketError(evt) {
    console.log("socket error");
    console.log(evt);
    var ctx = canvas.getContext("2d");
    ctx.fillStyle = "white";
    ctx.globalAlpha = 0.5;
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    ctx.font = "40px Oswald";
    ctx.textAlign = "center";
    ctx.globalAlpha = 1.0;
    ctx.textBaseline = "center";
    ctx.fillStyle = "black";
    ctx.fillText("Connection error :(", canvas.width / 2, (canvas.height / 2) - 30, canvas.width * 0.9);
}

function startGame() {
    $("#loading-container").addClass("hidden");
    $("#game-container").removeClass("hidden");

    $("body").css({"background": gridColorForId(playerId)});
    canvas = document.getElementById("drawing");

    canvas.width = MAP_WIDTH * Math.floor($("#game").width() / MAP_WIDTH);
    canvas.height = MAP_HEIGHT * Math.floor($("#game").height() / MAP_HEIGHT);

    drawGrid();
}

function endPlayer() {

}

function endGame(winnerId) {
    var winner = "Nobody";
    switch(winnerId) {
        case "0":
            winner = "Red";
            break;
        case "1":
            winner = "Blue";
            break;
        case "2":
            winner = "Green";
            break;
        case "3":
            winner = "Yellow";
            break;
    }

    var ctx = canvas.getContext("2d");
    ctx.fillStyle = "white";
    ctx.globalAlpha = 0.5;
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    ctx.font = "40px Oswald";
    ctx.textAlign = "center";
    ctx.globalAlpha = 1.0;
    ctx.textBaseline = "center";
    ctx.fillStyle = colorForId(winnerId);
    ctx.fillText(winner + " won this round!", canvas.width / 2, (canvas.height / 2) - 30, canvas.width * 0.9);
    ctx.fillStyle = "black";
    ctx.fillText("Refresh the page to play again.", canvas.width / 2, (canvas.height / 2) + 30, canvas.width * 0.9);
}

function sendKey(key) {
    ws.send("key " + key);
}

function clearCoord(coord) {
    var ctx = canvas.getContext("2d");
    ctx.fillStyle = gridColorForId(playerId);
    ctx.fillRect(coord.x * (canvas.width / MAP_WIDTH), coord.y * (canvas.height / MAP_HEIGHT), canvas.width / MAP_WIDTH, canvas.height / MAP_HEIGHT);
    r = toCtxRect({x: coord.x, y: coord.y});
    ctx.clearRect(r.x, r.y, r.w, r.h);
}

function drawCoord(coord) {
    var ctx = canvas.getContext("2d");
    ctx.fillStyle = colorForId(coord.id);
    r = toCtxRect(coord);
    ctx.fillRect(r.x, r.y, r.w, r.h);
}

function drawGrid() {
    for (var x = 0; x < MAP_WIDTH; x++) {
        for (var y = 0; y < MAP_HEIGHT; y++) {
           clearCoord({
                x: x,
                y: y
            });
        }
    }
}

function toCtxRect(coord) {
    return {
        x: coord.x * (canvas.width / MAP_WIDTH) + GRID_WIDTH,
        y: coord.y * (canvas.height / MAP_HEIGHT) + GRID_WIDTH,
        w: canvas.width / MAP_WIDTH - (2 * GRID_WIDTH),
        h: canvas.height / MAP_HEIGHT - (2 * GRID_WIDTH)
    }
}

function colorForId(id) {
    switch(id) {
        case "0":
            return "#FF3333";
        case "1":
            return "#3333FF";
        case "2":
            return "#33FF33";
        case "3":
            return "#FFFF33";
        case "x":
            return "#333333";
    }

    return "white";
}

function gridColorForId(id) {
    switch(id) {
        case "0":
            return "#FFAAAA";
        case "1":
            return "#AAAAFF";
        case "2":
            return "#AAFFAA";
        case "3":
            return "#FFFFAA";
    }

    return "#EEEEEE";
}
