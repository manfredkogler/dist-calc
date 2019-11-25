package requests

import (
	"fmt"
	"os"
)

var (
	hereAppID                  = os.Getenv("HERE_APP_ID")
	hereAppCode                = os.Getenv("HERE_APP_CODE")
	hereAPIstartingCredentials = fmt.Sprintf("app_id=%s&app_code=%s&", hereAppID, hereAppCode)
)
