go test -coverpkg=./... -coverprofile=cover.out 
go tool cover -html=cover.out -o cover.html
rm cover.out