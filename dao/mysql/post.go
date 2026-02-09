package mysql

import (
	"bluebell/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

func CreatePost(p *models.Post) error {
	sqlStr := `insert into post(post_id, title, content, author_id, community_id, create_time) values(?,?,?,?,?,?)`
	_, err := db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID, p.CreateTime)
	return err
}

// GetPostByID 根据帖子id到数据库里面查找帖子的详细信息
func GetPostByID(id int64) (data *models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time 
				from post
				where post_id = ?`
	data = new(models.Post)
	err = db.Get(data, sqlStr, id)
	return
}

// GetPostList 获取所有帖子列表mysql
func GetPostList(page int64, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time  from post
				limit ?,?`
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}

// GetPostListsByIDs 通过dis查询相应的帖子详情
func GetPostListsByIDs(ids []string) (posts []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time  
				from post
				where post_id in(?)
				order by FIND_IN_SET(post_id, ?)`

	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	err = db.Select(&posts, query, args...)
	return
}
