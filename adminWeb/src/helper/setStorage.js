export default (type = true, token, user) => {
  if (type) {
    localStorage.setItem('token', token)
    localStorage.setItem('user', user)
  }else{
    localStorage.clear()
  }
}