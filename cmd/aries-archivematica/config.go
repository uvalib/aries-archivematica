package main

import (
	"flag"
	"os"
	"strconv"
)

type configItem struct {
	flag string
	env  string
	desc string
}

type configStringItem struct {
	value string
	configItem
}

type configBoolItem struct {
	value bool
	configItem
}

type configData struct {
	listenPort         configStringItem
	dbHost              configStringItem
	dbName              configStringItem
	dbUser              configStringItem
	dbPass              configStringItem
	dbAllowOldPasswords configBoolItem
	useHttps           configBoolItem
	sslCrt             configStringItem
	sslKey             configStringItem
}

var config configData

func init() {
	config.listenPort = configStringItem{value: "", configItem: configItem{flag: "l", env: "ARIES_ARCHIVEMATICA_LISTEN_PORT", desc: "listen port"}}
	config.dbHost = configStringItem{value: "", configItem: configItem{flag: "h", env: "ARIES_ARCHIVEMATICA_DB_HOST", desc: "database host"}}
	config.dbName = configStringItem{value: "", configItem: configItem{flag: "n", env: "ARIES_ARCHIVEMATICA_DB_NAME", desc: "database name"}}
	config.dbUser = configStringItem{value: "", configItem: configItem{flag: "u", env: "ARIES_ARCHIVEMATICA_DB_USER", desc: "database user"}}
	config.dbPass = configStringItem{value: "", configItem: configItem{flag: "p", env: "ARIES_ARCHIVEMATICA_DB_PASS", desc: "database password"}}
	config.dbAllowOldPasswords = configBoolItem{value: false, configItem: configItem{flag: "o", env: "ARIES_ARCHIVEMATICA_DB_ALLOW_OLD_PASSWORDS", desc: "allow old database passwords"}}
	config.useHttps = configBoolItem{value: false, configItem: configItem{flag: "s", env: "ARIES_ARCHIVEMATICA_USE_HTTPS", desc: "use https"}}
	config.sslCrt = configStringItem{value: "", configItem: configItem{flag: "c", env: "ARIES_ARCHIVEMATICA_SSL_CRT", desc: "ssl crt"}}
	config.sslKey = configStringItem{value: "", configItem: configItem{flag: "k", env: "ARIES_ARCHIVEMATICA_SSL_KEY", desc: "ssl key"}}
}

func getBoolEnv(optEnv string) bool {
	value, _ := strconv.ParseBool(os.Getenv(optEnv))

	return value
}

func ensureConfigStringSet(item *configStringItem) bool {
	isSet := true

	if item.value == "" {
		isSet = false
		logger.Printf("[ERROR] %s is not set, use %s variable or -%s flag", item.desc, item.env, item.flag)
	}

	return isSet
}

func flagStringVar(item *configStringItem) {
	flag.StringVar(&item.value, item.flag, os.Getenv(item.env), item.desc)
}

func flagBoolVar(item *configBoolItem) {
	flag.BoolVar(&item.value, item.flag, getBoolEnv(item.env), item.desc)
}

func getConfigValues() {
	// get values from the command line first, falling back to environment variables
	flagStringVar(&config.listenPort)
	flagStringVar(&config.dbHost)
	flagStringVar(&config.dbName)
	flagStringVar(&config.dbUser)
	flagStringVar(&config.dbPass)
	flagBoolVar(&config.dbAllowOldPasswords)
	flagBoolVar(&config.useHttps)
	flagStringVar(&config.sslCrt)
	flagStringVar(&config.sslKey)

	flag.Parse()

	// check each required option, displaying a warning for empty values.
	// die if any of them are not set
	configOK := true
	configOK = ensureConfigStringSet(&config.listenPort) && configOK
	configOK = ensureConfigStringSet(&config.dbHost) && configOK
	configOK = ensureConfigStringSet(&config.dbName) && configOK
	configOK = ensureConfigStringSet(&config.dbUser) && configOK
	configOK = ensureConfigStringSet(&config.dbPass) && configOK
	if config.useHttps.value == true {
		configOK = ensureConfigStringSet(&config.sslCrt) && configOK
		configOK = ensureConfigStringSet(&config.sslKey) && configOK
	}

	if configOK == false {
		flag.Usage()
		os.Exit(1)
	}

	logger.Printf("[CONFIG] listenPort          = [%s]", config.listenPort.value)
	logger.Printf("[CONFIG] dbHost              = [%s]", config.dbHost.value)
	logger.Printf("[CONFIG] dbName              = [%s]", config.dbName.value)
	logger.Printf("[CONFIG] dbUser              = [%s]", config.dbUser.value)
	logger.Printf("[CONFIG] dbPass              = [REDACTED]")
	logger.Printf("[CONFIG] dbAllowOldPasswords = [%s]", strconv.FormatBool(config.dbAllowOldPasswords.value))
	logger.Printf("[CONFIG] useHttps            = [%s]", strconv.FormatBool(config.useHttps.value))
	logger.Printf("[CONFIG] sslCrt              = [%s]", config.sslCrt.value)
	logger.Printf("[CONFIG] sslKey              = [%s]", config.sslKey.value)
}
