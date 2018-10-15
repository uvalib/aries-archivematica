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
	applicationAPIUser    configStringItem
	applicationAPIKey     configStringItem
	applicationAPIUrlTemplate     configStringItem
	applicationDBProtocol configStringItem
	applicationDBHost     configStringItem
	applicationDBName     configStringItem
	applicationDBUser     configStringItem
	applicationDBPass     configStringItem
	storageAPIUser        configStringItem
	storageAPIKey         configStringItem
	storageAPIUrlTemplate         configStringItem
	storageDBProtocol     configStringItem
	storageDBHost         configStringItem
	storageDBName         configStringItem
	storageDBUser         configStringItem
	storageDBPass         configStringItem
	adminUrlTemplate      configStringItem
	useHttps              configBoolItem
	sslCrt                configStringItem
	sslKey                configStringItem
}

var config configData

func init() {
	config.listenPort = configStringItem{value: "", configItem: configItem{flag: "l", env: "ARIES_ARCHIVEMATICA_LISTEN_PORT", desc: "listen port"}}
	config.applicationAPIUser = configStringItem{value: "", configItem: configItem{flag: "a", env: "ARIES_ARCHIVEMATICA_APPLICATION_API_USER", desc: "application API user"}}
	config.applicationAPIKey = configStringItem{value: "", configItem: configItem{flag: "y", env: "ARIES_ARCHIVEMATICA_APPLICATION_API_KEY", desc: "application API key"}}
	config.applicationAPIUrlTemplate = configStringItem{value: "", configItem: configItem{flag: "w", env: "ARIES_ARCHIVEMATICA_APPLICATION_API_URL_TEMPLATE", desc: "application API url"}}
	config.applicationDBProtocol = configStringItem{value: "", configItem: configItem{flag: "r", env: "ARIES_ARCHIVEMATICA_APPLICATION_DB_PROT", desc: "application DB protocol"}}
	config.applicationDBHost = configStringItem{value: "", configItem: configItem{flag: "h", env: "ARIES_ARCHIVEMATICA_APPLICATION_DB_HOST", desc: "application DB host/file"}}
	config.applicationDBName = configStringItem{value: "", configItem: configItem{flag: "n", env: "ARIES_ARCHIVEMATICA_APPLICATION_DB_NAME", desc: "application DB name"}}
	config.applicationDBUser = configStringItem{value: "", configItem: configItem{flag: "u", env: "ARIES_ARCHIVEMATICA_APPLICATION_DB_USER", desc: "application DB user"}}
	config.applicationDBPass = configStringItem{value: "", configItem: configItem{flag: "p", env: "ARIES_ARCHIVEMATICA_APPLICATION_DB_PASS", desc: "application DB password"}}
	config.storageAPIUser = configStringItem{value: "", configItem: configItem{flag: "A", env: "ARIES_ARCHIVEMATICA_STORAGE_API_USER", desc: "storage service API user"}}
	config.storageAPIKey = configStringItem{value: "", configItem: configItem{flag: "Y", env: "ARIES_ARCHIVEMATICA_STORAGE_API_KEY", desc: "storage service API key"}}
	config.storageAPIUrlTemplate = configStringItem{value: "", configItem: configItem{flag: "W", env: "ARIES_ARCHIVEMATICA_STORAGE_API_URL_TEMPLATE", desc: "storage service API url"}}
	config.storageDBProtocol = configStringItem{value: "", configItem: configItem{flag: "R", env: "ARIES_ARCHIVEMATICA_STORAGE_DB_PROT", desc: "storage service DB protocol"}}
	config.storageDBHost = configStringItem{value: "", configItem: configItem{flag: "H", env: "ARIES_ARCHIVEMATICA_STORAGE_DB_HOST", desc: "storage service DB host/file"}}
	config.storageDBName = configStringItem{value: "", configItem: configItem{flag: "N", env: "ARIES_ARCHIVEMATICA_STORAGE_DB_NAME", desc: "storage service DB name"}}
	config.storageDBUser = configStringItem{value: "", configItem: configItem{flag: "U", env: "ARIES_ARCHIVEMATICA_STORAGE_DB_USER", desc: "storage service DB user"}}
	config.storageDBPass = configStringItem{value: "", configItem: configItem{flag: "P", env: "ARIES_ARCHIVEMATICA_STORAGE_DB_PASS", desc: "storage service DB password"}}
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
	flagStringVar(&config.applicationAPIUser)
	flagStringVar(&config.applicationAPIKey)
	flagStringVar(&config.applicationAPIUrlTemplate)
	flagStringVar(&config.applicationDBProtocol)
	flagStringVar(&config.applicationDBHost)
	flagStringVar(&config.applicationDBName)
	flagStringVar(&config.applicationDBUser)
	flagStringVar(&config.applicationDBPass)
	flagStringVar(&config.storageAPIUser)
	flagStringVar(&config.storageAPIKey)
	flagStringVar(&config.storageAPIUrlTemplate)
	flagStringVar(&config.storageDBProtocol)
	flagStringVar(&config.storageDBHost)
	flagStringVar(&config.storageDBName)
	flagStringVar(&config.storageDBUser)
	flagStringVar(&config.storageDBPass)
	flagStringVar(&config.adminUrlTemplate)
	flagBoolVar(&config.useHttps)
	flagStringVar(&config.sslCrt)
	flagStringVar(&config.sslKey)

	flag.Parse()

	// check each required option, displaying a warning for empty values.
	// die if any of them are not set
	configOK := true
	configOK = ensureConfigStringSet(&config.listenPort) && configOK
	configOK = ensureConfigStringSet(&config.applicationAPIUser) && configOK
	configOK = ensureConfigStringSet(&config.applicationAPIKey) && configOK
	configOK = ensureConfigStringSet(&config.applicationAPIUrlTemplate) && configOK
	configOK = ensureConfigStringSet(&config.applicationDBProtocol) && configOK
	configOK = ensureConfigStringSet(&config.applicationDBHost) && configOK
	configOK = ensureConfigStringSet(&config.applicationDBName) && configOK
	configOK = ensureConfigStringSet(&config.applicationDBUser) && configOK
	configOK = ensureConfigStringSet(&config.applicationDBPass) && configOK
	configOK = ensureConfigStringSet(&config.storageAPIUser) && configOK
	configOK = ensureConfigStringSet(&config.storageAPIKey) && configOK
	configOK = ensureConfigStringSet(&config.storageAPIUrlTemplate) && configOK
	configOK = ensureConfigStringSet(&config.storageDBProtocol) && configOK
	configOK = ensureConfigStringSet(&config.storageDBHost) && configOK
	configOK = ensureConfigStringSet(&config.storageDBName) && configOK
	configOK = ensureConfigStringSet(&config.storageDBUser) && configOK
	configOK = ensureConfigStringSet(&config.storageDBPass) && configOK
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
	logger.Printf("[CONFIG] applicationAPIUser        = [%s]", config.applicationAPIUser.value)
	logger.Printf("[CONFIG] applicationAPIKey         = [REDACTED]")
	logger.Printf("[CONFIG] applicationAPIUrlTemplate = [%s]", config.applicationAPIUrlTemplate.value)
	logger.Printf("[CONFIG] applicationDBProtocol     = [%s]", config.applicationDBProtocol.value)
	logger.Printf("[CONFIG] applicationDBHost         = [%s]", config.applicationDBHost.value)
	logger.Printf("[CONFIG] applicationDBName         = [%s]", config.applicationDBName.value)
	logger.Printf("[CONFIG] applicationDBUser         = [%s]", config.applicationDBUser.value)
	logger.Printf("[CONFIG] applicationDBPass         = [REDACTED]")
	logger.Printf("[CONFIG] storageAPIUser            = [%s]", config.storageAPIUser.value)
	logger.Printf("[CONFIG] storageAPIKey             = [REDACTED]")
	logger.Printf("[CONFIG] storageAPIUrlTemplate     = [%s]", config.storageAPIUrlTemplate.value)
	logger.Printf("[CONFIG] storageDBProtocol         = [%s]", config.storageDBProtocol.value)
	logger.Printf("[CONFIG] storageDBHost             = [%s]", config.storageDBHost.value)
	logger.Printf("[CONFIG] storageDBName             = [%s]", config.storageDBName.value)
	logger.Printf("[CONFIG] storageDBUser             = [%s]", config.storageDBUser.value)
	logger.Printf("[CONFIG] storageDBPass             = [REDACTED]")
	logger.Printf("[CONFIG] adminUrlTemplate          = [%s]", config.adminUrlTemplate.value)
	logger.Printf("[CONFIG] useHttps              = [%s]", strconv.FormatBool(config.useHttps.value))
	logger.Printf("[CONFIG] sslCrt                = [%s]", config.sslCrt.value)
	logger.Printf("[CONFIG] sslKey                = [%s]", config.sslKey.value)
}
