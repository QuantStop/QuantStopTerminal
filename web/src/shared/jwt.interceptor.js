import axios from "axios";
import * as appRouter from '../router'
import {userStore} from "../store/userStore";
const jwtInterceptor = axios.create({});
const defaultUserProfile = {
    id: 0,
    username: "",
    roles: "",
}

jwtInterceptor.interceptors.request.use((config) => {
  return config;
});

jwtInterceptor.interceptors.response.use((response) => {
        return response;
    },
    async (error) => {
        if (error.response.status === 401) {
            /*let response = await axios.get(
                "/api/refresh-token",
                {
                    withCredentials: true,
                    credentials: "include",
                    headers: {
                        'X-Requested-With': 'XMLHttpRequest' // CSRF prevention
                    },
                }
            )
            .catch((err) => {
                return Promise.reject(err);
            });
            if(response && response.data){
                return axios(error.config);
            }
            else {
                return Promise.reject(error);
            }*/



            //appRouter.push('/login')
            userStore.setLogout(true)
            userStore.setIsAuthed(false)
            userStore.setUserProfile(defaultUserProfile)


        } else {
            return Promise.reject(error);
        }
    }
);

export default jwtInterceptor;
