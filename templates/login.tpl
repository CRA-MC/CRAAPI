<!DOCTYPE html>
<html lang="zh-CN">

<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <link rel="stylesheet" href="https://cdn.staticfile.net/twitter-bootstrap/5.1.1/css/bootstrap.min.css" />
  <title>CRA-STUDiO - 登录</title>
</head>

<body>
  <div style="
        display: flex;
        align-items: center;
        position: absolute;
        top: 0;
        left: 0;
        bottom: 0;
        right: 0;
        justify-content: center;
      ">
    <div style="text-align: center">
      <h3 style="">CRA-STUDiO</h3>
      {{if .}} <h3 style="color: red">登录失败，用户名或者密码错误</h3> {{end}}
      <form action="/login" method="post">
        <div style="margin-top: 4vh;">
          <input style="width: 20rem;" type="text" name="username" placeholder="请输入用户名" />
        </div>
        <div style="margin-top: 1vh;">
          <input style="width: 20rem;" type="password" name="password" placeholder="请输入密码" />
        </div>
        <button class="btn btn-primary" style="width: 20rem;margin-top: 1vh;margin-bottom: 5rem" type="submit">
          登录
        </button>
      </form>
    </div>
  </div>
  <script src="https://cdn.staticfile.net/popper.js/2.9.3/umd/popper.min.js"></script>
  <script src="https://cdn.staticfile.net/twitter-bootstrap/5.1.1/js/bootstrap.min.js"></script>
</body>

</html>