<template>

  <div class="wg-cap-wrap" :style="{width: width}">
    <div class="wg-cap-wrap__header">
      <span>请在下图<em>依次</em>点击：</span>
      <img class="wg-cap-wrap__thumb" v-if="thumbBase64Code" :src="thumbBase64Code" alt=" ">
    </div>
    <div class="wg-cap-wrap__body">
      <img class="wg-cap-wrap__picture" v-if="imageBase64Code" :src="imageBase64Code" alt=" "
           @click="handleClickPos($event)">
      <img class="wg-cap-wrap__loading"
           src="data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz4KPHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zOnhsaW5rPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hsaW5rIiBzdHlsZT0ibWFyZ2luOiBhdXRvOyBiYWNrZ3JvdW5kOiByZ2JhKDI0MSwgMjQyLCAyNDMsIDApOyBkaXNwbGF5OiBibG9jazsgc2hhcGUtcmVuZGVyaW5nOiBhdXRvOyIgd2lkdGg9IjY0cHgiIGhlaWdodD0iNjRweCIgdmlld0JveD0iMCAwIDEwMCAxMDAiIHByZXNlcnZlQXNwZWN0UmF0aW89InhNaWRZTWlkIj4KICA8Y2lyY2xlIGN4PSI1MCIgY3k9IjM2LjgxMDEiIHI9IjEzIiBmaWxsPSIjM2U3Y2ZmIj4KICAgIDxhbmltYXRlIGF0dHJpYnV0ZU5hbWU9ImN5IiBkdXI9IjFzIiByZXBlYXRDb3VudD0iaW5kZWZpbml0ZSIgY2FsY01vZGU9InNwbGluZSIga2V5U3BsaW5lcz0iMC40NSAwIDAuOSAwLjU1OzAgMC40NSAwLjU1IDAuOSIga2V5VGltZXM9IjA7MC41OzEiIHZhbHVlcz0iMjM7Nzc7MjMiPjwvYW5pbWF0ZT4KICA8L2NpcmNsZT4KPC9zdmc+"
           alt="正在加载中...">
      <div v-for="(dot, key) in dots" :key="key" class="wg-cap-wrap__dot" :style="`top: ${dot.y}px; left:${dot.x}px;`">
        <span>{{ dot.index }}</span>
      </div>
    </div>
    <div class="wg-cap-wrap__footer">
      <div class="wg-cap-wrap__ico">
        <img @click="handleCloseEvent"
             src="data:image/svg+xml;base64,PHN2ZyB0PSIxNjI2NjE0NDM5NDIzIiBjbGFzcz0iaWNvbiIgdmlld0JveD0iMCAwIDEwMjQgMTAyNCIgdmVyc2lvbj0iMS4xIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHAtaWQ9Ijg2NzUiIHdpZHRoPSIyMDAiIGhlaWdodD0iMjAwIj48cGF0aCBkPSJNNTEyIDIzLjI3MjcyN2MyNjkuOTE3MDkxIDAgNDg4LjcyNzI3MyAyMTguODEwMTgyIDQ4OC43MjcyNzMgNDg4LjcyNzI3M2E0ODYuNjMyNzI3IDQ4Ni42MzI3MjcgMCAwIDEtODQuOTQ1NDU1IDI3NS40MDk0NTUgNDYuNTQ1NDU1IDQ2LjU0NTQ1NSAwIDAgMS03Ni44NDY1NDUtNTIuNTI2NTQ2QTM5My41NDE4MTggMzkzLjU0MTgxOCAwIDAgMCA5MDcuNjM2MzY0IDUxMmMwLTIxOC41MDc2MzYtMTc3LjEyODcyNy0zOTUuNjM2MzY0LTM5NS42MzYzNjQtMzk1LjYzNjM2NFMxMTYuMzYzNjM2IDI5My40OTIzNjQgMTE2LjM2MzYzNiA1MTJzMTc3LjEyODcyNyAzOTUuNjM2MzY0IDM5NS42MzYzNjQgMzk1LjYzNjM2NGEzOTUuMTcwOTA5IDM5NS4xNzA5MDkgMCAwIDAgMTI1LjQ0LTIwLjI5MzgxOSA0Ni41NDU0NTUgNDYuNTQ1NDU1IDAgMCAxIDI5LjQ4NjU0NSA4OC4yOTY3MjhBNDg4LjI2MTgxOCA0ODguMjYxODE4IDAgMCAxIDUxMiAxMDAwLjcyNzI3M0MyNDIuMDgyOTA5IDEwMDAuNzI3MjczIDIzLjI3MjcyNyA3ODEuOTE3MDkxIDIzLjI3MjcyNyA1MTJTMjQyLjA4MjkwOSAyMy4yNzI3MjcgNTEyIDIzLjI3MjcyN3ogbS0xMTUuMiAzMDcuNzEyTDUxMiA0NDYuMTM4MTgybDExNS4yLTExNS4yYTQ2LjU0NTQ1NSA0Ni41NDU0NTUgMCAxIDEgNjUuODE1MjczIDY1Ljg2MTgxOEw1NzcuODYxODE4IDUxMmwxMTUuMiAxMTUuMmE0Ni41NDU0NTUgNDYuNTQ1NDU1IDAgMSAxLTY1Ljg2MTgxOCA2NS44MTUyNzNMNTEyIDU3Ny44NjE4MThsLTExNS4yIDExNS4yYTQ2LjU0NTQ1NSA0Ni41NDU0NTUgMCAxIDEtNjUuODE1MjczLTY1Ljg2MTgxOEw0NDYuMTM4MTgyIDUxMmwtMTE1LjItMTE1LjJhNDYuNTQ1NDU1IDQ2LjU0NTQ1NSAwIDEgMSA2NS44NjE4MTgtNjUuODE1MjczeiIgcC1pZD0iODY3NiIgZmlsbD0iIzcwNzA3MCI+PC9wYXRoPjwvc3ZnPg=="
             alt="关闭">
        <img @click="handleRefreshEvent"
             src="data:image/svg+xml;base64,PHN2ZyB0PSIxNjI2NjE0NDk5NjM4IiBjbGFzcz0iaWNvbiIgdmlld0JveD0iMCAwIDEwMjQgMTAyNCIgdmVyc2lvbj0iMS4xIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHAtaWQ9IjEzNjAiIHdpZHRoPSIyMDAiIGhlaWdodD0iMjAwIj48cGF0aCBkPSJNMTg3LjQ1NiA0MjUuMDI0YTMzNiAzMzYgMCAwIDAgMzY4LjM4NCA0MjAuMjI0IDQ4IDQ4IDAgMCAxIDEyLjU0NCA5NS4xNjggNDMyIDQzMiAwIDAgMS00NzMuNjY0LTU0MC4xNmwtNTcuMjgtMTUuMzZhMTIuOCAxMi44IDAgMCAxLTYuMjcyLTIwLjkyOGwxNTkuMTY4LTE3OS40NTZhMTIuOCAxMi44IDAgMCAxIDIyLjE0NCA1Ljg4OGw0OC4wNjQgMjM1LjA3MmExMi44IDEyLjggMCAwIDEtMTUuODA4IDE0LjkxMmwtNTcuMjgtMTUuMzZ6TTgzNi40OCA1OTkuMDRhMzM2IDMzNiAwIDAgMC0zNjguMzg0LTQyMC4yMjQgNDggNDggMCAxIDEtMTIuNTQ0LTk1LjE2OCA0MzIgNDMyIDAgMCAxIDQ3My42NjQgNTQwLjE2bDU3LjI4IDE1LjM2YTEyLjggMTIuOCAwIDAgMSA2LjI3MiAyMC45MjhsLTE1OS4xNjggMTc5LjQ1NmExMi44IDEyLjggMCAwIDEtMjIuMTQ0LTUuODg4bC00OC4wNjQtMjM1LjA3MmExMi44IDEyLjggMCAwIDEgMTUuODA4LTE0LjkxMmw1Ny4yOCAxNS4zNnoiIGZpbGw9IiM3MDcwNzAiIHAtaWQ9IjEzNjEiPjwvcGF0aD48L3N2Zz4="
             alt="刷新">
      </div>
      <div class="wg-cap-wrap__btn">
        <el-button type="primary" @click="handleConfirmEvent">确认</el-button>
      </div>
    </div>
  </div>

</template>

<script>
export default {
  name: 'CaptchaPlus',
  mounted() {
    this.$emit('refresh')
  },
  props: {
    value: Boolean,
    width: {
      type: String,
      default: '300px'
    },
    calcPosType: {
      type: String,
      default: 'dom',
      validator: value => ['dom', 'screen'].includes(value)
    },
    maxDot: {
      type: Number,
      default: 5
      // validator: value => value > 10
    },
    imageBase64: String,
    thumbBase64: String
  },
  data() {
    return {
      dots: [],
      imageBase64Code: '',
      thumbBase64Code: ''
    }
  },
  watch: {
    value() {
      this.dots = []
      this.imageBase64Code = ''
      this.thumbBase64Code = ''
    },
    imageBase64(val) {
      this.dots = []
      this.imageBase64Code = val
    },
    thumbBase64(val) {
      this.dots = []
      this.thumbBase64Code = val
    }
  },
  methods: {
    /**
     * @Description: 处理关闭事件
     */
    handleCloseEvent() {
      this.$emit('close')
      // this.dots = []
      // this.imageBase64Code = ''
      // this.thumbBase64Code = ''
    },
    /**
     * @Description: 处理刷新事件
     */
    handleRefreshEvent() {
      this.dots = []
      this.$emit('refresh')
    },
    /**
     * @Description: 处理确认事件
     */
    handleConfirmEvent() {
      this.$emit('confirm', this.dots)
    },
    /**
     * @Description: 处理dot
     * @param ev
     */
    handleClickPos(ev) {
      if (this.dots.length >= this.maxDot) {
        return
      }
      const e = ev || window.event
      e.preventDefault()
      const dom = e.currentTarget

      const {domX, domY} = this.getDomXY(dom)
      // ===============================================
      // @notice 如 getDomXY 不准确可尝试使用 calcLocationLeft 或 calcLocationTop
      // const domX = this.calcLocationLeft(dom)
      // const domY = this.calcLocationTop(dom)
      // ===============================================
      let mouseX = (navigator.vendor === 'Netscape') ? e.pageX : e.x + document.body.offsetTop
      let mouseY = (navigator.vendor === 'Netscape') ? e.pageY : e.y + document.body.offsetTop

      if (this.calcPosType === 'screen') {
        mouseX = (navigator.vendor === 'Netscape') ? e.clientX : e.x
        mouseY = (navigator.vendor === 'Netscape') ? e.clientY : e.y
      }

      // 计算点击的相对位置
      const xPos = mouseX - domX
      const yPos = mouseY - domY

      // 转整形
      const xp = parseInt(xPos.toString())
      const yp = parseInt(yPos.toString())

      // 减去点的一半
      this.dots.push({
        x: xp - 11,
        y: yp - 11,
        index: this.dots.length + 1
      })
      return false
    },
    /**
     * @Description: 找到元素的屏幕位置
     * @param el
     */
    calcLocationLeft(el) {
      let tmp = el.offsetLeft
      let val = el.offsetParent
      while (val != null) {
        tmp += val.offsetLeft
        val = val.offsetParent
      }
      return tmp
    },
    /**
     * @Description: 找到元素的屏幕位置
     * @param el
     */
    calcLocationTop(el) {
      let tmp = el.offsetTop
      let val = el.offsetParent
      while (val != null) {
        tmp += val.offsetTop
        val = val.offsetParent
      }
      return tmp
    },
    /**
     * @Description: 找到元素的屏幕位置
     * @param dom
     */
    getDomXY(dom) {
      let x = 0
      let y = 0
      if (dom.getBoundingClientRect) {
        let box = dom.getBoundingClientRect();
        let D = document.documentElement;
        x = box.left + Math.max(D.scrollLeft, document.body.scrollLeft) - D.clientLeft;
        y = box.top + Math.max(D.scrollTop, document.body.scrollTop) - D.clientTop
      } else {
        while (dom !== document.body) {
          x += dom.offsetLeft
          y += dom.offsetTop
          dom = dom.offsetParent
        }
      }
      return {
        domX: x,
        domY: y
      }
    }
  }
}
</script>

<style scoped lang="stylus">
.wg-cap-wrap {
  background: #ffffff;

  -webkit-border-radius: 10px;
  -moz-border-radius: 10px;
  border-radius: 10px;

  -webkit-touch-callout: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;

  .wg-cap-wrap__header {
    height: 50px;
    width: 100%;
    font-size: 15px;

    display: -webkit-box;
    display: -webkit-flex;
    display: -ms-flexbox;
    display: flex;
    -webkit-box-align: center;
    -webkit-align-items: center;
    -ms-flex-align: center;
    align-items: center;

    span {
      padding-right: 5px;

      em {
        padding: 0 3px;
        font-weight: bold;
        color: #3e7cff;
        font-style: normal;
      }
    }

    .wg-cap-wrap__image {
      -webkit-border-radius: 5px;
      -moz-border-radius: 5px;
      border-radius: 5px;
      overflow: hidden;
      text-align: center;
      line-height: 1;
    }

    .wg-cap-wrap__thumb {
      min-width: 150px;
      text-align: center;
      line-height: 1;
      max-height: 100%;
    }

    .wg-cap-wrap__thumb.wg-cap-wrap__hidden {
      display: none;
    }

  }

  .wg-cap-wrap__body {
    position: relative;
    display: -webkit-box;
    display: -moz-box;
    display: -ms-flexbox;
    display: -webkit-flex;
    display: flex;
    background: #34383e;
    margin: auto;
    -webkit-border-radius: 5px;
    -moz-border-radius: 5px;
    border-radius: 5px;
    overflow: hidden;

    .wg-cap-wrap__picture {
      position: relative;
      z-index: 10;
      width: 100%;
    }

    .wg-cap-wrap__picture.wg-cap-wrap__hidden {
      display: none;
    }

    .wg-cap-wrap__loading {
      position: absolute;
      z-index: 9;
      top: 50%;
      left: 50%;
      width: 68px;
      height: 68px;
      margin-left: -34px;
      margin-top: -34px;
      line-height: 68px;
      text-align: center;
    }

    .wg-cap-wrap__dot {
      position: absolute;
      z-index: 10;
      width: 22px;
      height: 22px;
      color: #cedffe;
      background: #3e7cff;
      border: 2px solid #f7f9fb;
      line-height: 20px;
      text-align: center;
      -webkit-border-radius: 22px;
      -moz-border-radius: 22px;
      border-radius: 22px;
      cursor: default;
    }
  }

  .wg-cap-wrap__footer {
    width: 100%;
    height: 40px;
    color: #34383e;
    display: -webkit-box;
    display: -webkit-flex;
    display: -ms-flexbox;
    display: flex;
    -webkit-box-align: center;
    -webkit-align-items: center;
    -ms-flex-align: center;
    align-items: center;
    padding-top: 15px;

    .wg-cap-wrap__ico {
      flex: 1;

      img {
        width: 24px;
        height: 24px;
        color: #34383e;
        margin: 0 5px;
        cursor: pointer;
      }
    }

    .wg-cap-wrap__btn {
      display flex
      width: 120px;
      justify-content right
    }
  }
}
</style>
