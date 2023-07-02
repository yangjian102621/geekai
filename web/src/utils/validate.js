// 正则校验工具函数

export function validateEmail(email) {
    const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return regex.test(email);
}

export function validateMobile(mobile) {
    const regex = /^1[345789]\d{9}$/;
    return regex.test(mobile);
}