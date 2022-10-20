up-swagger-folders:
	mkdir ./internal/cloud/rest/generated

down-swagger-folders:
	rm -rf ./internal/cloud/rest/generated

swagger-folders: down-swagger-folders up-swagger-folders

swagger-cloud:
	swagger generate server -f ./api/swagger/cloud/swagger.yml --exclude-main -A clouder -t ./internal/cloud/rest/generated -s ops

swagger-generated:
	mkdir -p static/swagger-ui/generated

copy-swaggers:
	cp ./api/swagger/cloud/swagger.yml ./static/swagger-ui/generated/cloud.yml

swagger: swagger-folders swagger-generated copy-swaggers swagger-cloud