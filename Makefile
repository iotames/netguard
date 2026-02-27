APP_VERSION=v1.1.1

GEO_DB_URL = https://github.com/P3TERX/GeoLite.mmdb/releases/download/2026.02.25/GeoLite2-City.mmdb

# AMIS jssdk

JSSDK_URL = https://github.com/baidu/amis/releases/download/6.13.0/jssdk.tar.gz

# 根据操作系统设置目标文件名和链接库
ifeq ($(OS),Windows_NT)
BUILD_FILE_NAME=netguard.exe
BUILD_TIME:=$(shell powershell -Command "Get-Date -Format 'yyyy-MM-dd_HH_mm'")
# BUILD_TIME:=%date:~0,4%-%date:~5,2%-%date:~8,2%_%time:~0,2%_%time:~3,2%
COPY=copy
DIRSEP=\\
RM=del /Q
MKDIR=mkdir
RUN_FILE_NAME="Run.bat"
CMD_PRE=set CGO_ENABLED=1
else
BUILD_FILE_NAME=netguard
BUILD_TIME:=$(shell date +%Y-%m-%d_%H_%M)
COPY=cp -rf
DIRSEP=/
RM=rm -rf
MKDIR=mkdir -p
RUN_FILE_NAME="Run.sh"
CMD_PRE=export CGO_ENABLED=1
# apt install libpcap-dev	
endif

# Recipe（配方） 是 Makefile 中在 target（目标）下执行的命令
# Makefile 以 Tab 开头的命令，在它之前必须定义 target
GEO_DB_FILE_NAME = GeoLite2-City.mmdb
GEO_DB_FILE_PATH = main$(DIRSEP)$(GEO_DB_FILE_NAME)
JSSDK_TAR = main$(DIRSEP)jssdk.tar.gz
ZIP_FILE = main$(DIRSEP)NetGuard_$(APP_VERSION).zip
MAIN_HTML=main$(DIRSEP)amis.html
RELEASE_DIR=main$(DIRSEP)release
BUILD_FILE_PATH=main$(DIRSEP)$(BUILD_FILE_NAME)
MAIN_LOG_FILES=main$(DIRSEP)*.log
RELEASE_LOG_FILES=$(RELEASE_DIR)$(DIRSEP)*.log
SRC_AMIS_DIR=main$(DIRSEP)static$(DIRSEP)amis
SRC_PAGES_DIR=main$(DIRSEP)static$(DIRSEP)pages
RELEASE_AMIS_DIR=$(RELEASE_DIR)$(DIRSEP)static$(DIRSEP)amis
RELEASE_PAGES_DIR=$(RELEASE_DIR)$(DIRSEP)static$(DIRSEP)pages

# 编译文件
# apt install libpcap-dev
build: $(BUILD_FILE_PATH) $(SRC_AMIS_DIR)$(DIRSEP)sdk.js $(SRC_PAGES_DIR)$(DIRSEP)index.json

# 清理编译文件
clean:
	-$(RM) $(BUILD_FILE_PATH)
	-$(RM) $(RELEASE_DIR)$(DIRSEP)$(BUILD_FILE_NAME)
	-$(RM) $(RELEASE_DIR)$(DIRSEP)$(RUN_FILE_NAME)
	-$(RM) $(RELEASE_LOG_FILES)
	-$(RM) $(MAIN_LOG_FILES)
	-$(RM) $(ZIP_FILE)

# 执行编译程序
# 环境变量作用域问题：在 Makefile 中，每行命令都在独立子进程中执行，set CGO_ENABLED=1只影响当前行
$(BUILD_FILE_PATH):
	go mod tidy
	$(CMD_PRE) && go build -v -o $(BUILD_FILE_PATH) -trimpath -ldflags "-X 'main.BuildTime=$(BUILD_TIME)' -X 'main.Version=$(APP_VERSION)' " ./main

# 下载GEO数据库以解析IP地址
$(GEO_DB_FILE_PATH): 
ifeq ($(OS),Windows_NT)
# 	powershell -Command "Invoke-WebRequest -Uri $(GEO_DB_URL) -OutFile $(GEO_DB_FILE_PATH) -Verbose"
	powershell -Command "Start-BitsTransfer -Source $(GEO_DB_URL) -Destination $(GEO_DB_FILE_PATH) -DisplayName 'GEO DB download'"
else
	wget -c $(GEO_DB_URL) -O $(GEO_DB_FILE_PATH)
endif
	@echo "GEO_DB_FILE_PATH Download: $(GEO_DB_FILE_PATH)"

# 下载并解压AMIS jssdk
$(SRC_AMIS_DIR)$(DIRSEP)sdk.js: 
	-$(MKDIR) $(SRC_AMIS_DIR)
ifeq ($(OS),Windows_NT)
	if not exist $(JSSDK_TAR) powershell -Command "Invoke-WebRequest -Uri $(JSSDK_URL) -OutFile $(JSSDK_TAR)"
	powershell -Command "tar -xzf '$(JSSDK_TAR)' -C '$(SRC_AMIS_DIR)'"
else
	if [ ! -f $(JSSDK_TAR) ]; then wget -c $(JSSDK_URL) -O $(JSSDK_TAR); fi
	tar -xzf $(JSSDK_TAR) -C $(SRC_AMIS_DIR)
endif

run: $(BUILD_FILE_PATH)
ifeq ($(OS),Windows_NT)
	main$(DIRSEP)$(RUN_FILE_NAME)
else
	./man$(DIRSEP)$(RUN_FILE_NAME)
endif

# 整理发布包文件
release: $(BUILD_FILE_PATH) $(GEO_DB_FILE_PATH) $(SRC_AMIS_DIR)$(DIRSEP)sdk.js
	-$(MKDIR) $(RELEASE_DIR)
	-$(MKDIR) $(RELEASE_PAGES_DIR)
	-$(MKDIR) $(RELEASE_AMIS_DIR)
	$(COPY) $(BUILD_FILE_PATH) $(RELEASE_DIR)$(DIRSEP)$(BUILD_FILE_NAME)
	-$(COPY) $(GEO_DB_FILE_PATH) $(RELEASE_DIR)$(DIRSEP)$(GEO_DB_FILE_NAME)
	$(COPY) main$(DIRSEP)$(RUN_FILE_NAME) $(RELEASE_DIR)$(DIRSEP)$(RUN_FILE_NAME)
	$(COPY) $(MAIN_HTML) $(RELEASE_DIR)$(DIRSEP)amis.html
	$(COPY) $(SRC_PAGES_DIR)$(DIRSEP)* $(RELEASE_PAGES_DIR)$(DIRSEP)
	$(COPY) $(SRC_AMIS_DIR)$(DIRSEP)helper.css $(RELEASE_AMIS_DIR)$(DIRSEP)helper.css
	$(COPY) $(SRC_AMIS_DIR)$(DIRSEP)iconfont.css $(RELEASE_AMIS_DIR)$(DIRSEP)iconfont.css
	$(COPY) $(SRC_AMIS_DIR)$(DIRSEP)rest.js $(RELEASE_AMIS_DIR)$(DIRSEP)rest.js
	$(COPY) $(SRC_AMIS_DIR)$(DIRSEP)sdk.css $(RELEASE_AMIS_DIR)$(DIRSEP)sdk.css
	$(COPY) $(SRC_AMIS_DIR)$(DIRSEP)sdk.js $(RELEASE_AMIS_DIR)$(DIRSEP)sdk.js

# 生成ZIP包
zip: release
ifeq ($(OS),Windows_NT)
	powershell -Command "Compress-Archive -Path '$(RELEASE_DIR)' -DestinationPath '$(ZIP_FILE)' -Force"
else
	zip -r $(ZIP_FILE) $(RELEASE_DIR)
endif
	@echo "ZIP Generate Done: $(ZIP_FILE)"

.PHONY:	build run release zip clean
