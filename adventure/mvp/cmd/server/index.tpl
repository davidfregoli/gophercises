<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Choose Your Own Adventure</title>
  <style>
    body {
      padding-top: 60px;
    }
    * {
      max-width: 600px;
      margin: 0 auto 20px auto;
    }
  </style>
</head>
<body>
  <h1>{{.Title}}</h1>
  {{range .Paragraphs}}
  <p>{{.}}</p>
  {{end}}
  <hr>
  <ul>
  {{range .Options}}
  <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
  {{else}}
  <li><a href="/intro">Start over</a></li>
  {{end}}
  </ul>
</body>
</html>
