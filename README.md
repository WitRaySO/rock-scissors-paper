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



##
GET Method
/getAllUsers

##
GET Method
/ladderboard

##
GET Method
/user/ชื่อศัตรู/comparation
ส่ง JSON Body
{
"username" : "",
}
