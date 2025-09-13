package postgres

// 兼容性别名，用于保持向后兼容
// 这允许旧的代码继续工作，同时我们逐步迁移到新的结构

// GetClient 兼容旧的 pgdb.GetClient()
var GetPgClient = GetClient

// Init 兼容旧的 pgdb.Init()
var InitPg = Init