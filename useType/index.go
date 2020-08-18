package useType

import "time"

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

// ajax 获取的数据 对象
type AjaxPageData struct {
	Error bool `json:"error"`
	Body struct {
		IllustManga struct {
			Data []struct {
				IllustID string `json:"illustId,omitempty"`
				IllustTitle string `json:"illustTitle,omitempty"`
				ID string `json:"id,omitempty"`
				Title string `json:"title,omitempty"`
				IllustType int `json:"illustType,omitempty"`
				XRestrict int `json:"xRestrict,omitempty"`
				Restrict int `json:"restrict,omitempty"`
				Sl int `json:"sl,omitempty"`
				URL string `json:"url,omitempty"`
				Description string `json:"description,omitempty"`
				Tags []string `json:"tags,omitempty"`
				UserID string `json:"userId,omitempty"`
				UserName string `json:"userName,omitempty"`
				Width int `json:"width,omitempty"`
				Height int `json:"height,omitempty"`
				PageCount int `json:"pageCount,omitempty"`
				IsBookmarkable bool `json:"isBookmarkable,omitempty"`
				BookmarkData interface{} `json:"bookmarkData,omitempty"`
				Alt string `json:"alt,omitempty"`
				IsAdContainer bool `json:"isAdContainer"`
				TitleCaptionTranslation struct {
					WorkTitle interface{} `json:"workTitle"`
					WorkCaption interface{} `json:"workCaption"`
				} `json:"titleCaptionTranslation,omitempty"`
				CreateDate time.Time `json:"createDate,omitempty"`
				UpdateDate time.Time `json:"updateDate,omitempty"`
				ProfileImageURL string `json:"profileImageUrl,omitempty"`
			} `json:"data"`
			Total int `json:"total"`
			BookmarkRanges []struct {
				Min interface{} `json:"min"`
				Max interface{} `json:"max"`
			} `json:"bookmarkRanges"`
		} `json:"illustManga"`
		Popular struct {
			Recent []struct {
				IllustID string `json:"illustId"`
				IllustTitle string `json:"illustTitle"`
				ID string `json:"id"`
				Title string `json:"title"`
				IllustType int `json:"illustType"`
				XRestrict int `json:"xRestrict"`
				Restrict int `json:"restrict"`
				Sl int `json:"sl"`
				URL string `json:"url"`
				Description string `json:"description"`
				Tags []string `json:"tags"`
				UserID string `json:"userId"`
				UserName string `json:"userName"`
				Width int `json:"width"`
				Height int `json:"height"`
				PageCount int `json:"pageCount"`
				IsBookmarkable bool `json:"isBookmarkable"`
				BookmarkData interface{} `json:"bookmarkData"`
				Alt string `json:"alt"`
				IsAdContainer bool `json:"isAdContainer"`
				TitleCaptionTranslation struct {
					WorkTitle interface{} `json:"workTitle"`
					WorkCaption interface{} `json:"workCaption"`
				} `json:"titleCaptionTranslation"`
				CreateDate time.Time `json:"createDate"`
				UpdateDate time.Time `json:"updateDate"`
				ProfileImageURL string `json:"profileImageUrl"`
			} `json:"recent"`
			Permanent []struct {
				IllustID string `json:"illustId"`
				IllustTitle string `json:"illustTitle"`
				ID string `json:"id"`
				Title string `json:"title"`
				IllustType int `json:"illustType"`
				XRestrict int `json:"xRestrict"`
				Restrict int `json:"restrict"`
				Sl int `json:"sl"`
				URL string `json:"url"`
				Description string `json:"description"`
				Tags []string `json:"tags"`
				UserID string `json:"userId"`
				UserName string `json:"userName"`
				Width int `json:"width"`
				Height int `json:"height"`
				PageCount int `json:"pageCount"`
				IsBookmarkable bool `json:"isBookmarkable"`
				BookmarkData interface{} `json:"bookmarkData"`
				Alt string `json:"alt"`
				IsAdContainer bool `json:"isAdContainer"`
				TitleCaptionTranslation struct {
					WorkTitle interface{} `json:"workTitle"`
					WorkCaption interface{} `json:"workCaption"`
				} `json:"titleCaptionTranslation"`
				CreateDate time.Time `json:"createDate"`
				UpdateDate time.Time `json:"updateDate"`
				ProfileImageURL string `json:"profileImageUrl"`
			} `json:"permanent"`
		} `json:"popular"`
		RelatedTags []string `json:"relatedTags"`
		ZoneConfig struct {
			Header struct {
				URL string `json:"url"`
			} `json:"header"`
			Footer struct {
				URL string `json:"url"`
			} `json:"footer"`
			Infeed struct {
				URL string `json:"url"`
			} `json:"infeed"`
		} `json:"zoneConfig"`
		ExtraData struct {
			Meta struct {
				Title string `json:"title"`
				Description string `json:"description"`
				Canonical string `json:"canonical"`
				AlternateLanguages struct {
					Ja string `json:"ja"`
					En string `json:"en"`
				} `json:"alternateLanguages"`
				DescriptionHeader string `json:"descriptionHeader"`
			} `json:"meta"`
		} `json:"extraData"`
	} `json:"body"`
}


type DetailsInfoID struct {
IllustID string `json:"illustId"`
IllustTitle string `json:"illustTitle"`
IllustComment string `json:"illustComment"`
ID string `json:"id"`
Title string `json:"title"`
Description string `json:"description"`
IllustType int `json:"illustType"`
CreateDate time.Time `json:"createDate"`
UploadDate time.Time `json:"uploadDate"`
Restrict int `json:"restrict"`
XRestrict int `json:"xRestrict"`
Sl int `json:"sl"`
Urls struct {
Mini string `json:"mini"`
Thumb string `json:"thumb"`
Small string `json:"small"`
Regular string `json:"regular"`
Original string `json:"original"`
} `json:"urls"`
Tags struct {
AuthorID string `json:"authorId"`
IsLocked bool `json:"isLocked"`
Tags []struct {
Tag string `json:"tag"`
Locked bool `json:"locked"`
Deletable bool `json:"deletable"`
UserID string `json:"userId,omitempty"`
UserName string `json:"userName,omitempty"`
Translation struct {
En string `json:"en"`
} `json:"translation,omitempty"`
} `json:"tags"`
Writable bool `json:"writable"`
} `json:"tags"`
Alt string `json:"alt"`
StorableTags []string `json:"storableTags"`
UserID string `json:"userId"`
UserName string `json:"userName"`
UserAccount string `json:"userAccount"`
LikeData bool `json:"likeData"`
Width int `json:"width"`
Height int `json:"height"`
PageCount int `json:"pageCount"`
BookmarkCount int `json:"bookmarkCount"`
LikeCount int `json:"likeCount"`
CommentCount int `json:"commentCount"`
ResponseCount int `json:"responseCount"`
ViewCount int `json:"viewCount"`
IsHowto bool `json:"isHowto"`
IsOriginal bool `json:"isOriginal"`
ImageResponseOutData []interface{} `json:"imageResponseOutData"`
ImageResponseData []interface{} `json:"imageResponseData"`
ImageResponseCount int `json:"imageResponseCount"`
PollData interface{} `json:"pollData"`
SeriesNavData interface{} `json:"seriesNavData"`
DescriptionBoothID interface{} `json:"descriptionBoothId"`
DescriptionYoutubeID interface{} `json:"descriptionYoutubeId"`
ComicPromotion interface{} `json:"comicPromotion"`
FanboxPromotion interface{} `json:"fanboxPromotion"`
ContestBanners []interface{} `json:"contestBanners"`
IsBookmarkable bool `json:"isBookmarkable"`
BookmarkData interface{} `json:"bookmarkData"`
ContestData interface{} `json:"contestData"`
ZoneConfig struct {
Responsive struct {
URL string `json:"url"`
} `json:"responsive"`
Rectangle struct {
URL string `json:"url"`
} `json:"rectangle"`
Five00X500 struct {
URL string `json:"url"`
} `json:"500x500"`
Header struct {
URL string `json:"url"`
} `json:"header"`
Footer struct {
URL string `json:"url"`
} `json:"footer"`
ExpandedFooter struct {
URL string `json:"url"`
} `json:"expandedFooter"`
Logo struct {
URL string `json:"url"`
} `json:"logo"`
} `json:"zoneConfig"`
ExtraData struct {
Meta struct {
Title string `json:"title"`
Description string `json:"description"`
Canonical string `json:"canonical"`
AlternateLanguages struct {
Ja string `json:"ja"`
En string `json:"en"`
} `json:"alternateLanguages"`
DescriptionHeader string `json:"descriptionHeader"`
Ogp struct {
Description string `json:"description"`
Image string `json:"image"`
Title string `json:"title"`
Type string `json:"type"`
} `json:"ogp"`
Twitter struct {
Description string `json:"description"`
Image string `json:"image"`
Title string `json:"title"`
Card string `json:"card"`
} `json:"twitter"`
} `json:"meta"`
} `json:"extraData"`
TitleCaptionTranslation struct {
WorkTitle interface{} `json:"workTitle"`
WorkCaption interface{} `json:"workCaption"`
} `json:"titleCaptionTranslation"`
IsUnlisted bool `json:"isUnlisted"`
}


type DetailsInfoTwo struct {
	Timestamp time.Time `json:"timestamp"`
	Illust map[string] DetailsInfoID `json:"illust"`
}