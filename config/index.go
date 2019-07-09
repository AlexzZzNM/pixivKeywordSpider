package config

// 代理 列表
var ProxyList []string = []string{
	"socks5://127.0.0.1:1080",
}

// 爬取时， 最大并发 Goroutine 数目
const MaxGoroutine  = 4

// p站用户名 密码
const UserName string = ""
const Password string = ""

// 关键字 要搜索的关键字
const KeyWord string = "贞德" // 贞德天下第一  // 默认关键字 程序执行时会询问关键字

// 收藏少于 本选项的 不被爬取
const MinCollection int = 0

//数据库 s  暂时不用
const DBType string = "mysql"
const DBName string = "pixiv_keyword"
const DBUsername string = "root"
const DBPassword string = "123456"
const DBCity string = "127.0.0.1"
const DBCharacterSet string = "utf8mb4"
const MaxIdleConn int = 10
const MaxOpenConn int = 100
//数据库 e

