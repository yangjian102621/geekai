<template>
  <div class="container">
    <div class="wave">
      <div class="bar"></div>
      <div class="bar"></div>
      <div class="bar"></div>
      <div class="bar"></div>
      <div class="bar"></div>
      <div class="bar"></div>
      <div class="bar"></div>
      <div class="bar"></div>
      <div class="bar"></div>
      <div class="bar"></div>
    </div>
    <div class="text">正在生成歌曲</div>
  </div>
</template>

<style scoped lang="stylus">
.container {
  display: flex;
  flex-flow column
  justify-content: center;
  align-items: center;
  margin: 0;
  font-family: 'Arial', sans-serif;
  color: #e1e1e1;

  .wave {
    display: flex;
    justify-content: center;
    align-items: flex-end;
    height: 30px;
    margin-bottom: 5px

    .bar {
      width: 8px;
      margin: 0 2px;
      background-color: #919191;
      animation: wave 1.5s infinite;
    }

    .bar:nth-child(1) {
      animation-delay: 0s;
    }

    .bar:nth-child(2) {
      animation-delay: 0.1s;
    }

    .bar:nth-child(3) {
      animation-delay: 0.2s;
    }

    .bar:nth-child(4) {
      animation-delay: 0.3s;
    }

    .bar:nth-child(5) {
      animation-delay: 0.4s;
    }

    .bar:nth-child(6) {
      animation-delay: 0.5s;
    }

    .bar:nth-child(7) {
      animation-delay: 0.6s;
    }

    .bar:nth-child(8) {
      animation-delay: 0.7s;
    }

    .bar:nth-child(9) {
      animation-delay: 0.8s;
    }

    .bar:nth-child(10) {
      animation-delay: 0.9s;
    }
  }
}

@keyframes wave {
  0%, 100% {
    height:   10px;
  }
  50% {
    height: 30px;
  }
}

.text {
  font-size: 14px;
}
</style>