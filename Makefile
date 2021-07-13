# Add VERBOSE=1 flag when running make to echo all executed commands
ifndef VERBOSE
MAKEFLAGS += --silent
endif

# Not all make processors support the special .PHONY target, so mark FORCE as phony for those that do support it
# and hope this file does not exist for those that don't
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: FORCE
FORCE: ;

all: start-upstream start-downstream

build: build-upstream build-downstream

vendor: FORCE
	go mod vendor

build-upstream: vendor
	go build -o build/upstream ./main.go

build-downstream: vendor
	cd downstream
	go build -o build/downstream ./downstream/main.go

start-upstream: build-upstream
	@-pkill -f ./build/upstream
	./build/upstream  >> ./logs/upstream.log & disown

start-downstream: build-downstream
	@-pkill -f ./build/downstream
	./build/downstream >> ./logs/downstream.log & disown

