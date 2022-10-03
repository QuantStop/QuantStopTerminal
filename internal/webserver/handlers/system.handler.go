package handlers

/*import (
	"encoding/json"
	"github.com/quantstop/quantstopterminal/internal"
	"github.com/quantstop/quantstopterminal/internal/database/models"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/internal/webserver/errors"
	"github.com/quantstop/quantstopterminal/internal/webserver/write"
	"net/http"
	"net/url"
)

type setSubsystemRequest struct {
	Subsystem string `json:"subsystem"`
	Enable    bool   `json:"enable"`
}

type setSysConfigRequest struct {
	ApiUrl     string `json:"apiUrl"`
	GoMaxProcs string `json:"maxProcs"`
}



func SetSubsystem(bot internal.IEngine, user *models.User, w http.ResponseWriter, r *http.Request) http.HandlerFunc {

	decoder := json.NewDecoder(r.Body)
	req := setSubsystemRequest{}
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		return write.Error(errors.NoJSONBody)
	}

	if req.Subsystem == "" {
		return write.Error(errors.InvalidInput)
	}
	temp := ""
	if req.Enable {
		temp = "true"
	} else {
		temp = "false"
	}
	log.Debug(log.Webserver, "set "+req.Subsystem+" "+temp)

	// set subsystem status
	if err = bot.SetSubsystem(req.Subsystem, req.Enable); err != nil {
		return write.Error(err)
	}

	return write.Success()
}

func SetSystemConfig(bot internal.IEngine, user *models.User, w http.ResponseWriter, r *http.Request) http.HandlerFunc {

	decoder := json.NewDecoder(r.Body)
	req := setSysConfigRequest{}
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		return write.Error(errors.NoJSONBody)
	}

	if req.ApiUrl == "" {
		return write.Error(errors.InvalidInput)
	}

	if !isValidUrl(req.ApiUrl) {
		return write.Error(errors.InvalidInput)
	}

	// set config
	err = bot.SetConfig(req.ApiUrl, req.GoMaxProcs)
	if err != nil {
		return nil
	}

	return write.Success()
}

// isValidUrl tests a string to determine if it is a well-structured url or not.
func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
*/
