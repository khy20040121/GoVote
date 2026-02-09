package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// CreatePost创建帖子logic
func CreatePost(p *models.Post) error {
	// 1生成postID
	p.ID = snowflake.GenID()
	p.CreateTime = time.Now()

	// 2 保存到数据库
	if err := mysql.CreatePost(p); err != nil {
		zap.L().Error("mysql.CreatePost failed", zap.Error(err))
		return err
	}

	// 3 保存到redis
	if err := redis.CreatePost(p); err != nil {
		zap.L().Error("redis.CreatePost failed", zap.Error(err))
		return err
	}
	return nil
}

// GetPostByID 根据帖子的id来查询帖子的详细数据
func GetPostByID(id, userID int64) (data *models.ApiPostDetail, err error) {
	//查询帖子的基本信息
	post, err := mysql.GetPostByID(id)
	if err != nil {
		zap.L().Error("mysql.GetPostById failed", zap.Error(err))
		return nil, err
	}

	// 根据帖子的作者id查询作者的姓名
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID failed", zap.Error(err))
		return
	}

	// 根据社区id查询社区的详细信息
	communityDetail, err := mysql.GetCommunityDetail(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetail failed", zap.Error(err))
		return
	}

	// 从redis获取投票数
	voteData, err := redis.GetPostVoteList([]string{strconv.FormatInt(post.ID, 10)})
	if err != nil {
		zap.L().Error("redis.GetPostVoteList failed", zap.Error(err))
		// 不影响主流程，默认为0
		voteData = []int64{0}
	}

	// 获取当前用户的投票状态
	var voteStatus int32
	if userID > 0 {
		status, err := redis.GetPostVoteForUser(strconv.FormatInt(userID, 10), strconv.FormatInt(post.ID, 10))
		if err != nil && err != redis.Nil {
			zap.L().Error("redis.GetPostVoteForUser failed", zap.Error(err))
		} else {
			voteStatus = int32(status)
		}
	}

	// 将数据组合到模型中
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		VoteNum:         voteData[0],
		VoteStatus:      voteStatus,
		Post:            post,
		CommunityDetail: communityDetail,
	}
	return
}

// GetPostList 获取所有帖子的列表logic
func GetPostList(page int64, size int64) (data []*models.ApiPostDetail, err error) {
	var user *models.User
	var community *models.CommunityDetail
	var posts []*models.Post

	posts, err = mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList failed", zap.Error(err))
		return
	}

	data = make([]*models.ApiPostDetail, len(posts))
	for idx, post := range posts {

		// 根据作者id查找作者信息
		user, err = mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID failed", zap.Error(err))
			return
		}
		// 根据社区id查找社区信息
		community, err = mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetail failed", zap.Error(err))
			return
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data[idx] = postDetail
	}
	return
}

// GetPostList根据指定顺序获取帖子列表logic
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {

	// 去redis查询ids
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}

	// 根据ids去MYSQL中查询帖子的详细信息
	posts, err := mysql.GetPostListsByIDs(ids)
	if err != nil {
		return
	}

	// 根据ids查找每个帖子有多少赞成票
	voteData, err := redis.GetPostVoteList(ids)
	if err != nil {
		return
	}

	data = make([]*models.ApiPostDetail, len(posts))
	for idx, post := range posts {
		// 根据作者id查找作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID failed", zap.Error(err))
			return nil, err
		}

		// 根据社区id查找社区信息
		community, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetail failed", zap.Error(err))
			return nil, err
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data[idx] = postDetail
	}
	return
}

// GetCommunityList 按社区获取帖子的详情
func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 去redis查询ids
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}

	// 根据ids去MYSQL中查询帖子的详细信息
	posts, err := mysql.GetPostListsByIDs(ids)
	if err != nil {
		return
	}

	// 根据ids查找每个帖子有多少赞成票
	voteData, err := redis.GetPostVoteList(ids)
	if err != nil {
		return
	}

	data = make([]*models.ApiPostDetail, len(posts))
	for idx, post := range posts {
		// 根据作者id查找作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID failed", zap.Error(err))
			return nil, err
		}

		// 根据社区id查找社区信息
		community, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetail failed", zap.Error(err))
			return nil, err
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data[idx] = postDetail
	}
	return
}

// GetPostListNew 按社区按顺序查询所有帖子的详情
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 未按社区查询
	if p.CommunityID == 0 {
		data, err = GetPostList2(p)
	} else {
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		return nil, err
	}
	return data, err
}
