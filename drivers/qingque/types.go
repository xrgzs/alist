package qingque

import "github.com/alist-org/alist/v3/internal/model"

type Object struct {
	model.Object
	ShortcutID string
}

func (d *Object) GetShortcutID() string {
	return d.ShortcutID
}

type BaseResp struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Result     any    `json:"result"`
	Ok         bool   `json:"ok"`
	ResultCode string `json:"resultCode"`
}

//	type Entity struct {
//		EntityType      int    `json:"entityType"`
//		EntityID        string `json:"entityId"`
//		DomainID        string `json:"domainId"`
//		BusinessType    int    `json:"businessType"`
//		DisplayName     string `json:"displayName"`
//		CnName          string `json:"cnName"`
//		EnName          string `json:"enName"`
//		NationalityCode string `json:"nationalityCode"`
//		IsCurrentUser   bool   `json:"isCurrentUser"`
//		LoginID         string `json:"loginId"`
//		KwaiUserID      string `json:"kwaiUserId"`
//		KuimUserID      string `json:"kuimUserId"`
//		Avatar          string `json:"avatar"`
//		Detail          string `json:"detail"`
//		Resign          bool   `json:"resign"`
//	}
//
//	type Owner struct {
//		Entity Entity `json:"entity"`
//		Role   int    `json:"role"`
//	}
//
//	type Creator struct {
//		Entity Entity `json:"entity"`
//	}
//
//	type Ref struct {
//		ID              string      `json:"id"`
//		ShortcutID      string      `json:"shortcutId"`
//		IsOrigin        bool        `json:"isOrigin"`
//		InKnowledgeBase bool        `json:"inKnowledgeBase"`
//		KbaseType       interface{} `json:"kbaseType"`
//	}
type CosmoPaths struct {
	DocID     string `json:"docId"`
	DocName   string `json:"docName"`
	CosmoURL  string `json:"cosmoUrl"`
	DocTypeEn string `json:"docTypeEn"` //folder, yFile, doc, sheet, metaSheet
	// UserPhotoID string `json:"userPhotoId"`
	// SpaceName                    string `json:"spaceName"`
	// ShareType                    int    `json:"shareType"`
	// HasSubCosmo                  bool `json:"hasSubCosmo"`
	// IsOrigin                     bool `json:"isOrigin"`
	// ContainsExternalCollaborator bool `json:"containsExternalCollaborator"`
}
type FileList struct {
	CosmoPaths
	// DomainID                     string       `json:"domainId"`
	// LastViewTime                 int64    `json:"lastViewTime"`
	LastModifiedTime int64 `json:"lastModifiedTime"`
	// LastModifiedUserName         string       `json:"lastModifiedUserName"`
	// LastModifiedUserID           string       `json:"lastModifiedUserId"`
	CreateTime int64 `json:"createTime"`
	// Owner                        Owner        `json:"owner"`
	// Creator                      Creator      `json:"creator"`
	// AccessLevel                  int          `json:"accessLevel"`
	// IsCollect                    bool         `json:"isCollect"`
	// IsQuick                      bool         `json:"isQuick"`
	// Description                  string       `json:"description"`
	// IsCooperation                bool         `json:"isCooperation"`
	// IsRename                     bool         `json:"isRename"`
	// IsSafe                       bool         `json:"isSafe"`
	// IsHide                       bool         `json:"isHide"`
	// IsTag                        bool         `json:"isTag"`
	// LastViewTimeByMe             int64    `json:"lastViewTimeByMe"`
	// LastModifiedTimeByMe         int64    `json:"lastModifiedTimeByMe"`
	// IsTopSpace                   bool         `json:"isTopSpace"`
	// IsTopKnowledgeBase           bool         `json:"isTopKnowledgeBase"`
	// Unsolved                     int          `json:"unsolved"`
	// IsOriginDeleted              bool         `json:"isOriginDeleted"`
	// HasSubFoldersWithoutShortcut bool         `json:"hasSubFoldersWithoutShortcut"`
	// Role                         int          `json:"role"`
	// ClassifyLevel                int          `json:"classifyLevel"`
	ShortcutID string `json:"shortcutId"`
	// InKnowledgeBase              bool         `json:"inKnowledgeBase"`
	// Movable                      bool         `json:"movable"`
	// CanDelete                    bool         `json:"canDelete"`
	// CanAddSubNodes               bool         `json:"canAddSubNodes"`
	// Ref                          Ref          `json:"ref"`
	// CanCopy                      bool         `json:"canCopy"`
	// CanDuplicate                 bool         `json:"canDuplicate"`
	// PersonalKBase                bool         `json:"personalKBase"`
	FileSize int64 `json:"fileSize"`
	// InCollaboratorList           bool         `json:"inCollaboratorList"`
	// Collaborator                 int          `json:"collaborator"`
	// SubCosmo                     int          `json:"subCosmo"`
	// ViewCount                    int          `json:"viewCount"`
	// LikeCount                    int          `json:"likeCount"`
	// Purchased                    bool         `json:"purchased"`
	// IsRecommendation             bool         `json:"isRecommendation"`
	// CreateBrother                bool         `json:"createBrother"`
	// CosmoPathFirst               CosmoPaths   `json:"cosmoPathFirst"`
	// CosmoPaths                   []CosmoPaths `json:"cosmoPaths"`
	// OriginPath                   int          `json:"originPath"`
	// YfileType                    string       `json:"yfileType,omitempty"`
	// YfileWidth                   int          `json:"yfileWidth,omitempty"`
	// YfileHeight                  int          `json:"yfileHeight,omitempty"`
	// ThumbnailURL string `json:"thumbnailUrl,omitempty"` // only for online doc
}
type CosmoExtVoPage struct {
	// Total            int           `json:"total"`
	// Orders           []interface{} `json:"orders"`
	// OptimizeCountSQL bool        `json:"optimizeCountSql"`
	// HitCount         bool        `json:"hitCount"`
	// CountID          interface{} `json:"countId"`
	// MaxLimit         interface{} `json:"maxLimit"`
	List []FileList `json:"list"`
	// PageNum  int        `json:"pageNum"`
	// PageSize int        `json:"pageSize"`
	HasNext bool `json:"hasNext"`
	// SessionID string     `json:"sessionId"`
	// Pages     int        `json:"pages"`
}
type FileResp struct {
	CosmoExtVoPage CosmoExtVoPage `json:"cosmoExtVoPage"`
}
type DownloadResp struct {
	// CosmoID     string `json:"cosmoId"`
	// FileName    string `json:"fileName"`
	// FileSize    int    `json:"fileSize"`
	// FileType    string `json:"fileType"`
	FileURL string `json:"fileUrl"`
	// IsWrite     bool   `json:"isWrite"`
	// EncryptFlag bool   `json:"encryptFlag"`
}

type FolderNewResp struct {
	DocID string `json:"docId"`
	// ShortcutID string `json:"shortcutId"`
	DocTypeEn string `json:"docTypeEn"`
	// OpenDocURL string `json:"openDocUrl"`
	Cosmo any `json:"cosmo"`
}
