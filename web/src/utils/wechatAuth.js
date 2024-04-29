//微信相关
import wx from "weixin-js-sdk";

/**
 * 获取公众号授权URL
 * @returns {string}
 */
export function authUrl(appid, url = null) {
    if (url == null) {
        url = window.location.href
    }
    return `https://open.weixin.qq.com/connect/oauth2/authorize?appid=${appid}&redirect_uri=${encodeURIComponent(url)}&response_type=code&scope=snsapi_userinfo&state=STATE#wechat_redirect`
}

/**
 * 获取授权回调code
 */
export function getCode() {
    let url = window.location.href
    let urlStr = url.split('?')[1]
    const urlSearchParams = new URLSearchParams(urlStr)
    return Object.fromEntries(urlSearchParams.entries())
}

/**
 * 获取授权结果
 * @returns {boolean}
 */
export function authResult() {
    const queryBean = getCode()
    return queryBean !=null && queryBean.code !== undefined
}


export function getSignature (data, callback) {
    const {appId, nonceStr, paySign, timeStamp} = data
    // qryWxSignature 这个是调用后台获取签名的接口
    wx.config({
        beta: true,
        debug: false,
        appId: appId,
        timestamp: timeStamp,
        nonceStr: nonceStr,
        signature: paySign,
        // 这里是把所有的方法都写出来了 如果只需要一个方法可以只写一个
        jsApiList: [
            'checkJsApi',
            'onMenuShareTimeline',
            'onMenuShareAppMessage',
            'onMenuShareQQ',
            'onMenuShareWeibo',
            'hideMenuItems',
            'showMenuItems',
            'hideAllNonBaseMenuItem',
            'showAllNonBaseMenuItem',
            'translateVoice',
            'startRecord',
            'stopRecord',
            'onRecordEnd',
            'playVoice',
            'pauseVoice',
            'stopVoice',
            'uploadVoice',
            'downloadVoice',
            'chooseImage',
            'previewImage',
            'uploadImage',
            'downloadImage',
            'getNetworkType',
            'openLocation',
            'getLocation',
            'hideOptionMenu',
            'showOptionMenu',
            'closeWindow',
            'scanQRCode',
            'chooseWXPay',
            'openProductSpecificView',
            'addCard',
            'chooseCard',
            'openCard',
            'openWXDeviceLib',
            'closeWXDeviceLib',
            'configWXDeviceWiFi',
            'getWXDeviceInfos',
            'sendDataToWXDevice',
            'startScanWXDevice',
            'stopScanWXDevice',
            'connectWXDevice',
            'disconnectWXDevice',
            'getWXDeviceTicket',
            'WeixinJSBridgeReady',
            'onWXDeviceBindStateChange',
            'onWXDeviceStateChange',
            'onScanWXDeviceResult',
            'onReceiveDataFromWXDevice',
            'onWXDeviceBluetoothStateChange'
        ]
    })
    wx.ready(function () {
        console.log(callback, 'callback')
        if (callback) callback()
    })
}
