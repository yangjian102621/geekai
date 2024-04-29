<template>
  <div class="index-page" :style="{height: winHeight+'px'}">
    <div class="bg"></div>
    <div class="menu-box">
      <el-menu
          mode="horizontal"
          :ellipsis="false"
      >
        <div class="menu-item">
          <el-image :src="logo" alt="Geek-AI"/>
          <div class="title">{{ title }}</div>
        </div>
        <div class="menu-item">
          <a href="https://ai.r9it.com/docs/install/" target="_blank">
            <el-button type="primary" round>
              <i class="iconfont icon-book"></i>
              <span>部署文档</span>
            </el-button>
          </a>

          <a href="https://github.com/yangjian102621/chatgpt-plus" target="_blank">
            <el-button type="success" round>
              <i class="iconfont icon-github"></i>
              <span>项目源码</span>
            </el-button>
          </a>
          <el-button @click="router.push('/login')" round>登录</el-button>
          <el-button @click="router.push('/register')" round>注册</el-button>
        </div>
      </el-menu>
    </div>
    <div class="content">
      <h1>欢迎使用 {{ title }}</h1>
      <p>{{ slogan }}</p>
      <el-button @click="router.push('/chat')" color="#ffffff" style="color:#007bff" :dark="false">
        <i class="iconfont icon-chat"></i>
        <span>AI 对话</span>
      </el-button>
      <el-button @click="router.push('/mj')" color="#C4CCFD" style="color:#424282" :dark="false">
        <i class="iconfont icon-mj"></i>
        <span>MJ 绘画</span>
      </el-button>

      <el-button @click="router.push('/sd')" color="#4AE6DF" style="color:#424282" :dark="false">
        <i class="iconfont icon-sd"></i>
        <span>SD 绘画</span>
      </el-button>
      <el-button @click="router.push('/xmind')" color="#FFFD55" style="color:#424282" :dark="false">
        <i class="iconfont icon-xmind"></i>
        <span>思维导图</span>
      </el-button>
      <!--      <div id="animation-container"></div>-->
    </div>

    <div class="footer">
      <footer-bar />
    </div>
  </div>
</template>

<script setup>

// import * as THREE from 'three';
import {onMounted, ref} from "vue";
import {useRouter} from "vue-router";
import FooterBar from "@/components/FooterBar.vue";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";
import {isMobile} from "@/utils/libs";

const router = useRouter()

if (isMobile()) {
  router.push("/mobile")
}

const title = ref("Geek-AI 创作系统")
const logo = ref("/images/logo.png")
const slogan = ref("我辈之人，先干为敬，陪您先把 AI 用起来")
// const size = Math.max(window.innerWidth * 0.5, window.innerHeight * 0.8)
const winHeight = window.innerHeight - 150

onMounted(() => {
  httpGet("/api/config/get?key=system").then(res => {
    title.value = res.data.title
    logo.value = res.data.logo
  }).catch(e => {
    ElMessage.error("获取系统配置失败：" + e.message)
  })
  init()
})

const init = () => {
  // // 创建场景
  // // 创建场景
  // const scene = new THREE.Scene();
  //
  // // 创建相机
  // const camera = new THREE.PerspectiveCamera(30, 1, 0.1, 1000);
  // camera.position.z = 3.88;
  //
  // // 创建渲染器
  // const renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true });
  // renderer.setSize(size, size);
  // renderer.setClearColor(0x000000, 0);
  // const container = document.getElementById('animation-container');
  // container.appendChild(renderer.domElement);
  //
  // // 加载地球纹理
  // const loader = new THREE.TextureLoader();
  // loader.load(
  //     '/images/land_ocean_ice_cloud_2048.jpg',
  //     function (texture) {
  //       // 创建地球球体
  //       const geometry = new THREE.SphereGeometry(1, 32, 32);
  //       const material = new THREE.MeshPhongMaterial({
  //         map: texture,
  //         bumpMap: texture, // 使用同一张纹理作为凹凸贴图
  //         bumpScale: 0.05, // 调整凹凸贴图的影响程度
  //         specularMap: texture, // 高光贴图
  //         specular: new THREE.Color('#01193B'), // 高光颜色
  //       });
  //       const earth = new THREE.Mesh(geometry, material);
  //       scene.add(earth);
  //
  //       // 添加环境光和点光源
  //       const ambientLight = new THREE.AmbientLight(0xffffff, 0.3);
  //       scene.add(ambientLight);
  //       const pointLight = new THREE.PointLight(0xffffff, 0.8);
  //       pointLight.position.set(5, 5, 5);
  //       scene.add(pointLight);
  //
  //       // 创建动画
  //       const animate = function () {
  //         requestAnimationFrame(animate);
  //
  //         // 使地球自转和公转
  //         earth.rotation.y += 0.0006;
  //
  //         renderer.render(scene, camera);
  //       };
  //
  //       // 执行动画
  //       animate();
  //     }
  // );
}
</script>

<style lang="stylus" scoped>
@import '@/assets/iconfont/iconfont.css'
.index-page {
  margin: 0
  overflow hidden
  color #ffffff
  display flex
  justify-content center
  align-items baseline
  padding-top 150px

  .bg {
    position absolute
    top 0
    left 0
    width 100vw
    height 100vh
    background-image url("~@/assets/img/ai-bg.jpg")
    //filter: blur(8px);
    background-size: cover;
    background-position: center;
  }

  .menu-box {
    position absolute
    top 0
    width 100%
    display flex

    .el-menu {
      padding 0 30px
      width 100%
      display flex
      justify-content space-between
      background none
      border none

      .menu-item {
        display flex
        padding 20px 0

        color #ffffff

        .title {
          font-size 24px
          padding 10px 10px 0 10px
        }

        .el-image {
          height 50px
        }

        .el-button {
          margin-left 10px

          span {
            margin-left 5px
          }
        }
      }
    }
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

    .el-button {
      padding: 25px 20px;
      font-size: 1.3rem;
      transition: all 0.3s ease;

      .iconfont {
        font-size 1.6rem
        margin-right 10px
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
