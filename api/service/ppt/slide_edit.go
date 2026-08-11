package ppt

import (
	"context"
	"fmt"
	"geekai/core/types"
	"geekai/store/model"
	"geekai/store/vo"
	"strings"
)

// EditSlideImage 基于当前激活图做图生图，追加 image_history 并将 image_url 设为新版。
func (s *PptService) EditSlideImage(ctx context.Context, taskID string, userID uint, slideIndex int, prompt string, oss types.OSSConfig, app *types.AppConfig) ([]SlideData, error) {
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return nil, fmt.Errorf("请输入修改说明")
	}
	task, ok := s.GetTask(taskID)
	if !ok {
		return nil, ErrPptTaskNotFound
	}
	if task.UserID != userID {
		return nil, ErrPptTaskNotFound
	}
	refURL := ""
	for _, sl := range task.Slides {
		if sl.SlideIndex == slideIndex {
			normalizeSlideImageHistory(&sl)
			refURL = strings.TrimSpace(sl.ImageURL)
			break
		}
	}
	if refURL == "" {
		if slideExists(task.Slides, slideIndex) {
			return nil, ErrPptSlideNoImage
		}
		return nil, ErrPptSlideNotFound
	}

	cfg, err := s.loadPPTConfig()
	if err != nil {
		return nil, err
	}
	power, err := s.userPower(userID)
	if err != nil {
		return nil, err
	}
	if cfg.PowerCostPerSlide > 0 && power < cfg.PowerCostPerSlide {
		return nil, ErrInsufficientPower
	}

	generator, err := NewImageGenerator(cfg)
	if err != nil {
		return nil, err
	}

	refInputs, err := PrepareReferenceInputsForImg2Img(refURL, oss, app)
	if err != nil {
		return nil, fmt.Errorf("准备参考图失败：%w", err)
	}

	imgURL, err := generator.GenerateWithReference(ctx, prompt, refInputs)
	if err != nil {
		return nil, err
	}

	storedURL, err := s.uploadManager.GetUploadHandler().PutUrlFile(imgURL, ".png", false)
	if err != nil {
		return nil, fmt.Errorf("转存图片失败：%w", err)
	}

	if err := s.applySlideImageEdit(taskID, slideIndex, storedURL, prompt); err != nil {
		return nil, err
	}

	if cfg.PowerCostPerSlide > 0 {
		err = s.userService.DecreasePower(userID, cfg.PowerCostPerSlide, model.PowerLog{
			Type:   types.PowerConsume,
			Model:  generator.Provider(),
			Remark: fmt.Sprintf("PPT 任务 %s 第 %d 页图生图编辑", taskID, slideIndex),
		})
		if err != nil {
			return nil, fmt.Errorf("扣减算力失败：%v", err)
		}
	}

	task2, _ := s.GetTask(taskID)
	return task2.Slides, nil
}

func slideExists(slides []SlideData, slideIndex int) bool {
	for _, sl := range slides {
		if sl.SlideIndex == slideIndex {
			return true
		}
	}
	return false
}

func (s *PptService) applySlideImageEdit(taskID string, slideIndex int, newURL string, editPrompt string) error {
	s.slidesLock.Lock()
	defer s.slidesLock.Unlock()

	var job model.PPTJob
	if err := s.db.Where("task_id = ?", taskID).First(&job).Error; err != nil {
		return err
	}
	slides := job.Slides
	found := false
	for i := range slides {
		if slides[i].SlideIndex != slideIndex {
			continue
		}
		found = true
		sd := voToSlide(slides[i])
		normalizeSlideImageHistory(&sd)
		if strings.TrimSpace(sd.ImageURL) == "" {
			return ErrPptSlideNoImage
		}
		sd.ImageHistory = append(sd.ImageHistory, vo.PPTSlideImageVersion{
			ImageURL: newURL,
			Prompt:   editPrompt,
		})
		sd.ImageURL = newURL
		slides[i] = slideToVO(sd)
		break
	}
	if !found {
		return ErrPptSlideNotFound
	}
	job.Slides = slides
	biz := voSlidesToBiz(slides)
	return s.refreshJobMeta(&job, biz)
}

// SetActiveSlideVersion 将 image_url 切换为 image_history[versionIndex]。
func (s *PptService) SetActiveSlideVersion(taskID string, userID uint, slideIndex int, versionIndex int) ([]SlideData, error) {
	task, ok := s.GetTask(taskID)
	if !ok {
		return nil, ErrPptTaskNotFound
	}
	if task.UserID != userID {
		return nil, ErrPptTaskNotFound
	}

	s.slidesLock.Lock()
	defer s.slidesLock.Unlock()

	var job model.PPTJob
	if err := s.db.Where("task_id = ?", taskID).First(&job).Error; err != nil {
		return nil, err
	}
	slides := job.Slides
	found := false
	for i := range slides {
		if slides[i].SlideIndex != slideIndex {
			continue
		}
		found = true
		sd := voToSlide(slides[i])
		normalizeSlideImageHistory(&sd)
		hist := sd.ImageHistory
		if versionIndex < 0 || versionIndex >= len(hist) {
			return nil, ErrPptInvalidVersionIndex
		}
		sd.ImageURL = hist[versionIndex].ImageURL
		slides[i] = slideToVO(sd)
		break
	}
	if !found {
		return nil, ErrPptSlideNotFound
	}

	job.Slides = slides
	biz := voSlidesToBiz(slides)
	if err := s.refreshJobMeta(&job, biz); err != nil {
		return nil, err
	}

	out := voSlidesToBizNormalized(job.Slides)
	return out, nil
}
