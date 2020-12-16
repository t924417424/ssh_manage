var is_login = false;
const protocol = document.location.protocol.split(':')[0];
var ws_p = "ws";
if (protocol == "https") {
    ws_p = "wss";
}
const token = window.localStorage.getItem("token")
if (token == "") {
    if (window != top) {
        top.location.href = "/login";
    }
    window.location.href = "/login";
}
const auth = {
    type: "auth",
    token: token,
}
const socket = new WebSocket(ws_p + '://' + window.location.host + '/v1/term/' + GetQueryString("sid"));
const term = new Terminal({cols: 180, rows: 50, screenKeys: true, cursorBlink: true, cursorStyle: "block"});
term.open(document.getElementById('terms'));
window.onresize = function () {
    fit.fit(term);
};
socket.onopen = function () {
    socket.send(JSON.stringify(auth));  //验证权限
    term.write("正在验证\r\n");
    term.toggleFullscreen(true);
    fit.fit(term);
    term.on('data', function (data) {
        let sdata = {
            type: "cmd",
            cmd: data,
        }
        socket.send(JSON.stringify(sdata));
    });

    term.on('resize', size => {
        //console.log('resize', [size.cols, size.rows]);
        let sdata = {
            type: "resize",
            cols: size.cols,
            rows: size.rows,
        }
        socket.send(JSON.stringify(sdata));
    });

    socket.onmessage = function (msg) {
        if(!is_login){
            is_login = true
        }
        term.write(msg.data);
    };
    socket.onerror = function (e) {
        is_login = false
        console.log(e);
    };

    socket.onclose = function (e) {
        is_login = false
        console.log(e);
        term.write("连接已断开:" + e.reason + "\r\n");
        //term.destroy();
    };
};