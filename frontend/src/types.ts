export interface ApiResponse<T> {
  code: number;
  msg: string;
  data: T;
}

export interface User {
  user_id: string;
  username: string;
  token?: string;
}

export interface Community {
  id: number;
  name: string;
}

export interface CommunityDetail extends Community {
  introduction: string;
  create_time: string;
}

export interface Post {
  id: string;
  author_id: string;
  community_id: number;
  status: number;
  title: string;
  content: string;
  create_time: string;
}

export interface PostDetail extends Post {
  author_name: string;
  vote_num: number;
  vote_status?: number;
  community: Community;
}

export interface LoginParams {
  username: string;
  password: string;
}

export interface SignUpParams {
  username: string;
  password: string;
  re_password: string;
}

export interface CreatePostParams {
  title: string;
  content: string;
  community_id: number;
}

export interface VoteParams {
  post_id: string;
  direction: 1 | 0 | -1;
}

export interface PostListParams {
  page?: number;
  size?: number;
  order?: 'time' | 'score';
  community_id?: number;
}
