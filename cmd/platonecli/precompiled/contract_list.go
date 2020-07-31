package precompile

import "github.com/PlatONEnetwork/PlatONE-Go/common/syscontracts"

var (
	UserManagementAddress        = syscontracts.UserManagementAddress.String()        // The PlatONE Precompiled contract addr for user management
	NodeManagementAddress        = syscontracts.NodeManagementAddress.String()        // The PlatONE Precompiled contract addr for node management
	CnsManagementAddress         = syscontracts.CnsManagementAddress.String()         // The PlatONE Precompiled contract addr for CNS
	ParameterManagementAddress   = syscontracts.ParameterManagementAddress.String()   // The PlatONE Precompiled contract addr for parameter management
	FirewallManagementAddress    = syscontracts.FirewallManagementAddress.String()    // The PlatONE Precompiled contract addr for fire wall management
	GroupManagementAddress       = syscontracts.GroupManagementAddress.String()       // The PlatONE Precompiled contract addr for group management
	ContractDataProcessorAddress = syscontracts.ContractDataProcessorAddress.String() // The PlatONE Precompiled contract addr for group management
	CnsInvokeAddress             = syscontracts.CnsInvokeAddress.String()             // The PlatONE Precompiled contract addr for group management
)

// link the precompiled contract addresses with abi file bytes
var List = map[string]string{
	UserManagementAddress:        "../../release/linux/conf/contracts/userManager.cpp.abi.json",
	NodeManagementAddress:        "../../release/linux/conf/contracts/nodeManager.cpp.abi.json",
	CnsManagementAddress:         "../../release/linux/conf/contracts/cnsManager.cpp.abi.json",
	ParameterManagementAddress:   "../../release/linux/conf/contracts/paramManager.cpp.abi.json",
	FirewallManagementAddress:    "../../release/linux/conf/contracts/fireWall.abi.json",
	GroupManagementAddress:       "../../release/linux/conf/contracts/groupManager.cpp.abi.json",
	ContractDataProcessorAddress: "",
	CnsInvokeAddress:             "../../release/linux/conf/contracts/cnsinvoke.abi.json",
}
