import { useState, useEffect } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { postApi, communityApi } from '../utils/api';
import { formatTime } from '../utils/format';
import type { Post, Community } from '../types';

export default function PostList() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [posts, setPosts] = useState<Post[]>([]);
  const [communities, setCommunities] = useState<Community[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  
  // 从 URL 获取状态
  const order = (searchParams.get('order') as 'time' | 'score') || 'time';
  const communityIdParam = searchParams.get('community_id');
  const selectedCommunity = communityIdParam ? Number(communityIdParam) : undefined;

  const [hasMore, setHasMore] = useState(true);

  useEffect(() => {
    loadCommunities();
  }, []);

  useEffect(() => {
    loadPosts();
  }, [page, order, selectedCommunity]);

  const loadCommunities = async () => {
    try {
      const response = await communityApi.getList();
      if (response.data.code === 1000) {
        setCommunities(response.data.data);
      }
    } catch (error) {
      console.error('加载社区列表失败', error);
    }
  };

  const loadPosts = async () => {
    setLoading(true);
    try {
      const response = await postApi.getList({
        page,
        size: 10,
        order,
        community_id: selectedCommunity,
      });
      if (response.data.code === 1000) {
        const newPosts = response.data.data;
        if (page === 1) {
          setPosts(newPosts);
        } else {
          setPosts((prev) => [...prev, ...newPosts]);
        }
        setHasMore(newPosts.length === 10);
      }
    } catch (error) {
      console.error('加载帖子列表失败', error);
    } finally {
      setLoading(false);
    }
  };

  const handleFilterChange = (newOrder: 'time' | 'score', newCommunity?: number) => {
    // 只有当值真正改变时才更新状态并重新加载
    if (newOrder !== order || newCommunity !== selectedCommunity) {
      const params: Record<string, string> = { order: newOrder };
      if (newCommunity !== undefined) {
        params.community_id = String(newCommunity);
      }
      setSearchParams(params);
      setPage(1);
      setPosts([]);
    }
  };

  const loadMore = () => {
    if (!loading && hasMore) {
      setPage((prev) => prev + 1);
    }
  };

  return (
    <div className="space-y-6">
      {/* 筛选栏 */}
      <div className="card p-6">
        <div className="flex flex-wrap items-center gap-4">
          <div className="flex items-center space-x-2">
            <span className="text-gray-600 font-medium">排序：</span>
            <button
              onClick={() => {
                if (order !== 'time') {
                  handleFilterChange('time', selectedCommunity);
                }
              }}
              className={`px-4 py-2 rounded-lg transition-all ${
                order === 'time'
                  ? 'bg-blue-500 text-white shadow-md'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              时间
            </button>
            <button
              onClick={() => {
                if (order !== 'score') {
                  handleFilterChange('score', selectedCommunity);
                }
              }}
              className={`px-4 py-2 rounded-lg transition-all ${
                order === 'score'
                  ? 'bg-blue-500 text-white shadow-md'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              热度
            </button>
          </div>

          <div className="flex items-center space-x-2">
            <span className="text-gray-600 font-medium">社区：</span>
            <button
              onClick={() => handleFilterChange(order, undefined)}
              className={`px-4 py-2 rounded-lg transition-all ${
                selectedCommunity === undefined
                  ? 'bg-blue-500 text-white shadow-md'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              全部
            </button>
            {communities.map((community) => (
              <button
                key={community.id}
                onClick={() => handleFilterChange(order, community.id)}
                className={`px-4 py-2 rounded-lg transition-all ${
                  selectedCommunity === community.id
                    ? 'bg-blue-500 text-white shadow-md'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                {community.name}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* 帖子列表 */}
      <div className="space-y-4">
        {loading && posts.length === 0 ? (
          <div className="card p-12 text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
            <p className="mt-4 text-gray-500">加载中...</p>
          </div>
        ) : posts.length === 0 ? (
          <div className="card p-12 text-center">
            <p className="text-gray-500">暂无帖子</p>
          </div>
        ) : (
          posts.map((post) => (
            <Link
              key={post.id}
              to={`/post/${post.id}`}
              className="card p-6 hover:shadow-xl transition-all duration-200 block animate-slide-up"
            >
              <div className="flex items-start space-x-4">
                <div className="flex-1">
                  <h2 className="text-xl font-bold text-gray-900 mb-2 hover:text-blue-600 transition-colors">
                    {post.title}
                  </h2>
                  <p className="text-gray-600 line-clamp-2 mb-4">{post.content}</p>
                  <div className="flex items-center space-x-4 text-sm text-gray-500">
                    <span>{formatTime(post.create_time)}</span>
                    <span>•</span>
                    <span>社区 #{post.community_id}</span>
                  </div>
                </div>
              </div>
            </Link>
          ))
        )}

        {hasMore && posts.length > 0 && (
          <div className="text-center pt-4">
            <button
              onClick={loadMore}
              disabled={loading}
              className="btn-secondary disabled:opacity-50"
            >
              {loading ? '加载中...' : '加载更多'}
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

