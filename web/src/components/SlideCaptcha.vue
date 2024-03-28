<template>
  <div class="slide-captcha">
    <div class="bg-img">
      <el-image :src="backgroundImg" />
      <div :class="verifyMsgClass" v-if="checked !== 0">
        <span v-if="checked ===1">{{time}}s</span>
        {{verifyMsg}}
      </div>
      <div class="refresh" @click="emits('refresh')">
        <el-icon><Refresh /></el-icon>
      </div>
      <span class="block">
        <el-image :src="blockImg" :style="{left: blockLeft+'px'}" />
      </span>
    </div>

    <div class="verify">
      <div class="verify-bar-area">
        <span class="verify-msg">{{verifyText}}</span>

        <div :class="leftBarClass" :style="{width: leftBarWidth+'px'}">
          <div :class="blockClass" id="dragBlock"
               :style="{left: blockLeft+'px'}">
            <el-icon v-if="checked === 0"><ArrowRightBold /></el-icon>
            <el-icon v-if="checked === 1"><CircleCheckFilled /></el-icon>
            <el-icon v-if="checked === 2"><CircleCloseFilled /></el-icon>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
// eslint-disable-next-line no-undef
import {onMounted, ref, watch} from "vue";
import {ArrowRightBold, CircleCheckFilled, CircleCloseFilled, Refresh} from "@element-plus/icons-vue";

// eslint-disable-next-line no-undef
const props = defineProps({
  bgImg: String,
  bkImg: String,
  result: Number,
})

const verifyText = ref('向右滑动完成验证')
const verifyMsg = ref('')
const verifyMsgClass = ref("verify-text success")
const blockClass = ref('verify-move-block')
const leftBarClass = ref('verify-left-bar')
const backgroundImg = ref('')
const blockImg = ref('')
const leftBarWidth = ref(0)
const blockLeft = ref(0)
const checked = ref(0)
const time = ref('')

watch(() => props.bgImg, (newVal) => {
  backgroundImg.value = newVal;
});
watch(() => props.bkImg, (newVal) => {
  blockImg.value = newVal;
});
watch(() => props.result, (newVal) => {
  checked.value = newVal;
  if (newVal === 1) {
    verifyMsgClass.value = "verify-text success"
    blockClass.value = 'verify-move-block success'
    leftBarClass.value = 'verify-left-bar success'
    verifyMsg.value = '验证成功'
    setTimeout(() => emits('hide'), 1000)
  } else if (newVal ===2) {
    verifyMsgClass.value = "verify-text error"
    blockClass.value = 'verify-move-block error'
    leftBarClass.value = 'verify-left-bar error'
    verifyMsg.value = '验证失败'
    setTimeout(() => {
      reset()
      emits('refresh')
    }, 1000)
  } else {
    reset()
  }
});


// eslint-disable-next-line no-undef
const emits = defineEmits(['confirm','refresh','hide']);

let offsetX = 0, isDragging  = false
let start = 0
onMounted(() => {
  const dragBlock = document.getElementById('dragBlock');
  dragBlock.addEventListener('mousedown', (evt) => {
    blockClass.value = 'verify-move-block active'
    leftBarClass.value = 'verify-left-bar active'
    leftBarWidth.value = 32
    isDragging  = true
    verifyText.value = ""
    offsetX = evt.clientX
    start = new Date().getTime()
    evt.preventDefault();
  })

  document.body.addEventListener('mousemove',(evt) => {
    if (!isDragging) {
      return
    }
    const x = Math.max(evt.clientX - offsetX, 0)
    blockLeft.value = x;
    leftBarWidth.value = x + 32
  })

  document.body.addEventListener('mouseup', () => {
    if (!isDragging) {
      return
    }
    time.value = ((new Date().getTime() - start)/1000).toFixed(2)
    isDragging  = false
    emits('confirm', Math.floor(blockLeft.value))
  })

  // 触摸事件
  dragBlock.addEventListener('touchstart', function (e) {
    isDragging = true;
    blockClass.value = 'verify-move-block active'
    leftBarClass.value = 'verify-left-bar active'
    leftBarWidth.value = 32
    isDragging  = true
    verifyText.value = ""
    offsetX = e.touches[0].clientX - dragBlock.getBoundingClientRect().left;
    start = new Date().getTime()
    e.preventDefault();
  });

  document.addEventListener('touchmove', function (e) {
    if (!isDragging) {
      return
    }
    e.preventDefault();
    const x = Math.max(e.touches[0].clientX - offsetX, 0)
    blockLeft.value = x;
    leftBarWidth.value = x + 32
  });

  document.addEventListener('touchend', function () {
    if (!isDragging) {
      return
    }
    time.value = ((new Date().getTime() - start)/1000).toFixed(2)
    isDragging  = false
    emits('confirm', Math.floor(blockLeft.value))
  });
})


// 重置验证码
const reset = () => {
  blockClass.value = 'verify-move-block'
  leftBarClass.value = 'verify-left-bar'
  leftBarWidth.value = 0
  blockLeft.value = 0
  checked.value = 0
  verifyText.value = "向右滑动完成验证"
}
</script>

<style scoped lang="stylus">
@keyframes expandUp {
  0% {
    transform: scaleY(0);
  }
  100% {
    transform: scaleY(1);
  }
}

.slide-captcha {
  * {
    margin 0
    padding 0
  }

  .bg-img {
    position relative
    width 310px
    .verify-text {
      position absolute
      bottom 3px
      padding 5px 10px
      width 290px
      color #ffffff

      animation: expandUp 0.3s ease-in-out forwards;
      transform-origin: bottom center;
      transform: scaleY(0); /* 初始状态，元素高度为0 */
    }

    .verify-text.success {
      background-color rgba(92,184,92, 0.5)
    }

    .verify-text.error {
      background-color rgba(184,92,92, 0.5)
    }

    .refresh {
      position absolute
      right 5px
      top 5px
      font-size 20px
      cursor pointer
      color #ffffff
    }

    .block {
      position absolute
      top 0
      left 0
    }

  }

  .verify {
    .verify-bar-area {
      position relative
      border: 1px solid #dddddd
      overflow hidden
      height 34px

      .verify-msg {
        display flex
        line-height 32px
        width 100%
        justify-content center
      }

      .verify-left-bar {
        position absolute
        left 0
        top 0
        height 32px;

        .verify-move-block {
          position absolute
          width: 32px;
          height: 32px;
          background-color: rgb(255, 255, 255);
          border-top 1px solid #ffffff
          border-bottom 1px solid #ffffff
          border-right 1px solid #dddddd

          display flex
          justify-content center
          align-items center
          .el-icon {
            font-size 20px
            cursor pointer
          }
        }

        .verify-move-block.active {
          background #409eff
          color #ffffff
          border 1px solid #409eff
        }

        .verify-move-block.success {
          background #57AD57
          color #ffffff
          border 1px solid #57AD57
        }
        .verify-move-block.error {
          background #D9534F
          color #ffffff
          border 1px solid #D9534F
        }

      }

      .verify-left-bar.active {
        background-color #F0FFF0
        border 1px solid #409eff
      }
      .verify-left-bar.success {
        background-color #F0FFF0
        border 1px solid #57AD57
      }
      .verify-left-bar.error {
        background-color #F0FFF0
        border 1px solid #D9534F
      }
    }
  }
}
</style>