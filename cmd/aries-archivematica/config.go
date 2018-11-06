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
	listenPort            configStringItem
	storageAPIUser        configStringItem
	storageAPIKey         configStringItem
	storageAPIUrlTemplate configStringItem
	adminUrlTemplate      configStringItem
	useHttps              configBoolItem
	sslCrt                configStringItem
	sslKey                configStringItem
}

var config configData

func init() {
	config.listenPort = configStringItem{value: "", configItem: configItem{flag: "l", env: "ARIES_ARCHIVEMATICA_LISTEN_PORT", desc: "listen port"}}
	config.storageAPIUser = configStringItem{value: "", configItem: configItem{flag: "A", env: "ARIES_ARCHIVEMATICA_STORAGE_API_USER", desc: "storage service API user"}}
	config.storageAPIKey = configStringItem{value: "", configItem: configItem{flag: "Y", env: "ARIES_ARCHIVEMATICA_STORAGE_API_KEY", desc: "storage service API key"}}
	config.storageAPIUrlTemplate = configStringItem{value: "", configItem: configItem{flag: "W", env: "ARIES_ARCHIVEMATICA_STORAGE_API_URL_TEMPLATE", desc: "storage service API url"}}
	config.adminUrlTemplate = configStringItem{value: "", configItem: configItem{flag: "t", env: "ARIES_ARCHIVEMATICA_ADMIN_URL_TEMPLATE", desc: "admin url template"}}
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
	flagStringVar(&config.storageAPIUser)
	flagStringVar(&config.storageAPIKey)
	flagStringVar(&config.storageAPIUrlTemplate)
	flagStringVar(&config.adminUrlTemplate)
	flagBoolVar(&config.useHttps)
	flagStringVar(&config.sslCrt)
	flagStringVar(&config.sslKey)

	flag.Parse()

	// check each required option, displaying a warning for empty values.
	// die if any of them are not set
	configOK := true
	configOK = ensureConfigStringSet(&config.listenPort) && configOK
	configOK = ensureConfigStringSet(&config.storageAPIUser) && configOK
	configOK = ensureConfigStringSet(&config.storageAPIKey) && configOK
	configOK = ensureConfigStringSet(&config.storageAPIUrlTemplate) && configOK
	configOK = ensureConfigStringSet(&config.adminUrlTemplate) && configOK
	if config.useHttps.value == true {
		configOK = ensureConfigStringSet(&config.sslCrt) && configOK
		configOK = ensureConfigStringSet(&config.sslKey) && configOK
	}

	if configOK == false {
		flag.Usage()
		os.Exit(1)
	}

	logger.Printf("[CONFIG] listenPort                = [%s]", config.listenPort.value)
	logger.Printf("[CONFIG] storageAPIUser            = [%s]", config.storageAPIUser.value)
	logger.Printf("[CONFIG] storageAPIKey             = [REDACTED]")
	logger.Printf("[CONFIG] storageAPIUrlTemplate     = [%s]", config.storageAPIUrlTemplate.value)
	logger.Printf("[CONFIG] adminUrlTemplate          = [%s]", config.adminUrlTemplate.value)
	logger.Printf("[CONFIG] useHttps                  = [%s]", strconv.FormatBool(config.useHttps.value))
	logger.Printf("[CONFIG] sslCrt                    = [%s]", config.sslCrt.value)
	logger.Printf("[CONFIG] sslKey                    = [%s]", config.sslKey.value)
}
