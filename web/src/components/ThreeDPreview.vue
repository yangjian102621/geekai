<template>
  <div class="three-d-preview">
    <div ref="container" class="preview-container"></div>

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
import { Loading, Warning } from '@element-plus/icons-vue'
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
const modelType = computed(() => {
  if (props.modelType) {
    return props.modelType.toLowerCase()
  }
  // 从模型URL中获取类型
  if (props.modelUrl) {
    const url = new URL(props.modelUrl)
    return url.pathname.split('.').pop()
  }
  return 'glb'
})

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
  scene.background = new THREE.Color(0x2d2d2d) // 深灰色背景，匹配截图效果

  // 获取容器尺寸，完全自适应父容器
  const containerRect = container.value.getBoundingClientRect()
  const width = containerRect.width || 400
  const height = containerRect.height || 300

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
  renderer.shadowMap.enabled = false // 禁用阴影
  renderer.outputColorSpace = THREE.SRGBColorSpace
  // 提升曝光度让模型更加高亮
  renderer.toneMapping = THREE.ACESFilmicToneMapping
  renderer.toneMappingExposure = 2.2

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

// 添加光源 - 高亮显示模型，无阴影效果
const addLights = () => {
  // 强环境光 - 提供整体高亮照明
  const ambientLight = new THREE.AmbientLight(0xffffff, 1.0)
  scene.add(ambientLight)

  // 主方向光 - 从前上方照射，高亮度无阴影
  const directionalLight = new THREE.DirectionalLight(0xffffff, 1.8)
  directionalLight.position.set(5, 8, 5)
  directionalLight.castShadow = false
  scene.add(directionalLight)

  // 补充光源 - 从左侧照射，填充光照
  const fillLight = new THREE.DirectionalLight(0xffffff, 1.2)
  fillLight.position.set(-5, 4, 3)
  fillLight.castShadow = false
  scene.add(fillLight)

  // 背景光 - 从背后照射，增加轮廓高亮
  const rimLight = new THREE.DirectionalLight(0xffffff, 1.0)
  rimLight.position.set(0, 3, -5)
  rimLight.castShadow = false
  scene.add(rimLight)

  // 顶部光源 - 增加顶部高亮
  const topLight = new THREE.DirectionalLight(0xffffff, 0.8)
  topLight.position.set(0, 10, 0)
  topLight.castShadow = false
  scene.add(topLight)
}

// 添加地面网格 - 简洁网格，无阴影
const addGround = () => {
  // 创建网格辅助线 - 使用深色线条
  const gridHelper = new THREE.GridHelper(20, 20, 0x555555, 0x555555)
  gridHelper.position.y = 0
  scene.add(gridHelper)

  // 简单透明地面平面
  const groundGeometry = new THREE.PlaneGeometry(20, 20)
  const groundMaterial = new THREE.MeshBasicMaterial({
    color: 0x404040,
    transparent: true,
    opacity: 0.1,
  })
  const ground = new THREE.Mesh(groundGeometry, groundMaterial)
  ground.rotation.x = -Math.PI / 2
  ground.position.y = -0.01
  scene.add(ground)
}

// 添加坐标轴辅助线 - 匹配截图样式
const addAxesHelper = () => {
  const axesHelper = new THREE.AxesHelper(2)
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

    switch (modelType.value) {
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
        throw new Error(`不支持的模型格式: ${modelType.value}`)
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
      model.scale.setScalar(baseScale)

      // 根据模型大小调整相机距离 - 保持截图中的俯视角度
      const cameraDistance = maxDim > 0 ? maxDim * 2 : 5

      // 设置相机位置 - 匹配截图中的正面稍俯视角度
      camera.position.set(cameraDistance * 0.3, cameraDistance * 0.4, cameraDistance * 1.2)
      camera.lookAt(0, 0, 0)

      if (controls) {
        controls.target.set(0, 0, 0)
        controls.update()
      }

      // 移除阴影设置，让模型高亮显示
      model.traverse((child) => {
        if (child.isMesh) {
          child.castShadow = false
          child.receiveShadow = false
          // 如果材质支持，增加发光效果
          if (child.material) {
            child.material.emissive = new THREE.Color(0x111111) // 轻微发光
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

//

// 重试加载
const retryLoad = () => {
  loadModel()
}

// 窗口大小变化处理
const onWindowResize = () => {
  if (!container.value || !camera || !renderer) return

  const containerRect = container.value.getBoundingClientRect()
  const width = containerRect.width || 400
  const height = containerRect.height || 300

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
  display: flex;
  flex-direction: column;
}

.preview-container {
  width: 100%;
  height: 100%;
  position: relative;
  background: #2d2d2d;
  border-radius: 8px;
  overflow: hidden;
  // 移除min-height限制，让高度完全自适应

  // 确保在弹窗中能正确填充
  canvas {
    width: 100% !important;
    height: 100% !important;
    display: block;
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
</style>
