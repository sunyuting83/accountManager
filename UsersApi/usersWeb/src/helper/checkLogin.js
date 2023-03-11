import Config from '@/helper/config'
import Fetch from '@/helper/fetch'
export default async () => {
  const token = localStorage.getItem("token")
  const data = await Fetch(Config.Api.checklogin, {}, "get", token)
  if (data.status === 0) {
    return 0
  }else{
    localStorage.removeItem("token")
    localStorage.removeItem("user")
    return 1
  }
}