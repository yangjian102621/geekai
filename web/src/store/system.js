// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import Storage from "good-storage";

export function GetFileIcon(ext) {
    const files = {
        ".docx": "doc.png",
        ".doc": "doc.png",
        ".xls": "xls.png",
        ".xlsx": "xls.png",
        ".csv": "xls.png",
        ".ppt": "ppt.png",
        ".pptx": "ppt.png",
        ".md": "md.png",
        ".pdf": "pdf.png",
        ".sql": "sql.png",
        ".mp3": "mp3.png",
        ".wav": "mp3.png",
        ".mp4": "mp4.png",
        ".avi": "mp4.png",
    }
    if (files[ext]) {
        return '/images/ext/' + files[ext]
    }

    return '/images/ext/file.png'
}

// 获取文件类型
export function GetFileType (ext) {
    return ext.replace(".", "").toUpperCase()
}

// 将文件大小转成字符
export function FormatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

export function setRoute(path) {
    Storage.set(import.meta.env.VITE_APP_KEY_PREFIX + 'ROUTE_', path)
}

export function getRoute() {
    return Storage.get(import.meta.env.VITE_APP_KEY_PREFIX + 'ROUTE_')
}
