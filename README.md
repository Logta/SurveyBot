# SurveyBot

## 概要

Discord 上でアンケートをできるようにするための Bot

下記 URL で Bot を招待できます
https://discord.com/oauth2/authorize?client_id=868454195953561610&scope=bot&permissions=0

## Heroku へのデプロイ

- 基本的には main にマージされたら Heroku 側で変更が反映される

- うまくいかない場合は下記を実行
  - `heroku ps:scale woker=1`
