<template>
  <div class="discover-page">
    <div class="discover-content">
      <!-- AI工具列表 -->
      <div class="tools-section">
        <div class="section-header">
          <h2 class="section-title">AI 创作工具</h2>
          <p class="section-subtitle">探索强大的AI创作能力</p>
        </div>

        <div class="tools-grid">
          <div
            v-for="tool in aiTools"
            :key="tool.key"
            class="tool-item"
            @click="navigateTo(tool.url)"
          >
            <div class="tool-card">
              <div class="tool-icon-wrapper">
                <div class="tool-icon" :style="{ background: tool.gradient }">
                  <i class="iconfont" :class="tool.icon"></i>
                </div>
                <div class="tool-badge" v-if="tool.badge">{{ tool.badge }}</div>
              </div>
              <div class="tool-content">
                <h3 class="tool-name">{{ tool.name }}</h3>
                <p class="tool-desc">{{ tool.desc }}</p>
                <div class="tool-meta">
                  <span class="tool-tag" v-if="tool.tag">{{ tool.tag }}</span>
                  <span class="tool-status" :class="tool.status">{{ tool.statusText }}</span>
                </div>
              </div>
              <div class="tool-arrow">
                <van-icon name="arrow" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// AI工具配置
const aiTools = ref([
  {
    key: 'mj',
    name: 'MJ绘画',
    desc: 'Midjourney AI绘画创作',
    icon: 'icon-mj',
    gradient: 'linear-gradient(135deg, #8B5CF6, #A855F7)',
    badge: '热门',
    tag: 'AI绘画',
    status: 'active',
    statusText: '可用',
    url: '/mobile/create?tab=mj',
  },
  {
    key: 'sd',
    name: 'SD绘画',
    desc: 'Stable Diffusion本地化',
    icon: 'icon-sd',
    gradient: 'linear-gradient(135deg, #06B6D4, #0891B2)',
    tag: 'AI绘画',
    status: 'active',
    statusText: '可用',
    url: '/mobile/create?tab=sd',
  },
  {
    key: 'dalle',
    name: 'DALL·E',
    desc: 'OpenAI图像生成',
    icon: 'icon-dalle',
    gradient: 'linear-gradient(135deg, #F59E0B, #D97706)',
    tag: 'AI绘画',
    status: 'active',
    statusText: '可用',
    url: '/mobile/create?tab=dalle',
  },
  {
    key: 'suno',
    name: '音乐创作',
    desc: 'AI音乐生成与编辑',
    icon: 'icon-mp3',
    gradient: 'linear-gradient(135deg, #EF4444, #DC2626)',
    badge: '新功能',
    tag: 'AI音乐',
    status: 'active',
    statusText: '可用',
    url: '/mobile/suno-create',
  },
  {
    key: 'video',
    name: '视频生成',
    desc: 'AI视频创作工具',
    icon: 'icon-video',
    gradient: 'linear-gradient(135deg, #10B981, #059669)',
    tag: 'AI视频',
    status: 'beta',
    statusText: '测试版',
    url: '/mobile/video-create',
  },
  {
    key: 'jimeng',
    name: '即梦AI',
    desc: '即梦AI绘画平台',
    icon: 'icon-jimeng',
    gradient: 'linear-gradient(135deg, #F97316, #EA580C)',
    tag: 'AI绘画',
    status: 'active',
    statusText: '可用',
    url: '/mobile/jimeng-create',
  },
  {
    key: 'imgWall',
    name: 'AI画廊',
    desc: 'AI绘画作品展示',
    icon: 'icon-image-list',
    gradient: 'linear-gradient(135deg, #3B82F6, #2563EB)',
    tag: 'AI展示',
    status: 'active',
    statusText: '可用',
    url: '/mobile/imgWall',
  },
  {
    key: 'apps',
    name: '应用中心',
    desc: '更多AI应用工具',
    icon: 'icon-app',
    gradient: 'linear-gradient(135deg, #EC4899, #DB2777)',
    tag: '应用',
    status: 'active',
    statusText: '可用',
    url: '/mobile/apps',
  },
])

// 导航处理
const navigateTo = (url) => {
  if (url.startsWith('http')) {
    window.open(url, '_blank')
  } else {
    router.push(url)
  }
}
</script>

<style lang="scss" scoped>
.discover-page {
  min-height: 100vh;
  background: var(--van-background);

  .discover-content {
    padding: 20px 16px;

    .tools-section {
      .section-header {
        text-align: center;
        margin-bottom: 32px;

        .section-title {
          font-size: 24px;
          font-weight: 700;
          color: var(--van-text-color);
          margin: 0 0 8px 0;
          background: linear-gradient(135deg, var(--van-primary-color), #8b5cf6);
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          background-clip: text;
        }

        .section-subtitle {
          font-size: 14px;
          color: var(--van-gray-6);
          margin: 0;
        }
      }

      .tools-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 16px;

        .tool-item {
          .tool-card {
            background: var(--van-cell-background);
            border-radius: 16px;
            padding: 20px;
            box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
            cursor: pointer;
            position: relative;
            overflow: hidden;

            &::before {
              content: '';
              position: absolute;
              top: 0;
              left: 0;
              right: 0;
              height: 3px;
              background: linear-gradient(90deg, var(--van-primary-color), #8b5cf6);
              transform: scaleX(0);
              transition: transform 0.3s ease;
            }

            &:hover {
              transform: translateY(-4px);
              box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);

              &::before {
                transform: scaleX(1);
              }

              .tool-icon {
                transform: scale(1.1);
              }
            }

            &:active {
              transform: translateY(-2px);
            }

            .tool-icon-wrapper {
              position: relative;
              margin-bottom: 16px;

              .tool-icon {
                width: 56px;
                height: 56px;
                border-radius: 16px;
                display: flex;
                align-items: center;
                justify-content: center;
                margin: 0 auto;
                transition: transform 0.3s ease;
                box-shadow: 0 8px 20px rgba(0, 0, 0, 0.15);

                .iconfont {
                  font-size: 28px;
                  color: white;
                }
              }

              .tool-badge {
                position: absolute;
                top: -8px;
                right: -8px;
                background: linear-gradient(135deg, #ef4444, #dc2626);
                color: white;
                font-size: 10px;
                font-weight: 600;
                padding: 4px 8px;
                border-radius: 12px;
                box-shadow: 0 2px 8px rgba(239, 68, 68, 0.3);
              }
            }

            .tool-content {
              text-align: center;
              margin-bottom: 16px;

              .tool-name {
                font-size: 16px;
                font-weight: 600;
                color: var(--van-text-color);
                margin: 0 0 6px 0;
              }

              .tool-desc {
                font-size: 12px;
                color: var(--van-gray-6);
                margin: 0 0 12px 0;
                line-height: 1.4;
              }

              .tool-meta {
                display: flex;
                align-items: center;
                justify-content: center;
                gap: 8px;

                .tool-tag {
                  background: var(--van-primary-color);
                  color: white;
                  font-size: 10px;
                  font-weight: 500;
                  padding: 2px 8px;
                  border-radius: 10px;
                }

                .tool-status {
                  font-size: 10px;
                  font-weight: 500;
                  padding: 2px 8px;
                  border-radius: 10px;

                  &.active {
                    background: #10b981;
                    color: white;
                  }

                  &.beta {
                    background: #f59e0b;
                    color: white;
                  }

                  &.maintenance {
                    background: #6b7280;
                    color: white;
                  }
                }
              }
            }

            .tool-arrow {
              position: absolute;
              top: 20px;
              right: 20px;
              color: var(--van-gray-4);
              transition: all 0.3s ease;

              .van-icon {
                font-size: 16px;
              }
            }

            &:hover .tool-arrow {
              color: var(--van-primary-color);
              transform: translateX(2px);
            }
          }
        }
      }
    }
  }
}

// 深色主题优化
:deep(.van-theme-dark) {
  .discover-page {
    .tool-card {
      box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);

      &:hover {
        box-shadow: 0 8px 30px rgba(0, 0, 0, 0.3);
      }

      .tool-icon {
        box-shadow: 0 8px 20px rgba(0, 0, 0, 0.3);
      }
    }
  }
}

// 响应式优化
@media (max-width: 375px) {
  .discover-page {
    .discover-content {
      padding: 16px 12px;

      .tools-section {
        .section-header {
          margin-bottom: 24px;

          .section-title {
            font-size: 22px;
          }
        }

        .tools-grid {
          gap: 12px;

          .tool-item .tool-card {
            padding: 16px;

            .tool-icon-wrapper .tool-icon {
              width: 48px;
              height: 48px;

              .iconfont {
                font-size: 24px;
              }
            }

            .tool-content {
              .tool-name {
                font-size: 15px;
              }

              .tool-desc {
                font-size: 11px;
              }
            }
          }
        }
      }
    }
  }
}
</style>
