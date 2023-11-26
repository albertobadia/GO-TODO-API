package api

import "os"

var SECRET_KEY = os.Getenv("SECRET_KEY")
var IS_TESTING = os.Getenv("IS_TESTING") == "true"

var POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
var POSTGRES_PORT = os.Getenv("POSTGRES_PORT")
var POSTGRES_USER = os.Getenv("POSTGRES_USER")
var POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
var POSTGRES_DB = os.Getenv("POSTGRES_DB")
