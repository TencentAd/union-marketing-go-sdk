package ams
//
//import (
//	"encoding/base64"
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"os"
//	"testing"
//
//	"github.com/tencentad/union-marketing-go-sdk/api/sdk"
//	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/config"
//	"github.com/tencentad/marketing-api-go-sdk/pkg/errors"
//)
//
//func TestAddImage(t *testing.T) {
//	amsService := NewAMSService(&config.Config{})
//	file, err := os.Open("../../../test/addImageTest.jpeg")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	bytes, err := ioutil.ReadAll(file)
//	sEnc := base64.StdEncoding.EncodeToString(bytes)
//	input := sdk.ImageAddInput{
//		BaseInput: sdk.BaseInput{
//			AccountId:   25610,
//			AccountType: sdk.AccountTypeAMS,
//			AccessToken: "4b647781310e83b001408e3ce092e48e",
//		},
//		UploadType: sdk.UPLOAD_TYPE_BYTES,
//		Signature:  "b375ccc458c57ed5c79fd7bb076e87f6",
//		Bytes:      sEnc,
//		Desc:       "测试AddImage接口",
//	}
//
//	toutput, err := amsService.AddImage(&input)
//	if err != nil {
//		if resErr, ok := err.(errors.ResponseError); ok {
//			errStr, _ := json.Marshal(resErr)
//			fmt.Println("Response error:", string(errStr))
//		} else {
//			fmt.Println("Error:", err)
//		}
//	}
//	responseJson, _ := json.Marshal(toutput)
//	fmt.Println("Response data:", string(responseJson))
//}
//
//func TestAddImageByFile(t *testing.T) {
//	amsService := NewAMSService(&config.Config{})
//	file, err := os.Open("../../../test/addImageByFileTest.jpeg")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	input := sdk.ImageAddInput{
//		BaseInput: sdk.BaseInput{
//			AccountId:   25610,
//			AccountType: sdk.AccountTypeAMS,
//			AccessToken: "4b647781310e83b001408e3ce092e48e",
//		},
//		UploadType: sdk.UPLOAD_TYPE_FILE,
//		Signature:  "6f5320b2767e5a1c51079b5818505bb0",
//		File:      file,
//		Desc:       "测试AddImageByFile接口",
//	}
//
//	fmt.Println(input)
//
//	toutput, err := amsService.AddImage(&input)
//	if err != nil {
//		if resErr, ok := err.(errors.ResponseError); ok {
//			errStr, _ := json.Marshal(resErr)
//			fmt.Println("Response error:", string(errStr))
//		} else {
//			fmt.Println("Error:", err)
//		}
//	}
//	responseJson, _ := json.Marshal(toutput)
//	fmt.Println("Response data:", string(responseJson))
//}
//
//func TestGetImage(t *testing.T) {
//	amsService := NewAMSService(&config.Config{})
//
//	input := sdk.MaterialGetInput{
//		BaseInput: sdk.BaseInput{
//			AccountId:   25610,
//			AccountType: sdk.AccountTypeAMS,
//			AccessToken: "4b647781310e83b001408e3ce092e48e",
//		},
//		Filtering: &sdk.MaterialFiltering{
//			MaterialIds: []string{"752564005"},
//		},
//		Page:     0,
//		PageSize: 0,
//	}
//
//	toutput, err := amsService.GetImage(&input)
//	if err != nil {
//		if resErr, ok := err.(errors.ResponseError); ok {
//			errStr, _ := json.Marshal(resErr)
//			fmt.Println("Response error:", string(errStr))
//		} else {
//			fmt.Println("Error:", err)
//		}
//	}
//	responseJson, _ := json.Marshal(toutput)
//	fmt.Println("Response data:", string(responseJson))
//}
//
//func TestGetVideo(t *testing.T) {
//	amsService := NewAMSService(&config.Config{})
//
//	input := sdk.MaterialGetInput{
//		BaseInput: sdk.BaseInput{
//			AccountId:   25610,
//			AccountType: sdk.AccountTypeAMS,
//			AccessToken: "4b647781310e83b001408e3ce092e48e",
//		},
//	}
//
//	toutput, err := amsService.GetVideo(&input)
//	if err != nil {
//		if resErr, ok := err.(errors.ResponseError); ok {
//			errStr, _ := json.Marshal(resErr)
//			fmt.Println("Response error:", string(errStr))
//		} else {
//			fmt.Println("Error:", err)
//		}
//	}
//	responseJson, _ := json.Marshal(toutput)
//	fmt.Println("Response data:", string(responseJson))
//}
//
//func TestAddVideo(t *testing.T) {
//	amsService := NewAMSService(&config.Config{})
//	file, err := os.Open("../../../test/addVideoTest.mp4")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	input := sdk.VideoAddInput{
//		BaseInput: sdk.BaseInput{
//			AccountId:   25610,
//			AccountType: sdk.AccountTypeAMS,
//			AccessToken: "4b647781310e83b001408e3ce092e48e",
//		},
//		Signature:  "8c8fcf8f931319ea40175d7baedaf55a",
//		Desc:       "测试AddVideo接口",
//		File: file,
//	}
//
//	toutput, err := amsService.AddVideo(&input)
//	if err != nil {
//		if resErr, ok := err.(errors.ResponseError); ok {
//			errStr, _ := json.Marshal(resErr)
//			fmt.Println("Response error:", string(errStr))
//		} else {
//			fmt.Println("Error:", err)
//		}
//	}
//	responseJson, _ := json.Marshal(toutput)
//	fmt.Println("Response data:", string(responseJson))
//}
