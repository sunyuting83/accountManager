const CROSUrl = 'https://crossorigin.me/'
const RootU = 'http://localhost:13002'
const IMGUri = 'http://localhost:13005/image/'
const RootUrl = `${RootU}/admin/api/v1/`
const IndexUrl = `${CROSUrl}${RootUrl}`
const Name = '管理员帐号'
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
  'login': `${RootUrl}AdminLogin`,
  'checklogin': `${RootUrl}CheckLogin`,
  'repassword': `${RootUrl}RePassword`,
  'adminlist': `${RootUrl}AdminList`,
  'addadmin': `${RootUrl}AddAdmin`,
  'upstatus': `${RootUrl}UpStatus`,
  'deladmin': `${RootUrl}DelAdmin`,
  'userList': `${RootUrl}UserList`,
  'UsersAllList':`${RootUrl}UsersAllList`,
  'adduser': `${RootUrl}AddUser`,
  'upuserstatus': `${RootUrl}UpStatusUser`,
  'deluser': `${RootUrl}DelUser`,
  'reUserPassword': `${RootUrl}RePasswordUser`,
  'projectList': `${RootUrl}ProjectsList`,
  'userProjectList': `${RootUrl}UserProjectsList`,
  'addproject': `${RootUrl}AddProjects`,
  'delproject': `${RootUrl}DelProjects`,
  'upprojectstatus': `${RootUrl}UpStatusProjects`,
  'accountList': `${RootUrl}AccountList`,
  'UpdateProjects': `${RootUrl}UpdateProjects`,
  'AddGame': `${RootUrl}AddGame`,
  'DelGame': `${RootUrl}DelGame`,
  'GamesList': `${RootUrl}GamesList`,
  'GamesAllList':`${RootUrl}GamesAllList`,
  'setRemarks' : `${RootUrl}SetUserRemarks`,
  'DrawList' : `${RootUrl}DrawList`,
  'DrawData' : `${RootUrl}DrawData`,
  'AllCount' : `${RootUrl}AllCount`,
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
  IMGUri
}