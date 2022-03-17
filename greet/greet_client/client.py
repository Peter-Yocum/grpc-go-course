
import grpc
import greet_pb2 as greet_pb2
import greet_pb2_grpc as greet_pb2_grpc

def Greet(stub: greet_pb2_grpc.GreetServiceStub):
    req = greet_pb2.GreetRequest(
        greeting=greet_pb2.Greeting(first_name='Peter', last_name='Yocum')
    )

    result = stub.Greet(req)
    print(result)


def main():
    channel = grpc.insecure_channel('localhost:50051')
    greet_stub = greet_pb2_grpc.GreetServiceStub(channel)

    Greet(stub=greet_stub)


if __name__ == '__main__':
    main()
