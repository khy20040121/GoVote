import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { authApi } from '../utils/api';
import { setToken, setUser } from '../utils/auth';
import type { SignUpParams } from '../types';

interface SignUpProps {
  onSuccess?: () => void;
}

export default function SignUp({ onSuccess }: SignUpProps) {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<SignUpParams>({
    username: '',
    password: '',
    re_password: '',
  });
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (formData.password !== formData.re_password) {
      setError('两次输入的密码不一致');
      return;
    }

    setLoading(true);

    try {
      const response = await authApi.signup(formData);
      if (response.data.code === 1000) {
        // 注册成功后自动登录
        try {
          const loginResponse = await authApi.login({
            username: formData.username,
            password: formData.password,
          });
          if (loginResponse.data.code === 1000 && loginResponse.data.data) {
            setToken(loginResponse.data.data.token!);
            setUser(loginResponse.data.data);
            onSuccess?.();
            navigate('/');
            return;
          }
        } catch (loginErr) {
          console.warn('Auto login failed:', loginErr);
        }
        
        // 自动登录失败，跳转到登录页让用户手动登录
        navigate('/login');
      } else {
        setError(response.data.msg as string);
      }
    } catch (err: any) {
      setError(err.response?.data?.msg || '注册失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-[calc(100vh-200px)]">
      <div className="card p-8 w-full max-w-md animate-fade-in">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent mb-2">
            创建账户
          </h1>
          <p className="text-gray-500">加入 GoVote 社区</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          {error && (
            <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-xl">
              {error}
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              用户名
            </label>
            <input
              type="text"
              className="input-field"
              value={formData.username}
              onChange={(e) => setFormData({ ...formData, username: e.target.value })}
              required
              placeholder="请输入用户名"
              autoComplete="off"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              密码
            </label>
            <input
              type="password"
              className="input-field"
              value={formData.password}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
              required
              placeholder="请输入密码"
              autoComplete="new-password"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              确认密码
            </label>
            <input
              type="password"
              className="input-field"
              value={formData.re_password}
              onChange={(e) => setFormData({ ...formData, re_password: e.target.value })}
              required
              placeholder="请再次输入密码"
              autoComplete="new-password"
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="btn-primary w-full disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? '注册中...' : '注册'}
          </button>
        </form>

        <div className="mt-6 text-center text-sm text-gray-600">
          已有账户？{' '}
          <Link to="/login" className="text-blue-600 hover:text-blue-700 font-medium">
            立即登录
          </Link>
        </div>
      </div>
    </div>
  );
}

