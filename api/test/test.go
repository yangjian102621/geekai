package main

import (
	"chatplus/utils"
	"fmt"
	"os"
)

func main() {
	imgURL := "https://oaidalleapiprodscus.blob.core.windows.net/private/org-UJimNEKhVm07E58nxnjx5FeG/user-e5UAcPVbkm2nwD8urggRRM8q/img-zFXWyrJ9Z1HppI36dZMXNEaA.png?st=2023-11-26T09%3A57%3A49Z&se=2023-11-26T11%3A57%3A49Z&sp=r&sv=2021-08-06&sr=b&rscd=inline&rsct=image/png&skoid=6aaadede-4fb3-4698-a8f6-684d7786b067&sktid=a48cca56-e6da-484e-a814-9c849652bcb3&skt=2023-11-26T08%3A14%3A59Z&ske=2023-11-27T08%3A14%3A59Z&sks=b&skv=2021-08-06&sig=VmlU9didavbl02XYim2XuMmLMFJsLtCY/ULnzCjeO1g%3D"
	imageData, err := utils.DownloadImage(imgURL, "")
	if err != nil {
		panic(err)
	}
	newImagePath := "newimage.png"
	err = os.WriteFile(newImagePath, imageData, 0644)
	if err != nil {
		fmt.Println("Error writing image file:", err)
		return
	}

	fmt.Println("图片保存成功！")
}
