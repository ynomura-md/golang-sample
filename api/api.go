package api

import "go.uber.org/zap/zapcore"

type User struct {
    Name string
    Age  int
}


func (u User) MarshalLogObject(enc zapcore.ObjectEncoder) error {
    enc.AddString("name", u.Name)
    enc.AddInt("age", u.Age)
    return nil
}
