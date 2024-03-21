package tools

import (
	"Arkadiy_Servis_authorization/iternal/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

//func MultipartFormDataTopic(any interface{}, c *gin.Context) (multipart.File, error) {
//	t := reflect.TypeOf(any)
//	switch t.Name() {
//	case "Topic":
//		var topic domain.Topic
//
//		dataString := c.Request.FormValue("data")
//		data := []byte(dataString)
//
//		err := json.Unmarshal(data, &topic)
//		if err != nil {
//			return nil, err
//		}
//
//		file, _, err := c.Request.FormFile("logo")
//		if err != nil {
//			return nil, err
//		}
//		return file, nil
//	}
//
//	dataString := c.Request.FormValue("data")
//	data := []byte(dataString)
//
//	err := json.Unmarshal(data, &topic)
//	if err != nil {
//		return err
//	}
//
//	file, _, err := c.Request.FormFile("logo")
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func getType(myvar interface{}) string {
//	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
//		return "*" + t.Elem().Name()
//	} else {
//		return t.Name()
//	}
//}

func MultipartFormDataTopic(topic *models.Topic, c *gin.Context) (string, error) {
	var pathToLogo string
	dataString := c.Request.FormValue("data")
	data := []byte(dataString)

	err := json.Unmarshal(data, &topic)
	if err != nil {
		return "", err
	}

	for _, fileHeader := range c.Request.MultipartForm.File["logo"] {

		file, err2 := fileHeader.Open()
		if err2 != nil {
			return "", err2
		}

		path, err2 := MakeDir(file, "topic", fileHeader.Filename)
		pathToLogo = path
	}

	return pathToLogo, nil
}

func MultipartFormDataMassage(massage *models.Massage, c *gin.Context) ([]string, string, Claims, error) {
	var massDirPath []string

	topicID := c.Query("topic_id")

	claims, err := ParsTokenClaims(c.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, "", Claims{}, err
	}

	dataString := c.Request.FormValue("data")
	data := []byte(dataString)

	err = json.Unmarshal(data, &massage)
	if err != nil {
		return nil, "", Claims{}, err
	}

	for _, fileHeader := range c.Request.MultipartForm.File["file"] {
		file, err2 := fileHeader.Open()
		if err2 != nil {
			return nil, "", Claims{}, err2
		}

		dirPath, err2 := MakeDir(file, claims.UserID.String(), fileHeader.Filename)
		if err2 != nil {
			return nil, "", Claims{}, err2
		}

		massDirPath = append(massDirPath, dirPath)
	}

	return massDirPath, topicID, claims, nil
}
