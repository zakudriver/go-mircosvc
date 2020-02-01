
#!/usr/bin/env sh

# Install proto3 from source macOS only.
#  brew install autoconf automake libtool
#  git clone https://github.com/google/protobuf
#  ./autogen.sh ; ./configure ; make ; make install
#
# Update protoc Go bindings via
#  go get -usersvc github.com/golang/protobuf/{proto,protoc-gen-go}
#
# See also
#  https://github.com/grpc/grpc-go/tree/master/examples

protoc ./user/user.proto --go_out=plugins=grpc:.
protoc ./article/article.proto --go_out=plugins=grpc:.
