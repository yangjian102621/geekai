// 导入mitt包
import mitt from 'mitt'
// 创建EventBus实例对象
const bus = mitt()
// 共享出eventbus的实例对象
export default bus