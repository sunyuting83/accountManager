const CROSUrl = 'https://crossorigin.me/'
const RootU = 'http://localhost:13003'
const IMGUri = 'http://localhost:13005/image/'
const RootUrl = `${RootU}/api/v1/`
const IndexUrl = `${CROSUrl}${RootUrl}`
const Name = '帐号'
const Limit = 100
const GlobalTitle = `${Name}-后台管理`
const ImagesRoot = 'https://pic7.58cdn.com.cn/nowater/webim/'
const images = [
  `${ImagesRoot}big/n_v26a171fb1a3394f2abcfce3e1d0e2b662.jpg`,
  `${ImagesRoot}big/n_v224cd0671e0f4483d95f395494dd3a891.jpg`,
  `${ImagesRoot}small/n_v224cd0671e0f4483d95f395494dd3a891.jpg`,
  `https://image.suning.cn/uimg/ZR/share_order/158501870837440052.jpg`,
  `${ImagesRoot}big/n_v25205eb943f014624a20825a567ec7802.jpg`
]
const Api = {
  'login': `${RootUrl}Login`,
  'checklogin': `${RootUrl}CheckLogin`,
  'repassword': `${RootUrl}RePassword`,
  'projectList': `${RootUrl}ProjectsList`,
  'accountList': `${RootUrl}AccountList`
}

const MakeAccountListUri = (key) => {
  return `${RootUrl}${key}/AccountList`
}

const makePopeData = (e, message, active = true)=> {
  const target = e.target
  const data = {
    active: active,
    message: message,
    width: target.clientWidth,
    height: target.clientHeight,
    top: target.offsetTop,
    left: target.offsetLeft
  }
  return data
}

export default {
  IndexUrl,
  RootUrl,
  CROSUrl,
  GlobalTitle,
  images,
  Api,
  makePopeData,
  RootU,
  Limit,
  MakeAccountListUri,
  IMGUri
}