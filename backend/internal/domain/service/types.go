package service

// ContractInteractionData represents the data needed for frontend to interact with contracts
type ContractInteractionData struct {
	ContractAddress string `json:"contractAddress"`
	MethodName      string `json:"methodName"`
	Parameters      []any  `json:"parameters"`
	Value           string `json:"value,omitempty"` // For payable functions
	GasLimit        string `json:"gasLimit,omitempty"`
	GasPrice        string `json:"gasPrice,omitempty"`
}
