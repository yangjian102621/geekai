export function isWechat() {
  const ua = navigator.userAgent.toLowerCase();
  return ua.indexOf("micromessenger") !== -1;
}