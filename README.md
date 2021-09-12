# rock-scissors-paper

## ใส่ชี่อ User สู่ Database (sign up)
RESTful verb : `PUT`

path : `/signup`

JSON body : `{"username" : ""}`

## ส่งค่า User ที่มีอยู่ใน Database กลับมาทุกคน
RESTful verb : `GET`

path : `/getAllUsers`

JSON body : `-`

## ท้าผู้เล่นอีกคนในระบบที่มีชื่ออยู่ใน Database 
RESTful verb : `POST`

path : `/user/:ชื่อเรา/invitation`

JSON body : `{ "username" : "ชื่อศัตรู","choice" : ""}`

## เปรียบเทียบเรากับ User คนอิ่น
RESTful verb : `GET`

path : `/user/:ชื่อเรา/comparison`

JSON body : `{"username" : "ชื่อศัตรู"}`

## Leaderboard
RESTful verb : `GET`

path : `/leaderboard`

JSON body : `-`
