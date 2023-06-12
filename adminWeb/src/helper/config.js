const CROSUrl = 'https://crossorigin.me/'
const RootU = 'http://localhost:13002'
const IMGUri = 'http://localhost:18080/image/'
const RootUrl = `${RootU}/admin/api/v1/`
const IndexUrl = `${CROSUrl}${RootUrl}`
const Name = '管理员帐号'
const Limit = 100
const GlobalTitle = `${Name}-后台管理`
// const ImagesRoot = 'https://pic7.58cdn.com.cn/nowater/webim/'
const images = [
  ``,
  `https://kjimg10.360buyimg.com/ott/jfs/t20250611/133512/21/38815/17870/64872694F64c3f48a/42b7ef7069520cd9.png`,
  ``,
  `https://image.suning.cn/uimg/ZR/share_order/158501870837440052.jpg`,
  `https://img.pddpic.com/goods_mms/2022-08-17/3750025f-3037-4a2d-8ab2-33806e59aa86.jpeg`
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
  'ModifyCalc': `${RootUrl}ModifyCalc`,
  'GamesList': `${RootUrl}GamesList`,
  'GamesAllList':`${RootUrl}GamesAllList`,
  'setRemarks': `${RootUrl}SetUserRemarks`,
  'DrawList': `${RootUrl}DrawList`,
  'DrawData': `${RootUrl}DrawData`,
  'AllCount': `${RootUrl}AllCount`,
  'drawSelect': `${RootUrl}DrawSelect`,
  'drawSelectPull': `${RootUrl}DrawSelectPull`,
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