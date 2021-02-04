package namespace

import (
	"context"
	"fmt"
	"time"

	"github.com/vorteil/direktiv/pkg/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// List returns a list of namespaces on direktiv
func List(conn *grpc.ClientConn) ([]*protocol.GetNamespacesResponse_Namespace, error) {
	client := protocol.NewDirektivClient(conn)
	defer conn.Close()

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	request := protocol.GetNamespacesRequest{}

	resp, err := client.GetNamespaces(ctx, &request)
	if err != nil {
		// convert the error
		s := status.Convert(err)
		return nil, fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return resp.Namespaces, nil
}

// Delete deletes a namespace on direktiv
func Delete(name string, conn *grpc.ClientConn) (string, error) {
	client := protocol.NewDirektivClient(conn)
	defer conn.Close()

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	request := protocol.DeleteNamespaceRequest{
		Name: &name,
	}

	resp, err := client.DeleteNamespace(ctx, &request)
	if err != nil {
		// convert the error
		s := status.Convert(err)
		return "", fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return fmt.Sprintf("Deleted namespace: %s", resp.GetName()), nil
}

// Create creates a new namespace on direktiv
func Create(name string, conn *grpc.ClientConn) (string, error) {
	client := protocol.NewDirektivClient(conn)
	defer conn.Close()

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	request := protocol.AddNamespaceRequest{
		Name: &name,
	}

	resp, err := client.AddNamespace(ctx, &request)
	if err != nil {
		// convert the error
		s := status.Convert(err)
		return "", fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return fmt.Sprintf("Created namespace: %s", resp.GetName()), nil
}
