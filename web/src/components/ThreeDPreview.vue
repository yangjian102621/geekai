<template>
  <div class="three-d-preview">
    <div ref="container" class="preview-container"></div>

    <!-- 控制面板 -->
    <div class="control-panel">
      <div class="control-group">
        <label>缩放</label>
        <div class="scale-controls">
          <el-button size="small" @click="zoomOut" :disabled="scale <= 0.1">
            <el-icon><Minus /></el-icon>
          </el-button>
          <span class="scale-value">{{ scale.toFixed(1) }}x</span>
          <el-button size="small" @click="zoomIn" :disabled="scale >= 3">
            <el-icon><Plus /></el-icon>
          </el-button>
        </div>
      </div>

      <div class="control-group">
        <label>模型颜色</label>
        <div class="color-picker">
          <el-color-picker
            v-model="modelColor"
            @change="updateModelColor"
            :predefine="predefineColors"
            size="small"
          />
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-overlay">
      <div class="loading-content">
        <el-icon class="is-loading"><Loading /></el-icon>
        <p>加载3D模型中...</p>
        <div v-if="loadingProgress > 0" class="loading-progress">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: loadingProgress + '%' }"></div>
          </div>
          <span class="progress-text">{{ loadingProgress.toFixed(1) }}%</span>
        </div>
      </div>
    </div>

    <!-- 错误状态 -->
    <div v-if="error" class="error-overlay">
      <div class="error-content">
        <el-icon><Warning /></el-icon>
        <p>{{ error }}</p>
        <el-button size="small" @click="retryLoad">重试</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { Loading, Minus, Plus, Warning } from '@element-plus/icons-vue'
import { ElButton, ElIcon } from 'element-plus'
import * as THREE from 'three'
import { OrbitControls } from 'three/addons/controls/OrbitControls.js'
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js'
import { OBJLoader } from 'three/addons/loaders/OBJLoader.js'
import { STLLoader } from 'three/addons/loaders/STLLoader.js'
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'

// Props
const props = defineProps({
  modelUrl: {
    type: String,
    required: false,
  },
  modelType: {
    type: String,
    default: 'glb',
  },
})

// 响应式数据
const container = ref(null)
const loading = ref(true)
const error = ref('')
const loadingProgress = ref(0)
const scale = ref(1)
const modelColor = ref('#00ff88')
const predefineColors = ref([
  '#00ff88', // 亮绿色
  '#ff6b6b', // 亮红色
  '#4ecdc4', // 亮青色
  '#45b7d1', // 亮蓝色
  '#f9ca24', // 亮黄色
  '#f0932b', // 亮橙色
  '#eb4d4b', // 亮粉红
  '#6c5ce7', // 亮紫色
  '#a29bfe', // 亮靛蓝
  '#fd79a8', // 亮玫瑰
])

// Three.js 相关变量
let scene, camera, renderer, controls, model, mixer, clock
let animationId
let baseScale = 1 // 存储基础缩放值

// 初始化Three.js场景
const initThreeJS = () => {
  if (!container.value) {
    console.error('ThreeDPreview: 容器元素不存在')
    return
  }

  // 创建场景
  scene = new THREE.Scene()
  scene.background = new THREE.Color(0x2a2a2a) // 深灰色背景，类似截图

  // 获取容器尺寸，确保有最小尺寸
  const containerRect = container.value.getBoundingClientRect()
  const width = Math.max(containerRect.width || 400, 400)
  const height = Math.max(containerRect.height || 300, 300)

  // 创建相机 - 参考截图的视角（稍微俯视，从左上角观察）
  camera = new THREE.PerspectiveCamera(75, width / height, 0.1, 1000)
  camera.position.set(3, 3, 3) // 从左上角俯视角度

  // 创建渲染器
  renderer = new THREE.WebGLRenderer({
    antialias: true,
    alpha: true,
    preserveDrawingBuffer: true,
  })
  renderer.setSize(width, height)
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2))
  renderer.shadowMap.enabled = true
  renderer.shadowMap.type = THREE.PCFSoftShadowMap
  renderer.outputColorSpace = THREE.SRGBColorSpace

  // 添加到容器
  container.value.appendChild(renderer.domElement)

  // 创建控制器
  controls = new OrbitControls(camera, renderer.domElement)
  controls.enableDamping = true
  controls.dampingFactor = 0.05

  // 添加光源
  addLights()

  // 添加地面
  addGround()

  // 添加坐标轴辅助线
  addAxesHelper()

  // 创建时钟
  clock = new THREE.Clock()

  // 开始渲染循环
  animate()

  // 监听窗口大小变化
  window.addEventListener('resize', onWindowResize)
}

//

// 添加光源 - 参考截图的柔和光照效果
const addLights = () => {
  // 环境光 - 提供基础照明，参考截图的柔和效果
  const ambientLight = new THREE.AmbientLight(0x404040, 0.6)
  scene.add(ambientLight)

  // 主方向光 - 从左上角照射，模拟截图中的光照方向
  const directionalLight = new THREE.DirectionalLight(0xffffff, 0.8)
  directionalLight.position.set(5, 5, 3)
  directionalLight.castShadow = true
  directionalLight.shadow.mapSize.width = 2048
  directionalLight.shadow.mapSize.height = 2048
  directionalLight.shadow.camera.near = 0.5
  directionalLight.shadow.camera.far = 50
  directionalLight.shadow.camera.left = -10
  directionalLight.shadow.camera.right = 10
  directionalLight.shadow.camera.top = 10
  directionalLight.shadow.camera.bottom = -10
  scene.add(directionalLight)

  // 补充光源 - 从右侧照射，提供更均匀的光照
  const fillLight = new THREE.DirectionalLight(0xffffff, 0.4)
  fillLight.position.set(-3, 3, 3)
  scene.add(fillLight)

  // 背光 - 增加轮廓，但强度较低
  const backLight = new THREE.DirectionalLight(0xffffff, 0.15)
  backLight.position.set(0, 2, -5)
  scene.add(backLight)
}

// 添加地面网格 - 参考截图的深色背景和浅色网格线
const addGround = () => {
  // 创建网格辅助线 - 使用深色背景配浅色网格线，增加网格密度
  const gridHelper = new THREE.GridHelper(20, 40, 0x666666, 0x666666)
  gridHelper.position.y = -0.01 // 稍微向下一点，避免z-fighting
  scene.add(gridHelper)

  // 添加半透明地面 - 使用更深的颜色
  const groundGeometry = new THREE.PlaneGeometry(20, 20)
  const groundMaterial = new THREE.MeshLambertMaterial({
    color: 0x1a1a1a, // 更深的背景色
    transparent: true,
    opacity: 0.3,
  })
  const ground = new THREE.Mesh(groundGeometry, groundMaterial)
  ground.rotation.x = -Math.PI / 2
  ground.receiveShadow = true
  scene.add(ground)
}

// 添加坐标轴辅助线 - 参考截图的样式
const addAxesHelper = () => {
  const axesHelper = new THREE.AxesHelper(3) // 稍微小一点的坐标轴
  scene.add(axesHelper)
}

//

//

// 加载3D模型
const loadModel = async () => {
  if (!props.modelUrl) {
    console.warn('ThreeDPreview: 没有提供模型URL')
    return
  }

  try {
    loading.value = true
    loadingProgress.value = 0
    error.value = ''

    // 清除现有模型
    if (model) {
      scene.remove(model)
      model = null
    }

    let loadedModel

    switch (props.modelType.toLowerCase()) {
      case 'glb':
      case 'gltf':
        loadedModel = await loadGLTF(props.modelUrl)
        break
      case 'obj':
        loadedModel = await loadOBJ(props.modelUrl)
        break
      case 'stl':
        loadedModel = await loadSTL(props.modelUrl)
        break
      default:
        throw new Error(`不支持的模型格式: ${props.modelType}`)
    }

    if (loadedModel) {
      model = loadedModel
      scene.add(model)

      // 计算模型边界并调整相机位置
      const box = new THREE.Box3().setFromObject(model)
      const size = box.getSize(new THREE.Vector3())
      const center = box.getCenter(new THREE.Vector3())

      // 调整模型位置到原点
      model.position.sub(center)

      // 计算并保存基础缩放值
      const maxDim = Math.max(size.x, size.y, size.z)
      baseScale = maxDim > 0 ? 2 / maxDim : 1

      // 应用初始缩放
      model.scale.setScalar(baseScale * scale.value)

      // 根据模型大小调整相机距离 - 保持截图中的俯视角度
      const cameraDistance = maxDim > 0 ? maxDim * 2 : 5

      // 设置相机位置为左上角俯视角度
      camera.position.set(cameraDistance * 0.6, cameraDistance * 0.6, cameraDistance * 0.6)
      camera.lookAt(0, 0, 0)

      if (controls) {
        controls.target.set(0, 0, 0)
        controls.update()
      }

      // 设置阴影和材质
      model.traverse((child) => {
        if (child.isMesh) {
          child.castShadow = true
          child.receiveShadow = true

          // 将模型材质改为亮色
          if (child.material) {
            const colorHex = modelColor.value.replace('#', '0x')
            // 如果是数组材质
            if (Array.isArray(child.material)) {
              child.material.forEach((mat) => {
                if (mat.color) {
                  mat.color.setHex(colorHex)
                }
              })
            } else {
              // 单个材质
              if (child.material.color) {
                child.material.color.setHex(colorHex)
              }
            }
          }
        }
      })
    } else {
      console.warn('ThreeDPreview: 模型加载返回空值')
    }

    loading.value = false
    loadingProgress.value = 100
  } catch (err) {
    console.error('ThreeDPreview: 加载3D模型失败:', err)
    error.value = `加载模型失败: ${err.message}`
    loading.value = false
    loadingProgress.value = 0
  }
}

// 加载GLTF/GLB模型
const loadGLTF = (url) => {
  return new Promise((resolve, reject) => {
    const loader = new GLTFLoader()
    loader.load(
      url,
      (gltf) => {
        const model = gltf.scene

        // 处理动画
        if (gltf.animations && gltf.animations.length > 0) {
          mixer = new THREE.AnimationMixer(model)
          const action = mixer.clipAction(gltf.animations[0])
          action.play()
        }

        resolve(model)
      },
      (xhr) => {
        if (xhr.total > 0) {
          const percent = (xhr.loaded / xhr.total) * 100
          loadingProgress.value = percent
        }
      },
      (error) => {
        console.error('ThreeDPreview: GLTF模型加载失败', error)
        reject(error)
      }
    )
  })
}

// 加载OBJ模型
const loadOBJ = (url) => {
  return new Promise((resolve, reject) => {
    const loader = new OBJLoader()
    loader.load(
      url,
      (obj) => {
        // 为OBJ模型添加默认材质
        obj.traverse((child) => {
          if (child.isMesh) {
            child.material = new THREE.MeshLambertMaterial({
              color: 0x888888,
            })
          }
        })
        resolve(obj)
      },
      undefined,
      reject
    )
  })
}

// 加载STL模型
const loadSTL = (url) => {
  return new Promise((resolve, reject) => {
    const loader = new STLLoader()
    loader.load(
      url,
      (geometry) => {
        const material = new THREE.MeshLambertMaterial({
          color: 0x888888,
        })
        const mesh = new THREE.Mesh(geometry, material)
        resolve(mesh)
      },
      undefined,
      reject
    )
  })
}

// 放大
const zoomIn = () => {
  if (scale.value < 3) {
    scale.value = Math.min(scale.value + 0.1, 3)
    updateScale(scale.value)
  }
}

// 缩小
const zoomOut = () => {
  if (scale.value > 0.1) {
    scale.value = Math.max(scale.value - 0.1, 0.1)
    updateScale(scale.value)
  }
}

// 更新缩放
const updateScale = (value) => {
  if (model) {
    model.scale.setScalar(baseScale * value)
    console.log('ThreeDPreview: 更新缩放', { value, baseScale, finalScale: baseScale * value })
  }
}

// 更新模型颜色
const updateModelColor = (color) => {
  if (model && color) {
    model.traverse((child) => {
      if (child.isMesh && child.material) {
        if (Array.isArray(child.material)) {
          child.material.forEach((mat) => {
            if (mat.color) {
              mat.color.setHex(color.replace('#', '0x'))
            }
          })
        } else {
          if (child.material.color) {
            child.material.color.setHex(color.replace('#', '0x'))
          }
        }
      }
    })
  }
}

//

// 重试加载
const retryLoad = () => {
  loadModel()
}

// 窗口大小变化处理
const onWindowResize = () => {
  if (!container.value || !camera || !renderer) return

  const containerRect = container.value.getBoundingClientRect()
  const width = Math.max(containerRect.width || 400, 400)
  const height = Math.max(containerRect.height || 300, 300)

  camera.aspect = width / height
  camera.updateProjectionMatrix()
  renderer.setSize(width, height)
}

// 渲染循环
const animate = () => {
  animationId = requestAnimationFrame(animate)

  if (controls) {
    controls.update()
  }

  if (mixer) {
    const delta = clock.getDelta()
    mixer.update(delta)
  }

  if (renderer && scene && camera) {
    renderer.render(scene, camera)
  }
}

// 清理资源
const cleanup = () => {
  if (animationId) {
    cancelAnimationFrame(animationId)
  }

  if (mixer) {
    mixer.stopAllAction()
    mixer.uncacheRoot(model)
  }

  if (renderer) {
    renderer.dispose()
  }

  if (container.value && renderer) {
    container.value.removeChild(renderer.domElement)
  }

  window.removeEventListener('resize', onWindowResize)
}

// 监听模型URL变化
watch(
  () => props.modelUrl,
  (newUrl) => {
    if (newUrl) {
      loadModel()
    }
  }
)

// 监听模型类型变化
watch(
  () => props.modelType,
  () => {
    if (props.modelUrl) {
      loadModel()
    }
  }
)

// 生命周期
onMounted(() => {
  // 使用nextTick确保DOM完全渲染
  nextTick(() => {
    // 延迟初始化，确保容器有正确的尺寸
    setTimeout(() => {
      initThreeJS()
      if (props.modelUrl) {
        loadModel()
      }
    }, 100)
  })
})

onUnmounted(() => {
  cleanup()
})
</script>

<style lang="scss">
.three-d-preview {
  position: relative;
  width: 100%;
  height: 100%;
}

.preview-container {
  width: 100%;
  height: 100%;
  min-height: 400px;
  position: relative;
  background: #f0f0f0;
  border-radius: 8px;
  overflow: hidden;
}

.control-panel {
  position: absolute;
  top: 20px;
  right: 20px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);

  .control-group {
    margin-bottom: 16px;

    label {
      display: block;
      margin-bottom: 8px;
      font-size: 14px;
      color: #333;
      font-weight: 500;
    }
  }

  .color-picker {
    display: flex;
    justify-content: center;

    .el-color-picker--small {
      width: 100% !important;
    }
  }

  .scale-controls {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;

    .scale-value {
      min-width: 40px;
      text-align: center;
      font-size: 14px;
      color: #333;
      font-weight: 500;
    }
  }
}

.loading-overlay,
.error-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(5px);
}

.loading-content,
.error-content {
  text-align: center;

  .el-icon {
    font-size: 32px;
    margin-bottom: 12px;

    &.is-loading {
      animation: rotate 1s linear infinite;
    }
  }

  p {
    margin: 0 0 16px 0;
    color: #666;
    font-size: 14px;
  }
}

.loading-progress {
  width: 200px;
  margin-top: 16px;

  .progress-bar {
    width: 100%;
    height: 4px;
    background: rgba(0, 0, 0, 0.1);
    border-radius: 2px;
    overflow: hidden;
    margin-bottom: 8px;

    .progress-fill {
      height: 100%;
      background: #409eff;
      border-radius: 2px;
      transition: width 0.3s ease;
    }
  }

  .progress-text {
    font-size: 12px;
    color: #666;
  }
}

.error-content {
  .el-icon {
    color: #f56c6c;
  }
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

// 响应式设计
@media (max-width: 768px) {
  .control-panel {
    position: relative;
    top: auto;
    right: auto;
    margin: 16px;
    border-radius: 8px;
  }
}
</style>
