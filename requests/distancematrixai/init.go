package requests

import (
	"fmt"
	"os"
)

var (
	distanceMatrixAIaccessToken       = os.Getenv("DISTANCEMATRIXAI_ACCESS_TOKEN")
	distanceMatrixAIendingCredentials = fmt.Sprintf("&key=%s", distanceMatrixAIaccessToken)
)
