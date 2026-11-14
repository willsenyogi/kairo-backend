package response

type Response struct {
	Status         int    `json:"status"`
	Method         string `json:"http_method,omitempty"`
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	TotalDataIndex int    `json:"total_data_index"`
	Data           any    `json:"data,omitempty"`
}
