package main

import (
	"chatplus/store/model"
	"chatplus/utils"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"path"
)

func main() {
	MysqlDns := "root:12345678@tcp(localhost:3306)/chatgpt_plus?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(MysqlDns), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	_ = os.MkdirAll("static/upload/images", 0755)
	var jobs []model.MidJourneyJob
	db.Find(&jobs)
	for _, job := range jobs {
		basename := path.Base(job.ImgURL)
		imageData, err := utils.DownloadImage(job.ImgURL, "")
		if err != nil {
			fmt.Println("图片下载失败：" + job.ImgURL)
			continue
		}
		newImagePath := fmt.Sprintf("static/upload/images/%s", basename)
		err = os.WriteFile(newImagePath, imageData, 0644)
		if err != nil {
			fmt.Println("Error writing image file:", err)
			continue
		}
		fmt.Println("图片保存成功！", newImagePath)
		// 更新数据库
		job.ImgURL = fmt.Sprintf("http://localhost:5678/%s", newImagePath)
		db.Updates(&job)
	}

}
