up-swagger-folders:
	mkdir ./internal/cloud/rest/generated

down-swagger-folders:
	rm -rf ./internal/cloud/rest/generated

swagger-folders: down-swagger-folders up-swagger-folders
	echo "swagger folders made"

swagger-ui-folders:
	rm -rf static/swagger-ui/generated
	mkdir -p static/swagger-ui/generated

copy-swaggers:
	cp ./api/swagger/cloud/swagger.yml ./static/swagger-ui/generated/cloud.yml

swagger-ui: swagger-ui-folders copy-swaggers
	echo "swagger-ui generated"

up-mock-folders:
	mkdir internal/cloud/mocks

down-mocks-folders:
	rm -rf internal/cloud/mocks

mock-folders: down-mocks-folders up-mock-folders
	echo "mock folders made"

go-generate:
	go generate ./...

generate: swagger-folders mock-folders go-generate swagger-ui
	echo "project generated"