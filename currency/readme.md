To run on clience
grpcurl --plaintext localhost:9092 list
grpcurl --plaintext localhost:9092 Currency/SubscribeRate

grpcurl --plaintext -d @ localhost:9092 Currency/SubscribeRate 
{"Base":"CAD","Destination":"CAD"}