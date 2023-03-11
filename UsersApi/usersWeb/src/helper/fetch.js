export default async (url = '', params = {}, method = 'GET', token = '') => {
  method = method.toUpperCase()
  // 此处规定get请求的参数使用时放在data中，如同post请求
  if (method === 'GET') {
    let dataStr = ''
    Object.keys(params).forEach(key => {
      dataStr += key + '=' + params[key] + '&'
    })

    if (dataStr !== '') {
      dataStr = dataStr.substr(0, dataStr.lastIndexOf('&'))
      url = url + '?' + dataStr
    }
  }

  let requestConfig = {
    method: method,
  }

  if (method === 'POST' || method === 'PUT' || method === 'DELETE') {
    requestConfig.headers = {
      Accept: '*/*',
      'Content-Type': 'application/json;charset=UTF-8',
    }
    Object.defineProperty(requestConfig, 'body', {
      value: JSON.stringify(params)
    })
  }
  if (token !== null && token.length > 8) {
    requestConfig.headers = new Headers({
      Accept: '*/*',
      'Content-Type': 'application/json;charset=UTF-8',
    })
    requestConfig.headers.append('Authorization',`Bearer ${token}`)
  }
  return new Promise((resolve) => {
    fetch(url, requestConfig)
      .then(res => {
        if(res.ok) {
          return res.json()
        }else {
          resolve({
            status: 1,
            message: "访问出错"
          })
        }
      })
      .then(json => {
        resolve(json)
      })
      .catch((err) => {
        resolve({
          status: 1,
          message: err.message
        })
      })
  })
}