package api

import "os"

var SECRET_KEY = os.Getenv("SECRET_KEY")
var IS_TESTING = os.Getenv("IS_TESTING") == "true"
