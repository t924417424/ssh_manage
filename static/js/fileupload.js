let dropbox = document.getElementById("terms");       //要监听拖动上传的节点

let fileDrop = {
    startTime: 0,
    endTime: 0,
    uploadLength: 0, //上传数量
    //splitSize: 1024 * 1024 * 2, //文件上传分片大小
    filesList: [], // 文件列表数组
    errorLength: 0, //上传失败文件数量
    isUpload: true, //上传状态，是否可以上传
    //uploadSuspend:[],  //上传暂停参数
    isUploadNumber: 800,//限制单次上传数量
    uploadAllSize: 0, // 上传文件总大小
    uploadedSize: 0, // 已上传文件大小
    topUploadedSize: 0, // 上一次文件上传大小
    uploadExpectTime: 0, // 预计上传时间
    //initTimer:0, // 初始化计时
    speedInterval: null, //平局速度定时器
    timerSpeed: 0, //速度
    uploading: false,
    cancel: false,
}

dropbox.addEventListener("dragleave", function (e) {
    //e.stopPropagation();
    e.preventDefault();
}, false);

dropbox.addEventListener("dragenter", function (e) {
    //e.stopPropagation();
    e.preventDefault();
}, false);

dropbox.addEventListener("dragover", function (e) {
    //e.stopPropagation();
    e.preventDefault();
}, false);

dropbox.addEventListener("drop", changes, false);

function changes(e) {
    if(!is_login){
        layer.msg("请等待服务器连接！")
    }
    e.preventDefault();
    let items = e.dataTransfer.items, time, num = 0
    if (fileDrop.uploading) {
        layer.msg("已有文件队列上传中")
        return false
    }
    if (items && items.length && items[0].webkitGetAsEntry != null) {
        if (items[0].kind != 'file') return false;
    }
    if (fileDrop.filesList == null) fileDrop.filesList = []
    for (let i = fileDrop.filesList.length - 1; i >= 0; i--) {
        if (fileDrop.filesList[i].is_upload) fileDrop.filesList.splice(-i, 1)
    }

    function update_sync(s) {
        s.getFilesAndDirectories().then(function (subFilesAndDirs) {
            return iterateFilesAndDirs(subFilesAndDirs, s.path);
        });
    }

    let iterateFilesAndDirs = function (filesAndDirs, path) {
        for (let i = 0; i < filesAndDirs.length; i++) {
            if (typeof (filesAndDirs[i].getFilesAndDirectories) == 'function') {
                update_sync(filesAndDirs[i])
            } else {
                if (num > 100) {
                    //fileDrop.isUpload = false;
                    layer.msg(' '+ fileDrop.isUploadNumber +'份，无法上传,请压缩后上传!。',{icon:2,area:'405px'});
                    //clearTimeout(time);
                    return false;
                }
                fileDrop.filesList.push({
                    file: filesAndDirs[i],
                    path: path,
                    name: filesAndDirs[i].name.replace('//', '/'),
                    local: (path == "/" ? "" : path) + "/" + filesAndDirs[i].name.replace('//', '/'),
                    size: to_size(filesAndDirs[i].size),
                    upload: 0, //上传状态,未上传：0、上传中：1，已上传：2，上传失败：-1
                    is_upload: false
                });
                fileDrop.uploadAllSize += filesAndDirs[i].size
                fileDrop.uploadLength++;
            }
        }
    }
    if ('getFilesAndDirectories' in e.dataTransfer) {
        e.dataTransfer.getFilesAndDirectories().then(function (filesAndDirs) {
            return iterateFilesAndDirs(filesAndDirs, '/');
        });
    }
    //console.log(fileDrop.filesList)
    layer.load(1, {
        shade: [0.1,'#fff'] //0.1透明度的白色背景
    });
    setTimeout(function () {
        layer.closeAll('loading')
        open_upload_window()
    },3000)

}

function open_upload_window() {

    let template = `
    <table class="layui-table" lay-even="" lay-skin="row" id="file_upload" style="table-layout: fixed;padding-top: 0">
      <colgroup>
        <col width="250">
        <col width="150">
        <col width="150">
        <col>
      </colgroup>
      <thead>
         <tr>
          <th>文件路径</th>
          <th>文件大小</th>
          <th>状态</th>
        </tr>
      </thead>
      <tbody align="center">
        
      </tbody>
    </table>
    `
    layer.open({
        type: 1,
        closeBtn: 1,
        maxmin: true,
        area: ['550px', '455px'],
        btn: ['开始上传', '取消上传'],
        title: '上传文件',
        skin: 'file_dir_uploads',
        shade: 0.4,
        shadeClose: false,
        content: template,
        success: function () {
            for (let i = 0; i < fileDrop.filesList.length; i++) {
                $("#file_upload tbody").append(create_row(i, fileDrop.filesList[i]));
            }
        }
    });
}

function create_row(index, file) {
    console.log(file)
    return "<tr id='" + index + "'><td title='" + file.local + "' style=\"white-space:nowrap;overflow:hidden;text-overflow: ellipsis;\">" + file.local + "</td> <td>" + file.size + "</td> <td>" + getstatu(file.upload) + "</td></tr>"
}

function getstatu(statu) {
    //上传状态,未上传：0、上传中：1，已上传：2，上传失败：-1
    if (statu == -1) {
        return "<font color='red'>上传失败</font>"
    } else {
        if (statu == 0) {
            return "<font color='black'>未上传</font>"
        } else if (statu == 1) {
            return "<font color='#808080'>上传中</font>"
        } else {
            return "<font color='green'>已上传</font>"
        }
    }
}

function to_size(a) {
    var d = [" B", " KB", " MB", " GB", " TB", " PB"];
    var e = 1024;
    for (var b = 0; b < d.length; b += 1) {
        if (a < e) {
            var num = (b === 0 ? a : a.toFixed(2)) + d[b];
            return (!isNaN((b === 0 ? a : a.toFixed(2))) && typeof num != 'undefined') ? num : '0B';
        }
        a /= e
    }
}