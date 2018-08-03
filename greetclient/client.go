package main

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
	"io"
	"fmt"
	"grpc-practice/greet/greetpb"
	"log"
	"context"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("hello world")
	cc, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connection failed %v", err)
	}
	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)
	//doUnary(c)
	//doServerStreem(c)
	//doClientStreem(c)
	//dobidistream(c)
	//doErrorUnary(c)
	doUnaryWithDeadline(c,5*time.Second)	
	doUnaryWithDeadline(c,1*time.Second)
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient,seconds time.Duration){
	req:=&greetpb.GreetWithDeadRequest{
		Greeting: &greetpb.Greeting{
			FirstName:"Uday",
			SecondName:"Singh",
		},
	}
	ctx,cancel:=context.WithTimeout(context.Background(),seconds)
	defer cancel()

	res,err := c.GreetWithDeadline(ctx, req)
	if err!=nil{
		statusErr,ok:=status.FromError(err)
		if ok{
			if statusErr.Code() == codes.DeadlineExceeded{
				fmt.Println("Deadline was exceeded")
			}else{
				fmt.Println("Unexpected error",statusErr)
			}
			return
		}else{
		log.Fatalf("client failed %v",err)
		}
	}
	fmt.Println("the result is %v",res.Result)
}

func doErrorUnary(c greetpb.GreetServiceClient){
	number:=int64(-1)
	res,err:=c.SquarRoot(context.Background(),&greetpb.SquareRootRequest{
		Number:number,
	})
	if err!=nil{
		resError,ok:=status.FromError(err)
		if ok{
			fmt.Println(resError.Message())
			fmt.Println(resError.Code())
			if resError.Code()==codes.InvalidArgument{
				fmt.Println("Negative number sent to the server")
			}
		}else{
			log.Fatalf("Framework error thrown %v",err)
		}
	}else{
		fmt.Println("The square root of %d is %f",number,res.GetNumberRoot())
	}	
}

func doUnary(c greetpb.GreetServiceClient){
	req:=&greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName:"Uday",
			SecondName:"Singh",
		},
	}
	res,err := c.Greet(context.Background(), req)
	if err!=nil{
		log.Fatalf("client failed %v",err)
	}
	fmt.Println("the result is %v",res.Result)
}

func doServerStreem(c greetpb.GreetServiceClient){
	req:=&greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName:"Uday",
			SecondName:"Singh",
		},
	}
	resStream,err := c.GreetManyTimes(context.Background(), req)
	if err!=nil{
		log.Fatalf("client failed %v",err)
	}
	for{
		msg,err:=resStream.Recv()
		if err==io.EOF{
			break;
		}
		if err!=nil{
			log.Fatalf("Streaming error %v",err)
		}
		log.Println("Response ",msg.GetResult())
	}
}

func doClientStreem(c greetpb.GreetServiceClient){
	
req:=[]*greetpb.LongGreetrequest{
	&greetpb.LongGreetrequest{
		Greeting: &greetpb.Greeting{
			FirstName:"Uday",
			SecondName:"Singh",
		},
	},
	&greetpb.LongGreetrequest{
		Greeting: &greetpb.Greeting{
			FirstName:"Vidhu",
			SecondName:"Sharma",
		},
	},
	&greetpb.LongGreetrequest{
		Greeting: &greetpb.Greeting{
			FirstName:"Aditi",
			SecondName:"Gautam",
		},
	},
	&greetpb.LongGreetrequest{
		Greeting: &greetpb.Greeting{
			FirstName:"Ruchi",
			SecondName:"Sharma",
		},
	},
	&greetpb.LongGreetrequest{
		Greeting: &greetpb.Greeting{
			FirstName:"Deepika",
			SecondName:"Arora",
		},
	},
	&greetpb.LongGreetrequest{
		Greeting: &greetpb.Greeting{
			FirstName:"Ritika",
			SecondName:"Sharma",
		},
	},
}
stream,err:=c.LongGreet(context.Background())
if err!=nil{
	log.Fatalf("Error creating stream")
}
for _,request := range req{
	fmt.Println("Request is ",request)
	stream.Send(request)
	time.Sleep(1000*time.Millisecond)

}
res,err:=stream.CloseAndRecv()
if err!=nil{
	log.Fatalf("Error in stream",err)
}
fmt.Println("output",res)
}

func dobidistream(c greetpb.GreetServiceClient){
	stream,err:=c.GreetEveryone(context.Background())
	if err!=nil{
		log.Fatalf("Stream not being created %v",err)
	}
		
req:=[]*greetpb.GreetEveryoneRequest{
	&greetpb.GreetEveryoneRequest{
		Greeting: &greetpb.Greeting{
			FirstName:"Uday",
			SecondName:"Singh",
		},
	},
	&greetpb.GreetEveryoneRequest{
		Greeting: &greetpb.Greeting{
			FirstName:"Vidhu",
			SecondName:"Sharma",
		},
	},
	&greetpb.GreetEveryoneRequest{
		Greeting: &greetpb.Greeting{
			FirstName:"Aditi",
			SecondName:"Gautam",
		},
	},
	&greetpb.GreetEveryoneRequest{
		Greeting: &greetpb.Greeting{
			FirstName:"Ruchi",
			SecondName:"Sharma",
		},
	},
	&greetpb.GreetEveryoneRequest{
		Greeting: &greetpb.Greeting{
			FirstName:"Deepika",
			SecondName:"Arora",
		},
	},
	&greetpb.GreetEveryoneRequest{
		Greeting: &greetpb.Greeting{
			FirstName:"Ritika",
			SecondName:"Sharma",
		},
	},
}
	waitc:=make(chan struct{})
	go func(){
		for _,result := range req{
			fmt.Println("Request= ",result)
			stream.Send(result)
			time.Sleep(1000*time.Millisecond)
		}
		stream.CloseSend()
	}()
	go func(){
		for{
			response,err:=stream.Recv()
			if err==io.EOF{
				break
			}
			if err!=nil{
				log.Fatalf("Error collecting response %v",err)
				break
			}
			log.Println("the server says =",response.GetResult())
		}
		close(waitc)
	}()
	<-waitc
}