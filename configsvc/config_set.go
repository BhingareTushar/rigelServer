package configsvc

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/logharbour/logharbour"
	"github.com/remiges-tech/rigel"
)

type configset struct {
	App    string `json:"app" validate:"required"`
	Module string `json:"module" validate:"required"`
	Ver    int    `json:"ver" validate:"required"`
	Config string `json:"config" validate:"required"`
	Key    string `json:"key" validate:"required"`
	Value  any    `json:"value" validate:"required"`
}

type configupdate struct {
	App         string `json:"app" validate:"required"`
	Module      string `json:"module" validate:"required"`
	Ver         int    `json:"ver" validate:"required"`
	Config      string `json:"config" validate:"required"`
	Description string `json:"description" validate:"required"`
	Values      []struct {
		Name  string `json:"name" validate:"required"`
		Value string `json:"value" validate:"required"`
	} `json:"values" validate:"required"`
}

func Config_set(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Log("Starting execution of Config_set()")

	var configset configset
	err := wscutils.BindJSON(c, &configset)
	if err != nil {
		l.LogActivity("error while binding json", err)
		return
	}

	validationErrors := validateConfigset(configset, c)
	if len(validationErrors) > 0 {
		l.LogDebug("Validation errors:", logharbour.DebugInfo{Variables: map[string]any{"validationErrors": validationErrors}})
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	r := s.Dependencies["r"].(*rigel.Rigel)
	r.WithApp(configset.App).WithModule(configset.Module).WithVersion(configset.Ver).WithConfig(configset.Config)
	val := fmt.Sprintf("%#v", configset.Value)
	err = r.Set(c, configset.Key, val)
	if err != nil {
		l.LogActivity("error while setting value in etcd:", err)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse("unable_to_set"))
		return
	} else {
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: "data set successfully", Messages: []wscutils.ErrorMessage{}})
	}
}

// validateConfigset performs validation for the Configset.
func validateConfigset(config configset, c *gin.Context) []wscutils.ErrorMessage {
	// Validate the request body
	validationErrors := wscutils.WscValidate(config, config.getVals)

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return validationErrors
}

// getVals returns validation error details based on the field and tag.
func (config *configset) getVals(err validator.FieldError) []string {
	return nil
}
