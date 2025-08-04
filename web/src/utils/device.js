/**
 * 设备检测工具函数
 * 用于判断当前设备类型，支持PC端和移动端的智能识别
 */

/**
 * 检测设备类型
 * @returns {string} 'mobile' | 'desktop'
 */
export const detectDevice = () => {
  const userAgent = navigator.userAgent.toLowerCase()

  // 移动设备关键词检测
  const mobileKeywords = [
    'mobile',
    'android',
    'iphone',
    'ipad',
    'phone',
    'blackberry',
    'opera mini',
    'windows phone',
    'iemobile',
  ]

  // 平板设备关键词检测
  const tabletKeywords = ['tablet', 'ipad', 'playbook', 'silk', 'kindle']

  // 检查是否为移动设备
  const isMobile = mobileKeywords.some((keyword) => userAgent.includes(keyword))

  // 检查是否为平板设备
  const isTablet = tabletKeywords.some((keyword) => userAgent.includes(keyword))

  // 检查屏幕尺寸
  const screenWidth = window.innerWidth
  const screenHeight = window.innerHeight

  // 移动设备判断逻辑
  if (isMobile || isTablet || screenWidth <= 768) {
    return 'mobile'
  }

  return 'desktop'
}

/**
 * 检测是否为移动设备
 * @returns {boolean}
 */
export const isMobileDevice = () => {
  return detectDevice() === 'mobile'
}

/**
 * 检测是否为桌面设备
 * @returns {boolean}
 */
export const isDesktopDevice = () => {
  return detectDevice() === 'desktop'
}

/**
 * 获取设备跳转路径
 * @param {string} deviceType - 设备类型
 * @param {string} defaultPath - 默认路径
 * @returns {string} 跳转路径
 */
export const getDeviceRedirectPath = (deviceType, defaultPath = '/') => {
  if (deviceType === 'mobile') {
    return '/mobile'
  }
  return defaultPath
}

/**
 * 根据当前设备获取跳转路径
 * @param {string} defaultPath - 默认路径
 * @returns {string} 跳转路径
 */
export const getCurrentDeviceRedirectPath = (defaultPath = '/') => {
  const deviceType = detectDevice()
  return getDeviceRedirectPath(deviceType, defaultPath)
}

/**
 * 检测屏幕尺寸
 * @returns {object} { width, height, isSmall, isMedium, isLarge }
 */
export const getScreenInfo = () => {
  const width = window.innerWidth
  const height = window.innerHeight

  return {
    width,
    height,
    isSmall: width <= 768,
    isMedium: width > 768 && width <= 1024,
    isLarge: width > 1024,
  }
}

/**
 * 检测浏览器类型
 * @returns {string} 浏览器类型
 */
export const detectBrowser = () => {
  const userAgent = navigator.userAgent.toLowerCase()

  if (userAgent.includes('chrome')) return 'chrome'
  if (userAgent.includes('firefox')) return 'firefox'
  if (userAgent.includes('safari') && !userAgent.includes('chrome')) return 'safari'
  if (userAgent.includes('edge')) return 'edge'
  if (userAgent.includes('opera')) return 'opera'

  return 'unknown'
}

/**
 * 检测操作系统
 * @returns {string} 操作系统类型
 */
export const detectOS = () => {
  const userAgent = navigator.userAgent.toLowerCase()

  if (userAgent.includes('windows')) return 'windows'
  if (userAgent.includes('mac')) return 'macos'
  if (userAgent.includes('linux')) return 'linux'
  if (userAgent.includes('android')) return 'android'
  if (userAgent.includes('ios')) return 'ios'

  return 'unknown'
}
