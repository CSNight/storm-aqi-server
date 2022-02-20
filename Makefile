VERSION := v1.0.0
BUILD := $(shell git rev-parse --short HEAD)

image:
	docker build -t csnight/storm-aqi-server:$(VERSION)-$(BUILD) .

push: image
	docker push csnight/storm-aqi-server:$(VERSION)-$(BUILD)
