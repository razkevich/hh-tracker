package helper

// Data is used to static shape json return
type Data struct {
	Data interface{} `json:"data"`
} // @name Response.Data

// FullResponse is used to static shape json return
type FullResponse struct {
	Data  interface{} `json:"data"`
	Links interface{} `json:"links"`
} // @name Response.FullResponse

// FullResponseWithMeta is used to static shape json return
type FullResponseWithMeta struct {
	Meta  interface{} `json:"meta"`
	Data  interface{} `json:"data"`
	Links interface{} `json:"links"`
} // @name Response.FullResponseWithMeta

// WrapData method is to inject data value to dynamic success response
func WrapData(data interface{}) Data {
	res := Data{
		Data: data,
	}
	return res
} // @name Response

// BuildFullResponse method is to inject data value to dynamic success response with link
func BuildFullResponse(data interface{}, link interface{}) FullResponse {
	res := FullResponse{
		Data:  data,
		Links: link,
	}
	return res
}

// BuildFullResponseWithMeta method is to inject data value to dynamic success response with links and meta data
func BuildFullResponseWithMeta(data interface{}, links interface{}, meta interface{}) FullResponseWithMeta {
	res := FullResponseWithMeta{
		Meta:  meta,
		Data:  data,
		Links: links,
	}
	return res
}
