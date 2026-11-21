TAILWIND_VERSION := v4.1.17
DAISYUI_VERSION  := v5.5.5

OS   := $(shell uname -s | tr '[:upper:]' '[:lower:]' | sed 's/darwin/macos/')
ARCH := $(shell uname -m | sed -e 's/x86_64/x64/' -e 's/aarch64/arm64/')

TAILWIND := bin/tailwindcss-$(TAILWIND_VERSION)
DAISYUI  := web/assets/vendor/daisyui-$(DAISYUI_VERSION).stamp

.PHONY: run build test css css-watch tools clean

run:
	go run ./cmd/xmas-poll

build: css
	CGO_ENABLED=0 go build -trimpath -o bin/xmas-poll ./cmd/xmas-poll

test:
	go test ./...

# Rebuild web/static/app.css from web/assets/app.css and the classes used in web/templates/.
css: tools
	$(TAILWIND) -i web/assets/app.css -o web/static/app.css --minify

css-watch: tools
	$(TAILWIND) -i web/assets/app.css -o web/static/app.css --watch

tools: $(TAILWIND) $(DAISYUI)

$(TAILWIND):
	mkdir -p bin
	curl -sSfL -o $@ https://github.com/tailwindlabs/tailwindcss/releases/download/$(TAILWIND_VERSION)/tailwindcss-$(OS)-$(ARCH)
	chmod +x $@

$(DAISYUI):
	rm -rf web/assets/vendor && mkdir -p web/assets/vendor
	curl -sSfL -o web/assets/vendor/daisyui.mjs https://github.com/saadeghi/daisyui/releases/download/$(DAISYUI_VERSION)/daisyui.mjs
	curl -sSfL -o web/assets/vendor/daisyui-theme.mjs https://github.com/saadeghi/daisyui/releases/download/$(DAISYUI_VERSION)/daisyui-theme.mjs
	touch $@

clean:
	rm -rf bin web/assets/vendor
