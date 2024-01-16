package configsvc

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/logharbour/logharbour"
	"github.com/remiges-tech/rigel"
)

func Config_update(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Log("Starting execution of Config_set()")

	var configupdate configupdate
	err := wscutils.BindJSON(c, &configupdate)
	if err != nil {
		l.LogActivity("error while binding json", err)
		return
	}

	// Validate the user creation request
	validationErrors := validateConfigupdate(configupdate, c)
	if len(validationErrors) > 0 {
		l.LogDebug("Validation errors:", logharbour.DebugInfo{Variables: map[string]any{"validationErrors": validationErrors}})
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	r := s.Dependencies["r"].(*rigel.Rigel)
	r.WithApp(configupdate.App).WithModule(configupdate.Module).WithVersion(configupdate.Ver).WithConfig(configupdate.Config)

	for _, v := range configupdate.Values {
		err = r.Set(c, v.Name, v.Value)
		if err != nil {
			l.LogActivity("error while setting value in etcd:", err)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse("unable_to_set"))
			return
		}
	}
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: "data set successfully", Messages: []wscutils.ErrorMessage{}})
}

// validateConfigupdate performs validation for the Configupdate.
func validateConfigupdate(config configupdate, c *gin.Context) []wscutils.ErrorMessage {
	// Validate the request body
	validationErrors := wscutils.WscValidate(config, config.getValsForUser)

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return validationErrors
}

// getValsForUser returns validation error details based on the field and tag.
func (config *configupdate) getValsForUser(err validator.FieldError) []string {
	return nil
}
