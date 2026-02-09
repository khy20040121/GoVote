import { BrowserRouter, Routes, Route, Link, useNavigate } from 'react-router-dom';
import { useState, useEffect } from 'react';
import Login from './pages/Login';
import SignUp from './pages/SignUp';
import PostList from './pages/PostList';
import CreatePost from './pages/CreatePost';
import PostDetail from './pages/PostDetail';
import { getUser, logout } from './utils/auth';
import { User } from './types';

function NavBar() {
  const navigate = useNavigate();
  const [user, setUserState] = useState<User | null>(null);

  useEffect(() => {
    // 简单的用户状态同步，实际项目中可能需要 Context 或 Redux
    const checkUser = () => {
      setUserState(getUser());
    };
    
    checkUser();
    window.addEventListener('storage', checkUser);
    // 自定义事件用于组件间通信
    window.addEventListener('user-login', checkUser);
    window.addEventListener('user-logout', checkUser);
    
    return () => {
      window.removeEventListener('storage', checkUser);
      window.removeEventListener('user-login', checkUser);
      window.removeEventListener('user-logout', checkUser);
    };
  }, []);

  const handleLogout = () => {
    logout();
    setUserState(null);
    window.dispatchEvent(new Event('user-logout'));
    navigate('/login');
  };

  return (
    <nav className="bg-white shadow-sm sticky top-0 z-50">
      <div className="container mx-auto px-4 h-16 flex items-center justify-between">
        <Link to="/" className="text-2xl font-bold text-blue-600 hover:text-blue-700 transition-colors">
          GoVote
        </Link>

        <div className="flex items-center space-x-4">
          <Link to="/" className="text-gray-600 hover:text-blue-600 font-medium transition-colors">
            首页
          </Link>
          
          {user ? (
            <>
              <Link 
                to="/create-post" 
                className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors shadow-md hover:shadow-lg"
              >
                发布帖子
              </Link>
              <div className="flex items-center space-x-2 border-l pl-4 ml-2">
                <span className="text-gray-900 font-medium">{user.username}</span>
                <button 
                  onClick={handleLogout}
                  className="text-gray-500 hover:text-red-600 text-sm transition-colors"
                >
                  退出
                </button>
              </div>
            </>
          ) : (
            <div className="space-x-2">
              <Link 
                to="/login" 
                className="px-4 py-2 text-blue-600 hover:bg-blue-50 rounded-lg transition-colors font-medium"
              >
                登录
              </Link>
              <Link 
                to="/signup" 
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors shadow-md hover:shadow-lg font-medium"
              >
                注册
              </Link>
            </div>
          )}
        </div>
      </div>
    </nav>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <div className="min-h-screen bg-gray-50">
        <NavBar />
        <main className="container mx-auto px-4 py-8">
          <Routes>
            <Route path="/" element={<PostList />} />
            <Route path="/login" element={<Login onSuccess={() => window.dispatchEvent(new Event('user-login'))} />} />
            <Route path="/signup" element={<SignUp onSuccess={() => window.dispatchEvent(new Event('user-login'))} />} />
            <Route path="/create-post" element={<CreatePost />} />
            <Route path="/post/:id" element={<PostDetail />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  );
}
