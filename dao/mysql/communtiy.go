package mysql

import (
	"bluebell/models"
	"database/sql"

	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := `select community_id,community_name from community`

	if err = db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community", zap.Error(err))
			err = nil
		}
	}
	return
}

func GetCommunityDetail(id int64) (communityDetail *models.CommunityDetail, err error) {
	sqlStr := `select community_id,community_name,introduction, create_time
				from community
				where community_id = ?`

	communityDetail = new(models.CommunityDetail) // 需要分配内存
	if err = db.Get(communityDetail, sqlStr, int64(id)); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return communityDetail, err
}
