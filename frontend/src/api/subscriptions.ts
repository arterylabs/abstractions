import { axiosInstance } from "./index";

const subscriptions = {
  getSubscriptions: async () => {
    return axiosInstance.get("/subscriptions");
  },
  validateSubscription: async () => {
    return axiosInstance.get("/subscriptions/validate");
  },
};

export default subscriptions;
