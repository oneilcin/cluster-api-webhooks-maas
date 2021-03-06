package maas

import (
	"context"
	"fmt"

	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/juju/gomaasapi"
)

const (
	MAASMachineTag = "cluster.k8s.io/maas=" // Indicates a machine MAY be allocated if it has not been already
)

type Client struct {
	Controller gomaasapi.Controller
}

func NewClient(apiURL, apiVersion, apiKey string) (Client, error) {
	controller, err := gomaasapi.NewController(gomaasapi.ControllerArgs{
		BaseURL: apiURL,
		APIKey:  apiKey})
	if err != nil {
		return Client{}, fmt.Errorf("error creating controller with version: %v", err)
	}

	return Client{Controller: controller}, nil
}

// Create creates a machine
func (c Client) Create(ctx context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	allocateArgs := gomaasapi.AllocateMachineArgs{Tags: []string{MAASMachineTag + cluster.Namespace}}
	m, _, err := c.Controller.AllocateMachine(allocateArgs)
	if err != nil {
		return fmt.Errorf("error allocating machine: %v", err)
	}

	startArgs := gomaasapi.StartArgs{
		UserData:     "", // TODO(dwat): Use the Machine resource to set these?
		DistroSeries: "",
		Kernel:       "",
		Comment:      "",
	}
	err = m.Start(startArgs)
	if err != nil {
		return fmt.Errorf("error deploying machine (%v): %v", m.Hostname(), err)
	}

	return nil
}

// Delete deletes a machine
func (c Client) Delete(ctx context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	return nil
}

// Update updates a machine
func (c Client) Update(ctx context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	return nil
}

// Exists test for the existence of a machine
func (c Client) Exist(ctx context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) (bool, error) {
	// Get list of machines with tag
	machineArgs := gomaasapi.MachinesArgs{Hostnames: []string{MAASMachineTag + cluster.Namespace}}
	machines, err := c.Controller.Machines(machineArgs)
	if err != nil {
		return false, fmt.Errorf("error listing machines: %v", err)
	}
	if len(machines) != 1 {
		return false, fmt.Errorf("expected 1 machine (%s), found %d", machine.Name, len(machines))
	}

	return true, nil
}
