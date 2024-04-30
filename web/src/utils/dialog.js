/**
 * Util lib functions
 */
import {showConfirmDialog, showFailToast, showSuccessToast, showToast} from "vant";
import {isMobile} from "@/utils/libs";
import {ElMessage} from "element-plus";

export function showLoginDialog(router) {
    showConfirmDialog({
        title: '登录',
        message:
            '此操作需要登录才能进行，前往登录？',
    }).then(() => {
        router.push("/login")
    }).catch(() => {
        // on cancel
    });
}

export function showMessageOK(message) {
    if (isMobile()) {
        showSuccessToast(message)
    } else {
        ElMessage.success(message)
    }
}


export function showMessageInfo(message) {
    if (isMobile()) {
        showToast(message)
    } else {
        ElMessage.info(message)
    }
}

export function showMessageError(message) {
    if (isMobile()) {
        showFailToast(message)
    } else {
        ElMessage.error(message)
    }
}
