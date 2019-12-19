package requests

import (
	"fmt"
	"os"
)

var (
	hereAPIKey                 = os.Getenv("HERE_API_KEY")
	hereAPIstartingCredentials = fmt.Sprintf("apiKey=%s&", hereAPIKey)
)
