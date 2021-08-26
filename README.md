# SurveyBot

## 概要
Discord上でアンケートをできるようにするためのBot

下記URLでBotを招待できます
https://discord.com/oauth2/authorize?client_id=868454195953561610&scope=bot&permissions=0

## Herokuへのデプロイ
- 基本的にはmainにマージされたらHeroku川で変更が反映される

- うまくいかない場合は下記を実行
  - `heroku ps:scale woker=1`
