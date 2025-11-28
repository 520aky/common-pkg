# =========================================
# 自动 Tag +1 并推送（支持自动进位）
# =========================================

# 拉取远程标签（忽略错误）
FETCH_TAGS := $(shell git fetch --tags >/dev/null 2>&1 || true)

# 获取最新 tag（远程 + 本地）
LATEST_TAG := $(shell git tag -l "v*" --sort=-v:refname | head -n1 2>/dev/null || echo v0.0.0)

# 提取版本号
VER_NO_V    := $(patsubst v%,%,$(LATEST_TAG))
MAJOR       := $(word 1,$(subst ., ,$(VER_NO_V)))
MINOR       := $(word 2,$(subst ., ,$(VER_NO_V)))
PATCH       := $(word 3,$(subst ., ,$(VER_NO_V)))

# 计算新版本号（支持进位）
NEXT_MAJOR  := $(MAJOR)
NEXT_MINOR  := $(MINOR)
NEXT_PATCH  := $(PATCH)

define CALC_NEXT_TAG
MAJOR=$(MAJOR); MINOR=$(MINOR); PATCH=$(PATCH); \
if [ $$PATCH -ge 99 ]; then \
  PATCH=0; MINOR=$$((MINOR+1)); \
else \
  PATCH=$$((PATCH+1)); \
fi; \
if [ $$MINOR -ge 100 ]; then \
  MINOR=0; MAJOR=$$((MAJOR+1)); \
fi; \
printf "v%d.%d.%d" $$MAJOR $$MINOR $$PATCH
endef

NEXT_TAG := $(shell sh -c '$(CALC_NEXT_TAG)')

# 当前分支
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

tag:
	@echo "📦 当前分支: $(BRANCH)"
	@echo "🔍 同步远程标签..."
	@git fetch --tags >/dev/null 2>&1 || true
	@LATEST_TAG=$$(git tag -l "v*" --sort=-v:refname | head -n1 2>/dev/null || echo v0.0.0); \
	echo "🏷️  最新 Tag: $$LATEST_TAG"; \
	read -p "📝 请输入提交信息: " msg; \
	if [ -z "$$msg" ]; then echo "❌ 提交信息不能为空！"; exit 1; fi; \
	git add .; \
	git commit -m "$$msg" || true; \
	git push origin $(BRANCH); \
	MAJOR=$(MAJOR); MINOR=$(MINOR); PATCH=$(PATCH); \
	if [ $$PATCH -ge 99 ]; then \
		PATCH=0; MINOR=$$(expr $$MINOR + 1); \
	else \
		PATCH=$$(expr $$PATCH + 1); \
	fi; \
	if [ $$MINOR -ge 100 ]; then \
		MINOR=0; MAJOR=$$(expr $$MAJOR + 1); \
	fi; \
	NEXT_TAG=$$(printf "v%d.%d.%d" $$MAJOR $$MINOR $$PATCH); \
	echo "🔢 新 Tag: $$NEXT_TAG"; \
	git tag $$NEXT_TAG; \
	git push origin $$NEXT_TAG; \
	echo "✅ 已推送新 Tag: $$NEXT_TAG"