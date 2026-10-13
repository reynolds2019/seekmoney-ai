package errors

/**
 * desc: 全局错误码定义与错误封装
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */

const (
	CodeSuccess          = 200  // 操作成功
	CodeAccepted         = 202  // 异步接收
	CodeParamInvalid     = 4001 // 参数校验不通过
	CodeUnauthorized     = 4003 // 未授权或凭证失效
	CodeJobNotFound      = 4004 // 任务不存在
	CodeRateLimit        = 4029 // 请求频次受限
	CodeInternalError    = 5000 // 系统内部异常
	CodeCrawlFailed      = 5001 // 社交数据采集失败
	CodeLLMAnalyzeFailed = 5002 // 大模型分析异常
	CodeStorageError     = 5003 // 存储或数据库异常
)

// 获取错误码的中文描述
func GetErrorMessage(code int) string {
	switch code {
	case CodeSuccess:
		return "操作成功"
	case CodeAccepted:
		return "任务已接收并提交后台处理"
	case CodeParamInvalid:
		return "请求参数校验失败，请检查输入"
	case CodeUnauthorized:
		return "身份验证失败，请重新登录"
	case CodeJobNotFound:
		return "未找到指定的商机挖掘任务"
	case CodeRateLimit:
		return "操作过于频繁，请稍后再试"
	case CodeInternalError:
		return "服务器内部异常，请联系管理员"
	case CodeCrawlFailed:
		return "社交媒体数据抓取失败，请检查网络或更换关键词"
	case CodeLLMAnalyzeFailed:
		return "大模型深度商业分析服务调用失败"
	case CodeStorageError:
		return "数据持久化或对象存储操作异常"
	default:
		return "未知系统错误"
	}
}
