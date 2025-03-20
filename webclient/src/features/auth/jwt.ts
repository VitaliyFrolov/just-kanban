import { jwtDecode } from 'jwt-decode';

export interface AuthUserClaims {
  id: string;
  email: string;
  avatar: string;
  username: string;
  first_name: string;
  last_name: string;
}

export const getUserClaims = (token: string): AuthUserClaims => {
  return jwtDecode<AuthUserClaims>(token);
};
