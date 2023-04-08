/**
 * Util lib functions
 */

// generate a random string
export function randString(length) {
    const str = "0123456789abcdefghijklmnopqrstuvwxyz"
    const size = str.length
    let buf = []
    for (let i = 0; i < length; i++) {
        const rand = Math.random() * size
        buf.push(str.charAt(rand))
    }
    return buf.join("")
}

// 判断是否是移动设备
export function isMobile() {
    const userAgent = navigator.userAgent;
    const mobileRegex = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini|Mobile|mobile|CriOS/i;
    return mobileRegex.test(userAgent);
}

// 格式化日期
export function dateFormat(timestamp, format) {
    if (!timestamp) {
        return '';
    } else if (timestamp < 9680917502) {
        timestamp = timestamp * 1000;
    }
    let year, month, day, HH, mm, ss;
    let time = new Date(timestamp);
    let timeDate;
    year = time.getFullYear(); // 年
    month = time.getMonth() + 1; // 月
    day = time.getDate(); // 日
    HH = time.getHours(); // 时
    mm = time.getMinutes(); // 分
    ss = time.getSeconds(); // 秒

    month = month < 10 ? '0' + month : month;
    day = day < 10 ? '0' + day : day;
    HH = HH < 10 ? '0' + HH : HH; // 时
    mm = mm < 10 ? '0' + mm : mm; // 分
    ss = ss < 10 ? '0' + ss : ss; // 秒

    switch (format) {
        case 'yyyy':
            timeDate = String(year);
            break;
        case 'yyyy-MM':
            timeDate = year + '-' + month;
            break;
        case 'yyyy-MM-dd':
            timeDate = year + '-' + month + '-' + day;
            break;
        case 'yyyy/MM/dd':
            timeDate = year + '/' + month + '/' + day;
            break;
        case 'yyyy-MM-dd HH:mm:ss':
            timeDate = year + '-' + month + '-' + day + ' ' + HH + ':' + mm + ':' + ss;
            break;
        case 'HH:mm:ss':
            timeDate = HH + ':' + mm + ':' + ss;
            break;
        case 'MM':
            timeDate = String(month);
            break;
        default:
            timeDate = year + '-' + month + '-' + day + ' ' + HH + ':' + mm + ':' + ss;
            break;
    }
    return timeDate;
}
