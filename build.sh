protoc --proto_path ../../../ -I=./proto --go_out=plugins=grpc:./proto proto/frametracker.proto
mv proto/github.com/brotherlogic/frametracker/proto/* ./proto
