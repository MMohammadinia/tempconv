import grpc from '@grpc/grpc-js';
import protoLoader from '@grpc/proto-loader';

const packageDef = protoLoader.loadSync('./proto/temp.proto');
const grpcObj = grpc.loadPackageDefinition(packageDef);
const tempPackage = grpcObj.TempConverter;

const client = new tempPackage('localhost:50051', grpc.credentials.createInsecure());

client.ConvertTemperature({ value: 100, from: 'C' }, (err, response) => {
  if (err) console.error(err);
  else console.log(response);
});