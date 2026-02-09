import { useState } from 'react';
import { voteApi } from '../utils/api';
import { isAuthenticated } from '../utils/auth';

interface VoteButtonProps {
  postId: string;
  initialVoteNum: number;
  initialVoteStatus?: number;
  onVoteChange?: (newVoteNum: number) => void;
}

export default function VoteButton({ postId, initialVoteNum, initialVoteStatus = 0, onVoteChange }: VoteButtonProps) {
  const [voteNum, setVoteNum] = useState(initialVoteNum);
  const [voting, setVoting] = useState(false);
  const [userVote, setUserVote] = useState<number>(initialVoteStatus);

  const handleVote = async (direction: 1 | -1 | 0) => {
    if (!isAuthenticated()) {
      alert('请先登录');
      return;
    }

    if (voting) return;

    // 如果点击的是当前状态，则取消投票
    const newDirection = userVote === direction ? 0 : direction;
    
    setVoting(true);
    try {
      const response = await voteApi.vote({ post_id: postId, direction: newDirection });
      
      // 检查响应是否成功
      if (response.data.code === 1000) {
        // 计算新的投票数
        let newVoteNum = voteNum;
        if (userVote === 1 && newDirection === 0) {
          newVoteNum -= 1; // 取消赞成
        } else if (userVote === -1 && newDirection === 0) {
          newVoteNum += 1; // 取消反对
        } else if (userVote === 0) {
          newVoteNum += newDirection; // 新投票
        } else if (userVote === 1 && newDirection === -1) {
          newVoteNum -= 2; // 从赞成改为反对
        } else if (userVote === -1 && newDirection === 1) {
          newVoteNum += 2; // 从反对改为赞成
        }
        
        setVoteNum(newVoteNum);
        setUserVote(newDirection);
        onVoteChange?.(newVoteNum);
      } else {
        // 如果返回了错误码，显示错误消息
        const msg = response.data.msg || '投票失败';
        alert(msg);
        // 如果是token相关错误，跳转到登录页
        if (response.data.code === 1008 || response.data.code === 1009) {
          window.location.href = '/login';
        }
      }
    } catch (error: any) {
      // 处理网络错误或其他异常
      const msg = error.response?.data?.msg || error.message || '投票失败';
      alert(msg);
      // 如果是token相关错误，跳转到登录页
      if (error.response?.data?.code === 1008 || error.response?.data?.code === 1009) {
        window.location.href = '/login';
      }
    } finally {
      setVoting(false);
    }
  };

  return (
    <div className="flex flex-col items-center space-y-2">
      <button
        onClick={() => handleVote(1)}
        disabled={voting}
        className={`w-10 h-10 rounded-lg transition-all duration-200 ${
          userVote === 1
            ? 'bg-blue-500 text-white shadow-md'
            : 'bg-gray-100 text-gray-600 hover:bg-blue-100 hover:text-blue-600'
        } ${voting ? 'opacity-50 cursor-not-allowed' : ''}`}
      >
        <svg className="w-5 h-5 mx-auto" fill="currentColor" viewBox="0 0 20 20">
          <path fillRule="evenodd" d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z" clipRule="evenodd" />
        </svg>
      </button>
      <span className={`text-lg font-semibold ${voteNum > 0 ? 'text-blue-600' : voteNum < 0 ? 'text-red-600' : 'text-gray-600'}`}>
        {voteNum}
      </span>
      <button
        onClick={() => handleVote(-1)}
        disabled={voting}
        className={`w-10 h-10 rounded-lg transition-all duration-200 ${
          userVote === -1
            ? 'bg-red-500 text-white shadow-md'
            : 'bg-gray-100 text-gray-600 hover:bg-red-100 hover:text-red-600'
        } ${voting ? 'opacity-50 cursor-not-allowed' : ''}`}
      >
        <svg className="w-5 h-5 mx-auto" fill="currentColor" viewBox="0 0 20 20">
          <path fillRule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clipRule="evenodd" />
        </svg>
      </button>
    </div>
  );
}

