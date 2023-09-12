package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	imageURL := "https://cdn.discordapp.com/attachments/1151037077308325901/1151286701717733416/jiangjin_a_chrysanthemum_in_the_style_of_Van_Gogh_49b64011-6581-469d-9888-c285ab964e08.png"

	fmt.Println(filepath.Ext(filepath.Base(imageURL)))
}
