<template>
  <div class="tools-page">
    <van-nav-bar title="AI 工具" left-arrow @click-left="router.back()" fixed />

    <div class="tools-content">
      <!-- 工具分类 -->
      <van-tabs v-model:active="activeCategory" @change="onCategoryChange" sticky :offset-top="46">
        <van-tab title="全部" name="all" />
        <van-tab title="办公工具" name="office" />
        <van-tab title="创意工具" name="creative" />
        <van-tab title="学习工具" name="study" />
        <van-tab title="生活工具" name="life" />
      </van-tabs>

      <!-- 工具列表 -->
      <div class="tools-list">
        <div 
          v-for="tool in filteredTools" 
          :key="tool.key"
          class="tool-item"
          @click="openTool(tool)"
        >
          <div class="tool-header">
            <div class="tool-icon" :style="{ backgroundColor: tool.color }">
              <i class="iconfont" :class="tool.icon"></i>
            </div>
            <div class="tool-info">
              <div class="tool-name">{{ tool.name }}</div>
              <div class="tool-desc">{{ tool.desc }}</div>
            </div>
            <div class="tool-status">
              <van-tag :type="tool.status === 'available' ? 'success' : 'warning'" size="medium">
                {{ tool.status === 'available' ? '可用' : '开发中' }}
              </van-tag>
            </div>
          </div>
          
          <div class="tool-features" v-if="tool.features">
            <van-tag 
              v-for="feature in tool.features" 
              :key="feature"
              size="small"
              plain
              class="feature-tag"
            >
              {{ feature }}
            </van-tag>
          </div>
          
          <div class="tool-stats" v-if="tool.stats">
            <div class="stat-item">
              <span class="stat-label">使用次数：</span>
              <span class="stat-value">{{ tool.stats.usageCount }}次</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">好评率：</span>
              <span class="stat-value">{{ tool.stats.rating }}%</span>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <van-empty v-if="filteredTools.length === 0" description="该分类暂无工具" />
      </div>

      <!-- 推荐工具 -->
      <div class="recommend-section" v-if="activeCategory === 'all'">
        <h3 class="section-title">推荐工具</h3>
        <van-swipe :autoplay="3000" class="recommend-swipe">
          <van-swipe-item v-for="tool in recommendTools" :key="tool.key">
            <div class="recommend-card" @click="openTool(tool)">
              <div class="recommend-bg" :style="{ backgroundColor: tool.color }">
                <i class="iconfont" :class="tool.icon"></i>
              </div>
              <div class="recommend-content">
                <h4 class="recommend-title">{{ tool.name }}</h4>
                <p class="recommend-desc">{{ tool.desc }}</p>
                <van-button size="small" type="primary" plain round>
                  立即使用
                </van-button>
              </div>
            </div>
          </van-swipe-item>
        </van-swipe>
      </div>
    </div>

    <!-- 工具详情弹窗 -->
    <van-action-sheet v-model:show="showToolDetail" :title="selectedTool?.name">
      <div class="tool-detail" v-if="selectedTool">
        <div class="detail-header">
          <div class="detail-icon" :style="{ backgroundColor: selectedTool.color }">
            <i class="iconfont" :class="selectedTool.icon"></i>
          </div>
          <div class="detail-info">
            <h3 class="detail-name">{{ selectedTool.name }}</h3>
            <p class="detail-desc">{{ selectedTool.fullDesc || selectedTool.desc }}</p>
          </div>
        </div>
        
        <div class="detail-features" v-if="selectedTool.detailFeatures">
          <h4 class="features-title">功能特点</h4>
          <ul class="features-list">
            <li v-for="feature in selectedTool.detailFeatures" :key="feature">
              <van-icon name="checked" color="#07c160" />
              {{ feature }}
            </li>
          </ul>
        </div>
        
        <div class="detail-usage" v-if="selectedTool.usage">
          <h4 class="usage-title">使用说明</h4>
          <p class="usage-text">{{ selectedTool.usage }}</p>
        </div>
        
        <div class="detail-actions">
          <van-button 
            type="primary" 
            size="large" 
            round 
            block
            :disabled="selectedTool.status !== 'available'"
            @click="useTool(selectedTool)"
          >
            {{ selectedTool.status === 'available' ? '开始使用' : '开发中' }}
          </van-button>
        </div>
      </div>
    </van-action-sheet>

    <!-- 思维导图工具 -->
    <van-action-sheet v-model:show="showMindMap" title="思维导图" :close-on-click-overlay="false">
      <div class="mindmap-container">
        <div class="mindmap-toolbar">
          <van-button size="small" @click="createNewMap">新建</van-button>
          <van-button size="small" @click="saveMap">保存</van-button>
          <van-button size="small" @click="exportMap">导出</van-button>
          <van-button size="small" @click="closeMindMap">关闭</van-button>
        </div>
        <div class="mindmap-canvas" ref="mindmapCanvas">
          <!-- 这里会渲染思维导图 -->
          <div class="canvas-placeholder">
            <i class="iconfont icon-mind"></i>
            <p>思维导图工具</p>
            <p class="placeholder-desc">功能开发中，敬请期待</p>
          </div>
        </div>
      </div>
    </van-action-sheet>
  </div>
</template>

<script setup>
import { showNotify } from 'vant'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const activeCategory = ref('all')
const selectedTool = ref(null)
const showToolDetail = ref(false)
const showMindMap = ref(false)
const mindmapCanvas = ref()

// 工具列表配置
const tools = ref([
  {
    key: 'mindmap',
    name: '思维导图',
    desc: '智能生成思维导图，整理思路更清晰',
    fullDesc: '基于AI技术的智能思维导图生成工具，可以根据文本内容自动生成结构化的思维导图，支持多种导出格式。',
    icon: 'icon-mind',
    color: '#3B82F6',
    category: 'office',
    status: 'available',
    features: ['自动生成', '多种模板', '智能布局'],
    detailFeatures: [
      '支持文本自动转思维导图',
      '提供多种精美模板',
      '智能节点布局算法',
      '支持导出多种格式',
      '支持在线协作编辑'
    ],
    usage: '输入您的文本内容，AI会自动分析并生成对应的思维导图结构。您可以对生成的导图进行编辑、美化和导出。',
    stats: {
      usageCount: 1256,
      rating: 96
    }
  },
  {
    key: 'summary',
    name: '文档总结',
    desc: '快速提取文档要点，生成精准摘要',
    fullDesc: '智能文档总结工具，能够快速分析长文档并提取关键信息，生成简洁明了的摘要。',
    icon: 'icon-doc',
    color: '#10B981',
    category: 'office',
    status: 'available',
    features: ['关键词提取', '智能摘要', '多语言支持'],
    detailFeatures: [
      '支持多种文档格式',
      '智能关键词提取',
      '可控制摘要长度',
      '支持批量处理',
      '多语言文档支持'
    ],
    usage: '上传或粘贴文档内容，选择摘要长度和类型，AI会自动生成文档摘要。',
    stats: {
      usageCount: 2341,
      rating: 94
    }
  },
  {
    key: 'translation',
    name: '智能翻译',
    desc: '高质量多语言翻译，支持专业术语',
    fullDesc: '基于先进AI模型的多语言翻译工具，支持100+语言互译，特别适合专业文档翻译。',
    icon: 'icon-translate',
    color: '#8B5CF6',
    category: 'office',
    status: 'available',
    features: ['100+语言', '专业术语', '上下文理解'],
    detailFeatures: [
      '支持100多种语言互译',
      '专业术语库支持',
      '上下文语境理解',
      '批量文档翻译',
      '翻译质量评估'
    ],
    usage: '选择源语言和目标语言，输入需要翻译的内容，AI会提供高质量的翻译结果。',
    stats: {
      usageCount: 5678,
      rating: 98
    }
  },
  {
    key: 'poster',
    name: '海报设计',
    desc: '一键生成专业海报，多种风格可选',
    fullDesc: 'AI驱动的海报设计工具，提供丰富的模板和素材，轻松制作专业级海报。',
    icon: 'icon-design',
    color: '#F59E0B',
    category: 'creative',
    status: 'available',
    features: ['模板丰富', '一键生成', '高清输出'],
    detailFeatures: [
      '500+精美模板',
      '智能配色方案',
      '自动排版布局',
      '高清无水印导出',
      '支持自定义尺寸'
    ],
    usage: '选择海报类型和风格，输入文案内容，AI会自动生成专业海报设计。',
    stats: {
      usageCount: 3456,
      rating: 95
    }
  },
  {
    key: 'logo',
    name: 'Logo 设计',
    desc: 'AI 生成独特Logo，商用级品质',
    fullDesc: '专业的AI Logo设计工具，根据您的品牌理念生成独特的Logo设计方案。',
    icon: 'icon-logo',
    color: '#EF4444',
    category: 'creative',
    status: 'available',
    features: ['品牌风格', '矢量格式', '商用授权'],
    detailFeatures: [
      '多种设计风格选择',
      '矢量格式输出',
      '商用版权授权',
      '配色方案推荐',
      '标准化尺寸规范'
    ],
    usage: '描述您的品牌特点和期望风格，AI会生成多个Logo设计方案供您选择。',
    stats: {
      usageCount: 2234,
      rating: 93
    }
  },
  {
    key: 'study-plan',
    name: '学习计划',
    desc: '个性化学习路径规划，提升学习效率',
    fullDesc: '基于AI的个性化学习计划制定工具，根据您的学习目标和时间安排制定最优学习路径。',
    icon: 'icon-study',
    color: '#06B6D4',
    category: 'study',
    status: 'available',
    features: ['个性化', '进度跟踪', '智能调整'],
    detailFeatures: [
      '个性化学习路径',
      '学习进度跟踪',
      '智能计划调整',
      '学习效果评估',
      '多领域知识覆盖'
    ],
    usage: '输入您的学习目标、可用时间和当前水平，AI会为您制定详细的学习计划。',
    stats: {
      usageCount: 1890,
      rating: 97
    }
  },
  {
    key: 'recipe',
    name: '智能食谱',
    desc: '根据食材推荐美食，营养搭配建议',
    fullDesc: '智能食谱推荐系统，根据现有食材推荐美食制作方法，提供营养搭配建议。',
    icon: 'icon-food',
    color: '#F97316',
    category: 'life',
    status: 'development',
    features: ['食材识别', '营养分析', '制作指导'],
    detailFeatures: [
      '食材智能识别',
      '营养成分分析',
      '详细制作步骤',
      '口味偏好适配',
      '热量控制建议'
    ],
    usage: '拍照或输入现有食材，AI会推荐适合的菜谱并提供详细制作指导。',
    stats: {
      usageCount: 567,
      rating: 89
    }
  },
  {
    key: 'workout',
    name: '运动计划',
    desc: '定制化健身方案，科学训练指导',
    fullDesc: '个性化运动健身计划制定工具，根据身体状况和目标制定科学的训练方案。',
    icon: 'icon-sport',
    color: '#EC4899',
    category: 'life',
    status: 'development',
    features: ['个性定制', '科学指导', '进度跟踪'],
    detailFeatures: [
      '个性化训练计划',
      '科学运动指导',
      '训练进度跟踪',
      '饮食建议搭配',
      '健康数据分析'
    ],
    usage: '输入您的身体状况、运动目标和时间安排，AI会制定适合的运动计划。',
    stats: {
      usageCount: 234,
      rating: 91
    }
  }
])

// 推荐工具（取前3个可用的）
const recommendTools = computed(() => {
  return tools.value.filter(tool => tool.status === 'available').slice(0, 3)
})

// 根据分类筛选工具
const filteredTools = computed(() => {
  if (activeCategory.value === 'all') {
    return tools.value
  }
  return tools.value.filter(tool => tool.category === activeCategory.value)
})

onMounted(() => {
  // 检查URL参数，如果有指定工具则直接打开
  const urlParams = new URLSearchParams(window.location.search)
  const toolKey = urlParams.get('tool')
  if (toolKey) {
    const tool = tools.value.find(t => t.key === toolKey)
    if (tool) {
      openTool(tool)
    }
  }
})

// 分类切换
const onCategoryChange = (category) => {
  // 可以在这里添加数据加载逻辑
}

// 打开工具
const openTool = (tool) => {
  selectedTool.value = tool
  
  // 特殊工具直接打开对应界面
  if (tool.key === 'mindmap') {
    showMindMap.value = true
  } else {
    showToolDetail.value = true
  }
}

// 使用工具
const useTool = (tool) => {
  showToolDetail.value = false
  
  if (tool.status !== 'available') {
    showNotify({ type: 'warning', message: '该工具还在开发中，敬请期待' })
    return
  }
  
  // 根据工具类型跳转到对应页面或打开功能界面
  switch (tool.key) {
    case 'mindmap':
      showMindMap.value = true
      break
    case 'summary':
    case 'translation':
    case 'poster':
    case 'logo':
    case 'study-plan':
      showNotify({ type: 'primary', message: `正在启动${tool.name}工具...` })
      // 这里可以跳转到具体的工具页面
      break
    default:
      showNotify({ type: 'warning', message: '功能开发中' })
  }
}

// 思维导图相关方法
const createNewMap = () => {
  showNotify({ type: 'primary', message: '创建新的思维导图' })
}

const saveMap = () => {
  showNotify({ type: 'success', message: '思维导图已保存' })
}

const exportMap = () => {
  showNotify({ type: 'primary', message: '导出思维导图' })
}

const closeMindMap = () => {
  showMindMap.value = false
}
</script>

<style lang="scss" scoped>
.tools-page {
  min-height: 100vh;
  background: var(--van-background);
  
  .tools-content {
    padding-top: 46px;
    
    :deep(.van-tabs__nav) {
      background: var(--van-background);
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
    }
    
    .tools-list {
      padding: 16px;
      
      .tool-item {
        background: var(--van-cell-background);
        border-radius: 12px;
        padding: 16px;
        margin-bottom: 16px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
        cursor: pointer;
        transition: all 0.3s ease;
        
        &:active {
          transform: scale(0.98);
        }
        
        .tool-header {
          display: flex;
          align-items: center;
          margin-bottom: 12px;
          
          .tool-icon {
            width: 44px;
            height: 44px;
            border-radius: 12px;
            display: flex;
            align-items: center;
            justify-content: center;
            margin-right: 12px;
            
            .iconfont {
              font-size: 22px;
              color: white;
            }
          }
          
          .tool-info {
            flex: 1;
            
            .tool-name {
              font-size: 16px;
              font-weight: 600;
              color: var(--van-text-color);
              margin-bottom: 4px;
            }
            
            .tool-desc {
              font-size: 13px;
              color: var(--van-gray-6);
              line-height: 1.4;
            }
          }
          
          .tool-status {
            margin-left: 8px;
          }
        }
        
        .tool-features {
          display: flex;
          flex-wrap: wrap;
          gap: 8px;
          margin-bottom: 12px;
          
          .feature-tag {
            font-size: 11px;
          }
        }
        
        .tool-stats {
          display: flex;
          gap: 16px;
          font-size: 12px;
          color: var(--van-gray-6);
          
          .stat-label {
            margin-right: 4px;
          }
          
          .stat-value {
            color: var(--van-text-color);
            font-weight: 500;
          }
        }
      }
    }
    
    .recommend-section {
      padding: 0 16px 16px;
      
      .section-title {
        font-size: 18px;
        font-weight: 600;
        color: var(--van-text-color);
        margin: 0 0 16px 0;
      }
      
      .recommend-swipe {
        height: 160px;
        border-radius: 12px;
        overflow: hidden;
        
        .recommend-card {
          height: 100%;
          position: relative;
          display: flex;
          align-items: center;
          padding: 20px;
          background: linear-gradient(135deg, var(--van-primary-color), #8B5CF6);
          cursor: pointer;
          
          .recommend-bg {
            position: absolute;
            top: 16px;
            right: 16px;
            width: 80px;
            height: 80px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            background: rgba(255, 255, 255, 0.2);
            
            .iconfont {
              font-size: 40px;
              color: rgba(255, 255, 255, 0.8);
            }
          }
          
          .recommend-content {
            flex: 1;
            color: white;
            
            .recommend-title {
              font-size: 20px;
              font-weight: 700;
              margin: 0 0 8px 0;
            }
            
            .recommend-desc {
              font-size: 14px;
              opacity: 0.9;
              margin: 0 0 16px 0;
              line-height: 1.4;
            }
          }
        }
      }
    }
  }
  
  .tool-detail {
    padding: 20px;
    
    .detail-header {
      display: flex;
      margin-bottom: 20px;
      
      .detail-icon {
        width: 60px;
        height: 60px;
        border-radius: 16px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 16px;
        flex-shrink: 0;
        
        .iconfont {
          font-size: 28px;
          color: white;
        }
      }
      
      .detail-info {
        flex: 1;
        
        .detail-name {
          font-size: 20px;
          font-weight: 600;
          color: var(--van-text-color);
          margin: 0 0 8px 0;
        }
        
        .detail-desc {
          font-size: 14px;
          color: var(--van-gray-6);
          line-height: 1.5;
          margin: 0;
        }
      }
    }
    
    .detail-features,
    .detail-usage {
      margin-bottom: 20px;
      
      .features-title,
      .usage-title {
        font-size: 16px;
        font-weight: 600;
        color: var(--van-text-color);
        margin: 0 0 12px 0;
      }
      
      .features-list {
        padding: 0;
        margin: 0;
        list-style: none;
        
        li {
          display: flex;
          align-items: center;
          font-size: 14px;
          color: var(--van-text-color);
          margin-bottom: 8px;
          
          .van-icon {
            margin-right: 8px;
          }
        }
      }
      
      .usage-text {
        font-size: 14px;
        color: var(--van-gray-6);
        line-height: 1.5;
        margin: 0;
      }
    }
    
    .detail-actions {
      margin-top: 20px;
    }
  }
  
  .mindmap-container {
    height: 80vh;
    display: flex;
    flex-direction: column;
    
    .mindmap-toolbar {
      display: flex;
      gap: 8px;
      padding: 12px 16px;
      background: var(--van-background-2);
      border-bottom: 1px solid var(--van-border-color);
    }
    
    .mindmap-canvas {
      flex: 1;
      position: relative;
      background: var(--van-background);
      
      .canvas-placeholder {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        text-align: center;
        color: var(--van-gray-6);
        
        .iconfont {
          font-size: 48px;
          margin-bottom: 16px;
          color: var(--van-gray-5);
        }
        
        p {
          margin: 0 0 8px 0;
          font-size: 16px;
          
          &.placeholder-desc {
            font-size: 14px;
            opacity: 0.8;
          }
        }
      }
    }
  }
}

// 深色主题优化
:deep(.van-theme-dark) {
  .tools-page {
    .tool-item {
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
    }
  }
}
</style>