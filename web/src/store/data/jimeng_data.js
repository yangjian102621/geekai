import central_orbit from '@/assets/img/jimeng/yunjing/central_orbit.webp'
import clockwise_swivel from '@/assets/img/jimeng/yunjing/clockwise_swivel.webp'
import counterclockwise_swivel from '@/assets/img/jimeng/yunjing/counterclockwise_swivel.webp'
import crane_push from '@/assets/img/jimeng/yunjing/crane_push.webp'
import dynamic_orbit from '@/assets/img/jimeng/yunjing/dynamic_orbit.webp'
import handheld from '@/assets/img/jimeng/yunjing/handheld.webp'
import hitchcock_dolly_in from '@/assets/img/jimeng/yunjing/hitchcock_dolly_in.webp'
import hitchcock_dolly_out from '@/assets/img/jimeng/yunjing/hitchcock_dolly_out.webp'
import quick_pull_back from '@/assets/img/jimeng/yunjing/quick_pull_back.webp'
import rapid_push_pull from '@/assets/img/jimeng/yunjing/rapid_push_pull.webp'
import robo_arm from '@/assets/img/jimeng/yunjing/robo_arm.webp'

import acrylic_ornaments from '@/assets/img/jimeng/texiao/acrylic_ornaments.png'
import angel_figurine from '@/assets/img/jimeng/texiao/angel_figurine.png'
import birthday_photo_gorgeous from '@/assets/img/jimeng/texiao/birthday_photo_gorgeous.jpeg'
import birthday_photo_party from '@/assets/img/jimeng/texiao/birthday_photo_party.jpeg'
import birthday_photo_red from '@/assets/img/jimeng/texiao/birthday_photo_red.jpeg'
import car_miniature_ornaments from '@/assets/img/jimeng/texiao/car_miniature_ornaments.jpeg'
import Christmas_green_background from '@/assets/img/jimeng/texiao/Christmas_green_background.jpeg'
import Christmas_tree from '@/assets/img/jimeng/texiao/Christmas_tree.jpeg'
import claw_machine_style from '@/assets/img/jimeng/texiao/claw_machine_style.jpeg'
import earphone_case_style from '@/assets/img/jimeng/texiao/earphone_case_style.jpeg'
import electronic_pet_egg_style from '@/assets/img/jimeng/texiao/electronic_pet_egg_style.jpeg'
import felt_3d_polaroid from '@/assets/img/jimeng/texiao/felt_3d_polaroid.png'
import felt_keychain from '@/assets/img/jimeng/texiao/felt_keychain.png'
import furry_dream_doll from '@/assets/img/jimeng/texiao/furry_dream_doll.png'
import glass_ball from '@/assets/img/jimeng/texiao/glass_ball.png'
import graduation_photo from '@/assets/img/jimeng/texiao/graduation_photo.png'
import lofi_pixel_character_mini_card from '@/assets/img/jimeng/texiao/lofi_pixel_character_mini_card.png'
import lying_in_fluffy_belly from '@/assets/img/jimeng/texiao/lying_in_fluffy_belly.png'
import micro_landscape_mini_world from '@/assets/img/jimeng/texiao/micro_landscape_mini_world.png'
import micro_landscape_mini_world_professional from '@/assets/img/jimeng/texiao/micro_landscape_mini_world_professional.png'
import Mid_Autumn_Festival_individual from '@/assets/img/jimeng/texiao/Mid-Autumn_Festival_individual.jpeg'
import Mid_Autumn_Festival_new_chinese_style from '@/assets/img/jimeng/texiao/Mid-Autumn_Festival_new_chinese_style.jpeg'
import my_world from '@/assets/img/jimeng/texiao/my_world.png'
import my_world_universal from '@/assets/img/jimeng/texiao/my_world_universal.png'
import patchwork_collage_style from '@/assets/img/jimeng/texiao/patchwork_collage_style.jpeg'
import plastic_bubble_figure from '@/assets/img/jimeng/texiao/plastic_bubble_figure.png'
import plastic_bubble_figure_cartoon_text from '@/assets/img/jimeng/texiao/plastic_bubble_figure_cartoon_text.png'
import Spring_Festival_traditional_Chinese_architecture from '@/assets/img/jimeng/texiao/Spring_Festival_traditional_Chinese_architecture.png'

export const JimengParams = {
  image: [
    {
      name: '图片 4.0 文/图生图',
      version: '4.0',
      label: '支持文本、单图和多图输入，实现基于主体一致性的多图融合创作、图像编辑等多样玩法',
      key: 'doubao-seedream-4-0-250828',
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
          info: '建议不超过300个汉字或600个英文单词。字数过多信息容易分散，模型可能因此忽略细节。',
        },
        {
          name: 'image_urls',
          label: '参考图片',
          type: 'image',
          required: false,
          placeholder: '请上传图片',
          maxSize: 5,
          multiple: true,
          maxCount: 10,
          accept: '.png,.jpg,.jpeg',
          info: '支持编辑单张图片，或者一次融合多张图片',
        },
        {
          name: 'size',
          type: 'select',
          required: true,
          placeholder: '请选择尺寸',
          label: '图片尺寸',
          prefix: 'icon-resize',
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
          ],
        },
      ],
    },
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
          prefix: 'icon-resize',
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
          prefix: 'icon-resize',
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
          prefix: 'icon-resize',
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
          name: 'size',
          type: 'select',
          required: true,
          placeholder: '请选择尺寸',
          label: '图片尺寸',
          prefix: 'icon-resize',
          options: [
            {
              label: '1:1 (2048 * 2048)',
              value: '2048x2048',
            },
            {
              label: '4:3 (2304 * 1728)',
              value: '2304x1728',
            },
            {
              label: '3:4 (1728 * 2304)',
              value: '1728x2304',
            },
            {
              label: '16:9 (2560 * 1440)',
              value: '2560x1440',
            },
            {
              label: '9:16 (1440 * 2560)',
              value: '1440x2560',
            },
            {
              label: '3:2 (2496 * 1664)',
              value: '2496x1664',
            },
            {
              label: '2:3 (1664 * 2496)',
              value: '1664x2496',
            },
            {
              label: '21:9 (3024 * 1296)',
              value: '3024x1296',
            },
          ],
        },
      ],
    },

    {
      name: '图片 3.0 图像特效',
      version: '3.0',
      label: '将输入的单人写真图片，进行有创意的特效化处理。',
      key: 'i2i_multi_style_zx2x',
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
          info: '支持输入人像写真图片。',
        },
        {
          name: 'template_id',
          label: '特效模板ID',
          type: 'select',
          required: true,
          placeholder: '请选择特效模板ID',
          imgSize: '40px',
          popperClass: 'model-select',
          prefix: 'icon-sd',
          options: [
            {
              label: '毛毡3D拍立得风格',
              value: 'felt_3d_polaroid',
              image: felt_3d_polaroid,
            },
            {
              label: '像素世界风',
              value: 'my_world',
              image: my_world,
            },
            {
              label: '像素世界-万物通用版',
              value: 'my_world_universal',
              image: my_world_universal,
            },
            {
              label: '盲盒玩偶风',
              value: 'plastic_bubble_figure',
              image: plastic_bubble_figure,
            },
            {
              label: '塑料泡罩人偶-文字卡头版',
              value: 'plastic_bubble_figure_cartoon_text',
              image: plastic_bubble_figure_cartoon_text,
            },
            {
              label: '毛绒玩偶风',
              value: 'furry_dream_doll',
              image: furry_dream_doll,
            },
            {
              label: '迷你世界玩偶风',
              value: 'micro_landscape_mini_world',
              image: micro_landscape_mini_world,
            },
            {
              label: '微型景观小世界-职业版',
              value: 'micro_landscape_mini_world_professional',
              image: micro_landscape_mini_world_professional,
            },
            {
              label: '亚克力挂饰',
              value: 'acrylic_ornaments',
              image: acrylic_ornaments,
            },
            {
              label: '毛毡钥匙扣',
              value: 'felt_keychain',
              image: felt_keychain,
            },
            {
              label: 'Lofi 像素人物小卡',
              value: 'lofi_pixel_character_mini_card',
              image: lofi_pixel_character_mini_card,
            },
            {
              label: '天使形象手办',
              value: 'angel_figurine',
              image: angel_figurine,
            },
            {
              label: '躺在毛茸茸肚皮里',
              value: 'lying_in_fluffy_belly',
              image: lying_in_fluffy_belly,
            },
            {
              label: '玻璃球',
              value: 'glass_ball',
              image: glass_ball,
            },
            {
              label: '耳机盒',
              value: 'earphone_case_style',
              image: earphone_case_style,
            },
            {
              label: '电子宠物蛋',
              value: 'electronic_pet_egg_style',
              image: electronic_pet_egg_style,
            },
            {
              label: '拼贴缝布',
              value: 'patchwork_collage_style',
              image: patchwork_collage_style,
            },
            {
              label: '抓娃娃机',
              value: 'claw_machine_style',
              image: claw_machine_style,
            },
            {
              label: '车内微缩摆件',
              value: 'car_miniature_ornaments',
              image: car_miniature_ornaments,
            },
            {
              label: '中秋节-新中式',
              value: 'Mid-Autumn_Festival_new_chinese_style',
              image: Mid_Autumn_Festival_new_chinese_style,
            },
            {
              label: '中秋单人',
              value: 'Mid-Autumn_Festival_individual',
              image: Mid_Autumn_Festival_individual,
            },
            {
              label: '圣诞节绿背景',
              value: 'Christmas_green_background',
              image: Christmas_green_background,
            },
            {
              label: '圣诞节圣诞树',
              value: 'Christmas_tree',
              image: Christmas_tree,
            },
            {
              label: '春节红墙',
              value: 'Spring_Festival_traditional_Chinese_architecture',
              image: Spring_Festival_traditional_Chinese_architecture,
            },
            {
              label: '生日照华丽',
              value: 'birthday_photo_gorgeous',
              image: birthday_photo_gorgeous,
            },
            {
              label: '生日照红色',
              value: 'birthday_photo_red',
              image: birthday_photo_red,
            },
            {
              label: '生日照派对',
              value: 'birthday_photo_party',
              image: birthday_photo_party,
            },
            {
              label: '毕业照',
              value: 'graduation_photo',
              image: graduation_photo,
            },
          ],
        },
        {
          name: 'size',
          type: 'select',
          required: true,
          placeholder: '请选择尺寸',
          label: '图片尺寸',
          prefix: 'icon-resize',
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
          ],
        },
      ],
    },
  ],
  video: [
    // 视频 3.0 720P-文生视频
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
          required: true,
          placeholder: '请选择视频比例',
          prefix: 'icon-resize',
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
          prefix: 'icon-clock',
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
    // 视频 3.0 图生视频-首帧
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
          required: true,
          multiple: false,
          maxCount: 1,
          maxSize: 5,
          accept: '.png,.jpg,.jpeg',
        },
        {
          name: 'duration',
          type: 'select',
          label: '视频时长',
          prefix: 'icon-clock',
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
    // 视频 3.0 图生视频-首尾帧
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
          required: true,
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
          prefix: 'icon-clock',
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
    // 视频 3.0 图生视频-运镜
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
          required: true,
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
          prefix: 'icon-yunjing',
          imgSize: '54px',
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
          prefix: 'icon-clock',
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
    // 视频 3.0 1080P-文生视频
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
          prefix: 'icon-resize',
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
          prefix: 'icon-clock',
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
    // 视频 3.0 1080P-图生视频-首帧
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
          required: true,
          multiple: false,
          maxCount: 1,
          maxSize: 5,
          accept: '.png,.jpg,.jpeg',
        },
        {
          name: 'duration',
          type: 'select',
          label: '视频时长',
          prefix: 'icon-clock',
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
    // 视频 3.0 1080P-图生视频-首尾帧
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
          required: true,
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
          prefix: 'icon-clock',
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
    // 视频 3.0Pro 1080P-图生视频
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
          prefix: 'icon-resize',
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
          prefix: 'icon-clock',
          placeholder: '请选择视频时长',
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
    icon: 'icon-image',
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
