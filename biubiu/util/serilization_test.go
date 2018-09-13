package util

import (
	"testing"
	"encoding/gob"
	"os"
	"fmt"
	"bytes"
)

type User struct {
	Name string
	Function 	*func(string string)
}

func (this *User) Say() string {
	return this.Name + ` hello world ! `
}


//func (this *User) MarshalBinary() ([]byte, error) {
//	// A simple encoding: plain text.
//	var b bytes.Buffer
//	//var args []interface{}
//	//for _, value := range this.Values{
//	//	args = append(args, &value)
//	//}
//	//args = append(args, &this.Mutex)
//	//args = append(args)
//	a := this.a
//	id := this.id
//	name := this.name
//	a_bytes, err := json.Marshal(a)
//	id_bytes,err := json.Marshal(id)
//	name_bytes, err :=json.Marshal(name)
//
//	fmt.Println(id_bytes, "  ", name_bytes, "  ",a_bytes)
//	id_num, err   := b.Write(id_bytes)
//	name_num, err := b.Write(name_bytes)
//	a_num, err    := b.Write(a_bytes)
//	//fmt.Println("Marsh  result", res)
//	//_, err := fmt.Fprintln(&b, this.a)
//	if err != nil{
//		panic(err)
//	}
//	fmt.Println(id_num, name_num, a_num)
//	fmt.Println(b.Bytes())
//	return b.Bytes(), nil
//	//return res, nil
//}

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
//func (v *User) UnmarshalBinary(data []byte) error {
//	// A simple encoding: plain text.
//	b := bytes.NewBuffer(data)
//	//
//	//
//	//_, err := fmt.Fscanln(b, &v.a)
//	values := b.Bytes()
//	fmt.Println("values  ", values)
//	//err := json.Unmarshal(values[:0], &v.id)
//	var id, a int
//	err := json.Unmarshal(values[:0], &id)
//	//fmt.Println(v.id)
//	var name string
//	err = json.Unmarshal(values[1:6], &name)
//	//fmt.Println(values[1:7], v.name)
//	err = json.Unmarshal(values[7:], &a)
//	//fmt.Println(v.a)
//	//var a int
//	//fmt.Println("Name ", v.Name)
//	//err := json.Unmarshal(data, &a)
//	//v.a = a
//	//fmt.Println("a", a, v)
//	if err != nil{
//		panic(err)
//	}
//	v.id = id
//	v.a = a
//	v.name = name
//	return err
//}
func TestSerilization(t *testing.T){
//	序列化测试
//	利用gob进行序列化
	name := func(name_ string) {
		fmt.Println(name_)
	}
	user := User{Name: "Mike", Function: &name}
	fmt.Println(user)
	//user2 := User{Id: 3, Name: "Jack"}
	//u := []User{user, user2}
	file, err := os.Create("D://seritest.gob")
	if err != nil{
		fmt.Println(err)
		return
	}
	enc := gob.NewEncoder(file)
	err = enc.Encode(&user)
	if err != nil{
		fmt.Println(err)
		return
	}
}

func TestDeSerilization(t *testing.T) {
	file, err := os.Open("D://seritest.gob")
	if err != nil{
		fmt.Println(err)
		return
	}
	dec := gob.NewDecoder(file)
	var u User
	err = dec.Decode(&u)
	if err != nil{
		fmt.Println("反序列化失败", err)
		return
	}
	fmt.Println(u)
	//var wg sync.WaitGroup
	//
	//
	////for _, usr := range u{
	//	wg.Add(1)
	//	go func(user_ *User) {
	//		defer wg.Done()
	//		user_.Say()
	//	}(&u)
	////}
	//wg.Wait()
	fmt.Println(u.Say())
	name_tmp := u.Function
	(*name_tmp)("yes")
}
type Vector struct {
	x, y, z int
	B string
	Ch chan string
}

func (v Vector) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	fmt.Fprintln(&b, v.x, v.y, v.z, v.B, v.Ch)
	return b.Bytes(), nil
}

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (v *Vector) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &v.x, &v.y, &v.z, &v.B, &v.Ch)
	return err
}

func TestGob(t *testing.T){
	var network bytes.Buffer // Stand-in for the network.
	enc := gob.NewEncoder(&network)
	err := enc.Encode(Vector{x:3, y:4, z:5, B:"li", Ch:make(chan string)})
	if err != nil {
		fmt.Println("encode:", err)
	}

	// Create a decoder and receive a value.
	dec := gob.NewDecoder(&network)
	var v Vector
	err = dec.Decode(&v)
	if err != nil {
		fmt.Println("decode:", err)
	}
	fmt.Println(v)
}