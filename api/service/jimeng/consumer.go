package jimeng

import (
	"context"
	"time"

	"geekai/logger"
	"geekai/store/model"
)

var jimengLogger = logger.GetLogger()

// Consumer 即梦任务消费者
type Consumer struct {
	service *Service
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewConsumer 创建即梦任务消费者
func NewConsumer(service *Service) *Consumer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Consumer{
		service: service,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start 启动消费者
func (c *Consumer) Start() {
	jimengLogger.Info("Starting Jimeng task consumer...")
	go c.consume()
}

// Stop 停止消费者
func (c *Consumer) Stop() {
	jimengLogger.Info("Stopping Jimeng task consumer...")
	c.cancel()
}

// consume 消费任务
func (c *Consumer) consume() {
	for {
		select {
		case <-c.ctx.Done():
			jimengLogger.Info("Jimeng task consumer stopped")
			return
		default:
			c.processTask()
		}
	}
}

// processTask 处理任务
func (c *Consumer) processTask() {
	// 从队列中获取任务
	var task map[string]any
	if err := c.service.taskQueue.LPop(&task); err != nil {
		// 队列为空，等待1秒后重试
		time.Sleep(time.Second)
		return
	}

	// 解析任务
	jobIdFloat, ok := task["job_id"].(float64)
	if !ok {
		jimengLogger.Errorf("invalid job_id in task: %v", task)
		return
	}
	jobId := uint(jobIdFloat)

	taskType, ok := task["type"].(string)
	if !ok {
		jimengLogger.Errorf("invalid task type in task: %v", task)
		return
	}

	jimengLogger.Infof("Processing Jimeng task: job_id=%d, type=%s", jobId, taskType)

	// 处理任务
	if err := c.service.ProcessTask(jobId); err != nil {
		jimengLogger.Errorf("process jimeng task failed: job_id=%d, error=%v", jobId, err)

		// 任务失败，直接标记为失败状态，不进行重试
		c.service.UpdateJobStatus(jobId, model.JMTaskStatusFailed, err.Error())
	} else {
		jimengLogger.Infof("Jimeng task processed successfully: job_id=%d", jobId)
	}
}

// TaskQueueStatus 任务队列状态
type TaskQueueStatus struct {
	QueueLength int `json:"queue_length"`
	ActiveTasks int `json:"active_tasks"`
}

// GetQueueStatus 获取队列状态
func (c *Consumer) GetQueueStatus() (*TaskQueueStatus, error) {
	// 获取队列长度
	length, err := c.service.taskQueue.Size()
	if err != nil {
		return nil, err
	}

	// 获取活跃任务数（正在处理的任务）
	activeTasks, err := c.service.GetPendingTaskCount(0) // 0表示所有用户
	if err != nil {
		activeTasks = 0
	}

	return &TaskQueueStatus{
		QueueLength: int(length),
		ActiveTasks: int(activeTasks),
	}, nil
}

// MonitorQueue 监控队列状态
func (c *Consumer) MonitorQueue() {
	ticker := time.NewTicker(30 * time.Second) // 每30秒监控一次
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			status, err := c.GetQueueStatus()
			if err != nil {
				jimengLogger.Errorf("get queue status failed: %v", err)
				continue
			}

			if status.QueueLength > 0 || status.ActiveTasks > 0 {
				jimengLogger.Infof("Jimeng queue status: queue_length=%d, active_tasks=%d",
					status.QueueLength, status.ActiveTasks)
			}
		}
	}
}

// PushTaskToQueue 推送任务到队列（用于手动重试）
func (c *Consumer) PushTaskToQueue(task map[string]interface{}) error {
	return c.service.taskQueue.RPush(task)
}

// GetTaskStats 获取任务统计信息
func (c *Consumer) GetTaskStats() (map[string]any, error) {
	type StatResult struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	var stats []StatResult
	err := c.service.db.Model(&model.JimengJob{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&stats).Error
	if err != nil {
		return nil, err
	}

	result := map[string]any{
		"total":      int64(0),
		"completed":  int64(0),
		"processing": int64(0),
		"failed":     int64(0),
		"pending":    int64(0),
	}

	for _, stat := range stats {
		result["total"] = result["total"].(int64) + stat.Count
		result[stat.Status] = stat.Count
	}

	return result, nil
}
