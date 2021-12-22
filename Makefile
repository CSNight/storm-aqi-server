VERSION := v1.0.0
BUILD := $(shell git rev-parse --short HEAD)

image:
	sudo docker build -t csnight/aqi-server:$(VERSION)-$(BUILD) .

push: image
	sudo docker push csnight/aqi-server:$(VERSION)-$(BUILD)
