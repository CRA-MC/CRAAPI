<!DOCTYPE html>
<html lang="zh-CN">

<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <link rel="stylesheet" href="https://cdn.staticfile.net/twitter-bootstrap/5.1.1/css/bootstrap.min.css" />
  <title>CRA-STUDiO - 注册</title>
</head>
<body>
  <script src="https://cdn.staticfile.net/jquery/1.10.2/jquery.min.js">
    function is_username() {
      var username = document.getElementById("username");
    };

    function is_email() {
      var email = document.getElementById("email");
    };
  </script>
  <div>
    <h3>CRA-STOUDiO</h3>
    <form action="/register" method="post">
      <div><input type="text" placeholder="用户名" name="username" id="username" onblur="is_username()"></div>
  </div>
  <div><input type="text" placeholder="邮箱" name="email" id="email" onblur="is_email()"></div>
  <div>
    <input type="text" placeholder="验证码">
    <button class="btn btn-primary">获取验证码</button>
  </div>
  <input type="text" placeholder="密码" name="password">
  <input type="text" placeholder="确认密码" name="password_r">
  <div><button type="submit" class="btn btn-primary">注册</button></div>
  </form>
  </div>
  <script src="https://cdn.staticfile.net/popper.js/2.9.3/umd/popper.min.js"></script>
  <script src="https://cdn.staticfile.net/twitter-bootstrap/5.1.1/js/bootstrap.min.js"></script>
</body>
</html>
