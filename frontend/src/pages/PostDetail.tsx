import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { postApi } from '../utils/api';
import { formatTime } from '../utils/format';
import { PostDetail as IPostDetail } from '../types';
import VoteButton from '../components/VoteButton';

export default function PostDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [post, setPost] = useState<IPostDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (id) {
      loadPost(id);
    }
  }, [id]);

  const loadPost = async (postId: string) => {
    try {
      setLoading(true);
      const response = await postApi.getDetail(postId);
      if (response.data.code === 1000) {
        setPost(response.data.data);
      } else {
        setError(response.data.msg || '获取帖子详情失败');
      }
    } catch (err: any) {
      setError(err.response?.data?.msg || '获取帖子详情失败');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center py-12">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (error || !post) {
    return (
      <div className="card p-8 text-center">
        <p className="text-red-500 text-lg mb-4">{error || '帖子不存在'}</p>
        <button onClick={() => navigate('/')} className="btn-primary">
          返回首页
        </button>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Back Button */}
      <div className="flex justify-start">
        <button 
          onClick={() => navigate(-1)} 
          className="flex items-center text-gray-600 hover:text-blue-600 transition-colors"
        >
          <svg className="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          返回
        </button>
      </div>

      <div className="card p-0 overflow-hidden animate-slide-up">
        {/* Vote Section - Left Sidebar style for desktop */}
        <div className="flex">
          <div className="w-16 bg-gray-50 p-4 flex flex-col items-center border-r border-gray-100">
            <VoteButton 
              postId={post.id} 
              initialVoteNum={post.vote_num}
              initialVoteStatus={post.vote_status}
              onVoteChange={(newNum) => setPost({ ...post, vote_num: newNum })} 
            />
          </div>
          
          <div className="flex-1 p-6 md:p-8">
            <div className="flex items-center space-x-2 text-sm text-gray-500 mb-4">
              <span className="font-medium text-blue-600 hover:underline cursor-pointer">
                {post.community?.name}
              </span>
              <span>•</span>
              <span>由 {post.author_name} 发布</span>
              <span>•</span>
              <span>{formatTime(post.create_time)}</span>
            </div>

            <h1 className="text-2xl md:text-3xl font-bold text-gray-900 mb-6">
              {post.title}
            </h1>

            <div className="prose max-w-none text-gray-800 leading-relaxed whitespace-pre-wrap">
              {post.content}
            </div>
          </div>
        </div>
      </div>

      {/* Comments Section Placeholder - API doesn't specify comments endpoint yet */}
      <div className="card p-6">
        <h3 className="text-lg font-bold mb-4">评论</h3>
        <p className="text-gray-500 text-center py-8">暂无评论功能</p>
      </div>
    </div>
  );
}
