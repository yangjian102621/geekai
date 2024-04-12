<template>
  <div class="index-page" :style="{height: winHeight+'px'}">
    <div class="content">
      <h1>欢迎使用 {{ title }}</h1>
      <p>{{slogan}}</p>
      <button class="btn" @click="router.push('/chat')">立即使用</button>

      <div id="animation-container"></div>
    </div>

    <div class="footer">
      <footer-bar />
    </div>
  </div>
</template>

<script setup>

import * as THREE from 'three';
import {onMounted, ref} from "vue";
import {useRouter} from "vue-router";
import FooterBar from "@/components/FooterBar.vue";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";

const router = useRouter()

const title = ref("Geek-AI 创作系统")
const slogan = ref("我辈之人，先干为敬，陪您先把 AI 用起来")
const size = Math.max(window.innerWidth * 0.5, window.innerHeight * 0.8)
const winHeight = window.innerHeight - 150

onMounted(() => {
  httpGet("/api/config/get?key=system").then(res => {
    title.value = res.data['title']
  }).catch(e => {
    ElMessage.error("获取系统配置失败：" + e.message)
  })
  init()
})

const init = () => {
  // 创建场景
  // 创建场景
  const scene = new THREE.Scene();

  // 创建相机
  const camera = new THREE.PerspectiveCamera(30, 1, 0.1, 1000);
  camera.position.z = 3.88;

  // 创建渲染器
  const renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true });
  renderer.setSize(size, size);
  renderer.setClearColor(0x000000, 0);
  const container = document.getElementById('animation-container');
  container.appendChild(renderer.domElement);

  // 加载地球纹理
  const loader = new THREE.TextureLoader();
  loader.load(
      '/images/land_ocean_ice_cloud_2048.jpg',
      function (texture) {
        // 创建地球球体
        const geometry = new THREE.SphereGeometry(1, 32, 32);
        const material = new THREE.MeshPhongMaterial({
          map: texture,
          bumpMap: texture, // 使用同一张纹理作为凹凸贴图
          bumpScale: 0.05, // 调整凹凸贴图的影响程度
          specularMap: texture, // 高光贴图
          specular: new THREE.Color('#007bff'), // 高光颜色
        });
        const earth = new THREE.Mesh(geometry, material);
        scene.add(earth);

        // 添加环境光和点光源
        const ambientLight = new THREE.AmbientLight(0xffffff, 0.3);
        scene.add(ambientLight);
        const pointLight = new THREE.PointLight(0xffffff, 0.8);
        pointLight.position.set(5, 5, 5);
        scene.add(pointLight);

        // 创建动画
        const animate = function () {
          requestAnimationFrame(animate);

          // 使地球自转和公转
          earth.rotation.y += 0.001;

          renderer.render(scene, camera);
        };

        // 执行动画
        animate();
      }
  );
}
</script>

<style lang="stylus" scoped>
.index-page {
  margin: 0
  background-color #007bff /* 科技蓝色背景 */
  overflow hidden
  color #ffffff
  display flex
  justify-content center
  align-items baseline
  padding-top 150px

  .container {
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    color: #fff;
  }
  .content {
    text-align: center;
    position relative
    h1 {
      font-size: 5rem;
      margin-bottom: 1rem;
    }
    p {
      font-size: 1.5rem;
      margin-bottom: 2rem;
    }

    .btn {
      padding: 10px 20px;
      background-color: #fff;
      color: #007bff;
      border: none;
      border-radius: 5px;
      font-size: 1.3rem;
      cursor: pointer;
      transition: all 0.3s ease;

      &:hover {
        background-color: #e6e6e6;
      }
    }

    #animation-container {
      display flex
      justify-content center
      width 100%
      height: 300px;
      position: absolute;
      top: 350px

    }
  }

  .footer {
    .el-link__inner {
      color #ffffff
    }
  }

}
</style>
