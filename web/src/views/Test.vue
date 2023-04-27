<template>
  <div style="padding: 20px;" v-html="content" id="content"></div>
</template>

<script>
import {defineComponent, nextTick} from "vue"
import hl from 'highlight.js'
import 'highlight.js/styles/a11y-dark.css'

export default defineComponent({
  name: 'TestPage',
  data() {
    return {
      content: "测试页面",
    }
  },
  mounted() {

    let md = require('markdown-it')();
    this.content = md.render("```\n" +
        "const socket = new WebSocket('ws://localhost:8080');\n" +
        "\n" +
        "// 连接成功\n" +
        "socket.addEventListener('open', event => {\n" +
        "  console.log('WebSocket 连接成功！');\n" +
        "});\n" +
        "\n" +
        "// 接收消息\n" +
        "socket.addEventListener('message', event => {\n" +
        "  console.log('收到消息：', event.data);\n" +
        "});\n" +
        "\n" +
        "// 发送消息\n" +
        "socket.send('Hello, WebSocket!');\n" +
        "\n" +
        "```\n" +
        "\n" +
        "\n" +
        "以上代码创建了一个 WebSocket 连接，并在连接成功后输出一条提示信息。当收到消息时，会在控制台打印该消息。同时还演示了如何发送消息。在实际应用中，不同的框架和库可能会提供不同的 WebSocket 实现，代码可能会有所区别。");
    nextTick(() => {
      const blocks = document.getElementById('content').querySelectorAll('pre code');
      console.log(blocks)
      blocks.forEach((block) => {
        hl.highlightBlock(block)
      })
    })
  }
})
</script>
