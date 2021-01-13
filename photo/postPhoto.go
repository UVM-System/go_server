package photo

import (
	"bytes"
	"fmt"
	"go_server/config"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

func createFormDataFromBytes(images []ImageBytes) (string,*bytes.Buffer,error)  {
	bodyBuf := &bytes.Buffer{}

	bodyWriter :=multipart.NewWriter(bodyBuf)

	//写入文件
	for _,imageBytes :=range images {
		bufferWriter,_ :=bodyWriter.CreatePart(mimeHeader(imageBytes.FieldName,imageBytes.FileName,imageBytes.ContentType))
		bufferWriter.Write(imageBytes.Content)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	return contentType,bodyBuf,nil
}

func PostImage(images []ImageBytes) []byte  {
	contentType,bodyBuffer,_ := createFormDataFromBytes(images)
	response, err := http.Post(config.Config.DetectUrl, contentType, bodyBuffer)
	if err!=nil{
		panic(err.Error())
	}
	defer response.Body.Close()
	respBody, err := ioutil.ReadAll(response.Body)
	log.Println(string(respBody))
	return respBody
}

func mimeHeader(fieldname, filename,contenttype string) textproto.MIMEHeader {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", contenttype)
	return h
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
