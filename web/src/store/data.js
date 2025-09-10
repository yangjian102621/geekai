import central_orbit from '@/assets/img/jimeng/central_orbit.webp'
import clockwise_swivel from '@/assets/img/jimeng/clockwise_swivel.webp'
import counterclockwise_swivel from '@/assets/img/jimeng/counterclockwise_swivel.webp'
import crane_push from '@/assets/img/jimeng/crane_push.webp'
import dynamic_orbit from '@/assets/img/jimeng/dynamic_orbit.webp'
import handheld from '@/assets/img/jimeng/handheld.webp'
import hitchcock_dolly_in from '@/assets/img/jimeng/hitchcock_dolly_in.webp'
import hitchcock_dolly_out from '@/assets/img/jimeng/hitchcock_dolly_out.webp'
import quick_pull_back from '@/assets/img/jimeng/quick_pull_back.webp'
import rapid_push_pull from '@/assets/img/jimeng/rapid_push_pull.webp'
import robo_arm from '@/assets/img/jimeng/robo_arm.webp'

export const JimengParams = {
  image: [
    {
      name: '图片 2.1 文生图',
      version: '2.1',
      label: '平面绘感强，可生成文字海报',
      key: 'jimeng_high_aes_general_v21_L',
      params: [
        {
          name: 'prompt',
          label: '提示词',
          type: 'textarea',
          showWordLimit: true,
          maxlength: 800,
          autosize: { minRows: 3, maxRows: 5 },
          required: true,
          placeholder: '请输入提示词',
          info: '用于生成图像的提示词 ，中英文均可输入',
        },

        {
          name: 'size',
          type: 'select',
          required: true,
          placeholder: '请选择尺寸',
          label: '图片尺寸',
          options: [
            {
              label: '21:9 (1195 * 512)',
              value: '1195x512',
            },
            {
              label: '16:9 (1024 * 576)',
              value: '1024x576',
            },
            {
              label: '3:2 (1024 * 682)',
              value: '1024x682',
            },
            {
              label: '4:3 (1024 * 768)',
              value: '1024x768',
            },
            {
              label: '1:1 (1024 * 1024)',
              value: '1024x1024',
            },
            {
              label: '3:4 (768 * 1024)',
              value: '768x1024',
            },
            {
              label: '2:3 (682 * 1024)',
              value: '682x1024',
            },
            {
              label: '9:16 (576 * 1024)',
              value: '576x1024',
            },
          ],
        },
        {
          name: 'use_pre_llm',
          type: 'switch',
          required: false,
          label: '开启文本扩写',
          info: '开启后，系统会自动扩写提示词，提高生成质量',
          value: true,
        },
      ],
    },
    {
      name: '图片 3.0 文生图',
      version: '3.0',
      label: '影视质感，文字更准，直出2k高清图',
      key: 'jimeng_t2i_v30',
      params: [
        {
          name: 'prompt',
          label: '提示词',
          type: 'textarea',
          showWordLimit: true,
          maxlength: 800,
          autosize: { minRows: 3, maxRows: 5 },
          required: true,
          placeholder: '请输入提示词',
          info: '用于生成图像的提示词 ，中英文均可输入',
        },
        {
          name: 'size',
          type: 'select',
          required: true,
          placeholder: '请选择尺寸',
          label: '图片尺寸',

          options: [
            {
              label: '1:1 (1328 * 1328)',
              value: '1328x1328',
            },
            {
              label: '4:3 (1472 * 1104)',
              value: '1472x1104',
            },
            {
              label: '3:2 (1584 * 1056)',
              value: '1584x1056',
            },
            {
              label: '16:9 (1664 * 936)',
              value: '1664x936',
            },
            {
              label: '21:9 (2016 * 864)',
              value: '2016x864',
            },
            {
              label: '1:1 高清2K (2048 * 2048)',
              value: '2048x2048',
            },
            {
              label: '4:3 高清2K (2304 * 1728)',
              value: '2304x1728',
            },
            {
              label: '3:2 高清2K (2496 * 1664)',
              value: '2496x1664',
            },
            {
              label: '16:9 高清2K (2560 * 1440)',
              value: '2560x1440',
            },
            {
              label: '21:9 高清2K (3024 * 1296)',
              value: '3024x1296',
            },
          ],
        },
        {
          name: 'use_pre_llm',
          type: 'switch',
          required: true,
          label: '开启文本扩写',
          info: '开启后，系统会自动扩写提示词，提高生成质量',
          value: true,
        },
      ],
    },
    {
      name: '图片 3.1 文生图',
      version: '3.1',
      label: '丰富的美学多样性，画面更鲜明生动',
      key: 'jimeng_t2i_v31',
      params: [
        {
          name: 'prompt',
          label: '提示词',
          type: 'textarea',
          showWordLimit: true,
          maxlength: 800,
          autosize: { minRows: 3, maxRows: 5 },
          required: true,
          placeholder: '请输入提示词',
          info: '用于生成图像的提示词 ，中英文均可输入',
        },
        {
          name: 'size',
          type: 'select',
          required: true,
          placeholder: '请选择尺寸',
          label: '图片尺寸',

          options: [
            {
              label: '1:1 (1328 * 1328)',
              value: '1328x1328',
            },
            {
              label: '4:3 (1472 * 1104)',
              value: '1472x1104',
            },
            {
              label: '3:2 (1584 * 1056)',
              value: '1584x1056',
            },
            {
              label: '16:9 (1664 * 936)',
              value: '1664x936',
            },
            {
              label: '21:9 (2016 * 864)',
              value: '2016x864',
            },
            {
              label: '1:1 高清2K (2048 * 2048)',
              value: '2048x2048',
            },
            {
              label: '4:3 高清2K (2304 * 1728)',
              value: '2304x1728',
            },
            {
              label: '3:2 高清2K (2496 * 1664)',
              value: '2496x1664',
            },
            {
              label: '16:9 高清2K (2560 * 1440)',
              value: '2560x1440',
            },
            {
              label: '21:9 高清2K (3024 * 1296)',
              value: '3024x1296',
            },
          ],
        },
        {
          name: 'use_pre_llm',
          type: 'switch',
          required: true,
          label: '开启文本扩写',
          info: '开启后，系统会自动扩写提示词，提高生成质量',
          value: true,
        },
      ],
    },

    {
      name: '图片 3.0 图生图',
      version: '3.0',
      label: '精准执行编辑指令，保持图像内容完整性',
      key: 'jimeng_i2i_v30',
      params: [
        {
          name: 'prompt',
          label: '提示词',
          type: 'textarea',
          required: true,
          showWordLimit: true,
          maxlength: 800,
          autosize: { minRows: 3, maxRows: 5 },
          placeholder: '请输入用于编辑图像的提示词，如：把xxx改成xxx，删除xxx，添加xxx等',
          info: '建议长度<=120字符，最长不超过800字符',
        },
        {
          name: 'image_urls',
          label: '参考图片',
          type: 'image',
          required: true,
          placeholder: '请上传图片',
          maxSize: 5,
          accept: '.png,.jpg,.jpeg',
          info: '长边与短边比例在3以内，超出此比例或比例相对极端，会导致报错。',
        },
        {
          name: 'scale',
          label: '文本描述影响的程度',
          type: 'slider',
          min: 0,
          max: 1,
          step: 0.1,
          value: 0.5,
          info: '该值越大代表文本描述影响程度越大，且输入图片影响程度越小',
        },
        {
          name: 'size',
          type: 'select',
          required: true,
          placeholder: '请选择尺寸',
          label: '图片尺寸',

          options: [
            {
              label: '1：1 (1328 * 1328)',
              value: '1328x1328',
            },
            {
              label: '4：3 (1472 * 1104)',
              value: '1472x1104',
            },
            {
              label: '3：2 (1584 * 1056)',
              value: '1584x1056',
            },
            {
              label: '16：9 (1664 * 936)',
              value: '1664x936',
            },
            {
              label: '21：9 (2016 * 864)',
              value: '2016x864',
            },
          ],
        },
      ],
    },
    {
      name: '图片 4.0 文/图生图',
      version: '4.0',
      label:
        '支持文本、单图和多图输入，实现基于主体一致性的多图融合创作、图像编辑、组图生成等多样玩法',
      key: 'jimeng_i2i_v30',
      params: [
        {
          name: 'prompt',
          label: '提示词',
          type: 'textarea',
          required: true,
          showWordLimit: true,
          maxlength: 800,
          autosize: { minRows: 3, maxRows: 5 },
          placeholder: '请输入用于编辑图像的提示词，如：把xxx改成xxx，删除xxx，添加xxx等',
          info: '建议长度<=120字符，最长不超过800字符',
        },
        {
          name: 'image_urls',
          label: '参考图片',
          type: 'image',
          required: true,
          placeholder: '请上传图片',
          maxSize: 5,
          accept: '.png,.jpg,.jpeg',
          info: '长边与短边比例在3以内，超出此比例或比例相对极端，会导致报错。',
        },
        {
          name: 'scale',
          label: '文本描述影响的程度',
          type: 'slider',
          min: 0,
          max: 1,
          step: 0.1,
          value: 0.5,
          info: '该值越大代表文本描述影响程度越大，且输入图片影响程度越小',
        },
        {
          name: 'size',
          type: 'select',
          required: true,
          placeholder: '请选择尺寸',
          label: '图片尺寸',

          options: [
            {
              label: '1：1 (1328 * 1328)',
              value: '1328x1328',
            },
            {
              label: '4：3 (1472 * 1104)',
              value: '1472x1104',
            },
            {
              label: '3：2 (1584 * 1056)',
              value: '1584x1056',
            },
            {
              label: '16：9 (1664 * 936)',
              value: '1664x936',
            },
            {
              label: '21：9 (2016 * 864)',
              value: '2016x864',
            },
          ],
        },
      ],
    },
  ],
  video: [
    {
      name: '视频 3.0 720P-文生视频',
      version: '3.0',
      label: '生成效果与速度兼备',
      key: 'jimeng_t2v_v30',
      params: [
        {
          name: 'prompt',
          label: '提示词',
          type: 'textarea',
          showWordLimit: true,
          maxlength: 800,
          autosize: { minRows: 3, maxRows: 5 },
          required: true,
          placeholder: '请输入提示词',
          info: '用于生成视频的提示词 ，中英文均可输入',
        },
        {
          name: 'aspect_ratio',
          label: '视频比例',
          type: 'select',
          required: false,
          placeholder: '请选择视频比例',
          options: [
            {
              label: '16:9 (横版)',
              value: '16:9',
            },
            {
              label: '4:3 (标准)',
              value: '4:3',
            },
            {
              label: '1:1 (正方形)',
              value: '1:1',
            },
            {
              label: '3:4 (竖版)',
              value: '3:4',
            },
            {
              label: '9:16 (竖屏)',
              value: '9:16',
            },
            {
              label: '21:9 (超宽)',
              value: '21:9',
            },
          ],
        },
        {
          name: 'duration',
          type: 'select',
          label: '视频时长',
          options: [
            {
              label: '5秒',
              value: '5',
            },
            {
              label: '10秒',
              value: '10',
            },
          ],
        },
      ],
    },

    {
      name: '视频 3.0 720P-图生视频-首帧',
      version: '3.0',
      label: '根据提示词 + 首帧图片生成视频',
      key: 'jimeng_i2v_first_v30',
      params: [
        {
          name: 'prompt',
          label: '提示词',
          type: 'textarea',
          showWordLimit: true,
          maxlength: 800,
          autosize: { minRows: 3, maxRows: 5 },
          required: true,
          placeholder: '请输入提示词',
          info: '用于生成视频的提示词 ，中英文均可输入',
        },
        {
          name: 'image_urls',
          label: '首帧图片',
          type: 'image',
          required: false,
          multiple: false,
          maxCount: 1,
          maxSize: 5,
          accept: '.png,.jpg,.jpeg',
        },
        {
          name: 'duration',
          type: 'select',
          label: '视频时长',
          options: [
            {
              label: '5秒',
              value: '5',
            },
            {
              label: '10秒',
              value: '10',
            },
          ],
        },
      ],
    },

    {
      name: '视频 3.0 720P-图生视频-首尾帧',
      version: '3.0',
      label: '根据提示词 + 首尾帧图片生成视频',
      key: 'jimeng_i2v_first_tail_v30',
      params: [
        {
          name: 'prompt',
          label: '提示词',
          type: 'textarea',
          showWordLimit: true,
          maxlength: 800,
          autosize: { minRows: 3, maxRows: 5 },
          required: true,
          placeholder: '请输入提示词',
          info: '用于生成视频的提示词 ，中英文均可输入',
        },
        {
          name: 'image_urls',
          label: '首尾帧图片',
          type: 'image',
          required: false,
          multiple: true,
          maxCount: 2,
          maxSize: 5,
          accept: '.png,.jpg,.jpeg',
          info: '请上传两张图片，第一张为起始帧，第二张为尾帧',
        },
        {
          name: 'duration',
          type: 'select',
          label: '视频时长',
          options: [
            {
              label: '5秒',
              value: '5',
            },
            {
              label: '10秒',
              value: '10',
            },
          ],
        },
      ],
    },

    {
      name: '视频 3.0 720P-图生视频-运镜',
      version: '3.0',
      label: '根据提示词 + 运镜图片生成视频',
      key: 'jimeng_i2v_recamera_v30',
      params: [
        {
          name: 'prompt',
          label: '提示词',
          type: 'textarea',
          showWordLimit: true,
          maxlength: 800,
          autosize: { minRows: 3, maxRows: 5 },
          required: true,
          placeholder: '请输入提示词',
          info: '用于生成视频的提示词 ，中英文均可输入',
        },
        {
          name: 'image_urls',
          label: '运镜图片',
          type: 'image',
          required: false,
          placeholder: '请上传图片',
          maxSize: 5,
          multiple: true,
          maxCount: 1,
          accept: '.png,.jpg,.jpeg',
        },
        {
          name: 'template_id',
          label: '运镜控制',
          type: 'select',
          required: true,
          placeholder: '请选择运镜控制',
          popperClass: 'model-select',
          options: [
            {
              label: '希区柯克推进',
              value: 'hitchcock_dolly_in',
              image: hitchcock_dolly_in,
            },
            {
              label: '希区柯克拉远',
              value: 'hitchcock_dolly_out',
              image: hitchcock_dolly_out,
            },
            {
              label: '机械臂',
              value: 'robo_arm',
              image: robo_arm,
            },
            {
              label: '动感环绕',
              value: 'dynamic_orbit',
              image: dynamic_orbit,
            },
            {
              label: '中心环绕',
              value: 'central_orbit',
              image: central_orbit,
            },
            {
              label: '起重机',
              value: 'crane_push',
              image: crane_push,
            },
            {
              label: '超级拉远',
              value: 'quick_pull_back',
              image: quick_pull_back,
            },
            {
              label: '逆时针回旋',
              value: 'counterclockwise_swivel',
              image: counterclockwise_swivel,
            },
            {
              label: '顺时针回旋',
              value: 'clockwise_swivel',
              image: clockwise_swivel,
            },
            {
              label: '手持运镜',
              value: 'handheld',
              image: handheld,
            },
            {
              label: '快速推拉',
              value: 'rapid_push_pull',
              image: rapid_push_pull,
            },
          ],
          value: 'hitchcock_dolly_in',
        },
        {
          name: 'camera_strength',
          label: '运镜强度',
          type: 'select',
          required: true,
          placeholder: '请选择运镜强度',
          options: [
            {
              label: '弱',
              value: 'weak',
            },
            {
              label: '中',
              value: 'medium',
            },
            {
              label: '强',
              value: 'strong',
            },
          ],
          value: 'medium',
        },
        {
          name: 'duration',
          type: 'select',
          label: '视频时长',
          options: [
            {
              label: '5秒',
              value: '5',
            },
            {
              label: '10秒',
              value: '10',
            },
          ],
          value: '5',
        },
      ],
    },

    {
      name: '视频 3.0 1080P-文生视频',
      version: '3.0',
      label: '视觉表达流畅一致，支持1080P高清渲染',
      key: 'jimeng_t2v_v30_1080p',
      params: [
        {
          name: 'prompt',
          label: '提示词',
          type: 'textarea',
          showWordLimit: true,
          maxlength: 800,
          autosize: { minRows: 3, maxRows: 5 },
          required: true,
          placeholder: '请输入提示词',
          info: '用于生成视频的提示词 ，中英文均可输入',
        },
        {
          name: 'aspect_ratio',
          label: '视频比例',
          type: 'select',
          required: false,
          placeholder: '请选择视频比例',
          options: [
            {
              label: '16:9 (横版)',
              value: '16:9',
            },
            {
              label: '4:3 (标准)',
              value: '4:3',
            },
            {
              label: '1:1 (正方形)',
              value: '1:1',
            },
            {
              label: '3:4 (竖版)',
              value: '3:4',
            },
            {
              label: '9:16 (竖屏)',
              value: '9:16',
            },
            {
              label: '21:9 (超宽)',
              value: '21:9',
            },
          ],
        },
        {
          name: 'duration',
          type: 'select',
          label: '视频时长',
          options: [
            {
              label: '5秒',
              value: '5',
            },
            {
              label: '10秒',
              value: '10',
            },
          ],
        },
      ],
    },

    {
      name: '视频 3.0 1080P-图生视频-首帧',
      version: '3.0',
      label: '根据提示词 + 首帧图片生成1080P视频',
      key: 'jimeng_i2v_first_v30_1080',
      params: [
        {
          name: 'prompt',
          label: '提示词',
          type: 'textarea',
          showWordLimit: true,
          maxlength: 800,
          autosize: { minRows: 3, maxRows: 5 },
          required: true,
          placeholder: '请输入提示词',
          info: '用于生成视频的提示词 ，中英文均可输入',
        },
        {
          name: 'image_urls',
          label: '首帧图片',
          type: 'image',
          required: false,
          multiple: false,
          maxCount: 1,
          maxSize: 5,
          accept: '.png,.jpg,.jpeg',
        },
        {
          name: 'duration',
          type: 'select',
          label: '视频时长',
          options: [
            {
              label: '5秒',
              value: '5',
            },
            {
              label: '10秒',
              value: '10',
            },
          ],
        },
      ],
    },

    {
      name: '视频 3.0 1080P-图生视频-首尾帧',
      version: '3.0',
      label: '根据提示词 + 首尾帧图片生成1080P视频',
      key: 'jimeng_i2v_first_tail_v30_1080',
      params: [
        {
          name: 'prompt',
          label: '提示词',
          type: 'textarea',
          showWordLimit: true,
          maxlength: 800,
          autosize: { minRows: 3, maxRows: 5 },
          required: true,
          placeholder: '请输入提示词',
          info: '用于生成视频的提示词 ，中英文均可输入',
        },
        {
          name: 'image_urls',
          label: '首尾帧图片',
          type: 'image',
          required: false,
          multiple: true,
          maxCount: 2,
          maxSize: 5,
          accept: '.png,.jpg,.jpeg',
          info: '请上传两张图片，第一张为起始帧，第二张为尾帧',
        },
        {
          name: 'duration',
          type: 'select',
          label: '视频时长',
          options: [
            {
              label: '5秒',
              value: '5',
            },
            {
              label: '10秒',
              value: '10',
            },
          ],
        },
      ],
    },

    {
      name: '视频 3.0Pro 1080P-图生视频',
      version: '3.0',
      label: '根据提示词 + 首帧图片生成1080P视频',
      key: 'jimeng_ti2v_v30_pro',
      params: [
        {
          name: 'prompt',
          label: '提示词',
          type: 'textarea',
          showWordLimit: true,
          maxlength: 800,
          autosize: { minRows: 3, maxRows: 5 },
          required: true,
          placeholder: '请输入提示词',
          info: '用于生成视频的提示词 ，中英文均可输入',
        },
        {
          name: 'image_urls',
          label: '首帧图片',
          type: 'image',
          required: false,
          info: '只支持上传首帧图片',
          multiple: false,
          maxCount: 1,
          maxSize: 5,
          accept: '.png,.jpg,.jpeg',
        },
        // 比例
        {
          name: 'aspect_ratio',
          label: '视频比例',
          type: 'select',
          required: false,
          placeholder: '请选择视频比例',
          info: '只在文生视频场景下生效，图生视频场景会根据输入图的长宽比自动适配',
          options: [
            {
              label: '21:9 (2176 * 928)',
              value: '21:9',
            },
            {
              label: '16:9 (1920 * 1088)',
              value: '16:9',
            },
            {
              label: '4:3 (1664 * 1248)',
              value: '4:3',
            },
            {
              label: '1:1 (1440 * 1440)',
              value: '1:1',
            },
            {
              label: '3:4 (1248 * 1664)',
              value: '3:4',
            },
            {
              label: '9:16 (1088 * 1920)',
              value: '9:16',
            },
          ],
        },
        {
          name: 'duration',
          type: 'select',
          label: '视频时长',
          options: [
            {
              label: '5秒',
              value: '5',
            },
            {
              label: '10秒',
              value: '10',
            },
          ],
          value: '5',
        },
      ],
    },
  ],
  virtualHuman: [],
  actionTransfer: [],
}

export const JimengFunctions = [
  {
    key: 'image',
    icon: 'iconfont icon-image',
    name: '图片生成',
  },
  {
    key: 'video',
    icon: 'icon-video',
    name: '视频生成',
  },
  {
    key: 'virtualHuman',
    icon: 'icon-shuziren',
    name: '数字人',
  },
  {
    key: 'actionTransfer',
    icon: 'icon-action',
    name: '动作模仿',
  },
]
