<template>
  <div class="page-jimeng">
    <!-- 左侧参数设置面板 -->
    <div class="params-panel">
      <h2>即梦AI</h2>
      
      <!-- 功能分类按钮组 -->
      <div class="category-buttons">
        <div class="category-label">
          <el-icon><Star /></el-icon>
          功能分类
        </div>
        <div class="category-grid">
          <div 
            v-for="category in store.categories"
            :key="category.key"
            :class="['category-btn', { active: store.activeCategory === category.key }]"
            @click="store.switchCategory(category.key)"
          >
            <div class="category-icon">
              <i :class="getCategoryIcon(category.key)"></i>
            </div>
            <div class="category-name">{{ category.name }}</div>
          </div>
        </div>
      </div>

      <!-- 功能开关 -->
      <div class="function-switch" v-if="store.activeCategory === 'image_generation' || store.activeCategory === 'video_generation'">
        <div class="switch-label">
          <el-icon><Switch /></el-icon>
          生成模式
        </div>
        <div class="switch-container">
          <div class="switch-info">
            <div class="switch-title">
              {{ store.useImageInput ? (store.activeCategory === 'image_generation' ? '图生图' : '图生视频') : (store.activeCategory === 'image_generation' ? '文生图' : '文生视频') }}
            </div>
            <div class="switch-desc">
              {{ store.useImageInput ? '使用图片作为输入' : '使用文字作为输入' }}
            </div>
          </div>
          <el-switch 
            v-model="store.useImageInput" 
            @change="store.switchInputMode"
            size="large"
          />
        </div>
      </div>

      <!-- 参数容器 -->
      <div class="params-container">
        <!-- 文生图 -->
        <div v-if="store.activeFunction === 'text_to_image'" class="function-panel">
          <div class="param-line pt">
            <span class="label">提示词：</span>
            <el-tooltip content="输入你想要的图片内容描述" placement="right">
              <el-icon><InfoFilled /></el-icon>
            </el-tooltip>
          </div>
          <div class="param-line">
            <el-input
              v-model="store.textToImageParams.prompt"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 5 }"
              placeholder="请输入图片描述，越详细越好"
              maxlength="2000"
              show-word-limit
            />
          </div>
          
          <div class="param-line pt">
            <span class="label">图片尺寸：</span>
          </div>
          <div class="param-line">
            <el-select v-model="store.textToImageParams.size" placeholder="选择尺寸">
              <el-option label="1328x1328 (正方形)" value="1328x1328" />
              <el-option label="1024x1024 (正方形)" value="1024x1024" />
              <el-option label="1024x768 (横版)" value="1024x768" />
              <el-option label="768x1024 (竖版)" value="768x1024" />
            </el-select>
          </div>
          
          <div class="item-group">
            <span class="label">创意度：</span>
            <el-slider v-model="store.textToImageParams.scale" :min="1" :max="10" :step="0.5" />
          </div>
          
          <div class="item-group">
            <span class="label">种子值：</span>
            <el-input-number v-model="store.textToImageParams.seed" :min="-1" :max="999999" size="small" />
          </div>
          
          <div class="item-group flex justify-between">
            <span class="label">智能优化提示词</span>
            <el-switch v-model="store.textToImageParams.use_pre_llm" size="small" />
          </div>
        </div>

        <!-- 图生图 -->
        <div v-if="store.activeFunction === 'image_to_image_portrait'" class="function-panel">
          <div class="param-line pt">
            <span class="label">上传图片：</span>
          </div>
          <div class="param-line">
            <ImageUpload v-model="store.imageToImageParams.image_input" />
          </div>
          
          <div class="param-line pt">
            <span class="label">提示词：</span>
          </div>
          <div class="param-line">
            <el-input
              v-model="store.imageToImageParams.prompt"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 5 }"
              placeholder="描述你想要的图片效果"
              maxlength="2000"
              show-word-limit
            />
          </div>
          
          <div class="param-line pt">
            <span class="label">图片尺寸：</span>
          </div>
          <div class="param-line">
            <el-select v-model="store.imageToImageParams.size" placeholder="选择尺寸">
              <el-option label="1328x1328 (正方形)" value="1328x1328" />
              <el-option label="1024x1024 (正方形)" value="1024x1024" />
              <el-option label="1024x768 (横版)" value="1024x768" />
              <el-option label="768x1024 (竖版)" value="768x1024" />
            </el-select>
          </div>
          
          <div class="item-group">
            <span class="label">GPEN强度：</span>
            <el-slider v-model="store.imageToImageParams.gpen" :min="0" :max="1" :step="0.1" />
          </div>
          
          <div class="item-group">
            <span class="label">肌肤质感：</span>
            <el-slider v-model="store.imageToImageParams.skin" :min="0" :max="1" :step="0.1" />
          </div>
          
          <div class="item-group">
            <span class="label">种子值：</span>
            <el-input-number v-model="store.imageToImageParams.seed" :min="-1" :max="999999" size="small" />
          </div>
        </div>

        <!-- 图像编辑 -->
        <div v-if="store.activeFunction === 'image_edit'" class="function-panel">
          <div class="param-line pt">
            <span class="label">上传图片：</span>
          </div>
          <div class="param-line">
            <ImageUpload v-model="store.imageEditParams.image_urls" :multiple="true" />
          </div>
          
          <div class="param-line pt">
            <span class="label">编辑提示词：</span>
          </div>
          <div class="param-line">
            <el-input
              v-model="store.imageEditParams.prompt"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 5 }"
              placeholder="描述你想要的编辑效果"
              maxlength="2000"
              show-word-limit
            />
          </div>
          
          <div class="item-group">
            <span class="label">编辑强度：</span>
            <el-slider v-model="store.imageEditParams.scale" :min="0" :max="1" :step="0.1" />
          </div>
          
          <div class="item-group">
            <span class="label">种子值：</span>
            <el-input-number v-model="store.imageEditParams.seed" :min="-1" :max="999999" size="small" />
          </div>
        </div>

        <!-- 图像特效 -->
        <div v-if="store.activeFunction === 'image_effects'" class="function-panel">
          <div class="param-line pt">
            <span class="label">上传图片：</span>
          </div>
          <div class="param-line">
            <ImageUpload v-model="store.imageEffectsParams.image_input1" />
          </div>
          
          <div class="param-line pt">
            <span class="label">特效模板：</span>
          </div>
          <div class="param-line">
            <el-select v-model="store.imageEffectsParams.template_id" placeholder="选择特效模板">
              <el-option label="经典特效" value="classic" />
              <el-option label="艺术风格" value="artistic" />
              <el-option label="现代科技" value="modern" />
            </el-select>
          </div>
          
          <div class="param-line pt">
            <span class="label">输出尺寸：</span>
          </div>
          <div class="param-line">
            <el-select v-model="store.imageEffectsParams.size" placeholder="选择尺寸">
              <el-option label="1328x1328 (正方形)" value="1328x1328" />
              <el-option label="1024x1024 (正方形)" value="1024x1024" />
              <el-option label="1024x768 (横版)" value="1024x768" />
              <el-option label="768x1024 (竖版)" value="768x1024" />
            </el-select>
          </div>
        </div>

        <!-- 文生视频 -->
        <div v-if="store.activeFunction === 'text_to_video'" class="function-panel">
          <div class="param-line pt">
            <span class="label">提示词：</span>
          </div>
          <div class="param-line">
            <el-input
              v-model="store.textToVideoParams.prompt"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 5 }"
              placeholder="描述你想要的视频内容"
              maxlength="2000"
              show-word-limit
            />
          </div>
          
          <div class="param-line pt">
            <span class="label">视频比例：</span>
          </div>
          <div class="param-line">
            <el-select v-model="store.textToVideoParams.aspect_ratio" placeholder="选择比例">
              <el-option label="16:9 (横版)" value="16:9" />
              <el-option label="9:16 (竖版)" value="9:16" />
              <el-option label="1:1 (正方形)" value="1:1" />
            </el-select>
          </div>
          
          <div class="item-group">
            <span class="label">种子值：</span>
            <el-input-number v-model="store.textToVideoParams.seed" :min="-1" :max="999999" size="small" />
          </div>
        </div>

        <!-- 图生视频 -->
        <div v-if="store.activeFunction === 'image_to_video'" class="function-panel">
          <div class="param-line pt">
            <span class="label">上传图片：</span>
          </div>
          <div class="param-line">
            <ImageUpload v-model="store.imageToVideoParams.image_urls" :multiple="true" />
          </div>
          
          <div class="param-line pt">
            <span class="label">提示词：</span>
          </div>
          <div class="param-line">
            <el-input
              v-model="store.imageToVideoParams.prompt"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 5 }"
              placeholder="描述你想要的视频效果"
              maxlength="2000"
              show-word-limit
            />
          </div>
          
          <div class="param-line pt">
            <span class="label">视频比例：</span>
          </div>
          <div class="param-line">
            <el-select v-model="store.imageToVideoParams.aspect_ratio" placeholder="选择比例">
              <el-option label="16:9 (横版)" value="16:9" />
              <el-option label="9:16 (竖版)" value="9:16" />
              <el-option label="1:1 (正方形)" value="1:1" />
            </el-select>
          </div>
          
          <div class="item-group">
            <span class="label">种子值：</span>
            <el-input-number v-model="store.imageToVideoParams.seed" :min="-1" :max="999999" size="small" />
          </div>
        </div>

        <!-- 算力显示 -->
        <div class="text-info">
          <el-tag type="primary">当前算力: {{ store.userPower }}</el-tag>
          <el-tag type="warning">消耗: {{ store.currentPowerCost }}</el-tag>
        </div>

        <!-- 提交按钮 -->
        <div class="submit-btn">
          <el-button 
            type="primary" 
            @click="store.submitTask" 
            :loading="store.submitting"
            :disabled="!store.isLogin || store.userPower < store.currentPowerCost"
            size="large"
          >
            立即生成 ({{ store.currentPowerCost }}<i class="iconfont icon-vip2"></i>)
          </el-button>
        </div>
      </div>
    </div>

    <!-- 右侧任务列表 -->
    <div class="main-content" v-loading="store.loading">
      <div class="works-header">
        <h2 class="h-title">你的作品</h2>
        <div class="filter-buttons">
          <el-button-group>
            <el-button
              :type="store.taskFilter === 'all' ? 'primary' : 'default'"
              @click="store.switchTaskFilter('all')"
              size="small"
            >
              全部
            </el-button>
            <el-button
              :type="store.taskFilter === 'image' ? 'primary' : 'default'"
              @click="store.switchTaskFilter('image')"
              size="small"
            >
              图片
            </el-button>
            <el-button
              :type="store.taskFilter === 'video' ? 'primary' : 'default'"
              @click="store.switchTaskFilter('video')"
              size="small"
            >
              视频
            </el-button>
          </el-button-group>
        </div>
      </div>

      <div class="task-list">
        <div class="list-box" v-if="!store.noData">
          <div v-for="item in store.currentList" :key="item.id" class="task-item">
            <div class="task-left">
              <div class="task-preview">
                <el-image
                  v-if="item.img_url"
                  :src="item.img_url"
                  fit="cover"
                  class="preview-image"
                />
                <video
                  v-else-if="item.video_url"
                  :src="item.video_url"
                  class="preview-video"
                  preload="metadata"
                />
                <div v-else class="preview-placeholder">
                  <el-icon><Picture /></el-icon>
                  <span>{{ store.getTaskStatusText(item.status) }}</span>
                </div>
              </div>
            </div>
            
            <div class="task-center">
              <div class="task-info">
                <el-tag size="small" :type="store.getStatusType(item.status)">
                  {{ store.getTaskStatusText(item.status) }}
                </el-tag>
                <el-tag size="small">{{ store.getFunctionName(item.type) }}</el-tag>
              </div>
              <div class="task-prompt">
                {{ store.substr(item.prompt, 200) }}
              </div>
              <div class="task-meta">
                <span>{{ dateFormat(item.created_at) }}</span>
                <span v-if="item.power">{{ item.power }}算力</span>
              </div>
            </div>
            
            <div class="task-right">
              <div class="task-actions">
                <el-button
                  v-if="item.status === 'failed'"
                  type="primary"
                  size="small"
                  @click="store.retryTask(item.id)"
                >
                  重试
                </el-button>
                <el-button
                  v-if="item.video_url || item.img_url"
                  type="default"
                  size="small"
                  @click="store.downloadFile(item)"
                >
                  下载
                </el-button>
                <el-button
                  v-if="item.video_url"
                  type="default"
                  size="small"
                  @click="store.playVideo(item)"
                >
                  播放
                </el-button>
                <el-button
                  type="danger"
                  size="small"
                  @click="store.removeJob(item)"
                >
                  删除
                </el-button>
              </div>
            </div>
          </div>
        </div>

        <el-empty
          v-else
          :image="store.nodata"
          description="暂无任务，快去创建吧！"
        />

        <div class="pagination" v-if="store.total > store.pageSize">
          <el-pagination
            background
            layout="total, prev, pager, next"
            :current-page="store.page"
            :page-size="store.pageSize"
            :total="store.total"
            @current-change="store.fetchData"
          />
        </div>
      </div>
    </div>

    <!-- 视频预览对话框 -->
    <el-dialog
      v-model="store.showDialog"
      title="视频预览"
      width="70%"
      center
    >
      <video
        :src="store.currentVideoUrl"
        controls
        style="width: 100%; max-height: 60vh;"
      >
        您的浏览器不支持视频播放
      </video>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useJimengStore } from '@/store/jimeng'
import { dateFormat } from '@/utils/libs'
import ImageUpload from '@/components/ImageUpload.vue'
import { InfoFilled, Star, Switch, Picture } from '@element-plus/icons-vue'

const store = useJimengStore()

// 获取分类图标
const getCategoryIcon = (category) => {
  const iconMap = {
    'image_generation': 'iconfont icon-image',
    'image_editing': 'iconfont icon-edit',
    'image_effects': 'iconfont icon-magic',
    'video_generation': 'iconfont icon-video'
  }
  return iconMap[category] || 'iconfont icon-image'
}

onMounted(() => {
  store.init()
})

onUnmounted(() => {
  store.cleanup()
})
</script>

<style lang="stylus" scoped>
.page-jimeng {
  display: flex;
  min-height: 100vh;
  background: var(--chat-bg);
  
  // 左侧参数面板
  .params-panel {
    min-width: 380px;
    max-width: 380px;
    margin: 10px;
    padding: 20px;
    border-radius: 12px;
    background: #ffffff;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
    color: #333;
    font-size: 14px;
    overflow: auto;
    
    h2 {
      font-weight: bold;
      font-size: 20px;
      text-align: center;
      color: #333;
      margin-bottom: 30px;
    }
    
    // 功能分类按钮组
    .category-buttons {
      margin-bottom: 25px;
      
      .category-label {
        display: flex;
        align-items: center;
        margin-bottom: 15px;
        font-size: 16px;
        font-weight: 600;
        color: #333;
        
        .el-icon {
          margin-right: 8px;
          color: #5865f2;
        }
      }
      
      .category-grid {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 10px;
        
        .category-btn {
          display: flex;
          flex-direction: column;
          align-items: center;
          padding: 15px 10px;
          border: 2px solid #f0f0f0;
          border-radius: 12px;
          cursor: pointer;
          transition: all 0.3s ease;
          background: #fafafa;
          
          &:hover {
            border-color: #5865f2;
            background: #f8f9ff;
            transform: translateY(-2px);
          }
          
          &.active {
            border-color: #5865f2;
            background: linear-gradient(135deg, #5865f2 0%, #7289da 100%);
            color: white;
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(88, 101, 242, 0.3);
          }
          
          .category-icon {
            font-size: 20px;
            margin-bottom: 8px;
          }
          
          .category-name {
            font-size: 12px;
            font-weight: 500;
          }
        }
      }
    }
    
    // 功能开关
    .function-switch {
      margin-bottom: 25px;
      
      .switch-label {
        display: flex;
        align-items: center;
        margin-bottom: 15px;
        font-size: 16px;
        font-weight: 600;
        color: #333;
        
        .el-icon {
          margin-right: 8px;
          color: #5865f2;
        }
      }
      
      .switch-container {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 15px;
        border: 1px solid #e0e0e0;
        border-radius: 10px;
        background: #f9f9f9;
        
        .switch-info {
          flex: 1;
          
          .switch-title {
            font-size: 14px;
            font-weight: 600;
            color: #333;
            margin-bottom: 4px;
          }
          
          .switch-desc {
            font-size: 12px;
            color: #666;
          }
        }
      }
    }
    
    // 参数容器
    .params-container {
      .function-panel {
        .param-line {
          margin-bottom: 15px;
          
          &.pt {
            margin-top: 20px;
          }
          
          .label {
            display: flex;
            align-items: center;
            margin-bottom: 8px;
            font-weight: 600;
            color: #333;
          }
        }
        
        .item-group {
          display: flex;
          align-items: center;
          margin-bottom: 15px;
          
          .label {
            margin-right: 15px;
            font-weight: 600;
            color: #333;
            min-width: 80px;
          }
        }
        
        .text-info {
          margin: 20px 0;
          padding: 15px;
          background: #f0f8ff;
          border-radius: 8px;
          border-left: 4px solid #5865f2;
        }
        
        .submit-btn {
          margin-top: 30px;
          
          .el-button {
            width: 100%;
            height: 50px;
            font-size: 16px;
            font-weight: 600;
          }
        }
      }
    }
  }
  
  // 右侧主要内容区域
  .main-content {
    flex: 1;
    padding: 20px;
    background: var(--chat-bg);
    color: var(--text-theme-color);
    
    .works-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 20px;
      
      .h-title {
        font-size: 24px;
        font-weight: 600;
        color: var(--text-theme-color);
        margin: 0;
      }
    }
    
    .task-list {
      .list-box {
        .task-item {
          display: flex;
          align-items: center;
          padding: 20px;
          margin-bottom: 15px;
          background: var(--card-bg);
          border-radius: 12px;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
          
          .task-left {
            margin-right: 20px;
            
            .task-preview {
              width: 120px;
              height: 90px;
              border-radius: 8px;
              overflow: hidden;
              background: #f0f0f0;
              display: flex;
              align-items: center;
              justify-content: center;
              
              .preview-image, .preview-video {
                width: 100%;
                height: 100%;
                object-fit: cover;
              }
              
              .preview-placeholder {
                display: flex;
                flex-direction: column;
                align-items: center;
                color: #999;
                font-size: 12px;
                
                .el-icon {
                  font-size: 24px;
                  margin-bottom: 5px;
                }
              }
            }
          }
          
          .task-center {
            flex: 1;
            
            .task-info {
              display: flex;
              gap: 8px;
              margin-bottom: 10px;
            }
            
            .task-prompt {
              font-size: 14px;
              color: var(--text-theme-color);
              margin-bottom: 8px;
              line-height: 1.4;
            }
            
            .task-meta {
              display: flex;
              gap: 15px;
              font-size: 12px;
              color: #999;
            }
          }
          
          .task-right {
            .task-actions {
              display: flex;
              gap: 8px;
              flex-wrap: wrap;
            }
          }
        }
      }
      
      .pagination {
        margin-top: 30px;
        display: flex;
        justify-content: center;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .page-jimeng {
    flex-direction: column;
    
    .params-panel {
      min-width: 100%;
      max-width: 100%;
      margin: 10px 0;
    }
    
    .main-content {
      padding: 15px;
    }
  }
}
</style>