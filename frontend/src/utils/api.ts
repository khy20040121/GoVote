import axios, { AxiosResponse } from 'axios';
import type { 
  ApiResponse, 
  User, 
  Post, 
  PostDetail, 
  Community, 
  CommunityDetail,
  LoginParams,
  SignUpParams,
  CreatePostParams,
  VoteParams,
  PostListParams
} from '../types';

const api = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器：添加 token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// 响应拦截器：处理错误
api.interceptors.response.use(
  (response) => {
    // 检查响应体中的code字段，如果是token相关错误，清除token
    // 让组件自己处理错误响应，不要在这里直接跳转
    if (response.data?.code === 1008 || response.data?.code === 1009) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    }
    return response;
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/';
    }
    return Promise.reject(error);
  }
);

export const authApi = {
  login: (params: LoginParams): Promise<AxiosResponse<ApiResponse<User>>> => 
    api.post('/login', params),
  
  signup: (params: SignUpParams): Promise<AxiosResponse<ApiResponse<null>>> => 
    api.post('/signup', params),
};

export const postApi = {
  getList: (params?: PostListParams): Promise<AxiosResponse<ApiResponse<Post[]>>> => 
    api.get('/posts2', { params }),
  
  getDetail: (id: string): Promise<AxiosResponse<ApiResponse<PostDetail>>> => 
    api.get(`/post/${id}`),
  
  create: (params: CreatePostParams): Promise<AxiosResponse<ApiResponse<null>>> => 
    api.post('/post', params),
};

export const communityApi = {
  getList: (): Promise<AxiosResponse<ApiResponse<Community[]>>> => 
    api.get('/community'),
  
  getDetail: (id: number): Promise<AxiosResponse<ApiResponse<CommunityDetail>>> => 
    api.get(`/community/${id}`),
};

export const voteApi = {
  vote: (params: VoteParams): Promise<AxiosResponse<ApiResponse<null>>> => 
    api.post('/vote', params),
};

