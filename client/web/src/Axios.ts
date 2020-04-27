import Axios from "axios";
import config from './Config.json';

let axios = Axios.create({
    baseURL: config.server.addr,
})
axios.defaults.withCredentials=true

export default axios