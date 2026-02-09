import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { postApi, communityApi } from '../utils/api';
import { Community } from '../types';

export default function CreatePost() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [communities, setCommunities] = useState<Community[]>([]);
  const [formData, setFormData] = useState({
    title: '',
    content: '',
    community_id: 0,
  });
  const [error, setError] = useState('');

  useEffect(() => {
    loadCommunities();
  }, []);

  const loadCommunities = async () => {
    try {
      const response = await communityApi.getList();
      if (response.data.code === 1000) {
        setCommunities(response.data.data);
        if (response.data.data.length > 0) {
          setFormData(prev => ({ ...prev, community_id: response.data.data[0].id }));
        }
      }
    } catch (error) {
      console.error('加载社区列表失败', error);
      setError('加载社区列表失败');
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.title || !formData.content || !formData.community_id) {
      setError('请填写完整信息');
      return;
    }

    setLoading(true);
    setError('');

    try {
      const response = await postApi.create({
        title: formData.title,
        content: formData.content,
        community_id: Number(formData.community_id),
      });

      if (response.data.code === 1000) {
        navigate('/');
      } else {
        setError(response.data.msg || '发布失败');
      }
    } catch (err: any) {
      setError(err.response?.data?.msg || '发布失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto">
      <div className="card p-8 animate-fade-in">
        <h1 className="text-2xl font-bold mb-6">发布新帖子</h1>
        
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-xl mb-6">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              选择社区
            </label>
            <select
              className="input-field"
              value={formData.community_id}
              onChange={(e) => setFormData({ ...formData, community_id: Number(e.target.value) })}
              required
            >
              {communities.map((community) => (
                <option key={community.id} value={community.id}>
                  {community.name}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              标题
            </label>
            <input
              type="text"
              className="input-field"
              value={formData.title}
              onChange={(e) => setFormData({ ...formData, title: e.target.value })}
              placeholder="请输入标题"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              内容
            </label>
            <textarea
              className="input-field min-h-[200px]"
              value={formData.content}
              onChange={(e) => setFormData({ ...formData, content: e.target.value })}
              placeholder="请输入内容"
              required
            />
          </div>

          <div className="flex justify-end space-x-4">
            <button
              type="button"
              onClick={() => navigate('/')}
              className="px-6 py-2 rounded-lg bg-gray-100 text-gray-600 hover:bg-gray-200 transition-colors"
            >
              取消
            </button>
            <button
              type="submit"
              disabled={loading}
              className="btn-primary px-8 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? '发布中...' : '发布'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
