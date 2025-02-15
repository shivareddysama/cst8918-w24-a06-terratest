package test

import (
	"testing"
	"strings"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

var subscriptionID string = "fdcbdc66-325a-470d-957c-7977b0fb7718"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"labelPrefix": "sama0096",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`
	terraform.InitAndApply(t, terraformOptions)

	// Get output variables
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// 1. Confirm NIC exists and is connected to VM
	assert.True(t, azure.NetworkInterfaceExists(t, nicName, resourceGroupName, subscriptionID), "NIC should exist")

	vm := azure.GetVirtualMachine(t, vmName, resourceGroupName, subscriptionID)
	connectedNICs := vm.NetworkProfile.NetworkInterfaces

	isNicConnected := false
	for _, nic := range *connectedNICs {  // Dereference the pointer here
		if strings.Contains(*nic.ID, nicName) {
			isNicConnected = true
			break
		}
	}
	assert.True(t, isNicConnected, "NIC should be connected to the VM")

	// 2. Confirm the VM is running the correct Ubuntu version
	expectedOSVersion := "0001-com-ubuntu-server-jammy"  // Adjust to your expected version, e.g., "Ubuntu 20.04 LTS"
	actualOSVersion := *vm.StorageProfile.ImageReference.Offer
	t.Logf("Actual OS Version: %s", actualOSVersion)
	assert.True(t, strings.Contains(actualOSVersion, expectedOSVersion), "VM should be running the expected Ubuntu version")
}
