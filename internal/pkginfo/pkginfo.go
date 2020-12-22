package pkginfo

var (
	CommitHash  = "abcdefg"                                                         // 提交Hash
	CommitDate  = "YYYYMMDD"                                                        // 提交日期
	GoVersion   = "X.Y.Z"                                                           // Go版本号
	Name        = "jssh"                                                            // 命令名称
	Version     = "0.3"                                                             // 主版本号
	LongVersion = Version + "-" + CommitDate + "-" + CommitHash + "-go" + GoVersion // 完整版本号
)
