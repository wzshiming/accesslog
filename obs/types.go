package obs

import (
	"fmt"
	"net/url"

	"github.com/wzshiming/accesslog"
	"github.com/wzshiming/accesslog/unsafeutils"
)

// https://support.huaweicloud.com/intl/zh-cn/usermanual-obs/obs_03_0329.html
// 787f2f92b20943998a4fe2ab75eb09b8 bucket [13/Aug/2015:01:43:42 +0000] xx.xx.xx.xx 787f2f92b20943998a4fe2ab75eb09b8 281599BACAD9376ECE141B842B94535B  REST.GET.BUCKET.LOCATION - "GET /bucket?location HTTP/1.1" 200 - 211 - 6 6 "-"  "HttpClient" - -

type AccessLog struct {
	// 桶的ownerId。Bucket owner's ID.
	BucketOwner string
	// 桶名。Bucket name.
	Bucket string
	// 请求时间戳（UTC）。Request timestamp (UTC).
	Time string
	// 请求IP。Remote IP.
	RemoteIP string
	// 请求者ID。Requester ID.
	Requester string
	// 请求ID。Request ID.
	RequestID string
	// 操作名称。Operation name.
	Operation string
	// 对象名。Object name.
	Key string
	// 请求URI。Request URI.
	RequestURI string
	// 返回码。HTTP status code.
	HTTPStatus string
	// 错误码。Error code.
	ErrorCode string
	// HTTP响应的字节大小。Bytes sent in HTTP response.
	BytesSent string
	// 对象大小（bytes）。Object size in bytes.
	ObjectSize string
	// 服务端处理时间（ms）。Server processing time in ms.
	TotalTime string
	// 总请求时间（ms）。Total request time in ms.
	TurnAroundTime string
	// 请求的referrer头域。HTTP referer of the request.
	Referer string
	// 请求的user-agent头域。HTTP User-Agent header.
	UserAgent string
	// 请求中带的 versionId
	VersionID string
	// 联邦认证及委托授权信息
	STSLogUrn string
	// 当前的对象存储类别。Current object storage class.
	StorageClass string
	// 通过转换后的对象存储类别。Target object storage class.
	TargetStorageClass string
	// 对于并行文件系统，是文件/目录的内部标识，由父目录inode编号与文件/目录名称组成。对于对象桶，该字段为"-"。
	DentryName string
	// IAM用户ID。当使用匿名用户发起请求，记为Anonymous。
	IAMUserID string

	Unknown1 string
}

func (a AccessLog) Formatted() (f AccessLogFormatted, err error) {
	f.RemoteIP = accesslog.CleanupString(a.RemoteIP)
	f.Time = accesslog.CleanupString(a.Time)

	req, err := accesslog.ParseRequest(a.RequestURI)
	if err != nil {
		return f, fmt.Errorf("failed to formatted request URL: %w", err)
	}
	f.RequestMethod = req.Method
	f.RequestScheme = req.URL.Scheme
	f.RequestHost = req.URL.Host
	f.RequestPath = req.URL.Path
	f.RequestQuery = req.URL.RawQuery
	f.RequestProto = req.Proto

	f.Status = accesslog.CleanupString(a.HTTPStatus)
	f.BodySentBytes = accesslog.CleanupString(a.BytesSent)
	f.RequestTime = accesslog.CleanupString(a.TurnAroundTime)
	f.Referer = accesslog.CleanupString(a.Referer)
	f.UserAgent = accesslog.CleanupString(a.UserAgent)
	f.RequestID = accesslog.CleanupString(a.RequestID)
	f.BucketOwner = accesslog.CleanupString(a.BucketOwner)
	f.Operation = accesslog.CleanupString(a.Operation)
	f.BucketName = accesslog.CleanupString(a.Bucket)
	f.ObjectName = accesslog.CleanupString(a.Key)
	on, err := url.PathUnescape(f.ObjectName)
	if err == nil {
		f.ObjectName = on
	}
	f.ObjectSize = accesslog.CleanupString(a.ObjectSize)
	f.ServerCostTime = accesslog.CleanupString(a.TotalTime)
	f.ErrorCode = accesslog.CleanupString(a.ErrorCode)
	f.Requester = accesslog.CleanupString(a.Requester)
	f.VersionID = accesslog.CleanupString(a.VersionID)
	f.STSLogUrn = accesslog.CleanupString(a.STSLogUrn)
	f.StorageClass = accesslog.CleanupString(a.StorageClass)
	f.TargetStorageClass = accesslog.CleanupString(a.TargetStorageClass)
	f.DentryName = accesslog.CleanupString(a.DentryName)
	f.IAMUserID = accesslog.CleanupString(a.IAMUserID)
	return f, nil
}

type AccessLogFormatted struct {
	RemoteIP           string
	Time               string
	RequestMethod      string
	RequestScheme      string
	RequestHost        string
	RequestPath        string
	RequestQuery       string
	RequestProto       string
	Status             string
	BodySentBytes      string
	RequestTime        string
	Referer            string
	UserAgent          string
	RequestID          string
	LoggingFlag        string
	BucketOwner        string
	Operation          string
	BucketName         string
	ObjectName         string
	ObjectSize         string
	ServerCostTime     string
	ErrorCode          string
	Requester          string
	VersionID          string
	STSLogUrn          string
	StorageClass       string
	TargetStorageClass string
	DentryName         string
	IAMUserID          string
}

func (AccessLogFormatted) Fields() []string {
	return unsafeutils.Fields[AccessLogFormatted]()
}

func (e AccessLogFormatted) Values(fields []string) []string {
	accessLogEntryFieldsIndexMapping := unsafeutils.FieldsOffset[AccessLogFormatted]()
	out := make([]string, len(fields))
	for i, f := range fields {
		offset, ok := accessLogEntryFieldsIndexMapping[f]
		if !ok {
			continue
		}
		out[i] = unsafeutils.GetWithOffset[string](&e, offset)
	}
	return out
}
