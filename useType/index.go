package useType

// 用于解析 post
type PostKeyType struct {
	PostKey string `json:"pixivAccount.postKey"`
}

// page 页面出现的基础 作品信息
type BaseDataItem struct {
	IllustId 				string `json:"illustId"` // 作品 id
	IllustTitle				string `json:"illustTitle"` // 作品标题
	Url						string `json:"url"` // 缩略图地址， 直接访问会
	UserId					string `json:"userId"` // 作者id
	UserName				string `json:"userName"` // 作者昵称
	UserImage				string `json:"userImage"` // 作者头像
	Width					int `json:"width"` // 宽
	Height					int `json:"height"` // 高
	BookmarkCount			int `json:"bookmarkCount"` // 收藏数
	ImageBase64				string `json:"imageBase64"` // 缩略图的base64 位
}

type DetailsInfo struct {
	IllustId		string `json:"illustId"` // 作品id
	Urls			map[string]string `json:"urls"` // 不同 清晰度的 图片 mini  thumb  small  regular  original
	UserId			string `json:"userId"`
}

type Combination struct {
	BaseDataItem
	Urls			map[string]string `json:"urls"`
}

type CombinationSlice []Combination

func (c CombinationSlice) Len () int {
	return len(c)
}

func (c CombinationSlice) Swap(i, j int){     // 重写 Swap() 方法
	c[i], c[j] = c[j], c[i]
}
func (c CombinationSlice) Less(i, j int) bool {    // 重写 Less() 方法， 从大到小排序
	return c[j].BookmarkCount < c[i].BookmarkCount
}