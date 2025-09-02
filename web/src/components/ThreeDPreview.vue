<template>
  <div class="three-d-preview">
    <div ref="container" class="preview-container"></div>

    <!-- 控制面板 -->
    <div class="control-panel">
      <div class="control-group">
        <label>旋转速度</label>
        <el-slider
          v-model="rotationSpeed"
          :min="0"
          :max="0.1"
          :step="0.01"
          @change="updateRotationSpeed"
        />
      </div>

      <div class="control-group">
        <label>缩放</label>
        <el-slider v-model="scale" :min="0.1" :max="3" :step="0.1" @change="updateScale" />
      </div>

      <div class="control-buttons">
        <el-button size="small" @click="resetCamera">重置视角</el-button>
        <el-button size="small" @click="toggleAutoRotate">
          {{ autoRotate ? '停止旋转' : '自动旋转' }}
        </el-button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-overlay">
      <div class="loading-content">
        <el-icon class="is-loading"><Loading /></el-icon>
        <p>加载3D模型中...</p>
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
import { ElButton, ElIcon, ElSlider } from 'element-plus'
import * as THREE from 'three'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js'
import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader.js'
import { OBJLoader } from 'three/examples/jsm/loaders/OBJLoader.js'
import { STLLoader } from 'three/examples/jsm/loaders/STLLoader.js'
import { onMounted, onUnmounted, ref, watch } from 'vue'

// Props
const props = defineProps({
  modelUrl: {
    type: String,
    required: true,
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
const rotationSpeed = ref(0.02)
const scale = ref(1)
const autoRotate = ref(true)

// Three.js 相关变量
let scene, camera, renderer, controls, model, mixer, clock
let animationId

// 初始化Three.js场景
const initThreeJS = () => {
  if (!container.value) return

  // 创建场景
  scene = new THREE.Scene()
  scene.background = new THREE.Color(0xf0f0f0)

  // 创建相机
  const containerRect = container.value.getBoundingClientRect()
  camera = new THREE.PerspectiveCamera(75, containerRect.width / containerRect.height, 0.1, 1000)
  camera.position.set(0, 0, 5)

  // 创建渲染器
  renderer = new THREE.WebGLRenderer({ antialias: true })
  renderer.setSize(containerRect.width, containerRect.height)
  renderer.setPixelRatio(window.devicePixelRatio)
  renderer.shadowMap.enabled = true
  renderer.shadowMap.type = THREE.PCFSoftShadowMap

  // 添加到容器
  container.value.appendChild(renderer.domElement)

  // 创建控制器
  controls = new OrbitControls(camera, renderer.domElement)
  controls.enableDamping = true
  controls.dampingFactor = 0.05
  controls.autoRotate = autoRotate.value
  controls.autoRotateSpeed = rotationSpeed.value

  // 添加光源
  addLights()

  // 添加地面
  addGround()

  // 创建时钟
  clock = new THREE.Clock()

  // 开始渲染循环
  animate()

  // 监听窗口大小变化
  window.addEventListener('resize', onWindowResize)
}

// 添加光源
const addLights = () => {
  // 环境光
  const ambientLight = new THREE.AmbientLight(0x404040, 0.6)
  scene.add(ambientLight)

  // 方向光
  const directionalLight = new THREE.DirectionalLight(0xffffff, 0.8)
  directionalLight.position.set(10, 10, 5)
  directionalLight.castShadow = true
  directionalLight.shadow.mapSize.width = 2048
  directionalLight.shadow.mapSize.height = 2048
  scene.add(directionalLight)

  // 点光源
  const pointLight = new THREE.PointLight(0xffffff, 0.5)
  pointLight.position.set(-10, 10, -5)
  scene.add(pointLight)
}

// 添加地面
const addGround = () => {
  const groundGeometry = new THREE.PlaneGeometry(20, 20)
  const groundMaterial = new THREE.MeshLambertMaterial({
    color: 0xcccccc,
    transparent: true,
    opacity: 0.3,
  })
  const ground = new THREE.Mesh(groundGeometry, groundMaterial)
  ground.rotation.x = -Math.PI / 2
  ground.receiveShadow = true
  scene.add(ground)
}

// 加载3D模型
const loadModel = async () => {
  if (!props.modelUrl) return

  try {
    loading.value = true
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

      // 调整模型位置和大小
      centerModel()
      fitCameraToModel()

      // 设置阴影
      model.traverse((child) => {
        if (child.isMesh) {
          child.castShadow = true
          child.receiveShadow = true
        }
      })
    }

    loading.value = false
  } catch (err) {
    console.error('加载3D模型失败:', err)
    error.value = `加载模型失败: ${err.message}`
    loading.value = false
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
      undefined,
      reject
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

// 居中模型
const centerModel = () => {
  if (!model) return

  const box = new THREE.Box3().setFromObject(model)
  const center = box.getCenter(new THREE.Vector3())
  const size = box.getSize(new THREE.Vector3())

  // 居中
  model.position.sub(center)

  // 调整缩放
  const maxDim = Math.max(size.x, size.y, size.z)
  const scale = 2 / maxDim
  model.scale.setScalar(scale * props.scale)
}

// 调整相机以适应模型
const fitCameraToModel = () => {
  if (!model) return

  const box = new THREE.Box3().setFromObject(model)
  const size = box.getSize(new THREE.Vector3())
  const center = box.getCenter(new THREE.Vector3())

  const maxDim = Math.max(size.x, size.y, size.z)
  const fov = camera.fov * (Math.PI / 180)
  let cameraZ = Math.abs(maxDim / 2 / Math.tan(fov / 2))

  camera.position.set(center.x, center.y, center.z + cameraZ)
  camera.lookAt(center)

  controls.target.copy(center)
  controls.update()
}

// 更新旋转速度
const updateRotationSpeed = (value) => {
  if (controls) {
    controls.autoRotateSpeed = value
  }
}

// 更新缩放
const updateScale = (value) => {
  if (model) {
    const box = new THREE.Box3().setFromObject(model)
    const size = box.getSize(new THREE.Vector3())
    const maxDim = Math.max(size.x, size.y, size.z)
    const baseScale = 2 / maxDim
    model.scale.setScalar(baseScale * value)
  }
}

// 重置相机
const resetCamera = () => {
  if (camera && model) {
    fitCameraToModel()
  }
}

// 切换自动旋转
const toggleAutoRotate = () => {
  autoRotate.value = !autoRotate.value
  if (controls) {
    controls.autoRotate = autoRotate.value
  }
}

// 重试加载
const retryLoad = () => {
  loadModel()
}

// 窗口大小变化处理
const onWindowResize = () => {
  if (!container.value || !camera || !renderer) return

  const containerRect = container.value.getBoundingClientRect()
  camera.aspect = containerRect.width / containerRect.height
  camera.updateProjectionMatrix()
  renderer.setSize(containerRect.width, containerRect.height)
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
  initThreeJS()
  if (props.modelUrl) {
    loadModel()
  }
})

onUnmounted(() => {
  cleanup()
})
</script>

<style lang="scss" scoped>
.three-d-preview {
  position: relative;
  width: 100%;
  height: 100%;
}

.preview-container {
  width: 100%;
  height: 100%;
  position: relative;
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

  .control-buttons {
    display: flex;
    gap: 8px;

    .el-button {
      flex: 1;
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
