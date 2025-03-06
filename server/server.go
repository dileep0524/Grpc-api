package main

import (
	"Desktop/Emp/proto"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	_"github.com/lib/pq"
)
type server struct{
	proto.UnimplementedUserServiceServer
	db *sql.DB
}
func (s *server) CreateUser(ctx context.Context,req *proto.CreateUserReq)(*proto.Response,error){
	_,err:=s.db.Query("insert into emps (id,name,department,salary)values($1,$2,$3,$4)",req.Id,req.Name,req.Department,req.Salary)
	if err!=nil{
		return nil,err
	}
	return &proto.Response{Message: "Created"},nil

}
func (s *server) GetUser(ctx context.Context,req *proto.UserReq)(*proto.User,error){
	var user proto.User
	err:=s.db.QueryRow("select * from emps where id=$1",req.Id).Scan(&user.Id,&user.Name,&user.Department,&user.Salary)
	if err!=nil{
		return nil,err
	}
	return &user,nil
}
func (s *server) UpdateSal(ctx context.Context,req *proto.UpdateReq)(*proto.Response,error){
	_,err:=s.db.Exec("UPDATE emps SET salary=$1 where id=$2",req.Salary,req.Id)
	if err!=nil{
		return nil,err
	}
	return &proto.Response{Message: "Updated"},nil
}
func (s *server) DeleteUser(ctx context.Context,req *proto.UserReq)(*proto.Response,error){
	_,err:=s.db.Exec("delete from emps where id=$1",req.Id)
	if err!=nil{
		return nil,err
	}
	return &proto.Response{Message: "deleted"},nil
}
func (s *server) GetUsers(ctx context.Context,req *proto.Empty)(*proto.UsersResponse,error){
	rows,err:=s.db.Query("select * from emps")
	if err!=nil{
		return nil,err
	}
	var userss []*proto.User
	for rows.Next(){
		var user proto.User
		err:=rows.Scan(&user.Id,&user.Name,&user.Department,&user.Salary)
		if err!=nil{
			return nil,err
		}
		userss=append(userss, &user)
	}
	
	return &proto.UsersResponse{Users: userss},nil
}
func (s *server) ApplyBonus(ctx context.Context,req *proto.BonusReq)(*proto.BonusResponse,error){
	var sal int32
	err:=s.db.QueryRow("select salary from emps where id=$1",req.Id).Scan(&sal)
	if err!=nil{
		return nil,err
	}
	newsal:=sal+(sal*req.Percent/100)
	_,err=s.db.Query("update emps set salary=$1 where id=$2",newsal,req.Id)
	if err!=nil{
		return nil,err
	}
	return &proto.BonusResponse{Message: "Updated Successfully",Newsal: newsal},nil
}
func main(){
	db,err:=sql.Open("postgres","host=localhost port=5432 user=mrpsycho password=kusuma123 dbname=cruddb sslmode=disable")
	if err!=nil{
		log.Fatalf("Error Connecting to the DB:%v",err)
	}
	lis,err:=net.Listen("tcp",":50051")
	if err!=nil{
		log.Fatalf("Not Listening")
	}
	GrpcServer:=grpc.NewServer()
	proto.RegisterUserServiceServer(GrpcServer,&server{db: db})
	fmt.Println("Server Running")
	err=GrpcServer.Serve(lis)
	if err!=nil{
		log.Fatalf("Not Serving")
	}
}