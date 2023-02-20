# gobasicRedisgRPC
This is my repository for learning gRPC concept.

TODO:
- Rename and fix env var in file app copy.env
- Generate RSA key (2048 bits)
- Save priv key a file and fix the file name in file .env
- Change func startGrpcServer/startGrpcServer for start ginDefaultServer/gRPCServer respectively.
- Add logger
- Need Recover for gRPCServer from panic
- Posts route: setEx some requests to redis. Could load mongo to redis at the begin?
- Little work for you: implement:
```
func (uc *AuthServiceImpl) SignInUser(*models.SignInInput) (*models.DBResponse, error)
```
