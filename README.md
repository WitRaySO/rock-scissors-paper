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

path : `/user/:ชื่อศัตรู/invitation`

JSON body : `{ "username" : "","choice" : ""}`

## ท้าผู้เล่นอีกคนในระบบที่มีชื่ออยู่ใน Database 
RESTful verb : `GET`

path : `/user/:ชื่อศัตรู/comparison`

JSON body : `{"username" : ""}`

## ท้าผู้เล่นอีกคนในระบบที่มีชื่ออยู่ใน Database 
RESTful verb : `GET`

path : `/leadderboard`

JSON body : `-`
