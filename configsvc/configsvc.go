package configsvc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-aniket/etcd"
	"github.com/remiges-aniket/etcdls"
	"github.com/remiges-aniket/utils"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
)

const ORGANISATION_PREFIX = "/remiges/rigel/"

type getConfigResponse struct {
	App         *string  `json:"app,omitempty"`
	Module      *string  `json:"module,omitempty"`
	Version     *int     `json:"ver,omitempty"`
	Config      *string  `json:"config,omitempty"`
	Description string   `json:"description,omitempty"`
	Values      []values `json:"values,omitempty"`
}

type values struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Config_get: handles the GET /configget request
// func Config_get(c *gin.Context, s *service.Service) {
// 	lh := s.LogHarbour
// 	lh.Log("Config_get request received")

// 	client := s.Dependencies["client"].(*etcd.EtcdStorage)
// 	var response getConfigResponse
// 	var queryParams utils.GetConfigRequestParams
// 	err := c.ShouldBindQuery(&queryParams)
// 	if err != nil {
// 		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(wscutils.ERRCODE_INVALID_REQUEST, nil, err.Error())}))
// 		lh.Debug0().LogActivity("error while binding json request error:", err.Error)
// 		return
// 	}

// 	keyStr := ORGANISATION_PREFIX + *queryParams.App + "/" + *queryParams.Module + "/" + strconv.Itoa(queryParams.Version) + "/" + *queryParams.Config

// 	getValue, err := client.GetWithPrefix(c, keyStr)
// 	if err != nil {
// 		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(wscutils.ErrcodeMissing, nil, err.Error())}))
// 		lh.Debug0().LogActivity("error while get data from db error:", err.Error)
// 		return
// 	}

// 	// set response fields
// 	bindGetConfigResponse(&response, &queryParams, getValue)

// 	lh.Log(fmt.Sprintf("Record found: %v", map[string]any{"key with --prefix": keyStr, "value": response}))
// 	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(response))
// }

func Config_get(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("Config_get request received")

	client := s.Dependencies["client"].(*etcd.EtcdStorage)
	var response getConfigResponse
	var queryParams utils.GetConfigRequestParams
	err := c.ShouldBindQuery(&queryParams)
	if err != nil {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(wscutils.ERRCODE_INVALID_REQUEST, nil, err.Error())}))
		lh.Debug0().LogActivity("error while binding json request error:", err.Error)
		return
	}

	keyStr := ORGANISATION_PREFIX + *queryParams.App + "/" + *queryParams.Module + "/" + strconv.Itoa(queryParams.Version) + "/" + *queryParams.Config

	getValue, err := client.GetWithPrefix(c, keyStr)
	if err != nil {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(wscutils.ErrcodeMissing, nil, err.Error())}))
		lh.Debug0().LogActivity("error while get data from db error:", err.Error)
		return
	}

	fmt.Println(">>>>>>>>>>>>>")

	// set response fields
	bindGetConfigResponse(&response, &queryParams, getValue)

	lh.Log(fmt.Sprintf("Record found: %v", map[string]any{"key with --prefix": keyStr, "value": response}))
	te := make([]*etcdls.Node, 0)
	arr, _ := etcdls.BuildTree(te, getValue)
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(arr))
}

// Config_list: handles the GET /configlist request
func Config_list(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("v request received")

	client := s.Dependencies["client"].(*etcd.EtcdStorage)
	var response getConfigResponse
	var queryParams utils.GetConfigRequestParams
	err := c.ShouldBindQuery(&queryParams)
	if err != nil {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(wscutils.ERRCODE_INVALID_REQUEST, nil, err.Error())}))
		lh.Debug0().LogActivity("error while binding json request error:", err.Error)
		return
	}

	keyStr := ORGANISATION_PREFIX + *queryParams.App + "/" + *queryParams.Module + "/" + strconv.Itoa(queryParams.Version) + "/"

	getValue, err := client.GetWithPrefix(c, keyStr)
	if err != nil {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(wscutils.ErrcodeMissing, nil, err.Error())}))
		lh.Debug0().LogActivity("error while get data from db error:", err.Error)
		return
	}

	// set response fields
	bindGetConfigResponse(&response, &queryParams, getValue)

	lh.Log(fmt.Sprintf("Record found: %v", map[string]any{"key with --prefix": keyStr, "value": response}))
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(response))
}

// // Config_set: handles the POST /userset request
// func Config_set(c *gin.Context, s *service.Service) {
// 	lh := s.LogHarbour
// 	lh.Log("User_get request received")
// 	client := s.Dependencies["client"].(*clientv3.Client)

// 	var request utils.CreateConfigRequest

// 	// Unmarshal JSON request into user struct
// 	err := wscutils.BindJSON(c, &request)
// 	if err != nil {
// 		// l.LogActivity("Error Unmarshalling JSON to struct:", logharbour.DebugInfo{Variables: map[string]any{"Error": err.Error()}})
// 		return
// 	}

// 	valError := wscutils.WscValidate(request, getValsForConfigCreateReqError)
// 	if len(valError) > 0 {
// 		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, valError))
// 		return
// 	}

// 	// create key for DB using req
// 	keyStr := ORGANISATION_PREFIX + *request.App + "/" + *request.Module + "/" + strconv.Itoa(*request.Version) + "/" + *request.Config

// 	if &request.Description != nil || len(request.Description) != 0 {
// 		fmt.Println("INSIDE DESCRIP NILL >>> keyStr:", keyStr)
// 		ke, err := client.Put(c, keyStr+"/description", request.Description)
// 		fmt.Println("ke:", ke, ">>> err:", err)
// 	}
// 	keyStr += "/description/" + *request.Name
// 	_, err = client.Put(c, keyStr, *request.Value)

// 	if err != nil {
// 		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(wscutils.ErrcodeMissing, nil, err.Error())}))
// 		// lh.Debug0().Log()
// 		return
// 	}

// 	// lh.Log(fmt.Sprintf("User found: %v", map[string]any{"user": "userResp"}))
// 	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(map[string]any{"response": "config_set_success"}))
// }

// bindGetConfigResponse is specifically used in Cinfig_get to bing and set the response
func bindGetConfigResponse(response *getConfigResponse, queryParams *utils.GetConfigRequestParams, getValue map[string]string) {
	response.App = queryParams.App
	response.Module = queryParams.Module
	response.Version = &queryParams.Version
	response.Config = queryParams.Config
	for key, vals := range getValue {

		arry := strings.Split(key, "/")
		keyStr := arry[len(arry)-1]
		if strings.EqualFold(keyStr, "description") {
			response.Description = vals
			continue
		} else {

			response.Values = append(response.Values, values{
				Name:  keyStr,
				Value: vals,
			})
		}

	}
}

func getValsForConfigCreateReqError(err validator.FieldError) []string {
	validationErrorVals := utils.GetErrorValidationMapByAPIName("config_create")
	return utils.CommonValidation(validationErrorVals, err)
}

// func buildTree(c *gin.Context, client *etcd.EtcdStorage, path string) (*etcdls.Node, error) {

// 	resp, err := client.GetWithPrefix(c, "/")
// 	if err != nil {
// 		fmt.Println(">>>>>>> Error:", err.Error())
// 		return nil, err
// 	}

// 	tt := etcdls.BuildTree(resp)

// 	return tt, nil
// }
