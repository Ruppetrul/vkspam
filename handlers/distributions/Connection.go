package distributions

import (
	"google.golang.org/grpc"
	"sync"
)

var (
	connection *grpc.ClientConn
	once       sync.Once
)

func NewConnection() (*grpc.ClientConn, error) {
	var err error
	once.Do(func() {
		connection, err = grpc.Dial("vkspam_parser:10001", grpc.WithInsecure())
	})

	if err != nil {
		return nil, err
	}

	return connection, nil
}
