import { axiosInstance } from "./index";

interface SignupParams {
  email: string;
  password: string;
  name: string;
}

const auth = {
  loginWithGoogle: async () => {
    window.location.href = `${process.env.NEXT_PUBLIC_API_URL}/auth/google`;
  },
  getCurrentUser: async () => {
    return axiosInstance.get("/users/me");
  },
  logout: async () => {
    return axiosInstance.post("/auth/logout");
  },
  signup: async (data: SignupParams) => {
    return axiosInstance.post("/auth/signup", data);
  },
  login: async (email: string, password: string) => {
    return axiosInstance.post("/auth/login", { email, password });
  },
  verifyEmail: async (token: string) => {
    return axiosInstance.post(`/users/verify`, {
      token,
    });
  },
  resendVerificationEmail: async () => {
    return axiosInstance.post(`/users/verify/resend`);
  },
  resetPassword: async (email: string) => {
    return axiosInstance.post(`/users/password/reset`, { email });
  },
  resetPasswordConfirm: async (token: string, password: string, id: string) => {
    return axiosInstance.post(`/users/password/reset/confirm`, {
      token,
      password,
      id,
    });
  },
};

export default auth;
