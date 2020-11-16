var protocol = document.location.protocol.split(':')[0];
var ws_p = "ws";
if (protocol == "https") {
    ws_p = "wss";
}
var socket = new WebSocket(ws_p + '://' + window.location.host + '/v1/term/' + GetQueryString("sid"));
var term = new Terminal({cols: 180, rows: 50, screenKeys: true, cursorBlink: true, cursorStyle: "block"});
term.open(document.getElementById('terms'));
window.onresize = function () {
    fit.fit(term);
};
socket.onopen = function () {
    var token = window.localStorage.getItem("token")
    if (token == "") {
        if (window != top) {
            top.location.href = "/login";
        }
        window.location.href = "/login";
        return
    }
    var auth = {
        type: "auth",
        token: token,
    }
    socket.send(JSON.stringify(auth));  //验证权限
    term.write("正在验证\r\n");
    term.toggleFullscreen(true);
    fit.fit(term);
    term.on('data', function (data) {
        var sdata = {
            type: "cmd",
            cmd: data,
        }
        socket.send(JSON.stringify(sdata));
    });

    term.on('resize', size => {
        //console.log('resize', [size.cols, size.rows]);
        var sdata = {
            type: "resize",
            cols: size.cols,
            rows: size.rows,
        }
        socket.send(JSON.stringify(sdata));
    });

    socket.onmessage = function (msg) {
        term.write(msg.data);
    };
    socket.onerror = function (e) {
        console.log(e);
    };

    socket.onclose = function (e) {
        console.log(e);
        term.write("连接已断开:" + e.reason + "\r\n");
        //term.destroy();
    };
};