<template>
  <div>
   {{data}}
  </div>
</template>

<script setup>
import {ref, onMounted, onUpdated} from 'vue';
import {Markmap} from 'markmap-view';
import {loadJS, loadCSS} from 'markmap-common';
import {Transformer} from 'markmap-lib';
import {httpPost} from "@/utils/http";

const data=ref("")
httpPost("/api/test/sse",{
  "message":"你是什么模型",
  "user_id":123
}).then(res=>{
  // const source = new EventSource("http://localhost:5678/api/test/sse");
  // source.onmessage = function(event) {
  //   console.log(event.data)
  // };
})
</script>
