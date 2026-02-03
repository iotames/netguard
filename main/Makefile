APP_VERSION=v0.0.2
ZIP_NAME = NetGuard_$(APP_VERSION).zip
GEO_DB_FILE = GeoLite2-City.mmdb
GEO_DB_URL = https://github.com/P3TERX/GeoLite.mmdb/releases/download/2026.01.16/GeoLite2-City.mmdb

# 根据操作系统设置目标文件名和链接库
ifeq ($(OS),Windows_NT)
	BUILD_FILE_NAME=netguard.exe
	BUILD_TIME:=$(shell powershell -Command "Get-Date -Format 'yyyy-MM-dd_HH_mm'")
# 	BUILD_TIME:=%date:~0,4%-%date:~5,2%-%date:~8,2%_%time:~0,2%_%time:~3,2%
	COPY=copy
	DIRSEP=\\
	RM=del /Q
	MKDIR=mkdir
	RUN_FILE="Run.bat"
else
	BUILD_FILE_NAME=netguard
	BUILD_TIME:=$(shell date +%Y-%m-%d_%H_%M)
	COPY=cp -rf
	DIRSEP=/
	RM=rm -rf
	MKDIR=mkdir -p
	RUN_FILE="Run.sh"
endif

# 编译文件
# apt install libpcap-dev
build: $(BUILD_FILE_NAME)

# 清理编译文件
clean:
	-$(RM) $(BUILD_FILE_NAME)
	-$(RM) release$(DIRSEP)$(BUILD_FILE_NAME)
	-$(RM) release$(DIRSEP)$(RUN_FILE)
	-$(RM) $(ZIP_NAME)

# 执行编译程序
# CGO_ENABLED=1
$(BUILD_FILE_NAME):
	go build -v -o $(BUILD_FILE_NAME) -trimpath -ldflags "-X 'main.BuildTime=$(BUILD_TIME)' -X 'main.Version=$(APP_VERSION)' " .

# 下载GEO数据库以解析IP地址
$(GEO_DB_FILE): 
ifeq ($(OS),Windows_NT)
	powershell -Command "Invoke-WebRequest -Uri $(GEO_DB_URL) -OutFile $(GEO_DB_FILE)"
else
	wget -c $(GEO_DB_URL) -O $(GEO_DB_FILE)
endif
	@echo "GEO_DB_FILE Download: $(GEO_DB_FILE)"

run:
ifeq ($(OS),Windows_NT)
	$(RUN_FILE)
else
	./$(RUN_FILE)
endif

# 整理发布包文件
release: $(BUILD_FILE_NAME) $(GEO_DB_FILE)
	-$(MKDIR) release
	$(COPY) $(BUILD_FILE_NAME) release$(DIRSEP)$(BUILD_FILE_NAME)
	-$(COPY) $(GEO_DB_FILE) release$(DIRSEP)$(GEO_DB_FILE)
	$(COPY) $(RUN_FILE) release$(DIRSEP)$(RUN_FILE)

# 生成ZIP包
zip: release
ifeq ($(OS),Windows_NT)
	powershell -Command "Compress-Archive -Path 'release' -DestinationPath '$(ZIP_NAME)' -Force"
else
	zip -r $(ZIP_NAME) release
endif
	@echo "ZIP Generate Done: $(ZIP_NAME)"

.PHONY:	build run release zip clean
