export default (type = true, token, user, id) => {
  if (type) {
    localStorage.setItem('token', token)
    localStorage.setItem('user', user)
    localStorage.setItem('userid', id)
  }else{
    localStorage.clear()
  }
}