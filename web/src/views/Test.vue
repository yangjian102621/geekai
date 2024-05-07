<template>
  <div>
    <textarea v-model="value"/>
  </div>
  <svg ref="svgRef"/>
</template>

<script setup>
import {ref, onMounted, onUpdated} from 'vue';
import {Markmap} from 'markmap-view';
import {loadJS, loadCSS} from 'markmap-common';
import {Transformer} from 'markmap-lib';

const transformer = new Transformer();
const {scripts, styles} = transformer.getAssets();
loadCSS(styles);
loadJS(scripts);

const initValue = `# markmap

- beautiful
- useful
- easy
- interactive
`;

const value = ref(initValue);
const svgRef = ref(null);
let mm;

const update = () => {
  const {root} = transformer.transform(value.value);
  mm.setData(root);
  mm.fit();
};

onMounted(() => {
  mm = Markmap.create(svgRef.value);
  update();
});

onUpdated(update);
</script>
