package main

import (
	"math"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"time"
	"strconv"
	"fmt"
	"grpc-practice/greet/greetpb"
	"log"
	"net"
	"context"

	"google.golang.org/grpc"
)

type server struct{}

func (* server) Greet(ctx context.Context,req *greetpb.GreetRequest) (*greetpb.GreetResponse, error){
	firstName:=req.GetGreeting().GetFirstName()
	result:="hello"+firstName
	res:=&greetpb.GreetResponse{
		Result:result,
	}
	return res,nil
}
func (* server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadRequest) (*greetpb.GreetWithDeadResponse, error){
	for i:=0;i<3;i++ {
		if ctx.Err()==context.Canceled{
			fmt.Println("the client canceled the request")
			return nil,status.Error(codes.DeadlineExceeded,"the client")
		}
		time.Sleep(1000*time.Millisecond)
	}
	firstName:=req.GetGreeting().GetFirstName()
	result:="hello"+firstName
	res:=&greetpb.GreetWithDeadResponse{
		Result:result,
	}
	return res,nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest,stream greetpb.GreetService_GreetManyTimesServer) error{
	firstName:=req.GetGreeting().GetFirstName()
	secondName:=req.GetGreeting().GetSecondName()
	for i:=0;i<10;i++ {
		reply:="hello"+firstName+" "+secondName+" "+strconv.Itoa(i)
		res:=&greetpb.GreetManyTimesResponse{
			Result:reply,
		}
		stream.Send(res)
		time.Sleep(1000*time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) ( error){
	result:=""
	for{
		res,err:=stream.Recv()
		if err==io.EOF{
			return stream.SendAndClose(&greetpb.LongGreetresponse{
				Result:result,
			})
		}
		if err!=nil{
			log.Fatalf("Error sending stream %v",err)
		}
		result+="Hello "+res.GetGreeting().GetFirstName()+" "+res.GetGreeting().GetSecondName()+"!/n"
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error{
	result:=""
	for{
		request,err:=stream.Recv()
		if err==io.EOF{
			return nil
		}
		if err!=nil{
			log.Fatalf("Request can not be collected %v",err)
		}
		result="hello"+" "+request.GetGreeting().GetFirstName()+" "+request.GetGreeting().GetSecondName()+"!!!"
		err=stream.Send(
			&greetpb.GreetEveryoneresponse{
				Result:result,
			},
		)
		if err!=nil{
			log.Fatalf("Error while sending response %v",err)
			return err
		}
	}
}

func (*server) SquarRoot(ctx context.Context, req *greetpb.SquareRootRequest) (*greetpb.SquareRootResponse, error){
	num:=req.GetNumber()
	if num<0{
		return nil,status.Errorf(codes.InvalidArgument,fmt.Sprintf("Recieved negative number %v",num))
	}
	return &greetpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(num)),
	},nil
}
func main() {
	fmt.Println("Hello server")
	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}
