package apimodel

type LogInfo struct {
	LogName      string `json:"log_name"`      //日志文件名
	ErrorType    string `json:"error_type"`    // 状态
	Request      string `json:"request" `      // 请求信息
	Url          string `json:"url"`           //请求url
	ErrorCode    string `json:"error_code"`    // 状态码
	Msg          string `json:"msg"`           // 错误信息
	ErrorDetail  string `json:"error_detail"`  //详细错误数据
	OriginDetail string `json:"origin_detail"` //原始信息
	CreatedAt    string `json:"created_time"`
}

type LogInfoRequest struct {
	LogName   string `json:"name" form:"name"`
	StartTime string `json:"start_time" form:"start_time"`
	EndTime   string `json:"end_time" form:"end_time"`
	PaginationRequest
}

type LogListResponse struct {
	List []LogInfo `json:"list"`
	PaginationResponse
}

type LogInfoResponse struct {
	List []LogInfo `json:"list"`
	PaginationResponse
}

func (req LogInfoRequest) Valid(opt string) error {
	return nil
}

func (resp *LogListResponse) Load(total int, list []string) {
	resp.List = make([]LogInfo, 0, len(list))
	for _, v := range list {
		info := LogInfo{}
		info.LogName = v
		resp.List = append(resp.List, info)
	}
	resp.TotalSize = total
}

func (resp *LogInfoResponse) Load(total int, list []LogInfo) {
	resp.List = make([]LogInfo, 0, len(list))
	for _, v := range list {
		resp.List = append(resp.List, v)
	}
	resp.TotalSize = total
}
