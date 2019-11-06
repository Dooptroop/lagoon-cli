package cmd

import (
	"github.com/amazeeio/lagoon-cli/output"
)

// config vars
var lagoonHostname string
var lagoonPort string
var lagoonGraphQL string
var lagoonToken string

// variable vars
var variableValue string
var variableName string
var variableScope string

var jsonPatch string
var revealValue bool
var listAllProjects bool
var noHeader bool

var cmdProjectName string
var cmdProjectEnvironment string

var remoteID string

var notificationName string
var notificationNewName string
var notificationWebhook string
var notificationChannel string

var outputOptions = output.Options{
	Header: false,
	CSV:    false,
	JSON:   false,
	Pretty: false,
}
