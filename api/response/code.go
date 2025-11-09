package response

// https://google-cloud.gitbook.io/api-design-guide/errors

type responseData struct {
	Code      int         `json:"code"`
	Status    string      `json:"status"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	TraceID   string      `json:"trace_id,omitempty"`
	Timestamp int64       `json:"timestamp"`
	Total     *int        `json:"total,omitempty"`
}

// 通用的错误信息

// 没有错误。
var Success = responseData{
	Code:    200,
	Status:  "OK",
	Message: "请求成功",
}

// 客户端发送的数据包含非法参数。查看错误消息和错误详情来获取更多的信息。
var INVALID_ARGUMENT = responseData{
	Code:    400,
	Status:  "INVALID_ARGUMENT",
	Message: "请求参数错误",
}

// 现在的系统状态不可以执行当前的请求，例如删除一个非空的目录。
var FAILED_PRECONDITION = responseData{
	Code:    400,
	Status:  "FAILED_PRECONDITION",
	Message: "无法执行客户端请求",
}

// 客户端指定了一个非法的范围。
var OUT_OF_RANGE = responseData{
	Code:    400,
	Status:  "OUT_OF_RANGE",
	Message: "客户端越限访问",
}

// 因为缺失的，失效的或者过期的OAuth令牌，请求未能通过身份认证。
var UNAUTHENTICATED = responseData{
	Code:    401,
	Status:  "UNAUTHENTICATED",
	Message: "身份验证失败",
}

// 客户端没有足够的权限。这可能是因为OAuth令牌没有正确的作用域，或者客户端没有权限，或者是API对客户端代码禁用了。
var PERMISSION_DENIED = responseData{
	Code:    403,
	Status:  "PERMISSION_DENIED",
	Message: "客户端权限不足",
}

// 特定的资源没有被找到或者请求因为某些未被公开的原因拒绝（例如白名单）。
var NOT_FOUND = responseData{
	Code:    404,
	Status:  "NOT_FOUND",
	Message: "资源不存在",
}

// 并发冲突，如读 - 修改 - 写冲突。
var ABORTED = responseData{
	Code:    409,
	Status:  "ABORTED",
	Message: "数据处理冲突",
}

// 客户端尝试新建的资源已经存在了。
var ALREADY_EXISTS = responseData{
	Code:    409,
	Status:  "ALREADY_EXISTS",
	Message: "资源已存在",
}

// 资源配额不足或达不到速率限制。
var RESOURCE_EXHAUSTED = responseData{
	Code:    429,
	Status:  "RESOURCE_EXHAUSTED",
	Message: "资源配额不足或达不到速率限制",
}

// 请求被客户端取消了。
var CANCELLED = responseData{
	Code:    499,
	Status:  "CANCELLED",
	Message: "请求被客户端取消了",
}

// 不可恢复的数据丢失或数据损坏。 客户端应该向用户报告错误。
var DATA_LOSS = responseData{
	Code:    500,
	Status:  "DATA_LOSS",
	Message: "处理数据发生错误",
}

// 未知的服务端出错，通常是由于服务器出现bug了。
var UNKNOWN = responseData{
	Code:    500,
	Status:  "UNKNOWN",
	Message: "服务器未知错误",
}

// 服务器内部错误。通常是由于服务器出现bug了。
var INTERNAL = responseData{
	Code:    500,
	Status:  "INTERNAL",
	Message: "服务器内部错误",
}

// API方法没有被服务器实现。
var NOT_IMPLEMENTED = responseData{
	Code:    501,
	Status:  "NOT_IMPLEMENTED",
	Message: "API不存在",
}

// 服务不可用。通常是由于服务器宕机了。
var UNAVAILABLE = responseData{
	Code:    503,
	Status:  "UNAVAILABLE",
	Message: "服务不可用",
}

// 请求超过了截止日期。 只有当调用者设置的截止日期比方法的默认截止日期更短（服务器没能够在截止日期之前处理完请求）并且请求没有在截止日期内完成时，才会发生这种情况。
var DEALINE_EXCEED = responseData{
	Code:    504,
	Status:  "DEALINE_EXCEED",
	Message: "请求超时",
}

// 根据业务自定义的常用状态码
