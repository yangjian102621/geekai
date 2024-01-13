package storage

type uploadProgress struct {
	lastUploadedBytes int64
	progress          func(totalBytes, uploadedBytes int64)
}

func newUploadProgress(progressHandler func(totalBytes, uploadedBytes int64)) *uploadProgress {
	return &uploadProgress{
		lastUploadedBytes: 0,
		progress:          progressHandler,
	}
}

func (p *uploadProgress) onProgress(totalBytes, uploadedBytes int64) {
	if p.progress == nil {
		return
	}

	if p.lastUploadedBytes >= uploadedBytes {
		// 过滤重新上传的场景
		return
	}
	p.lastUploadedBytes = uploadedBytes
	p.progress(totalBytes, uploadedBytes)
}
