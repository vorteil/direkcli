package workflow

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/vorteil/direktiv/pkg/ingress"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// List returns an array of workflows for a given namespace
func List(conn *grpc.ClientConn, namespace string) ([]*ingress.GetWorkflowsResponse_Workflow, error) {
	client := ingress.NewDirektivIngressClient(conn)

	// set context with 3 second timeout
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	// prepare request
	request := ingress.GetWorkflowsRequest{
		Namespace: &namespace,
	}

	// send grpc request
	resp, err := client.GetWorkflows(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return nil, fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return resp.Workflows, nil
}

// // Execute runs the yaml provided from the workflow
func Execute(conn *grpc.ClientConn, namespace string, id string, input string) (string, error) {
	client := ingress.NewDirektivIngressClient(conn)

	// set context with 3 second timeout
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	var err error
	var b []byte
	if input != "" {
		b, err = ioutil.ReadFile(input)
		if err != nil {
			return "", err
		}
	}

	// prepare request
	request := ingress.InvokeWorkflowRequest{
		Namespace:  &namespace,
		Input:      b,
		WorkflowId: &id,
	}

	// send grpc request
	resp, err := client.InvokeWorkflow(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return "", fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return fmt.Sprintf("Successfully invoked, Instance ID: %s", resp.GetInstanceId()), nil
}

// getWorkflowUID returns uid of workflow so we can update/delete things related to it
func getWorkflowUID(conn *grpc.ClientConn, namespace, id string) (string, error) {
	client := ingress.NewDirektivIngressClient(conn)

	// set context with 3 second timeout
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	// prepare request
	request := ingress.GetWorkflowByIdRequest{
		Namespace: &namespace,
		Id:        &id,
	}

	// send grpc request
	resp, err := client.GetWorkflowById(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return "", fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}
	return resp.GetUid(), nil
}

// Get returns the YAML contents of the workflow
func Get(conn *grpc.ClientConn, namespace string, id string) (string, error) {
	client := ingress.NewDirektivIngressClient(conn)
	defer conn.Close()

	// set context with 3 second timeout
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	// prepare request
	request := ingress.GetWorkflowByIdRequest{
		Namespace: &namespace,
		Id:        &id,
	}

	// send grpc request
	resp, err := client.GetWorkflowById(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return "", fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return string(resp.GetWorkflow()), nil
}

// Update updates a workflow from the provided id
func Update(conn *grpc.ClientConn, namespace string, id string, filepath string) (string, error) {
	client := ingress.NewDirektivIngressClient(conn)
	defer conn.Close()

	// set context with 3 second timeout
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	uid, err := getWorkflowUID(conn, namespace, id)
	if err != nil {
		return "", err
	}

	// prepare request
	request := ingress.UpdateWorkflowRequest{
		Uid:      &uid,
		Workflow: b,
	}

	// send grpc request
	resp, err := client.UpdateWorkflow(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return "", fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return fmt.Sprintf("Successfully updated '%s'", resp.GetId()), nil
}

// Delete removes a workflow
func Delete(conn *grpc.ClientConn, namespace, id string) (string, error) {
	client := ingress.NewDirektivIngressClient(conn)
	defer conn.Close()

	// set context with 3 second timeout
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	uid, err := getWorkflowUid(conn, namespace, id)
	if err != nil {
		return "", err
	}

	// prepare request
	request := ingress.DeleteWorkflowRequest{
		Uid: &uid,
	}

	// send grpc request
	_, err = client.DeleteWorkflow(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return "", fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return fmt.Sprintf("Deleted workflow '%v'", id), nil
}

// Add creates a new workflow on a namespace
func Add(conn *grpc.ClientConn, namespace string, filepath string) (string, error) {
	client := ingress.NewDirektivIngressClient(conn)
	defer conn.Close()

	// set context with 3 second timeout
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	// prepare request
	request := ingress.AddWorkflowRequest{
		Namespace: &namespace,
		Workflow:  b,
	}

	// send grpc request
	resp, err := client.AddWorkflow(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return "", fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return fmt.Sprintf("Created workflow '%s'", resp.GetId()), nil
}
