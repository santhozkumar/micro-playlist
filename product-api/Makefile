check_swagger:
	which swagger || go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_swagger
	@echo Ensure you have the swagger cli installed
	@echo ...l	

	GO111MODULE=off	swagger generate spec -o ./swagger.yaml --scan-models



