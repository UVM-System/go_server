// Package photo 存储和照片有关的信息
package photo

type ImageBytes struct {
	FieldName string
	FileName string
	Content []byte
	ContentType string
}

//// 构造远程服务器检测的数据图像结构
//func CreateFormdata(mValue map[string][]string, images []ImageBytes) PostBytes {
//	fieldmap := make(map[string]string)
//	for key, value := range mValue  {
//		fieldmap[key] = value[0]
//	}
//	formData := PostBytes {
//		FileMap: images,
//		FieldMap: fieldmap,
//	}
//	return formData
//}

// CreateImageBytes 构造一张图像的结构，用于数据流传输
func CreateImageBytes(data []byte, filename string) ImageBytes {
	img := ImageBytes {
		FieldName: "image",
		FileName: filename,
		Content: data,
		ContentType: "image/jpeg",
	}
	return img
}
