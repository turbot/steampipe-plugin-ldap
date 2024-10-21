STEAMPIPE_INSTALL_DIR ?= ~/.steampipe
BUILD_TAGS = netgo
install:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/hub.steampipe.io/plugins/turbot/ldap@latest/steampipe-plugin-ldap.plugin -tags "${BUILD_TAGS}" *.go
