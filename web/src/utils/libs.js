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