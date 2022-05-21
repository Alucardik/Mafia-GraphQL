package main

import "MafiaGQL_server/utils"

const (
	DB_PORT        = 27017
	QUERY_ENDPOINT = "/query"
)

var (
	PORT        = utils.SetEnvVar("PORT", "8080")
	DB_HOST     = utils.SetEnvVar("DB_HOST", "localhost")
	DB_USERNAME = utils.SetEnvVar("DB_USERNAME", "root")
	DB_PASS     = utils.SetEnvVar("DB_PASS", "example")
)
